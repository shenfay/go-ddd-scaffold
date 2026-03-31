package common

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
}

// EntityMeta 聚合根元数据 (组合模式)
// 用于管理聚合根的通用状态，避免继承带来的复杂性
type EntityMeta struct {
	id                interface{}
	version           int
	uncommittedEvents []DomainEvent
	createdAt         time.Time
	updatedAt         time.Time
}

// NewEntityMeta 创建实体元数据
func NewEntityMeta(id interface{}, createdAt time.Time) *EntityMeta {
	return &EntityMeta{
		id:                id,
		version:           0,
		uncommittedEvents: make([]DomainEvent, 0),
		createdAt:         createdAt,
		updatedAt:         createdAt,
	}
}

// ID 返回聚合根 ID
func (m *EntityMeta) ID() interface{} {
	return m.id
}

// Version 返回当前版本号
func (m *EntityMeta) Version() int {
	return m.version
}

// IncrementVersion 增加版本号
func (m *EntityMeta) IncrementVersion() {
	m.version++
	m.updatedAt = time.Now()
}

// ApplyEvent 应用领域事件
func (m *EntityMeta) ApplyEvent(event DomainEvent) {
	m.uncommittedEvents = append(m.uncommittedEvents, event)
}

// GetUncommittedEvents 获取未提交的事件
func (m *EntityMeta) GetUncommittedEvents() []DomainEvent {
	return m.uncommittedEvents
}

// ClearUncommittedEvents 清除已提交的事件
func (m *EntityMeta) ClearUncommittedEvents() {
	m.uncommittedEvents = make([]DomainEvent, 0)
}

// CreatedAt 返回创建时间
func (m *EntityMeta) CreatedAt() time.Time {
	return m.createdAt
}

// UpdatedAt 返回更新时间
func (m *EntityMeta) UpdatedAt() time.Time {
	return m.updatedAt
}

// SetCreatedAt 设置创建时间
func (m *EntityMeta) SetCreatedAt(t time.Time) {
	m.createdAt = t
}

// SetUpdatedAt 设置更新时间
func (m *EntityMeta) SetUpdatedAt(t time.Time) {
	m.updatedAt = t
}

// SetID 设置聚合根 ID
func (m *EntityMeta) SetID(id interface{}) {
	m.id = id
}

// SetVersion 设置版本号
func (m *EntityMeta) SetVersion(version int) {
	m.version = version
}
