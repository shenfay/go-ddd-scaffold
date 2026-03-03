// Package repo 用户模块DAO仓储实现
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

// UserDAORepository 用户DAO仓储实现
type UserDAORepository struct {
	db        *gorm.DB
	querier   *dao.Query
	converter converter.Converter
}

// NewUserDAORepository 创建用户DAO仓储实例
func NewUserDAORepository(db *gorm.DB) repository.UserRepository {
	return &UserDAORepository{
		db:        db,
		querier:   dao.Use(db),
		converter: converter.NewConverter(),
	}
}

// Create 创建用户
func (r *UserDAORepository) Create(ctx context.Context, u *entity.User) error {
	id := u.ID.String()

	userModel := &model.User{
		ID:       &id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Avatar:   r.converter.ToStringPtr(u.Avatar),
		Status:   r.converter.ToStringPtr(string(u.Status)),
	}

	return r.querier.User.WithContext(ctx).Create(userModel)
}

// GetByID 根据ID获取用户
func (r *UserDAORepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	userModel, err := r.querier.User.WithContext(ctx).Where(r.querier.User.ID.Eq(id.String())).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return r.toEntity(userModel), nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserDAORepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	userModel, err := r.querier.User.WithContext(ctx).Where(r.querier.User.Email.Eq(email)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.toEntity(userModel), nil
}

// Update 更新用户
func (r *UserDAORepository) Update(ctx context.Context, u *entity.User) error {
	id := u.ID.String()

	userModel := &model.User{
		ID:       &id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Avatar:   r.converter.ToStringPtr(u.Avatar),
		Status:   r.converter.ToStringPtr(string(u.Status)),
	}

	_, err := r.querier.User.WithContext(ctx).Where(r.querier.User.ID.Eq(id)).Updates(userModel)
	return err
}

// Delete 删除用户
func (r *UserDAORepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.querier.User.WithContext(ctx).Where(r.querier.User.ID.Eq(id.String())).Delete()
	return err
}

// ListByTenant 列出租户下所有用户
func (r *UserDAORepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error) {
	// 通过tenant_members表获取租户下的用户
	memberModels, err := r.querier.TenantMember.WithContext(ctx).Where(
		r.querier.TenantMember.TenantID.Eq(tenantID.String()),
		r.querier.TenantMember.Status.Neq("removed"),
	).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to list members by tenant: %w", err)
	}

	users := make([]*entity.User, len(memberModels))
	for i, memberModel := range memberModels {
		userModel, err := r.querier.User.WithContext(ctx).Where(
			r.querier.User.ID.Eq(memberModel.UserID),
		).First()
		if err != nil {
			return nil, fmt.Errorf("failed to get user by id: %w", err)
		}
		users[i] = r.toEntityWithMembership(userModel, memberModel)
	}

	return users, nil
}

// ListChildrenByParent 列出父账户下的所有子女
func (r *UserDAORepository) ListChildrenByParent(ctx context.Context, parentID uuid.UUID) ([]*entity.User, error) {
	// 首先找到父用户的所有租户成员记录
	parentMemberModels, err := r.querier.TenantMember.WithContext(ctx).Where(
		r.querier.TenantMember.UserID.Eq(parentID.String()),
	).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get parent members: %w", err)
	}

	var users []*entity.User
	for _, parentMember := range parentMemberModels {
		// 在同一租户下寻找children
		childMemberModels, err := r.querier.TenantMember.WithContext(ctx).Where(
			r.querier.TenantMember.TenantID.Eq(parentMember.TenantID),
			r.querier.TenantMember.Role.Eq("child"),
			r.querier.TenantMember.Status.Neq("removed"),
		).Find()
		if err != nil {
			continue // 继续处理下一个租户
		}

		for _, childMemberModel := range childMemberModels {
			userModel, err := r.querier.User.WithContext(ctx).Where(
				r.querier.User.ID.Eq(childMemberModel.UserID),
			).First()
			if err != nil {
				continue // 继续处理下一个child
			}
			users = append(users, r.toEntityWithMembership(userModel, childMemberModel))
		}
	}

	return users, nil
}

// CountByTenant 统计租户用户数
func (r *UserDAORepository) CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	count, err := r.querier.TenantMember.WithContext(ctx).Where(
		r.querier.TenantMember.TenantID.Eq(tenantID.String()),
		r.querier.TenantMember.Status.Neq("removed"),
	).Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count members by tenant: %w", err)
	}
	return count, nil
}

// toEntity 将模型转换为实体
func (r *UserDAORepository) toEntity(userModel *model.User) *entity.User {
	id, _ := r.converter.ToUUID(*userModel.ID)

	avatar := ""
	if userModel.Avatar != nil {
		avatar = *userModel.Avatar
	}

	status := "active"
	if userModel.Status != nil {
		status = *userModel.Status
	}

	user := &entity.User{
		ID:       id,
		Email:    userModel.Email,
		Password: userModel.Password,
		Nickname: userModel.Nickname,
		Avatar:   avatar,
		Status:   entity.UserStatus(status),
	}

	return user
}

// toEntityWithMembership 将用户模型和成员模型转换为实体
func (r *UserDAORepository) toEntityWithMembership(userModel *model.User, memberModel *model.TenantMember) *entity.User {
	id, _ := r.converter.ToUUID(*userModel.ID)
	tenantID, _ := r.converter.ToUUID(memberModel.TenantID)

	avatar := ""
	if userModel.Avatar != nil {
		avatar = *userModel.Avatar
	}

	status := "active"
	if userModel.Status != nil {
		status = *userModel.Status
	}

	user := &entity.User{
		ID:       id,
		Email:    userModel.Email,
		Password: userModel.Password,
		Nickname: userModel.Nickname,
		Avatar:   avatar,
		Role:     entity.UserRole(memberModel.Role),
		Status:   entity.UserStatus(status),
		TenantID: &tenantID,
	}

	if memberModel.InvitedBy != nil {
		parentID, _ := r.converter.ToUUID(*memberModel.InvitedBy)
		user.ParentID = &parentID
	}

	return user
}
