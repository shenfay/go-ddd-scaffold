package entity

import (
	"time"

	sharedEntity "go-ddd-scaffold/internal/domain/shared/entity"
	"go-ddd-scaffold/internal/domain/tenant/event"

	"github.com/google/uuid"
)

// Tenant 租户聚合根（多租户 SaaS 场景）
// 纯领域对象，不包含任何基础设施标签
type Tenant struct {
	ID          uuid.UUID
	Name        string
	Description string

	// 租户配置
	MaxMembers int       // 最大成员数
	ExpiredAt  time.Time // 过期时间

	Members []TenantMember `json:"members,omitempty"` // 内部管理的成员列表

	// 领域事件（临时存储，由 Application Service 发布后清空）
	events []DomainEvent

	CreatedAt time.Time
	UpdatedAt time.Time
}

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetEventType() string
	GetEventID() string
	GetAggregateID() uuid.UUID
	GetOccurredAt() time.Time
	GetVersion() int
}

// NewTenant 创建新的租户实例
func NewTenant(name string, maxMembers int) *Tenant {
	return &Tenant{
		ID:          uuid.New(),
		Name:        name,
		Description: "",
		MaxMembers:  maxMembers,
		ExpiredAt:   time.Now().AddDate(1, 0, 0), // 默认一年后过期
		Members:     make([]TenantMember, 0),
	}
}

// IsValid 检查租户是否有效
func (t *Tenant) IsValid() bool {
	return t.ExpiredAt.After(time.Now()) && t.MaxMembers > 0
}

// IsExpired 检查租户是否过期
func (t *Tenant) IsExpired() bool {
	return t.ExpiredAt.Before(time.Now())
}

// CanAddMoreMembers 检查是否还能添加更多成员
func (t *Tenant) CanAddMoreMembers(currentCount int) bool {
	return currentCount < t.MaxMembers
}

// AddMember 添加成员到租户（聚合根方法）
func (t *Tenant) AddMember(userID uuid.UUID, role sharedEntity.UserRole, invitedBy *uuid.UUID) (*TenantMember, error) {
	if !t.IsValid() {
		return nil, ErrTenantInvalid
	}

	if len(t.Members) >= t.MaxMembers {
		return nil, ErrTenantMemberLimitExceeded
	}

	// 检查用户是否已是成员
	for _, member := range t.Members {
		if member.UserID == userID && member.Status == MemberStatusActive {
			return nil, ErrTenantMemberAlreadyExists
		}
	}

	member := TenantMember{
		ID:        uuid.New(),
		TenantID:  t.ID,
		UserID:    userID,
		Role:      role,
		Status:    MemberStatusActive,
		InvitedBy: invitedBy,
		JoinedAt:  time.Now(),
	}

	t.Members = append(t.Members, member)

	// 发布领域事件
	t.addEvent(event.NewTenantMemberAddedEvent(t.ID, userID, member.ID, string(role)))

	return &member, nil
}

// RemoveMember 移除租户成员（聚合根方法）
func (t *Tenant) RemoveMember(memberID uuid.UUID) error {
	for i, member := range t.Members {
		if member.ID == memberID {
			t.Members[i].Status = MemberStatusRemoved
			leftAt := time.Now()
			t.Members[i].LeftAt = &leftAt

			// 发布领域事件
			t.addEvent(event.NewTenantMemberRemovedEvent(t.ID, member.UserID, member.ID))

			return nil
		}
	}
	return ErrTenantMemberNotFound
}

// Events 获取待发布的领域事件
func (t *Tenant) Events() []DomainEvent {
	return t.events
}

// ClearEvents 清空已发布的领域事件
func (t *Tenant) ClearEvents() {
	t.events = nil
}

// addEvent 添加领域事件（内部方法）
func (t *Tenant) addEvent(event DomainEvent) {
	t.events = append(t.events, event)
}

// GetActiveMemberCount 获取活跃成员数量
func (t *Tenant) GetActiveMemberCount() int {
	count := 0
	for _, member := range t.Members {
		if member.Status == MemberStatusActive {
			count++
		}
	}
	return count
}
