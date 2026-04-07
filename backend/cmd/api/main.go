package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/listener"
	transhttp "github.com/shenfay/go-ddd-scaffold/internal/transport/http"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/handlers"
	pkglogger "github.com/shenfay/go-ddd-scaffold/pkg/logger"
	"github.com/shenfay/go-ddd-scaffold/pkg/metrics"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 导入生成的 Swagger 文档
	_ "github.com/shenfay/go-ddd-scaffold/api/swagger"
)

// @title           Go DDD Scaffold API
// @version         1.0
// @description     生产级 DDD 脚手架项目的 API 文档，包含用户认证、事件驱动等核心功能。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.github.com/shenfay/go-ddd-scaffold
// @contact.email  support@example.com

// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 使用 JWT Token，格式：Bearer {token}

func main() {
	// 1. 加载配置
	cfg, err := config.Load("development")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. 初始化日志系统
	if err := pkglogger.Init(cfg.Logger.Level); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer pkglogger.Sync() // 优雅关闭，确保所有日志写入磁盘

	pkglogger.Info("Starting application initialization...")

	// 3. 初始化数据库
	db, err := initDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. 初始化 Redis
	redisClient := initRedis(cfg.Redis)

	// 4. 初始化 Asynq Client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.Asynq.Addr,
	})
	defer asynqClient.Close()

	pkglogger.Info("✓ Asynq client initialized")

	// 创建 EventBus
	eventBus := messaging.NewEventBus(cfg.Redis.Addr, messaging.QueueConfig{
		Critical: "critical",
		Default:  "default",
	})
	pkglogger.Info("✓ Event Bus initialized")

	// 创建监听器并订阅事件
	auditLogListener := listener.NewAuditLogListener(eventBus)
	_ = auditLogListener // 保持引用，防止被 GC
	pkglogger.Info("✓ Audit Log Listener registered")

	activityLogListener := listener.NewActivityLogListener(eventBus)
	_ = activityLogListener
	pkglogger.Info("✓ Activity Log Listener registered")

	// 初始化服务依赖
	userRepo := repository.NewUserRepository(db)
	tokenService := authentication.NewTokenServiceImpl(
		redisClient,
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessExpire,
		cfg.JWT.RefreshExpire,
	)
	authService := authentication.NewService(userRepo, tokenService)
	authService.SetEventBus(eventBus)

	// 创建认证 Handler
	authHandler := handlers.NewAuthHandler(authService, tokenService)

	// 5. 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 6. 创建路由引擎
	engine := gin.New()

	// 7. 初始化 Prometheus 指标
	m := metrics.NewMetrics(prometheus.DefaultRegisterer)

	// 8. 注册中间件
	transhttp.Middlewares(engine, m)

	// 8. 创建并配置路由器
	apiRouter := transhttp.NewRouter(engine, authHandler)
	apiRouter.Setup()

	// 注册 Swagger UI 路由（开发环境）
	if gin.Mode() == gin.DebugMode {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("/swagger/doc.json"),
			ginSwagger.DefaultModelsExpandDepth(-1)))
		log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	}

	// 7. 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 8. 在 goroutine 中启动服务器
	go func() {
		log.Printf("Starting server on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 9. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 10. 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

// initDatabase 初始化数据库连接
func initDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 自动迁移表结构（开发环境）
	// 生产环境建议使用 migrations
	if err := db.AutoMigrate(&repository.UserPO{}); err != nil {
		return nil, err
	}

	log.Println("Database connection established and tables migrated")
	return db, nil
}

// initRedis 初始化 Redis 连接
func initRedis(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established")
	return client
}
