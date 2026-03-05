// Package repo 租户模块DAO仓储实现
package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/dao"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/model"
	cast "go-ddd-scaffold/pkg/uitl"
)

// TenantDAORepository 租户 DAO 仓储实现
type TenantDAORepository struct {
	db      *gorm.DB
	querier *dao.Query
}

// NewTenantDAORepository 创建租户 DAO 仓储实例
func NewTenantDAORepository(db *gorm.DB) repository.TenantRepository {
	return &TenantDAORepository{
		db:      db,
		querier: dao.Use(db),
	}
}

// Create 创建租户
func (r *TenantDAORepository) Create(ctx context.Context, t *entity.Tenant) error {
	id := t.ID.String()
	maxMembers := int32(t.MaxMembers)

	tenantModel := &model.Tenant{
		ID:          &id,
		Name:        t.Name,
		Description: cast.ToStringPtr(t.Description),
		ExpiredAt:   t.ExpiredAt,
		MaxMembers:  &maxMembers,
	}

	return r.querier.Tenant.WithContext(ctx).Create(tenantModel)
}

// GetByID 根据ID获取租户
func (r *TenantDAORepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error) {
	tenantModel, err := r.querier.Tenant.WithContext(ctx).Where(r.querier.Tenant.ID.Eq(id.String())).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to get tenant by id: %w", err)
	}

	return r.toEntity(tenantModel), nil
}

// Update 更新租户
func (r *TenantDAORepository) Update(ctx context.Context, t *entity.Tenant) error {
	id := t.ID.String()

	maxMembers := int32(t.MaxMembers)
	tenantModel := &model.Tenant{
		ID:          &id,
		Name:        t.Name,
		Description: cast.ToStringPtr(t.Description),
		ExpiredAt:   t.ExpiredAt,
		MaxMembers:  &maxMembers,
	}

	_, err := r.querier.Tenant.WithContext(ctx).Where(r.querier.Tenant.ID.Eq(id)).Updates(tenantModel)
	return err
}

// Delete 删除租户
func (r *TenantDAORepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.querier.Tenant.WithContext(ctx).Where(r.querier.Tenant.ID.Eq(id.String())).Delete()
	return err
}

// ListAll 列出所有租户
func (r *TenantDAORepository) ListAll(ctx context.Context) ([]*entity.Tenant, error) {
	models, err := r.querier.Tenant.WithContext(ctx).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to list all tenants: %w", err)
	}

	tenants := make([]*entity.Tenant, len(models))
	for i, m := range models {
		tenants[i] = r.toEntity(m)
	}

	return tenants, nil
}

// toEntity 将模型转换为实体
func (r *TenantDAORepository) toEntity(tenantModel *model.Tenant) *entity.Tenant {
	id, _ := cast.ToUUID(*tenantModel.ID)

	tenant := &entity.Tenant{
		ID:          id,
		Name:        tenantModel.Name,
		Description: "",
		ExpiredAt:   tenantModel.ExpiredAt,
		CreatedAt:   *tenantModel.CreatedAt,
		UpdatedAt:   *tenantModel.UpdatedAt,
	}

	if tenantModel.Description != nil {
		tenant.Description = *tenantModel.Description
	}

	if tenantModel.MaxMembers != nil {
		tenant.MaxMembers = int(*tenantModel.MaxMembers)
	}

	return tenant
}
