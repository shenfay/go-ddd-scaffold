package aggregate

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/valueobject"
	uservo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// Tenant 租户聚合根
// 使用组合模式替代继承
type Tenant struct {
	meta        *common.EntityMeta // 组合元数据
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
		meta:       common.NewEntityMeta(nil, time.Now()),
		name:       name,
		status:     valueobject.TenantStatusActive,
		ownerID:    ownerID,
		maxMembers: 100, // 默认最大成员数
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}

	// 设置初始 ID
	tenant.meta.SetID(valueobject.NewTenantID(1))

	// 验证和设置租户编码
	tc, err := valueobject.NewTenantCode(code)
	if err != nil {
		return nil, err
	}
	tenant.code = tc

	// 设置默认配置
	tenant.config = valueobject.NewDefaultTenantConfig()

	// 发布租户创建事件
	evt := event.NewTenantCreatedEvent(tenant.meta.ID().(valueobject.TenantID), code, name, ownerID)
	tenant.meta.ApplyEvent(evt)

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

// ID 获取租户 ID (实现 AggregateRoot 接口)
func (t *Tenant) ID() interface{} {
	return t.meta.ID()
}

// Version 获取版本号 (实现 AggregateRoot 接口)
func (t *Tenant) Version() int {
	return t.meta.Version()
}

// IncrementVersion 增加版本号 (实现 AggregateRoot 接口)
func (t *Tenant) IncrementVersion() {
	t.meta.IncrementVersion()
}

// ApplyEvent 应用领域事件 (实现 AggregateRoot 接口)
func (t *Tenant) ApplyEvent(event common.DomainEvent) {
	t.meta.ApplyEvent(event)
}

// GetUncommittedEvents 获取未提交事件 (实现 AggregateRoot 接口)
func (t *Tenant) GetUncommittedEvents() []common.DomainEvent {
	return t.meta.GetUncommittedEvents()
}

// ClearUncommittedEvents 清除已提交事件 (实现 AggregateRoot 接口)
func (t *Tenant) ClearUncommittedEvents() {
	t.meta.ClearUncommittedEvents()
}

// CreatedAt 获取创建时间
func (t *Tenant) CreatedAt() time.Time {
	return t.meta.CreatedAt()
}

// UpdatedAt 获取更新时间
func (t *Tenant) UpdatedAt() time.Time {
	return t.meta.UpdatedAt()
}

// SetCreatedAt 设置创建时间 (用于 Builder)
func (t *Tenant) SetCreatedAt(time time.Time) {
	t.meta.SetCreatedAt(time)
}

// SetUpdatedAt 设置更新时间 (用于 Builder)
func (t *Tenant) SetUpdatedAt(time time.Time) {
	t.meta.SetUpdatedAt(time)
}

// SetName 设置租户名称
func (t *Tenant) SetName(name string) {
	oldName := t.name
	t.name = name
	t.meta.SetUpdatedAt(time.Now())
	t.meta.IncrementVersion()

	evt := event.NewTenantNameChangedEvent(t.meta.ID().(valueobject.TenantID), oldName, name)
	t.meta.ApplyEvent(evt)
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
	t.meta.SetUpdatedAt(time.Now())
	t.meta.IncrementVersion()

	evt := event.NewTenantConfigChangedEvent(t.meta.ID().(valueobject.TenantID), "max_members", maxMembers)
	t.meta.ApplyEvent(evt)

	return nil
}

// Activate 激活租户
func (t *Tenant) Activate() error {
	if t.status == valueobject.TenantStatusActive {
		return common.NewBusinessError(CodeTenantAlreadyActive, "tenant is already active")
	}

	t.status = valueobject.TenantStatusActive
	t.meta.SetUpdatedAt(time.Now())
	t.meta.IncrementVersion()

	evt := event.NewTenantActivatedEvent(t.meta.ID().(valueobject.TenantID))
	t.meta.ApplyEvent(evt)

	return nil
}

// Deactivate 停用租户
func (t *Tenant) Deactivate(reason string) error {
	if t.status == valueobject.TenantStatusInactive {
		return common.NewBusinessError(CodeTenantAlreadyInactive, "tenant is already inactive")
	}

	t.status = valueobject.TenantStatusInactive
	t.meta.SetUpdatedAt(time.Now())
	t.meta.IncrementVersion()

	evt := event.NewTenantDeactivatedEvent(t.meta.ID().(valueobject.TenantID), reason)
	t.meta.ApplyEvent(evt)

	return nil
}

// Suspend 暂停租户
func (t *Tenant) Suspend(reason string) error {
	if t.status == valueobject.TenantStatusSuspended {
		return common.NewBusinessError(CodeTenantAlreadySuspended, "tenant is already suspended")
	}

	t.status = valueobject.TenantStatusSuspended
	t.meta.SetUpdatedAt(time.Now())
	t.meta.IncrementVersion()

	evt := event.NewTenantSuspendedEvent(t.meta.ID().(valueobject.TenantID), reason)
	t.meta.ApplyEvent(evt)

	return nil
}

// UpdateConfig 更新租户配置
func (t *Tenant) UpdateConfig(config *valueobject.TenantConfig) {
	t.config = config
	t.meta.SetUpdatedAt(time.Now())
	t.meta.IncrementVersion()

	evt := event.NewTenantConfigChangedEvent(t.meta.ID().(valueobject.TenantID), "config", config)
	t.meta.ApplyEvent(evt)
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
	t.meta.SetUpdatedAt(time.Now())
	t.meta.IncrementVersion()

	evt := event.NewTenantOwnershipTransferredEvent(t.meta.ID().(valueobject.TenantID), newOwnerID)
	t.meta.ApplyEvent(evt)

	return nil
}
