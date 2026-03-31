package common

import (
	"context"
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventName() string
	OccurredOn() time.Time
	AggregateID() interface{}
	Version() int
	Metadata() map[string]interface{}
}

// EventPublisher 领域事件发布器接口
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}

// BaseEvent 领域事件基础结构 (使用组合模式)
// 包含所有事件的通用字段，通过组合而非继承使用
type BaseEvent struct {
	eventName   string
	aggregateID interface{}
	version     int
	occurredOn  time.Time
	metadata    map[string]interface{}
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(eventName string, aggregateID interface{}, version int) *BaseEvent {
	return &BaseEvent{
		eventName:   eventName,
		aggregateID: aggregateID,
		version:     version,
		occurredOn:  time.Now(),
		metadata:    make(map[string]interface{}),
	}
}

// EventName 返回事件名称
func (e *BaseEvent) EventName() string {
	return e.eventName
}

// OccurredOn 返回事件发生时间
func (e *BaseEvent) OccurredOn() time.Time {
	return e.occurredOn
}

// AggregateID 返回聚合根 ID
func (e *BaseEvent) AggregateID() interface{} {
	return e.aggregateID
}

// Version 返回事件版本
func (e *BaseEvent) Version() int {
	return e.version
}

// Metadata 返回事件元数据
func (e *BaseEvent) Metadata() map[string]interface{} {
	if e.metadata == nil {
		e.metadata = make(map[string]interface{})
	}
	return e.metadata
}

// SetMetadata 设置事件元数据
func (e *BaseEvent) SetMetadata(key string, value interface{}) {
	if e.metadata == nil {
		e.metadata = make(map[string]interface{})
	}
	e.metadata[key] = value
}
