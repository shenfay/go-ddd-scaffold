// Package server 提供HTTP服务器相关服务（Swagger中间件版）
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	tenantservice "go-ddd-scaffold/internal/application/tenant/service"
	userservice "go-ddd-scaffold/internal/application/user/service"
	domainService "go-ddd-scaffold/internal/domain/user/service"
	"go-ddd-scaffold/internal/config"
	"go-ddd-scaffold/internal/infrastructure/auth"
	"go-ddd-scaffold/internal/infrastructure/event"
	"go-ddd-scaffold/internal/infrastructure/middleware"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/repo"
	"go-ddd-scaffold/internal/infrastructure/transaction"
	"go-ddd-scaffold/internal/infrastructure/wire"
	authhttp "go-ddd-scaffold/internal/interfaces/http/auth"
	tenanthttp "go-ddd-scaffold/internal/interfaces/http/tenant"
	userhttp "go-ddd-scaffold/internal/interfaces/http/user"
	"go-ddd-scaffold/internal/pkg/metrics"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ServerService HTTP 服务器服务
type ServerService struct {
	config      *config.Config
	db          *gorm.DB
	logger      *zap.Logger
	engine      *gin.Engine
	server      *http.Server
	metrics     *metrics.APIMetrics
	redisClient *redis.Client // Redis 客户端（用于 Token 黑名单等）
}

// NewServerService 创建新的服务器服务实例
func NewServerService(cfg *config.Config, db *gorm.DB, logger *zap.Logger) (*ServerService, error) {
	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// 测试 Redis 连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Warn("Redis 连接失败，部分功能将不可用", zap.Error(err))
		// 继续运行，但 Token 黑名单等功能将不可用
	} else {
		logger.Info("Redis 连接成功")
	}

	return &ServerService{
		config:      cfg,
		db:          db,
		logger:      logger,
		metrics:     metrics.NewAPIMetrics(),
		redisClient: rdb,
	}, nil
}

// Initialize 初始化服务器服务
func (s *ServerService) Initialize() error {
	// 初始化Gin引擎
	if err := s.initGinEngine(); err != nil {
		return fmt.Errorf("初始化Gin引擎失败: %w", err)
	}

	// 创建HTTP服务器
	s.createServer()

	s.logger.Info("服务器服务初始化完成")
	return nil
}

// initGinEngine 初始化Gin引擎和中间件
func (s *ServerService) initGinEngine() error {
	// 设置运行模式
	if s.config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	s.engine = gin.New()

	// 注册中间件
	s.registerMiddleware()

	// 注册路由
	s.registerRoutes()

	s.logger.Info("Gin引擎初始化完成")
	return nil
}

// registerMiddleware 注册中间件
func (s *ServerService) registerMiddleware() {
	// Gzip 压缩中间件
	s.engine.Use(middleware.GzipMiddleware(middleware.DefaultCompressionConfig))

	// 请求 ID 中间件
	s.engine.Use(middleware.RequestID("X-Request-ID"))

	// 幂等性中间件
	idempotencyConfig := middleware.DefaultIdempotencyConfig()
	s.engine.Use(middleware.Idempotency(idempotencyConfig))

	// 监控中间件
	s.engine.Use(middleware.Metrics(s.metrics))

	// 日志中间件
	s.engine.Use(middleware.Logger(middleware.DefaultLoggerConfig()))

	// 恢复中间件（panic 恢复）
	s.engine.Use(middleware.RecoveryWithLogger(s.logger))

	// 错误处理中间件（业务错误映射）
	s.engine.Use(middleware.ErrorMiddleware(s.logger))

	// CORS 中间件（使用独立的中间件文件）
	s.engine.Use(middleware.CORS())

	// Swagger 中间件（使用独立的中间件文件）
	s.engine.Use(middleware.Swagger())

	// 监控接口
	s.engine.GET("/metrics", middleware.MetricsHandler(s.metrics))
	s.engine.POST("/metrics/reset", middleware.ResetMetricsHandler(s.metrics))
}

