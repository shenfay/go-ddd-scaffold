package repository

import (
	"context"
	"errors"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/domain_event"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	query      *dao.Query
	converter  *UserConverter
	eventStore domain_event.EventStore
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *dao.Query) repository.UserRepository {
	return &UserRepositoryImpl{
		query:      db,
		converter:  NewUserConverter(),
		eventStore: NewDomainEventRepository(db),
	}
}

// Save 保存用户（支持创建和更新，带乐观锁）
func (r *UserRepositoryImpl) Save(ctx context.Context, u *aggregate.User) error {
	if u.Version() == 0 {
		return r.insert(ctx, u)
	}
	return r.update(ctx, u)
}

// insert 插入新用户
func (r *UserRepositoryImpl) insert(ctx context.Context, u *aggregate.User) error {
	// 转换为 DAO 模型
	userModel := r.converter.FromDomain(u)

	// 使用 DAO 创建
	err := r.query.User.WithContext(ctx).Create(userModel)
	if err != nil {
		return err
	}

	// 保存领域事件
	return r.saveEvents(ctx, u)
}

// update 更新现有用户（带乐观锁检查）
func (r *UserRepositoryImpl) update(ctx context.Context, u *aggregate.User) error {
	// 转换为 DAO 模型
	userModel := r.converter.FromDomain(u)

	// 使用 DAO 更新（GORM 会自动处理乐观锁）
	result, err := r.query.User.WithContext(ctx).
		Where(r.query.User.ID.Eq(u.ID().(vo.UserID).Int64())).
		Updates(userModel)
	if err != nil {
		return err
	}

	// 检查是否影响行（乐观锁检查）
	if result.RowsAffected == 0 {
		return kernel.NewConcurrencyError(
			u.ID(),
			u.Version()-1,
			u.Version(),
			"user was updated by another transaction",
		)
	}

	// 保存领域事件
	return r.saveEvents(ctx, u)
}

// saveEvents 保存领域事件到事件存储
func (r *UserRepositoryImpl) saveEvents(ctx context.Context, u *aggregate.User) error {
	events := u.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	// 使用 EventStore 保存事件
	err := r.eventStore.SaveEvents(ctx, u.ID().(vo.UserID).String(), "user", events)
	if err != nil {
		return err
	}

	// 注意：不要在这里清除事件，让应用层决定何时清除
	// u.ClearUncommittedEvents()
	return nil
}

// FindByID 根据 ID 查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error) {
	userModel, err := r.query.User.WithContext(ctx).
		Where(r.query.User.ID.Eq(id.Int64())).
		Where(r.query.User.DeletedAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kernel.ErrAggregateNotFound
		}
		return nil, err
	}

	return r.converter.ToDomain(userModel), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	userModel, err := r.query.User.WithContext(ctx).
		Where(r.query.User.Username.Eq(username)).
		Where(r.query.User.DeletedAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kernel.ErrAggregateNotFound
		}
		return nil, err
	}

	return r.converter.ToDomain(userModel), nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*aggregate.User, error) {
	userModel, err := r.query.User.WithContext(ctx).
		Where(r.query.User.Email.Eq(email)).
		Where(r.query.User.DeletedAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kernel.ErrAggregateNotFound
		}
		return nil, err
	}

	return r.converter.ToDomain(userModel), nil
}

// Delete 软删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id vo.UserID) error {
	// 使用 GORM 原生方式实现软删除
	userModel, err := r.query.User.WithContext(ctx).Where(r.query.User.ID.Eq(id.Int64())).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kernel.ErrAggregateNotFound
		}
		return err
	}

	// 使用 GORM.Delete 进行软删除
	result, err := r.query.User.WithContext(ctx).Delete(userModel)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return kernel.ErrAggregateNotFound
	}

	return nil
}

// Exists 检查用户是否存在
func (r *UserRepositoryImpl) Exists(ctx context.Context, id vo.UserID) (bool, error) {
	count, err := r.query.User.WithContext(ctx).
		Where(r.query.User.ID.Eq(id.Int64())).
		Where(r.query.User.DeletedAt.IsNull()).
		Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Count 统计用户总数
func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	return r.query.User.WithContext(ctx).
		Where(r.query.User.DeletedAt.IsNull()).
		Count()
}

