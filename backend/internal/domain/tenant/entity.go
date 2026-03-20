package tenant

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
)

// Tenant 租户聚合根
type Tenant struct {
	kernel.BaseEntity

	code        *TenantCode
	name        string
	description string
	status      TenantStatus
	config      *TenantConfig
	ownerID     vo.UserID
	maxMembers  int
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTenant 创建新租户
func NewTenant(code, name string, ownerID vo.UserID) (*Tenant, error) {
	tenant := &Tenant{
		name:       name,
		status:     TenantStatusActive,
		ownerID:    ownerID,
		maxMembers: 100, // 默认最大成员数
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}

	// 设置初始ID
	tenant.SetID(NewTenantID(1))

	// 验证和设置租户编码
	tc, err := NewTenantCode(code)
	if err != nil {
		return nil, err
	}
	tenant.code = tc

	// 设置默认配置
	tenant.config = NewDefaultTenantConfig()

	// 发布租户创建事件
	event := NewTenantCreatedEvent(tenant.ID().(TenantID), code, name, ownerID)
	tenant.ApplyEvent(event)

	return tenant, nil
}

// Code 获取租户编码
func (t *Tenant) Code() *TenantCode {
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
func (t *Tenant) Status() TenantStatus {
	return t.status
}

// Config 获取租户配置
func (t *Tenant) Config() *TenantConfig {
	return t.config
}

// OwnerID 获取所有者ID
func (t *Tenant) OwnerID() vo.UserID {
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

	event := NewTenantNameChangedEvent(t.ID().(TenantID), oldName, name)
	t.ApplyEvent(event)
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
		return kernel.NewBusinessError(kernel.CodeInvalidMaxMembers, "max members must be greater than 0")
	}

	t.maxMembers = maxMembers
	t.updatedAt = time.Now()
	t.IncrementVersion()

	event := NewTenantConfigChangedEvent(t.ID().(TenantID), "max_members", maxMembers)
	t.ApplyEvent(event)

	return nil
}

// Activate 激活租户
func (t *Tenant) Activate() error {
	if t.status == TenantStatusActive {
		return kernel.NewBusinessError(kernel.CodeTenantAlreadyActive, "tenant is already active")
	}

	t.status = TenantStatusActive
	t.updatedAt = time.Now()
	t.IncrementVersion()

	event := NewTenantActivatedEvent(t.ID().(TenantID))
	t.ApplyEvent(event)

	return nil
}

// Deactivate 停用租户
func (t *Tenant) Deactivate(reason string) error {
	if t.status == TenantStatusInactive {
		return kernel.NewBusinessError(kernel.CodeTenantAlreadyInactive, "tenant is already inactive")
	}

	t.status = TenantStatusInactive
	t.updatedAt = time.Now()
	t.IncrementVersion()

	event := NewTenantDeactivatedEvent(t.ID().(TenantID), reason)
	t.ApplyEvent(event)

	return nil
}

// Suspend 暂停租户
func (t *Tenant) Suspend(reason string) error {
	if t.status == TenantStatusSuspended {
		return kernel.NewBusinessError(kernel.CodeTenantAlreadySuspended, "tenant is already suspended")
	}

	t.status = TenantStatusSuspended
	t.updatedAt = time.Now()
	t.IncrementVersion()

	event := NewTenantSuspendedEvent(t.ID().(TenantID), reason)
	t.ApplyEvent(event)

	return nil
}

// UpdateConfig 更新租户配置
func (t *Tenant) UpdateConfig(config *TenantConfig) {
	t.config = config
	t.updatedAt = time.Now()
	t.IncrementVersion()

	event := NewTenantConfigChangedEvent(t.ID().(TenantID), "config", config)
	t.ApplyEvent(event)
}

// IsActive 检查租户是否活跃
func (t *Tenant) IsActive() bool {
	return t.status == TenantStatusActive
}

// CanAddMember 检查是否可以添加成员
func (t *Tenant) CanAddMember(currentMemberCount int) bool {
	if !t.IsActive() {
		return false
	}
	return currentMemberCount < t.maxMembers
}
