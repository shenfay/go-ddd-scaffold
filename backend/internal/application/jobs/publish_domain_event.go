package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	"go.uber.org/zap"
)

// PublishDomainEventJob 发布领域事件作业
// 职责：从 Asynq 队列接收事件并发布到 EventBus
// 注意：实际的事件发布由 OutboxProcessor 完成，这个 Job 是作为备用机制
type PublishDomainEventJob struct {
	eventBus common.EventBus
	logger   *zap.Logger
}

// NewPublishDomainEventJob 创建发布领域事件作业
func NewPublishDomainEventJob(eventBus common.EventBus, logger *zap.Logger) *PublishDomainEventJob {
	return &PublishDomainEventJob{
		eventBus: eventBus,
		logger:   logger.Named("publish_domain_event"),
	}
}

// Execute 执行发布领域事件任务
func (j *PublishDomainEventJob) Execute(ctx context.Context, payload map[string]interface{}) error {
	j.logger.Debug("Processing domain event publication", zap.Any("payload", payload))

	// 1. 从 payload 中提取事件数据
	eventType, ok := payload["event_type"].(string)
	if !ok {
		j.logger.Error("Invalid event type in payload")
		return nil
	}

	dataStr, ok := payload["data"].(string)
	if !ok {
		j.logger.Error("Invalid data in payload")
		return nil
	}

	// 2. 反序列化事件数据
	var eventData map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &eventData); err != nil {
		j.logger.Error("Failed to unmarshal event data", zap.Error(err))
		return nil
	}

	// 3. 根据事件类型重建事件对象
	domainEvent := j.reconstructEvent(eventType, eventData)
	if domainEvent == nil {
		j.logger.Error("Unknown event type", zap.String("event_type", eventType))
		return nil
	}

	// 4. 通过 EventBus 发布事件
	if j.eventBus != nil {
		if err := j.eventBus.Publish(ctx, domainEvent); err != nil {
			j.logger.Error("Failed to publish event", zap.Error(err))
			return err
		}
	}

	j.logger.Info("Domain event published successfully",
		zap.String("event_type", eventType),
		zap.Any("aggregate_id", domainEvent.AggregateID()))

	return nil
}

// reconstructEvent 根据事件类型重建事件对象
func (j *PublishDomainEventJob) reconstructEvent(eventType string, data map[string]interface{}) common.DomainEvent {
	switch eventType {
	case "UserRegistered":
		return &userEvent.UserRegisteredEvent{
			// 实际项目中应该从 data 中恢复所有字段
		}
	case "UserActivated":
		return &userEvent.UserActivatedEvent{}
	case "UserLoggedIn":
		return &userEvent.UserLoggedInEvent{}
	case "UserDeactivated":
		return &userEvent.UserDeactivatedEvent{}
	case "UserPasswordChanged":
		return &userEvent.UserPasswordChangedEvent{}
	case "UserProfileUpdated":
		return &userEvent.UserProfileUpdatedEvent{}
	case "UserEmailChanged":
		return &userEvent.UserEmailChangedEvent{}
	case "UserLocked":
		return &userEvent.UserLockedEvent{}
	case "UserUnlocked":
		return &userEvent.UserUnlockedEvent{}
	default:
		return nil
	}
}

// Queue 返回队列名称
func (j *PublishDomainEventJob) Queue() string {
	return "jobs_realtime"
}

// MaxRetry 返回最大重试次数
func (j *PublishDomainEventJob) MaxRetry() int {
	return 3
}

// Timeout 返回超时时间
func (j *PublishDomainEventJob) Timeout() time.Duration {
	return 30 * time.Second
}
