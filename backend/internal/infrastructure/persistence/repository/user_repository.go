package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
	"gorm.io/gorm"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Save 保存用户（支持创建和更新，带乐观锁）
func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
	if u.Version() == 0 {
		return r.insert(ctx, u)
	}
	return r.update(ctx, u)
}

// insert 插入新用户
func (r *UserRepositoryImpl) insert(ctx context.Context, u *user.User) error {
	// 转换为 DAO 模型
	userModel := r.fromDomain(u)

	// 使用 DAO 创建
	err := dao.User.WithContext(ctx).Create(userModel)
	if err != nil {
		return err
	}

	// 保存领域事件
	return r.saveEvents(ctx, u)
}

// update 更新现有用户（带乐观锁检查）
func (r *UserRepositoryImpl) update(ctx context.Context, u *user.User) error {
	// 转换为 DAO 模型
	userModel := r.fromDomain(u)

	// 使用 DAO 更新（GORM 会自动处理乐观锁）
	result, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(u.ID().(user.UserID).Int64())).
		Updates(userModel)
	if err != nil {
		return err
	}

	// 检查是否影响行（乐观锁检查）
	if result.RowsAffected == 0 {
		return ddd.NewConcurrencyError(
			u.ID(),
			u.Version()-1,
			u.Version(),
			"user was updated by another transaction",
		)
	}

	// 保存领域事件
	return r.saveEvents(ctx, u)
}

// saveEvents 保存领域事件到事件存储（保持原生 SQL）
func (r *UserRepositoryImpl) saveEvents(ctx context.Context, u *user.User) error {
	events := u.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return err
		}

		domainEvent := &model.DomainEvent{
			AggregateID:   u.ID().(user.UserID).String(),
			AggregateType: "user",
			EventType:     event.EventName(),
			EventVersion:  int32(event.Version()),
			EventData:     string(eventData),
			OccurredOn:    event.OccurredOn(),
			Processed:     util.Bool(false),
		}

		err = dao.DomainEvent.WithContext(ctx).Create(domainEvent)
		if err != nil {
			return err
		}
	}

	// 清除已保存的事件
	u.ClearUncommittedEvents()
	return nil
}

// FindByID 根据 ID 查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	userModel, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(id.Int64())).
		Where(dao.User.DeletedAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ddd.ErrAggregateNotFound
		}
		return nil, err
	}

	return r.toDomain(userModel), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	userModel, err := dao.User.WithContext(ctx).
		Where(dao.User.Username.Eq(username)).
		Where(dao.User.DeletedAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ddd.ErrAggregateNotFound
		}
		return nil, err
	}

	return r.toDomain(userModel), nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	userModel, err := dao.User.WithContext(ctx).
		Where(dao.User.Email.Eq(email)).
		Where(dao.User.DeletedAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ddd.ErrAggregateNotFound
		}
		return nil, err
	}

	return r.toDomain(userModel), nil
}

// Delete 软删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id user.UserID) error {
	// 使用 GORM 原生方式实现软删除
	userModel, err := dao.User.WithContext(ctx).Where(dao.User.ID.Eq(id.Int64())).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ddd.ErrAggregateNotFound
		}
		return err
	}

	// 使用 GORM.Delete 进行软删除
	result, err := dao.User.WithContext(ctx).Delete(userModel)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ddd.ErrAggregateNotFound
	}

	return nil
}

// Exists 检查用户是否存在
func (r *UserRepositoryImpl) Exists(ctx context.Context, id user.UserID) (bool, error) {
	count, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(id.Int64())).
		Where(dao.User.DeletedAt.IsNull()).
		Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Count 统计用户总数
func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	return dao.User.WithContext(ctx).
		Where(dao.User.DeletedAt.IsNull()).
		Count()
}

// CountByStatus 按状态统计用户数
func (r *UserRepositoryImpl) CountByStatus(ctx context.Context, status user.UserStatus) (int64, error) {
	return dao.User.WithContext(ctx).
		Where(dao.User.Status.Eq(int16(status))).
		Where(dao.User.DeletedAt.IsNull()).
		Count()
}

// FindByStatus 根据状态查找用户
func (r *UserRepositoryImpl) FindByStatus(ctx context.Context, status user.UserStatus) ([]*user.User, error) {
	userModels, err := dao.User.WithContext(ctx).
		Where(dao.User.Status.Eq(int16(status))).
		Where(dao.User.DeletedAt.IsNull()).
		Find()
	if err != nil {
		return nil, err
	}

	users := make([]*user.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = r.toDomain(userModel)
	}
	return users, nil
}

