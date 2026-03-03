// Package event 通用领域事件定义
package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DomainEvent 领域事件接口（所有领域事件的通用接口）
type DomainEvent interface {
	// GetEventType 获取事件类型
	GetEventType() string
	// GetEventID 获取事件 ID
	GetEventID() string
	// GetAggregateID 获取聚合根 ID
	GetAggregateID() uuid.UUID
	// GetOccurredAt 获取事件发生时间
	GetOccurredAt() time.Time
	// GetVersion 获取事件版本
	GetVersion() int
}

// BaseDomainEvent 基础领域事件（方便实现 DomainEvent 接口）
type BaseDomainEvent struct {
	EventType   string
	EventID     string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// GetEventType 获取事件类型
func (e *BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetEventID 获取事件 ID
func (e *BaseDomainEvent) GetEventID() string {
	return e.EventID
}

// GetAggregateID 获取聚合根 ID
func (e *BaseDomainEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

// GetOccurredAt 获取事件发生时间
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetVersion 获取事件版本
func (e *BaseDomainEvent) GetVersion() int {
	return e.Version
}

// EventHandler 领域事件处理器函数类型
type EventHandler func(ctx context.Context, event DomainEvent) error

// EventPublisher 事件发布器
type EventPublisher struct {
	subscribers map[string][]EventHandler
}

// NewEventPublisher 创建事件发布器
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (ep *EventPublisher) Subscribe(eventType string, handler EventHandler) {
	ep.subscribers[eventType] = append(ep.subscribers[eventType], handler)
}

// Unsubscribe 取消订阅（简化实现：清空该类型的所有订阅者）
func (ep *EventPublisher) Unsubscribe(eventType string, handler EventHandler) {
	// 由于 Go 中函数无法直接比较，这里简单清空该事件类型的所有订阅
	// 如果需要精确移除特定处理器，需要使用处理器 ID 或其他机制
	delete(ep.subscribers, eventType)
}

// Publish 异步发布事件（在当前 goroutine 中执行）
func (ep *EventPublisher) Publish(event DomainEvent) error {
	handlers := ep.subscribers[event.GetEventType()]
	for _, handler := range handlers {
		go handler(context.Background(), event)
	}
	return nil
}

// PublishSync 同步发布事件
func (ep *EventPublisher) PublishSync(event DomainEvent) error {
	handlers := ep.subscribers[event.GetEventType()]
	for _, handler := range handlers {
		if err := handler(context.Background(), event); err != nil {
			return err
		}
	}
	return nil
}

// GetSubscriberCount 获取订阅者数量
func (ep *EventPublisher) GetSubscriberCount(eventType string) int {
	return len(ep.subscribers[eventType])
}