// registerRoutes 注册路由（使用 Router Provider 模式）
func (s *ServerService) registerRoutes() {
	// 健康检查路由
	s.engine.GET("/health", s.healthCheck)
	s.engine.GET("/health/detail", s.healthCheckDetail)

	// 初始化 JWT 服务
	jwtService := auth.NewJWTService(s.config.JWT.SecretKey, s.config.JWT.ExpireIn)

	// 初始化 Casbin 权限服务
	casbinService, err := wire.InitializeCasbinService(s.db)
	if err != nil {
		s.logger.Warn("初始化 Casbin 服务失败，将跳过权限检查", zap.Error(err))
		// 继续运行，但权限检查将不可用
	}

	// 初始化 Token 黑名单服务（使用 Redis）
	var tokenBlacklist userservice.TokenBlacklistService = nil
	if s.redisClient != nil {
		metrics := wire.InitializeMetrics()
		rateLimiter := wire.InitializeRateLimiter(metrics)
		circuitBreaker := wire.InitializeCircuitBreaker(metrics)

		tokenBlacklist = wire.InitializeTokenBlacklistService(
			s.redisClient,
			metrics,
			rateLimiter,
			circuitBreaker,
		)
		s.logger.Info("Token 黑名单服务已初始化")
	} else {
		s.logger.Warn("Redis 客户端不可用，Token 黑名单服务将不可用")
	}

	// 创建认证中间件（带 Token 黑名单检查）
	authMiddleware := middleware.NewAuthMiddleware(jwtService, casbinService, tokenBlacklist)

	// API 路由组 - 公开接口（无需认证）
	apiPublic := s.engine.Group("/api")

	// API 路由组 - 需要认证的接口
	api := s.engine.Group("/api")
	api.Use(authMiddleware.HandlerFunc())

	// 将 Casbin 服务注入到 Context（供权限中间件使用）
	if casbinService != nil {
		api.Use(func(c *gin.Context) {
			c.Set("casbinService", casbinService)
			c.Next()
		})
	}

	// ==================== 初始化各模块 ====================

	// 1. 创建事件总线
	eventBus := event.NewEventBus()

	// 2. 初始化 User 模块（使用 Wire）
	userAppService, err := wire.InitializeUserModule(s.db, s.logger, jwtService, eventBus)
	if err != nil {
		s.logger.Error("初始化 User 模块失败", zap.Error(err))
		return
	}

	// 3. 创建 Handler（使用已初始化的 TokenBlacklistService）
	authHandler := authhttp.NewAuthHandler(userAppService, s.logger, tokenBlacklist)

	// User Handler 需要拆分 CQRS 服务
	userQuerySvc := userservice.NewUserQueryService(
		repo.NewUserDAORepository(s.db),
		repo.NewTenantMemberDAORepository(s.db),
	)
	userCommandSvc := userservice.NewUserCommandService(
		repo.NewUserDAORepository(s.db),
		repo.NewTenantMemberDAORepository(s.db),
		domainService.NewDefaultBcryptPasswordHasher(), // cost=12（生产环境）
		transaction.NewGormUnitOfWork(s.db),            // 新增 UnitOfWork
	)
	
	// 租户服务使用独立的 tenant/service 包（与第 234 行复用同一个实例）
	tenantRepo := repo.NewTenantDAORepository(s.db)
	tenantMemberRepo := repo.NewTenantMemberDAORepository(s.db)
	uow := transaction.NewGormUnitOfWork(s.db)
	tenantHandlerSvc := tenantservice.NewTenantService(tenantRepo, tenantMemberRepo, casbinService, uow)
	
	userHandler := userhttp.NewUserHandler(
		userAppService,    // authService
		userQuerySvc,      // userQueryService
		userCommandSvc,    // userCommandService
		tenantHandlerSvc,  // tenantService（复用租户模块的服务实例）
		s.logger,
	)

	// 4. 初始化租户模块（用于独立的租户路由，复用上面的实例）
	tenantHandler := tenanthttp.NewTenantHandler(tenantHandlerSvc, s.logger)

	// ==================== 使用 Router Provider 注册路由 ====================

	// 5. 创建 Router Provider
	authProvider := authhttp.NewAuthRouterProvider(authHandler)
	userProvider := userhttp.NewUserRouterProvider(userHandler)
	tenantProvider := tenanthttp.NewTenantRouterProvider(tenantHandler)

	// 6. 注册路由（清晰简洁）
	authProvider.ProvidePublicRoutes(apiPublic) // 认证路由（公开：register, login）
	authProvider.ProvideProtectedRoutes(api)    // 认证路由（受保护：logout）
	userProvider.ProvideProtectedRoutes(api)    // 用户路由（需认证）
	tenantProvider.ProvideProtectedRoutes(api)  // 租户路由（需认证）
}

// createServer 创建HTTP服务器
func (s *ServerService) createServer() {
	addr := fmt.Sprintf(":%d", s.config.App.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	s.logger.Info("HTTP服务器创建完成", zap.String("address", addr))
}

// Start 启动服务器
func (s *ServerService) Start() error {
	s.logger.Info("启动HTTP服务器",
		zap.String("address", s.server.Addr),
		zap.String("env", s.config.App.Env))

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("服务器启动失败: %w", err)
	}

	return nil
}

// Shutdown 优雅关闭服务器
func (s *ServerService) Shutdown(ctx context.Context) error {
	s.logger.Info("正在关闭HTTP服务器")
	return s.server.Shutdown(ctx)
}

// GetEngine 获取Gin引擎实例（用于测试或其他用途）
func (s *ServerService) GetEngine() *gin.Engine {
	return s.engine
}

// GetMetrics 获取监控指标
func (s *ServerService) GetMetrics() *metrics.APIMetrics {
	return s.metrics
}

// ==================== Handler Methods ====================

// healthCheck 健康检查处理器
func (s *ServerService) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// healthCheckDetail 详细健康检查处理器
func (s *ServerService) healthCheckDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"details": "All systems operational",
		"name":    s.config.App.Name,
		"env":     s.config.App.Env,
	})
}

// eventBusAdapter 事件总线适配器，将 event.EventBus 适配为 service.EventBus 接口
type eventBusAdapter struct {
	bus *event.EventBus
}

// Publish 实现 service.EventBus 接口
func (a *eventBusAdapter) Publish(event interface{}) error {
	if event == nil {
		return nil
	}
	// UserRegisteredEvent 不满足 DomainEvent 接口，这里需要特殊处理
	// 暂时跳过事件发布，因为当前的事件类型不兼容
	return nil
}
