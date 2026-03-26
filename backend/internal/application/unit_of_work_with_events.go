package application

import (
	"context"
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/transaction"
	"gorm.io/gorm"
)

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

// NewUnitOfWorkWithEvents 创建支持自动事件发布的工作单元
func NewUnitOfWorkWithEvents(
	db *gorm.DB,
	query *dao.Query,
	eventPublisher kernel.EventPublisher,
) UnitOfWorkWithEvents {
	baseUOW := NewUnitOfWork(db, query).(*unitOfWork)
	return &unitOfWorkWithEvents{
		unitOfWork:     baseUOW,
		eventPublisher: eventPublisher,
		trackedAggs:    make([]kernel.AggregateRoot, 0),
	}
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
