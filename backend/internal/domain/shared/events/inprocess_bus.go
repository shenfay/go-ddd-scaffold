package events

import (
	"context"
	"sync"
)

// InProcessBus 进程内事件总线
// 基于 sync.RWMutex 实现线程安全的同步事件发布/订阅
type InProcessBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

// NewInProcessBus 创建进程内事件总线
func NewInProcessBus() *InProcessBus {
	return &InProcessBus{
		handlers: make(map[string][]Handler),
	}
}

// Publish 发布领域事件
// 同步调用所有订阅了该事件类型的 Handler，按注册顺序执行
func (b *InProcessBus) Publish(ctx context.Context, event DomainEvent) error {
	b.mu.RLock()
	handlers := b.handlers[event.EventName()]
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅指定类型的领域事件
// 支持运行时动态注册，线程安全
func (b *InProcessBus) Subscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}