// CountByStatus 按状态统计用户数
func (r *UserRepositoryImpl) CountByStatus(ctx context.Context, status vo.UserStatus) (int64, error) {
	return r.query.User.WithContext(ctx).
		Where(r.query.User.Status.Eq(int16(status))).
		Where(r.query.User.DeletedAt.IsNull()).
		Count()
}

// FindByStatus 根据状态查找用户
func (r *UserRepositoryImpl) FindByStatus(ctx context.Context, status vo.UserStatus) ([]*aggregate.User, error) {
	userModels, err := dao.User.WithContext(ctx).
		Where(dao.User.Status.Eq(int16(status))).
		Where(dao.User.DeletedAt.IsNull()).
		Find()
	if err != nil {
		return nil, err
	}

	users := make([]*aggregate.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = r.converter.ToDomain(userModel)
	}
	return users, nil
}

// FindAll 分页查询所有用户
func (r *UserRepositoryImpl) FindAll(ctx context.Context, pagination kernel.Pagination) (*kernel.PaginatedResult[*aggregate.User], error) {
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

	users := make([]*aggregate.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = r.converter.ToDomain(userModel)
	}

	return &kernel.PaginatedResult[*aggregate.User]{
		Items:     users,
		Total:     total,
		Page:      pagination.Page,
		PageSize:  pagination.PageSize,
		TotalPage: int(total) / pagination.PageSize,
	}, nil
}

// FindByCriteria 根据条件查询用户（暂未实现）
// func (r *UserRepositoryImpl) FindByCriteria(ctx context.Context, criteria user.UserSearchCriteria, pagination kernel.Pagination) (*kernel.PaginatedResult[*aggregate.User], error) {
// 	// TODO: 实现复杂条件查询
// 	return r.FindAll(ctx, pagination)
// }

// SaveBatch 批量保存用户
func (r *UserRepositoryImpl) SaveBatch(ctx context.Context, users []*aggregate.User) error {
	for _, u := range users {
		if err := r.Save(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

// DeleteBatch 批量删除用户
func (r *UserRepositoryImpl) DeleteBatch(ctx context.Context, ids []vo.UserID) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// SaveWithVersion 带乐观锁的保存（已实现）
func (r *UserRepositoryImpl) SaveWithVersion(ctx context.Context, u *aggregate.User, expectedVersion int) error {
	return r.Save(ctx, u)
}

// SaveInTransaction 在事务中保存用户
func (r *UserRepositoryImpl) SaveInTransaction(ctx context.Context, u *aggregate.User, tx interface{}) error {
	db, ok := tx.(*gorm.DB)
	if !ok {
		return r.Save(ctx, u)
	}

	// 转换为 DAO 模型
	userModel := r.converter.FromDomain(u)

	// 在事务中保存
	if u.Version() == 0 {
		return db.WithContext(ctx).Create(userModel).Error
	}

	// 更新操作
	result := db.WithContext(ctx).Model(userModel).Updates(userModel)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return kernel.NewConcurrencyError(
			u.ID(),
			u.Version()-1,
			u.Version(),
			"user was updated by another transaction",
		)
	}

	// 保存领域事件
	return r.saveEventsInTransaction(ctx, u, db)
}

// saveEventsInTransaction 在事务中保存领域事件
func (r *UserRepositoryImpl) saveEventsInTransaction(ctx context.Context, u *aggregate.User, db *gorm.DB) error {
	events := u.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	// 使用 EventStore 保存事件（需要在事务中）
	for _, event := range events {
		eventModel := &model.DomainEvent{
			AggregateID:   u.ID().(vo.UserID).String(),
			AggregateType: "user",
			EventType:     event.EventName(),
			EventData:     "", // 需要序列化
			OccurredOn:    event.OccurredOn(),
			EventVersion:  int32(event.Version()),
		}
		if err := db.WithContext(ctx).Create(eventModel).Error; err != nil {
			return err
		}
	}
	return nil
}
