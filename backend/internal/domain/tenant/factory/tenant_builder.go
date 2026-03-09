package factory

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/user/entity"
)

// TenantBuilder 租户构建器（Factory 模式）
type TenantBuilder struct {
	name        string
	description string
	ownerID     uuid.UUID
	maxMembers int
	expiredAt   time.Time
}

// NewTenantBuilder 创建租户构建器
func NewTenantBuilder(ownerID uuid.UUID, name string) *TenantBuilder {
	return &TenantBuilder{
		ownerID:    ownerID,
		name:       name,
		maxMembers: 10,                          // 默认最大成员数
		expiredAt:  time.Now().AddDate(1, 0, 0), // 默认 1 年有效期
	}
}

// WithDescription 设置描述信息
func (b *TenantBuilder) WithDescription(desc string) *TenantBuilder {
	b.description = desc
	return b
}

// WithMaxMembers 设置最大成员数
func (b *TenantBuilder) WithMaxMembers(count int) *TenantBuilder {
	if count <= 0 {
		panic("max members must be positive")
	}
	b.maxMembers = count
	return b
}

// WithExpiredAt 设置过期时间
func (b *TenantBuilder) WithExpiredAt(expiredAt time.Time) *TenantBuilder {
	b.expiredAt = expiredAt
	return b
}

// Build 构建租户实例
// 保证聚合根的不变性和有效性
func (b *TenantBuilder) Build() (*entity.Tenant, error) {
	// 验证不变量
	if b.name == "" {
		return nil, errors.New("tenant name is required")
	}

	if b.ownerID == uuid.Nil {
		return nil, errors.New("owner ID is required")
	}

	now := time.Now()
	if b.expiredAt.Before(now) {
		return nil, errors.New("expired_at must be in the future")
	}

	return &entity.Tenant{
		ID:          uuid.New(),
		Name:        b.name,
		Description: b.description,
		MaxMembers:  b.maxMembers,
		ExpiredAt:  b.expiredAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// MustBuild 构建租户实例，失败则 panic
func (b *TenantBuilder) MustBuild() *entity.Tenant {
	tenant, err := b.Build()
	if err != nil {
		panic(err)
	}
	return tenant
}
