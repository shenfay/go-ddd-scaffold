// Package repo 用户模块DAO仓储实现
package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/dao"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/model"
	cast "go-ddd-scaffold/pkg/uitl"
)

// UserDAORepository 用户 DAO 仓储实现
type UserDAORepository struct {
	db      *gorm.DB
	querier *dao.Query
}

// NewUserDAORepository 创建用户 DAO 仓储实例
func NewUserDAORepository(db *gorm.DB) repository.UserRepository {
	return &UserDAORepository{
		db:      db,
		querier: dao.Use(db),
	}
}

// Create 创建用户
func (r *UserDAORepository) Create(ctx context.Context, u *entity.User) error {
	id := u.ID.String()

	userModel := &model.User{
		ID:       &id,
		Email:    u.Email.String(),
		Password: u.Password.String(),
		Nickname: u.Nickname.String(),
		Avatar:   u.Avatar,
		Status:   cast.ToStringPtr(string(u.Status)),
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
		Email:    u.Email.String(),
		Password: u.Password.String(),
		Nickname: u.Nickname.String(),
		Avatar:   u.Avatar,
		Status:   cast.ToStringPtr(string(u.Status)),
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
	id, _ := cast.ToUUID(*userModel.ID)

	status := "active"
	if userModel.Status != nil {
		status = *userModel.Status
	}

	// 创建值对象 - 由于是从数据库读取的已验证数据，我们使用 NewXxxFromString 方法
	email := valueobject.NewEmailFromString(userModel.Email)
	nickname := valueobject.NewNicknameFromString(userModel.Nickname)
	password := entity.HashedPassword(userModel.Password)

	user := &entity.User{
		ID:        id,
		Email:     email,
		Password:  password,
		Nickname:  nickname,
		Avatar:    userModel.Avatar,
		Phone:     userModel.Phone,
		Bio:       userModel.Bio,
		Status:    entity.UserStatus(status),
		CreatedAt: *userModel.CreatedAt,
		UpdatedAt: *userModel.UpdatedAt,
	}

	return user
}

// toEntityWithMembership 将用户模型和成员模型转换为实体（仅返回用户基础信息）
func (r *UserDAORepository) toEntityWithMembership(userModel *model.User, memberModel *model.TenantMember) *entity.User {
	id, _ := cast.ToUUID(*userModel.ID)

	status := "active"
	if userModel.Status != nil {
		status = *userModel.Status
	}

	email := valueobject.NewEmailFromString(userModel.Email)
	nickname := valueobject.NewNicknameFromString(userModel.Nickname)
	password := entity.HashedPassword(userModel.Password)

	user := &entity.User{
		ID:        id,
		Email:     email,
		Password:  password,
		Nickname:  nickname,
		Avatar:    userModel.Avatar,
		Phone:     userModel.Phone,
		Bio:       userModel.Bio,
		Status:    entity.UserStatus(status),
		CreatedAt: *userModel.CreatedAt,
		UpdatedAt: *userModel.UpdatedAt,
	}

	return user
}
