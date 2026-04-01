package factory

import (
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	asynq_pkg "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/event_bus/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/worker"
	"go.uber.org/zap"
)

// CreateEventPublisher 创建事件发布器
func CreateEventPublisher(
	cfg *config.AppConfig,
	query *dao.Query,
	redisClient *redis.Client,
	logger *zap.Logger,
) (common.EventPublisher, error) {
	// 创建 Asynq 客户端
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 配置事件发布器
	publisherConfig := asynq_pkg.EventPublisherConfig{
		DefaultQueue:      "default",
		HighPriorityQueue: "high_priority",
		LowPriorityQueue:  "low_priority",
		DeduplicationTTL:  cfg.Events.Deduplication.TTL,
		MaxRetries:        cfg.Events.Asynq.Retry.MaxAttempts,
		BaseDelay:         cfg.Events.Asynq.Retry.MinBackoff,
		MaxDelay:          cfg.Events.Asynq.Retry.MaxBackoff,
		EnableMetrics:     cfg.Events.Monitoring.EnableMetrics,
	}

	// 创建优化的事件发布器
	publisher := asynq_pkg.NewEventPublisher(
		query,
		asynqClient,
		redisClient,
		logger,
		publisherConfig,
	)

	return publisher, nil
}

// CreateAsynqServer 创建 Asynq 服务器
func CreateAsynqServer(
	cfg *config.AppConfig,
	redisClient *redis.Client,
	logger *zap.Logger,
) *asynq.Server {
	// 构建队列配置
	queues := make(map[string]int)
	for queueName, priority := range cfg.Events.Asynq.Queues {
		queues[queueName] = priority
	}

	// 创建 Asynq 服务器
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
		asynq.Config{
			Concurrency: cfg.Events.Asynq.Concurrency,
			Queues:      queues,
			// Logger:      logger.Named("asynq_server"), // 暂时注释掉，需要适配 Asynq 的 Logger 接口
		},
	)

	return server
}

// CreateSubscriberManager 创建订阅者管理器
func CreateSubscriberManager(
	cfg *config.AppConfig,
	logger *zap.Logger,
) *worker.SubscriberManager {
	subscriberConfig := worker.SubscriberConfig{
		MaxConcurrency:  cfg.Events.Asynq.Concurrency,
		BufferSize:      100,
		HandlerTimeout:  cfg.Events.Asynq.Retry.MaxBackoff,
		MaxRetries:      cfg.Events.Asynq.Retry.MaxAttempts,
		RetryDelay:      cfg.Events.Asynq.Retry.MinBackoff,
		DeadLetterQueue: "dead_letters",
		EnableMetrics:   cfg.Events.Monitoring.EnableMetrics,
		MetricsPrefix:   "events",
	}

	return worker.NewSubscriberManager(logger, subscriberConfig)
}
