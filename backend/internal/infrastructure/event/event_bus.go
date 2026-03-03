package event

import (
	"context"
	"fmt"
	"sync"
)

// EventBus 事件总线
// 基础设施层实现，协调事件的发布和订阅
type EventBus struct {
	publisher *EventPublisher
	handlers  map[string][]EventHandler
	mutex     sync.RWMutex
}

// NewEventBus 创建事件总线实例
func NewEventBus() *EventBus {
	return &EventBus{
		publisher: NewEventPublisher(),
		handlers:  make(map[string][]EventHandler),
	}
}

// RegisterHandler 注册事件处理器
func (eb *EventBus) RegisterHandler(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.publisher.Subscribe(eventType, handler)
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// UnregisterHandler 注销事件处理器（简化实现：清空该类型的所有处理器）
func (eb *EventBus) UnregisterHandler(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.publisher.Unsubscribe(eventType, handler)

	// 由于 Go 中函数无法直接比较，这里简单清空该事件类型的所有处理器
	delete(eb.handlers, eventType)
}

// Publish 发布领域事件（异步）
func (eb *EventBus) Publish(ctx context.Context, event DomainEvent) error {
	if event == nil {
		return fmt.Errorf("事件不能为空")
	}

	return eb.publisher.Publish(event)
}

// PublishSync 同步发布领域事件
func (eb *EventBus) PublishSync(ctx context.Context, event DomainEvent) error {
	if event == nil {
		return fmt.Errorf("事件不能为空")
	}

	return eb.publisher.PublishSync(event)
}

// GetHandlerCount 获取指定事件类型的处理器数量
func (eb *EventBus) GetHandlerCount(eventType string) int {
	return eb.publisher.GetSubscriberCount(eventType)
}

// ListRegisteredEvents 列出所有已注册事件类型
func (eb *EventBus) ListRegisteredEvents() []string {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	events := make([]string, 0, len(eb.handlers))
	for eventType := range eb.handlers {
		events = append(events, eventType)
	}

	return events
}
