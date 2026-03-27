package aggregate

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/valueobject"
	uservo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// Tenant 租户聚合根
type Tenant struct {
	common.BaseEntity

	code        *valueobject.TenantCode
	name        string
	description string
	status      valueobject.TenantStatus
	config      *valueobject.TenantConfig
	ownerID     uservo.UserID
	maxMembers  int
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTenant 创建新租户
func NewTenant(code, name string, ownerID uservo.UserID) (*Tenant, error) {
	tenant := &Tenant{
		name:       name,
		status:     valueobject.TenantStatusActive,
		ownerID:    ownerID,
		maxMembers: 100, // 默认最大成员数
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}

	// 设置初始ID
	tenant.SetID(valueobject.NewTenantID(1))

	// 验证和设置租户编码
	tc, err := valueobject.NewTenantCode(code)
	if err != nil {
		return nil, err
	}
	tenant.code = tc

	// 设置默认配置
	tenant.config = valueobject.NewDefaultTenantConfig()

	// 发布租户创建事件
	evt := event.NewTenantCreatedEvent(tenant.ID().(valueobject.TenantID), code, name, ownerID)
	tenant.ApplyEvent(evt)

	return tenant, nil
}

// Code 获取租户编码
func (t *Tenant) Code() *valueobject.TenantCode {
	return t.code
}

// Name 获取租户名称
func (t *Tenant) Name() string {
	return t.name
}

// Description 获取租户描述
func (t *Tenant) Description() string {
	return t.description
}

// Status 获取租户状态
func (t *Tenant) Status() valueobject.TenantStatus {
	return t.status
}

// Config 获取租户配置
func (t *Tenant) Config() *valueobject.TenantConfig {
	return t.config
}

// OwnerID 获取所有者ID
func (t *Tenant) OwnerID() uservo.UserID {
	return t.ownerID
}

// MaxMembers 获取最大成员数
func (t *Tenant) MaxMembers() int {
	return t.maxMembers
}

// SetName 设置租户名称
func (t *Tenant) SetName(name string) {
	oldName := t.name
	t.name = name
	t.updatedAt = time.Now()
	t.IncrementVersion()

	evt := event.NewTenantNameChangedEvent(t.ID().(valueobject.TenantID), oldName, name)
	t.ApplyEvent(evt)
}

// SetDescription 设置租户描述
func (t *Tenant) SetDescription(description string) {
	t.description = description
	t.updatedAt = time.Now()
	t.IncrementVersion()
}

// SetMaxMembers 设置最大成员数
func (t *Tenant) SetMaxMembers(maxMembers int) error {
	if maxMembers <= 0 {
		return common.NewBusinessError(CodeInvalidMaxMembers, "max members must be greater than 0")
	}

	t.maxMembers = maxMembers
	t.updatedAt = time.Now()
	t.IncrementVersion()

	evt := event.NewTenantConfigChangedEvent(t.ID().(valueobject.TenantID), "max_members", maxMembers)
	t.ApplyEvent(evt)

	return nil
}

// Activate 激活租户
func (t *Tenant) Activate() error {
	if t.status == valueobject.TenantStatusActive {
		return common.NewBusinessError(CodeTenantAlreadyActive, "tenant is already active")
	}

	t.status = valueobject.TenantStatusActive
	t.updatedAt = time.Now()
	t.IncrementVersion()

	evt := event.NewTenantActivatedEvent(t.ID().(valueobject.TenantID))
	t.ApplyEvent(evt)

	return nil
}

// Deactivate 停用租户
func (t *Tenant) Deactivate(reason string) error {
	if t.status == valueobject.TenantStatusInactive {
		return common.NewBusinessError(CodeTenantAlreadyInactive, "tenant is already inactive")
	}

	t.status = valueobject.TenantStatusInactive
	t.updatedAt = time.Now()
	t.IncrementVersion()

	evt := event.NewTenantDeactivatedEvent(t.ID().(valueobject.TenantID), reason)
	t.ApplyEvent(evt)

	return nil
}

// Suspend 暂停租户
func (t *Tenant) Suspend(reason string) error {
	if t.status == valueobject.TenantStatusSuspended {
		return common.NewBusinessError(CodeTenantAlreadySuspended, "tenant is already suspended")
	}

	t.status = valueobject.TenantStatusSuspended
	t.updatedAt = time.Now()
	t.IncrementVersion()

	evt := event.NewTenantSuspendedEvent(t.ID().(valueobject.TenantID), reason)
	t.ApplyEvent(evt)

	return nil
}

// UpdateConfig 更新租户配置
func (t *Tenant) UpdateConfig(config *valueobject.TenantConfig) {
	t.config = config
	t.updatedAt = time.Now()
	t.IncrementVersion()

	evt := event.NewTenantConfigChangedEvent(t.ID().(valueobject.TenantID), "config", config)
	t.ApplyEvent(evt)
}

// IsActive 检查租户是否活跃
func (t *Tenant) IsActive() bool {
	return t.status == valueobject.TenantStatusActive
}

// CanAddMember 检查是否可以添加成员
func (t *Tenant) CanAddMember(currentMemberCount int) bool {
	if !t.IsActive() {
		return false
	}
	return currentMemberCount < t.maxMembers
}

// TransferOwnership 转移租户所有权
func (t *Tenant) TransferOwnership(newOwnerID uservo.UserID) error {
	if t.ownerID.Equals(newOwnerID) {
		return common.NewBusinessError(CodeInvalidOperation, "new owner is already the current owner")
	}

	t.ownerID = newOwnerID
	t.updatedAt = time.Now()
	t.IncrementVersion()

	evt := event.NewTenantOwnershipTransferredEvent(t.ID().(valueobject.TenantID), newOwnerID)
	t.ApplyEvent(evt)

	return nil
}
