package container

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL 驱动
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/audit"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/loginlog"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	infraRepo "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
)

// CacheClient 缓存客户端接口（解耦具体实现）
type CacheClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
}

// redisCacheAdapter 适配现有的 RedisClient 到 CacheClient 接口
type redisCacheAdapter struct {
	client *redis.Client
}

func (a *redisCacheAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return a.client.Set(ctx, key, value, expiration).Err()
}

func (a *redisCacheAdapter) Get(ctx context.Context, key string) (string, error) {
	return a.client.Get(ctx, key).Result()
}

func (a *redisCacheAdapter) Del(ctx context.Context, keys ...string) error {
	return a.client.Del(ctx, keys...).Err()
}

func (a *redisCacheAdapter) Exists(ctx context.Context, keys ...string) (int64, error) {
	return a.client.Exists(ctx, keys...).Result()
}

// Container 应用容器接口（轻量级基础设施 + 路由管理）
type Container interface {
	// === 基础设施访问 ===
	GetDB() *sql.DB
	GetGormDB() *gorm.DB
	GetRedis() *redis.Client
	GetCache() CacheClient
	GetLogger(name string) *zap.Logger
	GetConfig() *config.AppConfig

	// === HTTP 路由访问 ===
	GetRouter() *gin.Engine

	// === ID 生成器访问 ===
	GetSnowflake() *snowflake.Node

	// === Repository 访问 ===
	GetUserRepo() repository.UserRepository
	GetLoginStatsRepo() repository.LoginStatsRepository
	GetAuditLogRepo() audit.AuditLogRepository
	GetLoginLogRepo() loginlog.LoginLogRepository

	// === Unit of Work 访问 ===
	GetUnitOfWork() application.UnitOfWork

	// === 生命周期管理 ===
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// ContainerInternal 内部接口，用于注册生命周期钩子
//
// export
//
//lint:ignore U1000 This is used for type assertion in other packages
type ContainerInternal interface {
	Container
	OnStart(fn func(context.Context) error)
	OnStop(fn func(context.Context) error)
}

// ContainerImpl 容器实现
type ContainerImpl struct {
	// 基础设施（一次性初始化）
	db        *sql.DB
	gormDB    *gorm.DB
	redis     *redis.Client
	cache     CacheClient
	logger    *zap.Logger
	config    *config.AppConfig
	router    *gin.Engine
	snowflake *snowflake.Node

	// Repository
	userRepo       repository.UserRepository
	loginStatsRepo repository.LoginStatsRepository
	auditLogRepo   audit.AuditLogRepository
	loginLogRepo   loginlog.LoginLogRepository

	// Unit of Work
	uow application.UnitOfWork

	// 生命周期钩子
	onStart []func(context.Context) error
	onStop  []func(context.Context) error
	mu      sync.Mutex // 保护钩子切片
}

// GetDB 获取数据库连接
func (c *ContainerImpl) GetDB() *sql.DB {
	return c.db
}

// GetGormDB 获取 GORM 数据库连接
func (c *ContainerImpl) GetGormDB() *gorm.DB {
	return c.gormDB
}

// GetRedis 获取 Redis 客户端
func (c *ContainerImpl) GetRedis() *redis.Client {
	return c.redis
}

// GetCache 获取缓存客户端
func (c *ContainerImpl) GetCache() CacheClient {
	return c.cache
}

// GetLogger 获取命名 logger
func (c *ContainerImpl) GetLogger(name string) *zap.Logger {
	return c.logger.Named(name)
}

// GetConfig 获取配置
func (c *ContainerImpl) GetConfig() *config.AppConfig {
	return c.config
}

// GetRouter 获取路由引擎
func (c *ContainerImpl) GetRouter() *gin.Engine {
	return c.router
}

// GetSnowflake 获取 Snowflake ID 生成器
func (c *ContainerImpl) GetSnowflake() *snowflake.Node {
	return c.snowflake
}

// GetUserRepo 获取用户 Repository
func (c *ContainerImpl) GetUserRepo() repository.UserRepository {
	return c.userRepo
}

// GetLoginStatsRepo 获取登录统计 Repository
func (c *ContainerImpl) GetLoginStatsRepo() repository.LoginStatsRepository {
	return c.loginStatsRepo
}

// GetAuditLogRepo 获取审计日志 Repository
func (c *ContainerImpl) GetAuditLogRepo() audit.AuditLogRepository {
	return c.auditLogRepo
}

// GetLoginLogRepo 获取登录日志 Repository
func (c *ContainerImpl) GetLoginLogRepo() loginlog.LoginLogRepository {
	return c.loginLogRepo
}

// GetUnitOfWork 获取 Unit of Work
func (c *ContainerImpl) GetUnitOfWork() application.UnitOfWork {
	return c.uow
}

// OnStart 注册启动钩子
func (c *ContainerImpl) OnStart(fn func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onStart = append(c.onStart, fn)
}

// OnStop 注册停止钩子
func (c *ContainerImpl) OnStop(fn func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onStop = append(c.onStop, fn)
}

// Start 启动容器（调用所有启动钩子）
func (c *ContainerImpl) Start(ctx context.Context) error {
	c.mu.Lock()
	hooks := make([]func(context.Context) error, len(c.onStart))
	copy(hooks, c.onStart)
	c.mu.Unlock()

	for _, fn := range hooks {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Stop 停止容器（调用所有停止钩子）
func (c *ContainerImpl) Stop(ctx context.Context) error {
	c.mu.Lock()
	hooks := make([]func(context.Context) error, len(c.onStop))
	copy(hooks, c.onStop)
	c.mu.Unlock()

	for _, fn := range hooks {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ContainerOption 容器选项
type ContainerOption func(*ContainerImpl)

// WithRouter 设置自定义路由引擎
func WithRouter(router *gin.Engine) ContainerOption {
	return func(c *ContainerImpl) {
		c.router = router
	}
}

// NewContainer 创建新容器
func NewContainer(
	cfg *config.AppConfig,
	logger *zap.Logger,
	opts ...ContainerOption,
) (Container, error) {
	// 1. 初始化数据库（复用现有代码）
	db, err := initDatabase(cfg.Database, logger)
	if err != nil {
		return nil, err
	}

	// 2. 初始化 GORM DB
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 3. 初始化 Redis（复用现有代码）
	redisClient, err := initRedisClient(cfg.Redis, logger)
	if err != nil {
		return nil, err
	}

	// 4. 初始化缓存（使用 Redis 客户端）
	cacheClient := &redisCacheAdapter{client: redisClient}

	// 5. 初始化路由
	router := gin.New()

	// 6. 初始化 Snowflake ID 生成器
	nodeID := cfg.GetSnowflakeNodeID()
	snowflakeNode, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	logger.Info("snowflake node initialized", zap.Int64("node_id", nodeID))

	// 应用选项
	c := &ContainerImpl{
		db:        db,
		gormDB:    gormDB,
		redis:     redisClient,
		cache:     cacheClient,
		logger:    logger,
		config:    cfg,
		router:    router,
		snowflake: snowflakeNode,
	}

	// 初始化 Repository
	daoQuery := dao.Use(gormDB)
	c.userRepo = infraRepo.NewUserRepository(daoQuery)
	c.loginStatsRepo = infraRepo.NewLoginStatsRepository(daoQuery)
	c.auditLogRepo = infraRepo.NewAuditLogRepository(daoQuery)
	c.loginLogRepo = infraRepo.NewLoginLogRepository(daoQuery)

	// 初始化 Unit of Work
	c.uow = application.NewUnitOfWork(gormDB, daoQuery)

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// initDatabase 初始化数据库连接
func initDatabase(cfg config.DatabaseConfig, logger *zap.Logger) (*sql.DB, error) {
	// 这里直接复用现有的 NewDatabase 函数
	// 为了简化，我们直接在这里实现
	dsn := cfg.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// 配置连接池
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	logger.Info("database connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name))

	return db, nil
}

// initRedisClient 初始化 Redis 客户端
func initRedisClient(cfg config.RedisConfig, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Info("redis connection established",
		zap.String("addr", cfg.Addr),
		zap.Int("db", cfg.DB))

	return client, nil
}
