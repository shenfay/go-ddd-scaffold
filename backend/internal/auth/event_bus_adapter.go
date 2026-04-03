package auth

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// eventBusAdapter 适配器：将旧的 event.EventBus 转换为新的 messaging.EventBus
type eventBusAdapter struct {
	eventBus messaging.EventBus
}

// NewEventBusAdapter 创建事件总线适配器
func NewEventBusAdapter(eventBus messaging.EventBus) event.EventBus {
	return &eventBusAdapter{eventBus: eventBus}
}

// Publish 发布事件（实现 event.EventBus 接口）
func (a *eventBusAdapter) Publish(ctx context.Context, evt event.Event) error {
	return a.eventBus.Publish(ctx, evt)
}

// Subscribe 订阅事件（实现 event.EventBus 接口）
func (a *eventBusAdapter) Subscribe(eventType string, handler event.EventHandler) {
	a.eventBus.Subscribe(eventType, handler)
}
