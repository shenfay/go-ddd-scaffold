package domain_event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/task_queue"
	"go.uber.org/zap"
)

// Publisher 领域事件发布器接口
type Publisher interface {
	Publish(ctx context.Context, event kernel.DomainEvent) error
}

// AsynqPublisher 基于 asynq 的事件发布器实现
type AsynqPublisher struct {
	publisher *task_queue.Publisher
	logger    *zap.Logger
}

// NewAsynqPublisher 创建 asynq 事件发布器
func NewAsynqPublisher(publisher *task_queue.Publisher, logger *zap.Logger) *AsynqPublisher {
	return &AsynqPublisher{
		publisher: publisher,
		logger:    logger,
	}
}

// Publish 发布领域事件到 asynq 队列
func (p *AsynqPublisher) Publish(ctx context.Context, event kernel.DomainEvent) error {
	// 安全地获取 aggregate_id
	var aggregateID string
	if id, ok := event.AggregateID().(string); ok {
		aggregateID = id
	} else {
		// 如果不是字符串类型，尝试转换为字符串
		aggregateID = fmt.Sprintf("%v", event.AggregateID())
	}

	p.logger.Info("Publishing event to asynq",
		zap.String("event_type", event.EventName()),
		zap.String("aggregate_id", aggregateID),
	)

	// 序列化事件数据
	eventData, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event", zap.Error(err))
		return err
	}

	// 创建任务负载
	payload := task_queue.DomainEventPayload{
		AggregateID:   aggregateID, // 使用之前安全转换的变量
		AggregateType: getAggregateType(event),
		EventType:     event.EventName(),
		EventVersion:  int32(event.Version()),
		EventData:     eventData,
		OccurredOn:    event.OccurredOn().Format(time.RFC3339),
		Metadata:      convertMetadata(event.Metadata()),
	}

	// 根据事件类型选择优先级队列
	queue := p.getQueueForEvent(event.EventName())

	// 将任务加入队列
	err = p.publisher.PublishDomainEvent(ctx, payload, queue)
	if err != nil {
		return err
	}

	p.logger.Debug("Event published to asynq",
		zap.String("event_type", event.EventName()),
		zap.String("queue", queue),
	)

	return nil
}

// getQueueForEvent 根据事件类型返回相应的队列
func (p *AsynqPublisher) getQueueForEvent(eventType string) string {
	// 高优先级事件
	criticalEvents := map[string]bool{
		"UserRegistered":   true, // 用户注册
		"PaymentCompleted": true, // 支付完成
	}

	if criticalEvents[eventType] {
		return "critical"
	}

	// 低优先级事件
	lowEvents := map[string]bool{
		"UserLoggedIn":     true, // 用户登录（日志类）
		"NotificationSent": true, // 通知发送
	}

	if lowEvents[eventType] {
		return "low"
	}

	return "default"
}

// getAggregateType 从事件中获取聚合类型（简单实现，实际项目中可以更复杂）
func getAggregateType(event kernel.DomainEvent) string {
	// 这里可以根据事件名称推断聚合类型
	// 例如：UserCreated -> User, TenantCreated -> Tenant
	return "Unknown"
}

// convertMetadata 转换元数据类型
func convertMetadata(meta map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range meta {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}
	return result
}
