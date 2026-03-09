// Package repository 用户领域仓储接口定义
package repository

import (
	"context"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/user/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	BaseRepository[entity.User, uuid.UUID]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error)
	CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// TenantRepository 租户仓储接口
type TenantRepository interface {
	BaseRepository[entity.Tenant, uuid.UUID]
}

// TenantMemberRepository 租户成员仓储接口
type TenantMemberRepository interface {
	BaseRepository[entity.TenantMember, uuid.UUID]
	GetByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) (*entity.TenantMember, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.TenantMember, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.TenantMember, error)
	DeleteByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) error
}
