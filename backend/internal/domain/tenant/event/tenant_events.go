package event

import (
	"time"

	"github.com/google/uuid"
)

// TenantCreatedEvent 租户创建事件
type TenantCreatedEvent struct{
	TenantID   uuid.UUID `json:"tenantId"`
	OwnerID    uuid.UUID `json:"ownerId"`
	TenantName string    `json:"tenantName"`
	CreatedAt  time.Time `json:"createdAt"`
	EventID    string    `json:"eventId"`
	EventType  string    `json:"eventType"`
	AggregateID uuid.UUID `json:"aggregateId"`
	OccurredAt  time.Time `json:"occurredAt"`
	Version    int       `json:"version"`
}

// NewTenantCreatedEvent 创建租户创建事件
func NewTenantCreatedEvent(tenantID, ownerID uuid.UUID, tenantName string) *TenantCreatedEvent {
	return &TenantCreatedEvent{
		TenantID:   tenantID,
		OwnerID:    ownerID,
		TenantName: tenantName,
		CreatedAt:  time.Now(),
		EventID:     uuid.New().String(),
		EventType:   "TenantCreated",
		AggregateID: tenantID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

// MemberJoinedEvent 成员加入租户事件
type MemberJoinedEvent struct {
	TenantID  uuid.UUID      `json:"tenantId"`
	UserID    uuid.UUID      `json:"userId"`
	Role    string         `json:"role"`
	JoinedAt  time.Time     `json:"joinedAt"`
	EventID    string    `json:"eventId"`
	EventType  string    `json:"eventType"`
	AggregateID uuid.UUID `json:"aggregateId"`
	OccurredAt  time.Time `json:"occurredAt"`
	Version    int       `json:"version"`
}

// NewMemberJoinedEvent 创建成员加入事件
func NewMemberJoinedEvent(tenantID, userID uuid.UUID, role string) *MemberJoinedEvent {
	return &MemberJoinedEvent{
		TenantID: tenantID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
		EventID:     uuid.New().String(),
		EventType:   "MemberJoined",
		AggregateID: tenantID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

// MemberLeftEvent 成员离开租户事件
type MemberLeftEvent struct {
	TenantID  uuid.UUID `json:"tenantId"`
	UserID    uuid.UUID `json:"userId"`
	LeftAt    time.Time `json:"leftAt"`
	EventID    string    `json:"eventId"`
	EventType  string    `json:"eventType"`
	AggregateID uuid.UUID `json:"aggregateId"`
	OccurredAt  time.Time `json:"occurredAt"`
	Version    int       `json:"version"`
}

// NewMemberLeftEvent 创建成员离开事件
func NewMemberLeftEvent(tenantID, userID uuid.UUID) *MemberLeftEvent {
	return &MemberLeftEvent{
		TenantID: tenantID,
		UserID:   userID,
		LeftAt:   time.Now(),
		EventID:     uuid.New().String(),
		EventType:   "MemberLeft",
		AggregateID: tenantID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

// RoleChangedEvent 角色变更事件
type RoleChangedEvent struct {
	TenantID    uuid.UUID `json:"tenantId"`
	UserID      uuid.UUID `json:"userId"`
	OldRole   string    `json:"oldRole"`
	NewRole   string    `json:"newRole"`
	ChangedAt   time.Time `json:"changedAt"`
	EventID    string    `json:"eventId"`
	EventType  string    `json:"eventType"`
	AggregateID uuid.UUID `json:"aggregateId"`
	OccurredAt  time.Time `json:"occurredAt"`
	Version    int       `json:"version"`
}

// NewRoleChangedEvent 创建角色变更事件
func NewRoleChangedEvent(tenantID, userID uuid.UUID, oldRole, newRole string) *RoleChangedEvent {
	return &RoleChangedEvent{
		TenantID:  tenantID,
		UserID:    userID,
		OldRole:   oldRole,
		NewRole:   newRole,
		ChangedAt: time.Now(),
		EventID:     uuid.New().String(),
		EventType:   "RoleChanged",
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

func (e *MemberJoinedEvent) GetEventType() string      { return e.EventType }
func (e *MemberJoinedEvent) GetEventID() string        { return e.EventID }
func (e *MemberJoinedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *MemberJoinedEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *MemberJoinedEvent) GetVersion() int           { return e.Version }

func (e *MemberLeftEvent) GetEventType() string      { return e.EventType }
func (e *MemberLeftEvent) GetEventID() string        { return e.EventID }
func (e *MemberLeftEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *MemberLeftEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *MemberLeftEvent) GetVersion() int           { return e.Version }

func (e *RoleChangedEvent) GetEventType() string      { return e.EventType }
func (e *RoleChangedEvent) GetEventID() string        { return e.EventID }
func (e *RoleChangedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *RoleChangedEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *RoleChangedEvent) GetVersion() int           { return e.Version }

// 兼容旧版本的事件构造函数
func NewTenantMemberAddedEvent(tenantID, userID, memberID uuid.UUID, role string) *MemberJoinedEvent {
	return &MemberJoinedEvent{
		TenantID:   tenantID,
		UserID:     userID,
		Role:       role,
		JoinedAt:   time.Now(),
		EventID:   uuid.New().String(),
		EventType:  "TenantMemberAdded",
		AggregateID: tenantID,
		OccurredAt: time.Now(),
		Version:    1,
	}
}

func NewTenantMemberRemovedEvent(tenantID, userID, memberID uuid.UUID) *MemberLeftEvent {
	return &MemberLeftEvent{
		TenantID:   tenantID,
		UserID:     userID,
		LeftAt:     time.Now(),
		EventID:   uuid.New().String(),
		EventType:  "TenantMemberRemoved",
		AggregateID: tenantID,
		OccurredAt: time.Now(),
		Version:    1,
	}
}
