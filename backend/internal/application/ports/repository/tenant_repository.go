package repository

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/aggregate"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/valueobject"
)

// TenantRepository 租户仓储端口
type TenantRepository interface {
	// Save 保存租户
	Save(ctx context.Context, tenant *aggregate.Tenant) error

	// FindByID 根据 ID 查找租户
	FindByID(ctx context.Context, id vo.TenantID) (*aggregate.Tenant, error)

	// FindByCode 根据租户编码查找
	FindByCode(ctx context.Context, code vo.TenantCode) (*aggregate.Tenant, error)

	// Update 更新租户
	Update(ctx context.Context, tenant *aggregate.Tenant) error

	// Delete 删除租户
	Delete(ctx context.Context, id vo.TenantID) error
}
