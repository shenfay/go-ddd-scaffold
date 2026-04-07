package messaging

import (
	"context"
	"sync"

	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// EventBus 事件总线接口
type EventBus interface {
	Publish(ctx context.Context, evt event.Event) error
	Subscribe(eventType string, handler event.EventHandler)
}

// InProcessEventBus 内存事件总线（同步分发）
type InProcessEventBus struct {
	handlers map[string][]event.EventHandler
	mu       sync.RWMutex
}

// NewInProcessEventBus 创建内存事件总线
func NewInProcessEventBus() EventBus {
	return &InProcessEventBus{
		handlers: make(map[string][]event.EventHandler),
	}
}

// Publish 发布事件到订阅者（同步调用）
func (b *InProcessEventBus) Publish(ctx context.Context, evt event.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers := b.handlers[evt.GetType()]
	for _, handler := range handlers {
		if err := handler(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅事件
func (b *InProcessEventBus) Subscribe(eventType string, handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}
