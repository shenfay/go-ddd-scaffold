// Package repo 租户成员模块DAO仓储实现
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
	"go-ddd-scaffold/pkg/converter"
)

// TenantMemberDAORepository 租户成员DAO仓储实现
type TenantMemberDAORepository struct {
	db        *gorm.DB
	querier   *dao.Query
	converter converter.Converter
}

// NewTenantMemberDAORepository 创建租户成员DAO仓储实例
func NewTenantMemberDAORepository(db *gorm.DB) repository.TenantMemberRepository {
	return &TenantMemberDAORepository{
		db:        db,
		querier:   dao.Use(db),
		converter: converter.NewConverter(),
	}
}

// Create 创建租户成员关系
func (r *TenantMemberDAORepository) Create(ctx context.Context, member *entity.TenantMember) error {
	memberID := member.ID.String()
	tenantID := member.TenantID.String()
	userID := member.UserID.String()
	var invitedBy *string
	if member.InvitedBy != nil {
		invitedByStr := member.InvitedBy.String()
		invitedBy = &invitedByStr
	}

	memberModel := &model.TenantMember{
		ID:        &memberID,
		TenantID:  tenantID,
		UserID:    userID,
		Role:      string(member.Role),
		Status:    r.converter.ToStringPtr(string(member.Status)),
		InvitedBy: invitedBy,
	}

	return r.querier.TenantMember.WithContext(ctx).Create(memberModel)
}

// GetByID 根据ID获取租户成员关系
func (r *TenantMemberDAORepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TenantMember, error) {
	memberModel, err := r.querier.TenantMember.WithContext(ctx).Where(r.querier.TenantMember.ID.Eq(id.String())).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant member not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to get tenant member by id: %w", err)
	}

	return r.toEntity(memberModel), nil
}

// GetByUserAndTenant 根据用户ID和租户ID获取租户成员关系
func (r *TenantMemberDAORepository) GetByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) (*entity.TenantMember, error) {
	memberModel, err := r.querier.TenantMember.WithContext(ctx).Where(
		r.querier.TenantMember.UserID.Eq(userID.String()),
		r.querier.TenantMember.TenantID.Eq(tenantID.String()),
	).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant member not found for user %s and tenant %s", userID.String(), tenantID.String())
		}
		return nil, fmt.Errorf("failed to get tenant member by user and tenant: %w", err)
	}

	return r.toEntity(memberModel), nil
}

// Update 更新租户成员关系
func (r *TenantMemberDAORepository) Update(ctx context.Context, member *entity.TenantMember) error {
	memberID := member.ID.String()
	tenantID := member.TenantID.String()
	userID := member.UserID.String()
	var invitedBy *string
	if member.InvitedBy != nil {
		invitedByStr := member.InvitedBy.String()
		invitedBy = &invitedByStr
	}

	memberModel := &model.TenantMember{
		ID:        &memberID,
		TenantID:  tenantID,
		UserID:    userID,
		Role:      string(member.Role),
		Status:    r.converter.ToStringPtr(string(member.Status)),
		InvitedBy: invitedBy,
	}

	_, err := r.querier.TenantMember.WithContext(ctx).Where(r.querier.TenantMember.ID.Eq(memberID)).Updates(memberModel)
	return err
}

// Delete 删除租户成员关系
func (r *TenantMemberDAORepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.querier.TenantMember.WithContext(ctx).Where(r.querier.TenantMember.ID.Eq(id.String())).Delete()
	return err
}

// ListByTenant 列出租户下所有成员
func (r *TenantMemberDAORepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.TenantMember, error) {
	models, err := r.querier.TenantMember.WithContext(ctx).Where(
		r.querier.TenantMember.TenantID.Eq(tenantID.String()),
		r.querier.TenantMember.Status.Neq("removed"),
	).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to list members by tenant: %w", err)
	}

	members := make([]*entity.TenantMember, len(models))
	for i, model := range models {
		members[i] = r.toEntity(model)
	}

	return members, nil
}

// ListByUser 列出用户在所有租户中的成员关系
func (r *TenantMemberDAORepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.TenantMember, error) {
	models, err := r.querier.TenantMember.WithContext(ctx).Where(
		r.querier.TenantMember.UserID.Eq(userID.String()),
		r.querier.TenantMember.Status.Neq("removed"),
	).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to list members by user: %w", err)
	}

	members := make([]*entity.TenantMember, len(models))
	for i, model := range models {
		members[i] = r.toEntity(model)
	}

	return members, nil
}

// DeleteByUserAndTenant 删除用户在指定租户中的成员关系
func (r *TenantMemberDAORepository) DeleteByUserAndTenant(ctx context.Context, userID, tenantID uuid.UUID) error {
	_, err := r.querier.TenantMember.WithContext(ctx).Where(
		r.querier.TenantMember.UserID.Eq(userID.String()),
		r.querier.TenantMember.TenantID.Eq(tenantID.String()),
	).Delete()
	return err
}

// toEntity 将模型转换为实体
func (r *TenantMemberDAORepository) toEntity(model *model.TenantMember) *entity.TenantMember {
	id, _ := r.converter.ToUUID(*model.ID)
	tenantID, _ := r.converter.ToUUID(model.TenantID)
	userID, _ := r.converter.ToUUID(model.UserID)

	member := &entity.TenantMember{
		ID:       id,
		TenantID: tenantID,
		UserID:   userID,
		Role:     entity.UserRole(model.Role),
		Status:   entity.MemberStatus(*model.Status),
	}

	if model.InvitedBy != nil {
		invitedByID, _ := r.converter.ToUUID(*model.InvitedBy)
		member.InvitedBy = &invitedByID
	}

	return member
}
