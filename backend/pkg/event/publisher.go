package event

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/pkg/logger"
	"go.uber.org/zap"
)

// Publisher 事件发布器
type Publisher struct {
	bus EventBus
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(ctx context.Context, evt Event) error
}

// NewPublisher 创建事件发布器
func NewPublisher(bus EventBus) *Publisher {
	return &Publisher{bus: bus}
}

// Publish 发布事件（自动处理错误日志）
func (p *Publisher) Publish(ctx context.Context, evt Event) {
	if p.bus == nil || evt == nil {
		return
	}

	if err := p.bus.Publish(ctx, evt); err != nil {
		logger.Warn("Failed to publish domain event",
			zap.String("event_type", evt.GetType()),
			zap.Error(err),
		)
	}
}
