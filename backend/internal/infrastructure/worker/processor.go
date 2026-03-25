package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	asynq_pkg "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/messaging/asynq"
	"go.uber.org/zap"
)

// Handler 领域事件处理器接口
type Handler interface {
	// CanHandle 是否可以处理该事件
	CanHandle(eventType string) bool
	// Handle 处理事件
	Handle(ctx context.Context, event kernel.DomainEvent) error
}

// Processor asynq 任务处理器
type Processor struct {
	logger   *zap.Logger
	handlers []Handler
}

// NewProcessor 创建 asynq 任务处理器
func NewProcessor(logger *zap.Logger, handlers ...Handler) *Processor {
	return &Processor{
		logger:   logger,
		handlers: handlers,
	}
}

// ProcessTask 处理 asynq 任务（实现 asynq.HandlerFunc）
func (p *Processor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	switch task.Type() {
	case asynq_pkg.TaskTypeDomainEvent:
		return p.processDomainEvent(ctx, task)
	default:
		return fmt.Errorf("unknown task type: %s", task.Type())
	}
}

// processDomainEvent 处理领域事件
func (p *Processor) processDomainEvent(ctx context.Context, task *asynq.Task) error {
	payload, err := asynq_pkg.ExtractDomainEventPayload(task)
	if err != nil {
		return fmt.Errorf("failed to extract payload: %w", err)
	}

	p.logger.Info("Processing domain event",
		zap.String("event_type", payload.EventType),
		zap.String("aggregate_id", payload.AggregateID),
	)

	// 根据事件类型反序列化为具体的领域事件
	domainEvent, err := p.deserializeEvent(payload)
	if err != nil {
		p.logger.Error("Failed to deserialize event",
			zap.String("event_type", payload.EventType),
			zap.Error(err),
		)
		return err
	}

	// 调用相应的处理器
	for _, handler := range p.handlers {
		if handler.CanHandle(payload.EventType) {
			if err := handler.Handle(ctx, domainEvent); err != nil {
				p.logger.Error("Handler failed",
					zap.String("event_type", payload.EventType),
					zap.Error(err),
				)
				return err
			}
		}
	}

	return nil
}

// deserializeEvent 根据事件类型反序列化为具体的领域事件
func (p *Processor) deserializeEvent(payload *asynq_pkg.DomainEventPayload) (kernel.DomainEvent, error) {
	// 使用映射表获取事件创建器
	eventCreator := eventCreators[payload.EventType]
	if eventCreator == nil {
		// 未知事件类型，返回通用适配器
		return &domainEventAdapter{
			eventName:   payload.EventType,
			aggregateID: payload.AggregateID,
			version:     int(payload.EventVersion),
			occurredOn:  parseTime(payload.OccurredOn),
			eventData:   payload.EventData,
		}, nil
	}

	// 创建事件实例并反序列化
	ev := eventCreator()
	if err := json.Unmarshal(payload.EventData, ev); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event %s: %w", payload.EventType, err)
	}

	// 初始化 BaseEvent
	setBaseEvent(ev, payload.EventType, payload.AggregateID, int(payload.EventVersion))

	return ev, nil
}

// eventCreators 事件创建器映射表
var eventCreators = map[string]func() kernel.DomainEvent{
	"UserRegistered":      func() kernel.DomainEvent { return &event.UserRegisteredEvent{} },
	"UserActivated":       func() kernel.DomainEvent { return &event.UserActivatedEvent{} },
	"UserDeactivated":     func() kernel.DomainEvent { return &event.UserDeactivatedEvent{} },
	"UserLoggedIn":        func() kernel.DomainEvent { return &event.UserLoggedInEvent{} },
	"UserPasswordChanged": func() kernel.DomainEvent { return &event.UserPasswordChangedEvent{} },
	"UserEmailChanged":    func() kernel.DomainEvent { return &event.UserEmailChangedEvent{} },
	"UserLocked":          func() kernel.DomainEvent { return &event.UserLockedEvent{} },
	"UserUnlocked":        func() kernel.DomainEvent { return &event.UserUnlockedEvent{} },
	"UserProfileUpdated":  func() kernel.DomainEvent { return &event.UserProfileUpdatedEvent{} },
}

// setBaseEvent 设置事件的 BaseEvent 字段（通过反射或类型断言）
func setBaseEvent(ev kernel.DomainEvent, eventType string, aggregateID interface{}, version int) {
	baseEvent := kernel.NewBaseEvent(eventType, aggregateID, version)

	// 根据具体事件类型设置 BaseEvent
	switch e := ev.(type) {
	case *event.UserRegisteredEvent:
		e.BaseEvent = baseEvent
	case *event.UserActivatedEvent:
		e.BaseEvent = baseEvent
	case *event.UserDeactivatedEvent:
		e.BaseEvent = baseEvent
	case *event.UserLoggedInEvent:
		e.BaseEvent = baseEvent
	case *event.UserPasswordChangedEvent:
		e.BaseEvent = baseEvent
	case *event.UserEmailChangedEvent:
		e.BaseEvent = baseEvent
	case *event.UserLockedEvent:
		e.BaseEvent = baseEvent
	case *event.UserUnlockedEvent:
		e.BaseEvent = baseEvent
	case *event.UserProfileUpdatedEvent:
		e.BaseEvent = baseEvent
	}
}

// domainEventAdapter 适配器，将 asynq 任务转换为 DomainEvent 接口
type domainEventAdapter struct {
	eventName   string
	aggregateID string
	version     int
	occurredOn  time.Time
	eventData   json.RawMessage
}

func (e *domainEventAdapter) EventName() string                { return e.eventName }
func (e *domainEventAdapter) OccurredOn() time.Time            { return e.occurredOn }
func (e *domainEventAdapter) AggregateID() interface{}         { return e.aggregateID }
func (e *domainEventAdapter) Version() int                     { return e.version }
func (e *domainEventAdapter) Metadata() map[string]interface{} { return make(map[string]interface{}) }

// parseTime 解析时间字符串
func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return t
}
