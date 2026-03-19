package domain_event

import (
	"context"
	"sync"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event kernel.DomainEvent) error
	// Subscribe 订阅事件
	Subscribe(eventType string, handler EventHandler) error
	// Unsubscribe 取消订阅
	Unsubscribe(eventType string, handler EventHandler) error
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event kernel.DomainEvent) error
}

// EventHandlerFunc 事件处理器函数类型
type EventHandlerFunc func(ctx context.Context, event kernel.DomainEvent) error

// Handle 实现 EventHandler 接口
func (f EventHandlerFunc) Handle(ctx context.Context, event kernel.DomainEvent) error {
	return f(ctx, event)
}

// SimpleEventBus 简单事件总线实现
type SimpleEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewSimpleEventBus 创建简单事件总线
func NewSimpleEventBus() EventBus {
	return &SimpleEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish 发布事件
func (b *SimpleEventBus) Publish(ctx context.Context, event kernel.DomainEvent) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	eventType := event.EventName()
	handlers := b.handlers[eventType]

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

// Subscribe 订阅事件
func (b *SimpleEventBus) Subscribe(eventType string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

// Unsubscribe 取消订阅
func (b *SimpleEventBus) Unsubscribe(eventType string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventType]
	for i, h := range handlers {
		if h == handler {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

// AsyncEventBus 异步事件总线（带 worker 池）
type AsyncEventBus struct {
	SimpleEventBus
	eventChan chan eventWithContext
	workers   int
	wg        sync.WaitGroup
}

type eventWithContext struct {
	ctx   context.Context
	event kernel.DomainEvent
}

// NewAsyncEventBus 创建异步事件总线
func NewAsyncEventBus(workers int) EventBus {
	bus := &AsyncEventBus{
		eventChan: make(chan eventWithContext, 100),
		workers:   workers,
	}
	bus.handlers = make(map[string][]EventHandler)
	bus.startWorkers()
	return bus
}

func (b *AsyncEventBus) startWorkers() {
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			for e := range b.eventChan {
				b.SimpleEventBus.Publish(e.ctx, e.event)
			}
		}()
	}
}

// Publish 发布事件（异步）
func (b *AsyncEventBus) Publish(ctx context.Context, event kernel.DomainEvent) error {
	b.eventChan <- eventWithContext{ctx: ctx, event: event}
	return nil
}

// Stop 停止事件总线
func (b *AsyncEventBus) Stop() {
	close(b.eventChan)
	b.wg.Wait()
}
