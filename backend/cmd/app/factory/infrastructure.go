package factory

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/model"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/event_bus/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	httpinfra "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/http"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
)

// Infrastructure 基础设施组件集合
type Infrastructure struct {
	DB             *gorm.DB
	Redis          *redis.Client
	Logger         *zap.Logger
	Config         *config.AppConfig
	EventPublisher common.EventPublisher
	EventBus       common.EventBus
	TaskPublisher  *asynq.EventPublisher
	ErrorMapper    *httpinfra.ErrorMapper
}

// NewInfrastructure 创建基础设施组件
func NewInfrastructure(cfg *config.AppConfig, logger *zap.Logger) (*Infrastructure, func(), error) {
	var cleanups []func()

	// 1. 初始化 PostgreSQL (GORM)
	gormDB, err := initGormDB(cfg.Database, logger)
	if err != nil {
		return nil, nil, err
	}
	cleanups = append(cleanups, func() {
		if sqlDB, err := gormDB.DB(); err == nil {
			_ = sqlDB.Close()
			logger.Info("database connection closed")
		}
	})

	// 2. 初始化 Redis
	redisClient, err := initRedis(cfg.Redis, logger)
	if err != nil {
		runCleanups(cleanups)
		return nil, nil, err
	}
	cleanups = append(cleanups, func() {
		_ = redisClient.Close()
		logger.Info("redis connection closed")
	})

	// 3. 初始化 Snowflake ID 生成器（使用 yitter/idgenerator-go）
	nodeID := cfg.GetSnowflakeNodeID()
	idgen.Initialize(uint64(nodeID), 10) // WorkerIdBitLength=10，支持 1024 个节点
	logger.Info("snowflake id generator initialized", zap.Int64("worker_id", nodeID))

	// 4. 初始化 Asynq Client
	asynqClient := asynq.NewClient(asynq.Config{
		RedisAddr:     cfg.Redis.Addr,
		RedisPassword: cfg.Redis.Password,
		RedisDB:       cfg.Redis.DB,
	})

	// 5. 初始化 Repository（用于 ActivityLog 和 EventLog）
	query := dao.Use(gormDB)

	// 6. 初始化 Asynq Publisher（任务发布器）
	asynqPublisher := asynq.NewEventPublisher(
		query,
		asynqClient,
		redisClient,
		logger,
		asynq.EventPublisherConfig{
			DefaultQueue: "default",
		},
	)

	// 7. 初始化 EventPublisher（使用新的工厂函数）
	eventPub, err := CreateEventPublisher(
		cfg,
		query,
		redisClient,
		logger.Named("event_publisher"),
	)
	if err != nil {
		runCleanups(cleanups)
		return nil, nil, err
	}

	// 8. 初始化 ErrorMapper（移至 infrastructure 层）
	errorMapper := httpinfra.NewErrorMapper()

	// 9. 初始化 EventBus（同步事件总线）
	eventBus := common.NewSimpleEventBus()
	logger.Info("event bus initialized")

	// 构建 cleanup 函数（按逆序执行）
	cleanup := func() {
		runCleanups(cleanups)
	}

	return &Infrastructure{
		DB:             gormDB,
		Redis:          redisClient,
		Logger:         logger,
		Config:         cfg,
		EventPublisher: eventPub,
		EventBus:       eventBus,
		TaskPublisher:  asynqPublisher,
		ErrorMapper:    errorMapper,
	}, cleanup, nil
}

// initGormDB 初始化 GORM 数据库连接
func initGormDB(cfg config.DatabaseConfig, logger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 获取底层 sql.DB 配置连接池
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	logger.Info("database connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name))

	return gormDB, nil
}

// initRedis 初始化 Redis 客户端
func initRedis(cfg config.RedisConfig, logger *zap.Logger) (*redis.Client, error) {
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

// ActivityLogRepository 获取活动日志仓储实例
func (i *Infrastructure) ActivityLogRepository() model.ActivityLogRepository {
	query := dao.Use(i.DB)
	return repository.NewActivityLogRepository(query)
}

// runCleanups 按逆序执行所有清理函数
func runCleanups(cleanups []func()) {
	for i := len(cleanups) - 1; i >= 0; i-- {
		cleanups[i]()
	}
}
