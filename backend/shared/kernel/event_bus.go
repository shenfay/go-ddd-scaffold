package kernel

import (
	"context"
	"sync"
)

// EventHandler 事件处理器函数类型
type EventHandler func(ctx context.Context, event DomainEvent) error

// EventBus 事件总线接口
type EventBus interface {
	Subscribe(eventName string, handler EventHandler)
	Publish(ctx context.Context, event DomainEvent) error
}

// SimpleEventBus 简单内存事件总线实现
type SimpleEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewSimpleEventBus 创建简单事件总线
func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (b *SimpleEventBus) Subscribe(eventName string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// Publish 发布事件
func (b *SimpleEventBus) Publish(ctx context.Context, event DomainEvent) error {
	b.mu.RLock()
	handlers := b.handlers[event.EventName()]
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
