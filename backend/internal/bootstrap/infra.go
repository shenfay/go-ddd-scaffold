package bootstrap

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	domain_event "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/eventstore"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
	task_queue "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/taskqueue"
)

// Infra 基础设施组件集合
// 纯数据结构体，用于存放所有基础设施组件
type Infra struct {
	DB             *gorm.DB
	Redis          *redis.Client
	Logger         *zap.Logger
	Config         *config.AppConfig
	Snowflake      *snowflake.Node
	EventPublisher kernel.EventPublisher
	EventBus       kernel.EventBus // 同步事件总线，用于领域事件订阅
	AsynqClient    *asynq.Client
	ErrorMapper    *kernel.ErrorMapper
}

// NewInfra 创建基础设施组件
// 返回 Infra 实例、cleanup 函数和可能的错误
// cleanup 函数按创建逆序释放资源
func NewInfra(cfg *config.AppConfig, logger *zap.Logger) (*Infra, func(), error) {
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

	// 3. 初始化 Snowflake ID 生成器
	nodeID := cfg.GetSnowflakeNodeID()
	snowflakeNode, err := snowflake.NewNode(nodeID)
	if err != nil {
		runCleanups(cleanups)
		return nil, nil, err
	}
	logger.Info("snowflake node initialized", zap.Int64("node_id", nodeID))

	// 4. 初始化 Asynq Client
	asynqClient := task_queue.NewClient(task_queue.Config{
		RedisAddr:     cfg.Redis.Addr,
		RedisPassword: cfg.Redis.Password,
		RedisDB:       cfg.Redis.DB,
	})
	cleanups = append(cleanups, func() {
		_ = asynqClient.Close()
		logger.Info("asynq client closed")
	})

	// 5. 初始化 Repository（用于 ActivityLog 和 EventLog）
	query := dao.Use(gormDB)

	// 6. 初始化 Asynq Publisher（任务发布器）
	asynqPublisher := task_queue.NewPublisher(asynqClient)

	// 7. 初始化 EventPublisher（使用适配器模式）
	eventPub := domain_event.NewEventPublisherAdapter(
		query,
		asynqPublisher,
		logger.Named("event_publisher"),
	)

	// 8. 初始化 ErrorMapper
	errorMapper := kernel.NewErrorMapper()

	// 9. 初始化 EventBus（同步事件总线）
	eventBus := kernel.NewSimpleEventBus()
	logger.Info("event bus initialized")

	// 构建 cleanup 函数（按逆序执行）
	cleanup := func() {
		runCleanups(cleanups)
	}

	return &Infra{
		DB:             gormDB,
		Redis:          redisClient,
		Logger:         logger,
		Config:         cfg,
		Snowflake:      snowflakeNode,
		EventPublisher: eventPub,
		EventBus:       eventBus,
		AsynqClient:    asynqClient,
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

// runCleanups 按逆序执行所有清理函数
func runCleanups(cleanups []func()) {
	for i := len(cleanups) - 1; i >= 0; i-- {
		cleanups[i]()
	}
}
