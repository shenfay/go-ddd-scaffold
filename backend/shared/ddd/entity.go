package ddd

import (
	"time"
)

// AggregateRoot 聚合根接口定义
type AggregateRoot interface {
	ID() interface{}
	Version() int
	IncrementVersion()
	ApplyEvent(event DomainEvent)
	GetUncommittedEvents() []DomainEvent
	ClearUncommittedEvents()
	LoadFromHistory(events []DomainEvent) error
}

// BaseEntity 聚合根基础结构
type BaseEntity struct {
	id                interface{}
	version           int
	uncommittedEvents []DomainEvent
	createdAt         time.Time
	updatedAt         time.Time
}

// ID 返回聚合根ID
func (e *BaseEntity) ID() interface{} {
	return e.id
}

// Version 返回当前版本号
func (e *BaseEntity) Version() int {
	return e.version
}

// IncrementVersion 增加版本号
func (e *BaseEntity) IncrementVersion() {
	e.version++
	e.updatedAt = time.Now()
}

// ApplyEvent 应用领域事件
func (e *BaseEntity) ApplyEvent(event DomainEvent) {
	e.uncommittedEvents = append(e.uncommittedEvents, event)
}

// GetUncommittedEvents 获取未提交的事件
func (e *BaseEntity) GetUncommittedEvents() []DomainEvent {
	return e.uncommittedEvents
}

// ClearUncommittedEvents 清除已提交的事件
func (e *BaseEntity) ClearUncommittedEvents() {
	e.uncommittedEvents = []DomainEvent{}
}

// CreatedAt 返回创建时间
func (e *BaseEntity) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt 返回更新时间
func (e *BaseEntity) UpdatedAt() time.Time {
	return e.updatedAt
}

// SetCreatedAt 设置创建时间
func (e *BaseEntity) SetCreatedAt(t time.Time) {
	e.createdAt = t
}

// SetUpdatedAt 设置更新时间
func (e *BaseEntity) SetUpdatedAt(t time.Time) {
	e.updatedAt = t
}

// SetID 设置聚合根ID
func (e *BaseEntity) SetID(id interface{}) {
	e.id = id
}

// SetVersion 设置版本号
func (e *BaseEntity) SetVersion(version int) {
	e.version = version
}