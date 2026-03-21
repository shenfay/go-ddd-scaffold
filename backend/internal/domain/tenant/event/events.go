package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/valueobject"
	uservo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// TenantCreatedEvent 租户创建事件
type TenantCreatedEvent struct {
	*kernel.BaseEvent
	TenantID  valueobject.TenantID `json:"tenant_id"`
	Code      string               `json:"code"`
	Name      string               `json:"name"`
	OwnerID   uservo.UserID        `json:"owner_id"`
	CreatedAt time.Time            `json:"created_at"`
}

// NewTenantCreatedEvent 创建租户创建事件
func NewTenantCreatedEvent(tenantID valueobject.TenantID, code, name string, ownerID uservo.UserID) *TenantCreatedEvent {
	event := &TenantCreatedEvent{
		BaseEvent: kernel.NewBaseEvent("TenantCreated", tenantID, 1),
		TenantID:  tenantID,
		Code:      code,
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantActivatedEvent 租户激活事件
type TenantActivatedEvent struct {
	*kernel.BaseEvent
	TenantID    valueobject.TenantID `json:"tenant_id"`
	ActivatedAt time.Time            `json:"activated_at"`
}

// NewTenantActivatedEvent 创建租户激活事件
func NewTenantActivatedEvent(tenantID valueobject.TenantID) *TenantActivatedEvent {
	event := &TenantActivatedEvent{
		BaseEvent:   kernel.NewBaseEvent("TenantActivated", tenantID, 1),
		TenantID:    tenantID,
		ActivatedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantDeactivatedEvent 租户停用事件
type TenantDeactivatedEvent struct {
	*kernel.BaseEvent
	TenantID      valueobject.TenantID `json:"tenant_id"`
	Reason        string               `json:"reason"`
	DeactivatedAt time.Time            `json:"deactivated_at"`
}

// NewTenantDeactivatedEvent 创建租户停用事件
func NewTenantDeactivatedEvent(tenantID valueobject.TenantID, reason string) *TenantDeactivatedEvent {
	event := &TenantDeactivatedEvent{
		BaseEvent:     kernel.NewBaseEvent("TenantDeactivated", tenantID, 1),
		TenantID:      tenantID,
		Reason:        reason,
		DeactivatedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantSuspendedEvent 租户暂停事件
type TenantSuspendedEvent struct {
	*kernel.BaseEvent
	TenantID    valueobject.TenantID `json:"tenant_id"`
	Reason      string               `json:"reason"`
	SuspendedAt time.Time            `json:"suspended_at"`
}

// NewTenantSuspendedEvent 创建租户暂停事件
func NewTenantSuspendedEvent(tenantID valueobject.TenantID, reason string) *TenantSuspendedEvent {
	event := &TenantSuspendedEvent{
		BaseEvent:   kernel.NewBaseEvent("TenantSuspended", tenantID, 1),
		TenantID:    tenantID,
		Reason:      reason,
		SuspendedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantNameChangedEvent 租户名称变更事件
type TenantNameChangedEvent struct {
	*kernel.BaseEvent
	TenantID  valueobject.TenantID `json:"tenant_id"`
	OldName   string               `json:"old_name"`
	NewName   string               `json:"new_name"`
	ChangedAt time.Time            `json:"changed_at"`
}

// NewTenantNameChangedEvent 创建租户名称变更事件
func NewTenantNameChangedEvent(tenantID valueobject.TenantID, oldName, newName string) *TenantNameChangedEvent {
	event := &TenantNameChangedEvent{
		BaseEvent: kernel.NewBaseEvent("TenantNameChanged", tenantID, 1),
		TenantID:  tenantID,
		OldName:   oldName,
		NewName:   newName,
		ChangedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantConfigChangedEvent 租户配置变更事件
type TenantConfigChangedEvent struct {
	*kernel.BaseEvent
	TenantID    valueobject.TenantID `json:"tenant_id"`
	ConfigKey   string               `json:"config_key"`
	ConfigValue interface{}          `json:"config_value"`
	ChangedAt   time.Time            `json:"changed_at"`
}

// NewTenantConfigChangedEvent 创建租户配置变更事件
func NewTenantConfigChangedEvent(tenantID valueobject.TenantID, key string, value interface{}) *TenantConfigChangedEvent {
	event := &TenantConfigChangedEvent{
		BaseEvent:   kernel.NewBaseEvent("TenantConfigChanged", tenantID, 1),
		TenantID:    tenantID,
		ConfigKey:   key,
		ConfigValue: value,
		ChangedAt:   time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantMemberAddedEvent 租户成员添加事件
type TenantMemberAddedEvent struct {
	*kernel.BaseEvent
	TenantID valueobject.TenantID `json:"tenant_id"`
	UserID   uservo.UserID        `json:"user_id"`
	Role     string               `json:"role"`
	AddedBy  uservo.UserID        `json:"added_by"`
	AddedAt  time.Time            `json:"added_at"`
}

// NewTenantMemberAddedEvent 创建租户成员添加事件
func NewTenantMemberAddedEvent(tenantID valueobject.TenantID, userID, addedBy uservo.UserID, role valueobject.TenantRole) *TenantMemberAddedEvent {
	event := &TenantMemberAddedEvent{
		BaseEvent: kernel.NewBaseEvent("TenantMemberAdded", tenantID, 1),
		TenantID:  tenantID,
		UserID:    userID,
		Role:      role.String(),
		AddedBy:   addedBy,
		AddedAt:   time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantMemberRemovedEvent 租户成员移除事件
type TenantMemberRemovedEvent struct {
	*kernel.BaseEvent
	TenantID  valueobject.TenantID `json:"tenant_id"`
	UserID    uservo.UserID        `json:"user_id"`
	RemovedBy uservo.UserID        `json:"removed_by"`
	RemovedAt time.Time            `json:"removed_at"`
}

// NewTenantMemberRemovedEvent 创建租户成员移除事件
func NewTenantMemberRemovedEvent(tenantID valueobject.TenantID, userID, removedBy uservo.UserID) *TenantMemberRemovedEvent {
	event := &TenantMemberRemovedEvent{
		BaseEvent: kernel.NewBaseEvent("TenantMemberRemoved", tenantID, 1),
		TenantID:  tenantID,
		UserID:    userID,
		RemovedBy: removedBy,
		RemovedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantMemberRoleChangedEvent 租户成员角色变更事件
type TenantMemberRoleChangedEvent struct {
	*kernel.BaseEvent
	TenantID  valueobject.TenantID `json:"tenant_id"`
	UserID    uservo.UserID        `json:"user_id"`
	OldRole   string               `json:"old_role"`
	NewRole   string               `json:"new_role"`
	ChangedBy uservo.UserID        `json:"changed_by"`
	ChangedAt time.Time            `json:"changed_at"`
}

// NewTenantMemberRoleChangedEvent 创建租户成员角色变更事件
func NewTenantMemberRoleChangedEvent(tenantID valueobject.TenantID, userID, changedBy uservo.UserID, oldRole, newRole valueobject.TenantRole) *TenantMemberRoleChangedEvent {
	event := &TenantMemberRoleChangedEvent{
		BaseEvent: kernel.NewBaseEvent("TenantMemberRoleChanged", tenantID, 1),
		TenantID:  tenantID,
		UserID:    userID,
		OldRole:   oldRole.String(),
		NewRole:   newRole.String(),
		ChangedBy: changedBy,
		ChangedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}

// TenantOwnershipTransferredEvent 租户所有权转移事件
type TenantOwnershipTransferredEvent struct {
	*kernel.BaseEvent
	TenantID      valueobject.TenantID `json:"tenant_id"`
	NewOwnerID    uservo.UserID        `json:"new_owner_id"`
	TransferredAt time.Time            `json:"transferred_at"`
}

// NewTenantOwnershipTransferredEvent 创建租户所有权转移事件
func NewTenantOwnershipTransferredEvent(tenantID valueobject.TenantID, newOwnerID uservo.UserID) *TenantOwnershipTransferredEvent {
	event := &TenantOwnershipTransferredEvent{
		BaseEvent:     kernel.NewBaseEvent("TenantOwnershipTransferred", tenantID, 1),
		TenantID:      tenantID,
		NewOwnerID:    newOwnerID,
		TransferredAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "tenant")
	return event
}
