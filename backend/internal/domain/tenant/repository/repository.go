// Package repository 租户领域仓储接口定义
package repository

import (
	"context"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/user/entity"
)

// TenantRepository 租户仓储接口
type TenantRepository interface {
	Create(ctx context.Context, tenant *entity.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
	Update(ctx context.Context, tenant *entity.Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListAll(ctx context.Context) ([]*entity.Tenant, error)
}
