package task_queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
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
	case TaskTypeDomainEvent:
		return p.processDomainEvent(ctx, task)
	default:
		return fmt.Errorf("unknown task type: %s", task.Type())
	}
}

// processDomainEvent 处理领域事件
func (p *Processor) processDomainEvent(ctx context.Context, task *asynq.Task) error {
	payload, err := ExtractDomainEventPayload(task)
	if err != nil {
		return fmt.Errorf("failed to extract payload: %w", err)
	}

	p.logger.Info("Processing domain event",
		zap.String("event_type", payload.EventType),
		zap.String("aggregate_id", payload.AggregateID),
	)

	// 创建一个简单的事件实现用于传递给处理器
	event := &domainEventAdapter{
		eventName:   payload.EventType,
		aggregateID: payload.AggregateID,
		version:     int(payload.EventVersion),
		occurredOn:  parseTime(payload.OccurredOn),
		eventData:   payload.EventData,
	}

	// 调用相应的处理器
	for _, handler := range p.handlers {
		if handler.CanHandle(payload.EventType) {
			if err := handler.Handle(ctx, event); err != nil {
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
