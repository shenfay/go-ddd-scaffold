// Package repository 租户领域仓储接口定义
package repository

import (
	"context"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/tenant/entity"
)

// TenantRepository 租户仓储接口
type TenantRepository interface {
	Create(ctx context.Context, tenant *entity.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
	Update(ctx context.Context, tenant *entity.Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListAll(ctx context.Context) ([]*entity.Tenant, error)
}

// TenantMemberRepository 租户成员仓储接口
type TenantMemberRepository interface {
	Create(ctx context.Context, member *entity.TenantMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TenantMember, error)
	GetByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) (*entity.TenantMember, error)
	Update(ctx context.Context, member *entity.TenantMember) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.TenantMember, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.TenantMember, error)
	DeleteByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) error
}
