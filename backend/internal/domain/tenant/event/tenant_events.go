package event

import (
	"time"

	"github.com/google/uuid"
)

// TenantCreatedEvent 租户创建事件
type TenantCreatedEvent struct {
	TenantID    uuid.UUID
	Name        string
	MaxMembers  int
	ExpiredAt   time.Time
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewTenantCreatedEvent 创建租户创建事件
func NewTenantCreatedEvent(tenantID uuid.UUID, name string, maxMembers int, expiredAt time.Time) *TenantCreatedEvent {
	return &TenantCreatedEvent{
		TenantID:    tenantID,
		Name:        name,
		MaxMembers:  maxMembers,
		ExpiredAt:   expiredAt,
		EventID:     uuid.New().String(),
		EventType:   "TenantCreated",
		AggregateID: tenantID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

func (e *TenantCreatedEvent) GetEventType() string      { return e.EventType }
func (e *TenantCreatedEvent) GetEventID() string        { return e.EventID }
func (e *TenantCreatedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *TenantCreatedEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *TenantCreatedEvent) GetVersion() int           { return e.Version }

// TenantMemberAddedEvent 租户成员添加事件
type TenantMemberAddedEvent struct {
	TenantID    uuid.UUID
	UserID      uuid.UUID
	MemberID    uuid.UUID
	Role        string
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewTenantMemberAddedEvent 创建租户成员添加事件
func NewTenantMemberAddedEvent(tenantID, userID, memberID uuid.UUID, role string) *TenantMemberAddedEvent {
	return &TenantMemberAddedEvent{
		TenantID:    tenantID,
		UserID:      userID,
		MemberID:    memberID,
		Role:        role,
		EventID:     uuid.New().String(),
		EventType:   "TenantMemberAdded",
		AggregateID: tenantID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

func (e *TenantMemberAddedEvent) GetEventType() string      { return e.EventType }
func (e *TenantMemberAddedEvent) GetEventID() string        { return e.EventID }
func (e *TenantMemberAddedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *TenantMemberAddedEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *TenantMemberAddedEvent) GetVersion() int           { return e.Version }

// TenantMemberRemovedEvent 租户成员移除事件
type TenantMemberRemovedEvent struct {
	TenantID    uuid.UUID
	UserID      uuid.UUID
	MemberID    uuid.UUID
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewTenantMemberRemovedEvent 创建租户成员移除事件
func NewTenantMemberRemovedEvent(tenantID, userID, memberID uuid.UUID) *TenantMemberRemovedEvent {
	return &TenantMemberRemovedEvent{
		TenantID:    tenantID,
		UserID:      userID,
		MemberID:    memberID,
		EventID:     uuid.New().String(),
		EventType:   "TenantMemberRemoved",
		AggregateID: tenantID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

func (e *TenantMemberRemovedEvent) GetEventType() string      { return e.EventType }
func (e *TenantMemberRemovedEvent) GetEventID() string        { return e.EventID }
func (e *TenantMemberRemovedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *TenantMemberRemovedEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *TenantMemberRemovedEvent) GetVersion() int           { return e.Version }
