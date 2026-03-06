// Package wire 依赖注入配置 - 事件模块
package wire

import (
	"go-ddd-scaffold/internal/infrastructure/event"

	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// EventModuleSet 事件模块的依赖集合
var EventModuleSet = wire.NewSet(
	InitializeOutboxRepository,
	InitializeTransactionalEventPublisher,
	InitializeEventManager,
)

// InitializeOutboxRepository 初始化发件箱仓储
func InitializeOutboxRepository(db *gorm.DB) event.OutboxRepository {
	return event.NewGormOutboxRepository(db)
}

// InitializeTransactionalEventPublisher 初始化事务性事件发布器
func InitializeTransactionalEventPublisher(
	outboxRepo event.OutboxRepository,
	eventBus *event.EventBus,
	logger *zap.Logger,
) *event.TransactionalEventPublisher {
	return event.NewTransactionalEventPublisher(
		outboxRepo,
		eventBus,
		logger,
		event.TransactionalEventPublisherConfig{
			PollInterval: 5,    // 5 秒轮询一次
			BatchSize:    100,  // 每批处理 100 个事件
			MaxWorkers:   3,    // 3 个并发协程
			EnableWorker: true, // 启用后台 Worker
		},
	)
}

// InitializeEventManager 初始化事件管理器（带重试支持）
func InitializeEventManager(
	eventBus *event.EventBus,
	logger *zap.Logger,
) *event.EventManager {
	return event.NewEventManager(eventBus, logger)
}
