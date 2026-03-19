package messaging

import (
	"context"
	"log"
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
	handlers, exists := b.handlers[eventType]
	if !exists {
		return nil // 没有订阅者是正常的
	}

	// 异步处理事件
	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			if err := h.Handle(ctx, event); err != nil {
				log.Printf("Error handling event %s: %v", eventType, err)
				errChan <- err
			}
		}(handler)
	}

	// 等待所有处理器完成
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 收集错误（不中断主流程）
	for err := range errChan {
		_ = err // 错误已记录，这里只是消费通道
	}

	return nil
}

// Subscribe 订阅事件
func (b *SimpleEventBus) Subscribe(eventType string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.handlers[eventType] == nil {
		b.handlers[eventType] = make([]EventHandler, 0)
	}

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

// Unsubscribe 取消订阅
func (b *SimpleEventBus) Unsubscribe(eventType string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers, exists := b.handlers[eventType]
	if !exists {
		return nil
	}

	// 查找并移除处理器
	for i, h := range handlers {
		if h == handler {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

// AsyncEventBus 异步事件总线
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
func NewAsyncEventBus(bufferSize int, workers int) *AsyncEventBus {
	bus := &AsyncEventBus{
		SimpleEventBus: SimpleEventBus{
			handlers: make(map[string][]EventHandler),
		},
		eventChan: make(chan eventWithContext, bufferSize),
		workers:   workers,
	}

	// 启动工作协程
	for i := 0; i < workers; i++ {
		bus.wg.Add(1)
		go bus.worker()
	}

	return bus
}

// worker 事件处理工作协程
func (b *AsyncEventBus) worker() {
	defer b.wg.Done()

	for ec := range b.eventChan {
		b.processEvent(ec.ctx, ec.event)
	}
}

// processEvent 处理单个事件
func (b *AsyncEventBus) processEvent(ctx context.Context, event kernel.DomainEvent) {
	b.mu.RLock()
	handlers, exists := b.handlers[event.EventName()]
	b.mu.RUnlock()

	if !exists {
		return
	}

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			log.Printf("Error handling event %s: %v", event.EventName(), err)
		}
	}
}

// Publish 异步发布事件
func (b *AsyncEventBus) Publish(ctx context.Context, event kernel.DomainEvent) error {
	select {
	case b.eventChan <- eventWithContext{ctx: ctx, event: event}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop 停止事件总线
func (b *AsyncEventBus) Stop() {
	close(b.eventChan)
	b.wg.Wait()
}

// EventPublisherAdapter 事件发布适配器
type EventPublisherAdapter struct {
	bus EventBus
}

// NewEventPublisherAdapter 创建事件发布适配器
func NewEventPublisherAdapter(bus EventBus) *EventPublisherAdapter {
	return &EventPublisherAdapter{bus: bus}
}

// Publish 发布事件（适配 commands.EventPublisher 接口）
func (a *EventPublisherAdapter) Publish(ctx context.Context, event kernel.DomainEvent) error {
	return a.bus.Publish(ctx, event)
}
