// Package aggregate 租户聚合根
package aggregate

import (
	"time"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/user/entity"
)

// TenantAggregate 租户聚合根
type TenantAggregate struct {
	tenant *entity.Tenant
}

// NewTenantAggregate 创建租户聚合根
func NewTenantAggregate(tenant *entity.Tenant) *TenantAggregate {
	return &TenantAggregate{tenant: tenant}
}

// GetTenant 获取租户实体
func (a *TenantAggregate) GetTenant() *entity.Tenant {
	return a.tenant
}

// GetID 获取聚合根ID
func (a *TenantAggregate) GetID() uuid.UUID {
	return a.tenant.ID
}

// IsValid 检查租户是否有效
func (a *TenantAggregate) IsValid() bool {
	return a.tenant.IsValid()
}

// IsExpired 检查租户是否过期
func (a *TenantAggregate) IsExpired() bool {
	return a.tenant.IsExpired()
}

// CanAddMoreMembers 检查是否还能添加更多成员
func (a *TenantAggregate) CanAddMoreMembers(currentCount int) bool {
	return a.tenant.CanAddMoreMembers(currentCount)
}

// UpdateInfo 更新租户信息
func (a *TenantAggregate) UpdateInfo(name, description string) {
	a.tenant.Name = name
	a.tenant.Description = description
	a.tenant.UpdatedAt = time.Now()
}

// ExtendExpiration 延长租户过期时间
func (a *TenantAggregate) ExtendExpiration(newExpiredAt time.Time) {
	a.tenant.ExpiredAt = newExpiredAt
	a.tenant.UpdatedAt = time.Now()
}