// FindAll 分页查询所有用户
func (r *UserRepositoryImpl) FindAll(ctx context.Context, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.User], error) {
	// 获取总数
	total, err := dao.User.WithContext(ctx).
		Where(dao.User.DeletedAt.IsNull()).
		Count()
	if err != nil {
		return nil, err
	}

	// 获取数据
	offset := pagination.Offset()
	userModels, err := dao.User.WithContext(ctx).
		Where(dao.User.DeletedAt.IsNull()).
		Order(dao.User.CreatedAt.Desc()).
		Limit(pagination.PageSize).
		Offset(offset).
		Find()
	if err != nil {
		return nil, err
	}

	users := make([]*user.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = r.toDomain(userModel)
	}

	return &ddd.PaginatedResult[*user.User]{
		Items:      users,
		TotalCount: total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.CalculateTotalPages(total),
	}, nil
}

// FindByCriteria 根据条件查询用户（暂未实现）
// func (r *UserRepositoryImpl) FindByCriteria(ctx context.Context, criteria user.UserSearchCriteria, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.User], error) {
// 	// TODO: 实现复杂条件查询
// 	return r.FindAll(ctx, pagination)
// }

// SaveBatch 批量保存用户
func (r *UserRepositoryImpl) SaveBatch(ctx context.Context, users []*user.User) error {
	for _, u := range users {
		if err := r.Save(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

// DeleteBatch 批量删除用户
func (r *UserRepositoryImpl) DeleteBatch(ctx context.Context, ids []user.UserID) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// SaveWithVersion 带乐观锁的保存（已实现）
func (r *UserRepositoryImpl) SaveWithVersion(ctx context.Context, u *user.User, expectedVersion int) error {
	return r.Save(ctx, u)
}

// toDomain 将数据库模型转换为领域对象
func (r *UserRepositoryImpl) toDomain(model *model.User) *user.User {
	// 处理可能的 nil 值
	loginCount := 0
	if model.LoginCount != nil {
		loginCount = int(*model.LoginCount)
	}

	failedAttempts := 0
	if model.FailedAttempts != nil {
		failedAttempts = int(*model.FailedAttempts)
	}

	version := 0
	if model.Version != nil {
		version = int(*model.Version)
	}

	gender := user.UserGender(0)
	if model.Gender != nil {
		gender = user.UserGender(*model.Gender)
	}

	var displayName string
	if model.DisplayName != nil {
		displayName = *model.DisplayName
	}

	var phoneNumber string
	if model.PhoneNumber != nil {
		phoneNumber = *model.PhoneNumber
	}

	var avatarURL string
	if model.AvatarURL != nil {
		avatarURL = *model.AvatarURL
	}

	createdAt := time.Time{}
	if model.CreatedAt != nil {
		createdAt = *model.CreatedAt
	}

	updatedAt := time.Time{}
	if model.UpdatedAt != nil {
		updatedAt = *model.UpdatedAt
	}

	// 使用 Builder 模式优雅地重建用户对象
	return user.NewUserBuilder().
		WithID(model.ID).
		WithUsername(model.Username).
		WithEmail(model.Email).
		WithPasswordHash(model.PasswordHash).
		WithStatus(user.UserStatus(model.Status)).
		WithGender(gender).
		WithDisplayName(displayName).
		WithPhoneNumber(phoneNumber).
		WithAvatarURL(avatarURL).
		WithLastLoginAt(model.LastLoginAt).
		WithLoginCount(loginCount).
		WithFailedAttempts(failedAttempts).
		WithLockedUntil(model.LockedUntil).
		WithVersion(version).
		WithTimestamps(createdAt, updatedAt).
		Build()
}

// fromDomain 将领域对象转换为数据库模型
func (r *UserRepositoryImpl) fromDomain(u *user.User) *model.User {
	displayName := u.DisplayName()
	phoneNumber := u.PhoneNumber()
	avatarURL := u.AvatarURL()
	loginCount := int(u.LoginCount())
	failedAttempts := int(u.FailedAttempts())
	version := int(u.Version())

	return &model.User{
		ID:             u.ID().(user.UserID).Int64(),
		Username:       u.Username().Value(),
		Email:          u.Email().Value(),
		PasswordHash:   u.Password().Value(),
		Status:         int16(u.Status()),
		DisplayName:    util.StringPtrNilIfEmpty(displayName),
		Gender:         util.Int16PtrNilIfZero(int16(u.Gender())),
		PhoneNumber:    util.StringPtrNilIfEmpty(phoneNumber),
		AvatarURL:      util.StringPtrNilIfEmpty(avatarURL),
		LastLoginAt:    u.LastLoginAt(),
		LoginCount:     util.Int32PtrNilIfZero(int32(loginCount)),
		FailedAttempts: util.Int32PtrNilIfZero(int32(failedAttempts)),
		LockedUntil:    u.LockedUntil(),
		Version:        util.Int32PtrNilIfZero(int32(version)),
		CreatedAt:      util.Time(u.CreatedAt()),
		UpdatedAt:      util.Time(u.UpdatedAt()),
	}
}
