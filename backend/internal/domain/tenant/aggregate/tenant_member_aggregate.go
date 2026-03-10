// Package aggregate 租户领域聚合根定义
package aggregate

import (
	"time"

	"github.com/google/uuid"

	shared_entity"go-ddd-scaffold/internal/domain/shared/entity"
	tenant_entity"go-ddd-scaffold/internal/domain/tenant/entity"
	tenant_event "go-ddd-scaffold/internal/domain/tenant/event"
)

// TenantMemberAggregate 租户成员聚合根
type TenantMemberAggregate struct {
	*tenant_entity.TenantMember
	events []tenant_entity.DomainEvent
}

// NewTenantMemberAggregate 创建租户成员聚合根
func NewTenantMemberAggregate(tenantID, userID uuid.UUID, role shared_entity.UserRole, invitedBy *uuid.UUID) *TenantMemberAggregate {
	member := &tenant_entity.TenantMember{
		ID:        uuid.New(),
		TenantID: tenantID,
		UserID:  userID,
		Role:      role,
		Status:    tenant_entity.MemberStatusActive,
		InvitedBy: invitedBy,
		JoinedAt:  time.Now(),
	}

	aggregate := &TenantMemberAggregate{
		TenantMember: member,
		events:       make([]tenant_entity.DomainEvent, 0),
	}

	// 记录成员加入事件
	aggregate.recordEvent(tenant_event.NewMemberJoinedEvent(tenantID, userID, string(role)))

	return aggregate
}

// TenantMemberAggregateFromExisting 从已有成员创建聚合根（用于从数据库加载）
func TenantMemberAggregateFromExisting(member *tenant_entity.TenantMember) *TenantMemberAggregate {
	return &TenantMemberAggregate{
		TenantMember: member,
		events:       make([]tenant_entity.DomainEvent, 0),
	}
}

// Activate 激活成员
func (a *TenantMemberAggregate) Activate() error {
	if a.Status == tenant_entity.MemberStatusActive {
		return nil
	}

	a.Status = tenant_entity.MemberStatusActive
	a.UpdatedAt = time.Now()

	return nil
}

// Deactivate 停用成员
func (a *TenantMemberAggregate) Deactivate() error {
	if a.Status == tenant_entity.MemberStatusInactive {
		return nil
	}

	a.Status = tenant_entity.MemberStatusInactive
	a.UpdatedAt = time.Now()

	return nil
}

// Remove 移除成员
func(a *TenantMemberAggregate) Remove() error {
	if a.IsRemoved() {
		return nil
	}

	a.Status = tenant_entity.MemberStatusRemoved
	a.LeftAt = func() *time.Time {
		t := time.Now()
		return &t
	}()
	a.UpdatedAt = time.Now()

	a.recordEvent(tenant_event.NewMemberLeftEvent(a.TenantID, a.UserID))

	return nil
}

// ChangeRole 变更成员角色
func (a *TenantMemberAggregate) ChangeRole(newRole shared_entity.UserRole) error {
	if a.Role == newRole {
		return nil
	}

	oldRole := a.Role
	a.Role = newRole
	a.UpdatedAt = time.Now()

	a.recordEvent(tenant_event.NewRoleChangedEvent(a.TenantID, a.UserID, string(oldRole), string(newRole)))

	return nil
}

// InviteMember 邀请新成员
func (a *TenantMemberAggregate) InviteMember(targetUserID uuid.UUID, role shared_entity.UserRole) (*TenantMemberAggregate, error) {
	if !a.CanInviteMembers() {
		return nil, tenant_entity.ErrUnauthorized
	}

	newMember := NewTenantMemberAggregate(a.TenantID, targetUserID, role, &a.UserID)
	return newMember, nil
}

// CanInviteMembers 检查是否可以邀请成员
func(a *TenantMemberAggregate) CanInviteMembers() bool {
	return a.IsActive() && (a.Role == shared_entity.RoleOwner || a.Role == shared_entity.RoleAdmin)
}

// IsActive 检查成员是否活跃
func(a *TenantMemberAggregate) IsActive() bool {
	return a.Status == tenant_entity.MemberStatusActive
}

// IsOwner 检查是否是所有者
func(a *TenantMemberAggregate) IsOwner() bool {
	return a.Role == shared_entity.RoleOwner
}

// recordEvent 记录领域事件
func(a *TenantMemberAggregate) recordEvent(evt tenant_entity.DomainEvent) {
	a.events = append(a.events, evt)
}

// GetEvents 获取所有未处理的领域事件
func(a *TenantMemberAggregate) GetEvents() []tenant_entity.DomainEvent {
	return a.events
}

// ClearEvents 清空已处理的领域事件
func (a *TenantMemberAggregate) ClearEvents() {
	a.events = make([]tenant_entity.DomainEvent, 0)
}
