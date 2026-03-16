package ddd

import (
	"context"
	"reflect"
	"runtime"
	"sync"
)

// EventHandler 事件处理器函数类型
type EventHandler func(ctx context.Context, event DomainEvent) error

// EventBus 事件总线接口
type EventBus interface {
	// Subscribe 订阅指定类型的事件
	Subscribe(eventType string, handler EventHandler)

	// Publish 发布事件给所有订阅者
	Publish(ctx context.Context, event DomainEvent) error

	// Unsubscribe 取消订阅（可选）
	Unsubscribe(eventType string, handler EventHandler)
}

// SimpleEventBus 简单事件总线实现（同步模式）
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

// Subscribe 订阅指定类型的事件
func (b *SimpleEventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish 发布事件给所有订阅者（同步执行）
func (b *SimpleEventBus) Publish(ctx context.Context, event DomainEvent) error {
	b.mu.RLock()
	handlers := b.handlers[event.EventName()]
	b.mu.RUnlock()

	// 同步执行所有处理器
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			// 记录错误但不中断其他处理器
			// TODO: 添加日志记录
			continue
		}
	}

	return nil
}

// Unsubscribe 取消订阅
func (b *SimpleEventBus) Unsubscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventType]
	// 由于函数类型无法比较，我们重新创建一个不包含该 handler 的切片
	newHandlers := make([]EventHandler, 0, len(handlers))
	for _, h := range handlers {
		// 这里通过函数地址来判断（Go 中同一个函数变量的地址相同）
		// 注意：这个方法有局限性，实际使用中可能不需要 Unsubscribe
		if getFuncName(h) != getFuncName(handler) {
			newHandlers = append(newHandlers, h)
		}
	}
	b.handlers[eventType] = newHandlers
}

// getFuncName 获取函数名（用于比较）
func getFuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// AsyncEventBus 异步事件总线实现（可选）
type AsyncEventBus struct {
	*SimpleEventBus
	queue chan asyncEvent
}

type asyncEvent struct {
	ctx   context.Context
	event DomainEvent
}

// NewAsyncEventBus 创建异步事件总线
func NewAsyncEventBus(queueSize int) *AsyncEventBus {
	bus := &AsyncEventBus{
		SimpleEventBus: NewSimpleEventBus(),
		queue:          make(chan asyncEvent, queueSize),
	}

	// 启动后台消费者
	go bus.consume()

	return bus
}

// consume 消费事件队列（后台 goroutine）
func (b *AsyncEventBus) consume() {
	for ae := range b.queue {
		// 异步处理，不阻塞发布者
		_ = b.SimpleEventBus.Publish(ae.ctx, ae.event)
	}
}

// Publish 发布事件到队列（异步执行）
func (b *AsyncEventBus) Publish(ctx context.Context, event DomainEvent) error {
	select {
	case b.queue <- asyncEvent{ctx: ctx, event: event}:
		return nil
	default:
		// 队列已满，降级为同步执行
		return b.SimpleEventBus.Publish(ctx, event)
	}
}
