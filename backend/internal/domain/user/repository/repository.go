// Package repository 用户领域仓储接口定义
package repository

import (
	"context"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/user/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error)
	ListChildrenByParent(ctx context.Context, parentID uuid.UUID) ([]*entity.User, error)
	CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// TenantRepository 租户仓储接口
type TenantRepository interface {
	Create(ctx context.Context, tenant *entity.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
	Update(ctx context.Context, tenant *entity.Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
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
