package application

import (
	"context"
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	persistenceRepo "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/transaction"
	"gorm.io/gorm"
)

// =============================================================================
// UnitOfWork Option 模式
// =============================================================================

// unitOfWorkOptions UnitOfWork 的配置选项
type unitOfWorkOptions struct {
	eventPublisher kernel.EventPublisher
}

// UnitOfWorkOption UnitOfWork 选项接口
type UnitOfWorkOption interface {
	apply(*unitOfWorkOptions)
}

// optionFunc 函数式选项实现
type optionFunc func(*unitOfWorkOptions)

func (f optionFunc) apply(opts *unitOfWorkOptions) {
	f(opts)
}

// WithEventPublisher 配置事件发布器，启用领域事件自动发布功能
func WithEventPublisher(publisher kernel.EventPublisher) UnitOfWorkOption {
	return optionFunc(func(opts *unitOfWorkOptions) {
		opts.eventPublisher = publisher
	})
}

// =============================================================================
// 基础 UnitOfWork
// =============================================================================

// UnitOfWork 工作单元接口，用于管理事务和仓储
type UnitOfWork interface {
	// Transaction 在事务中执行函数
	Transaction(ctx context.Context, fn func(context.Context) error) error
	// UserRepository 返回用户仓储
	UserRepository() repository.UserRepository
	// LoginStatsRepository 返回登录统计仓储
	LoginStatsRepository() repository.LoginStatsRepository
	// ActivityLogRepository 返回活动日志仓储
	ActivityLogRepository() aggregate.ActivityLogRepository
}

// unitOfWork 工作单元实现
type unitOfWork struct {
	db              *gorm.DB
	query           *dao.Query
	userRepo        repository.UserRepository
	loginStatsRepo  repository.LoginStatsRepository
	activityLogRepo aggregate.ActivityLogRepository
}

// NewUnitOfWork 创建工作单元实例（支持选项模式）
//
// 基础用法：
//
//	uow := application.NewUnitOfWork(db, query)
//
// 带领域事件：
//
//	uow := application.NewUnitOfWork(db, query,
//	    application.WithEventPublisher(eventPublisher))
func NewUnitOfWork(db *gorm.DB, query *dao.Query, opts ...UnitOfWorkOption) UnitOfWork {
	// 应用所有选项
	options := &unitOfWorkOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}

	// 创建基础 unitOfWork 实例
	baseUOW := &unitOfWork{
		db:              db,
		query:           query,
		userRepo:        persistenceRepo.NewUserRepository(query),
		loginStatsRepo:  persistenceRepo.NewLoginStatsRepository(query),
		activityLogRepo: persistenceRepo.NewActivityLogRepository(query),
	}

	// 如果配置了事件发布器，返回增强版本
	if options.eventPublisher != nil {
		return &unitOfWorkWithEvents{
			unitOfWork:     baseUOW,
			eventPublisher: options.eventPublisher,
			trackedAggs:    make([]kernel.AggregateRoot, 0),
		}
	}

	return baseUOW
}

// Transaction 在事务中执行函数
func (u *unitOfWork) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务注入上下文
		txCtx := transaction.WithTransaction(ctx, tx)
		return fn(txCtx)
	})
}

// UserRepository 返回用户仓储
func (u *unitOfWork) UserRepository() repository.UserRepository {
	return u.userRepo
}

// LoginStatsRepository 返回登录统计仓储
func (u *unitOfWork) LoginStatsRepository() repository.LoginStatsRepository {
	return u.loginStatsRepo
}

// ActivityLogRepository 返回活动日志仓储
func (u *unitOfWork) ActivityLogRepository() aggregate.ActivityLogRepository {
	return u.activityLogRepo
}

// =============================================================================
// 扩展 UnitOfWork（支持领域事件）
// =============================================================================

// UnitOfWorkWithEvents 扩展的工作单元接口，支持自动事件发布
type UnitOfWorkWithEvents interface {
	UnitOfWork
	// TrackAggregate 跟踪聚合根以自动发布事件
	TrackAggregate(aggregate kernel.AggregateRoot)
	// TransactionWithEvents 在事务中执行函数，并自动发布领域事件
	TransactionWithEvents(ctx context.Context, fn func(context.Context) error) error
}

// unitOfWorkWithEvents 扩展的工作单元实现
type unitOfWorkWithEvents struct {
	*unitOfWork
	eventPublisher kernel.EventPublisher
	trackedAggs    []kernel.AggregateRoot
}

// TrackAggregate 跟踪聚合根以自动发布事件
func (u *unitOfWorkWithEvents) TrackAggregate(aggregate kernel.AggregateRoot) {
	u.trackedAggs = append(u.trackedAggs, aggregate)
}

// TransactionWithEvents 在事务中执行函数，并自动发布领域事件
func (u *unitOfWorkWithEvents) TransactionWithEvents(ctx context.Context, fn func(context.Context) error) error {
	// 清除之前跟踪的聚合根
	u.trackedAggs = make([]kernel.AggregateRoot, 0)

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务注入上下文
		txCtx := transaction.WithTransaction(ctx, tx)

		// 执行业务逻辑
		if err := fn(txCtx); err != nil {
			return err
		}

		// 自动发布所有跟踪的聚合根的事件
		for _, aggregate := range u.trackedAggs {
			events := aggregate.GetUncommittedEvents()
			for _, event := range events {
				if err := u.eventPublisher.Publish(txCtx, event); err != nil {
					return fmt.Errorf("failed to publish event %s: %w", event.EventName(), err)
				}
			}
			// 清除已发布的事件
			aggregate.ClearUncommittedEvents()
		}

		return nil
	})
}

// Reset 重置跟踪的聚合根列表
func (u *unitOfWorkWithEvents) Reset() {
	u.trackedAggs = make([]kernel.AggregateRoot, 0)
}

// TrackedAggregates 返回当前跟踪的聚合根列表（用于测试）
func (u *unitOfWorkWithEvents) TrackedAggregates() []kernel.AggregateRoot {
	return u.trackedAggs
}

// Compile-time interface check
var _ UnitOfWorkWithEvents = (*unitOfWorkWithEvents)(nil)
