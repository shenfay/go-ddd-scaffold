// Package transaction 提供事务管理相关的基础设施
package transaction

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// gormUnitOfWork GORM 工作单元实现
type gormUnitOfWork struct {
	db *gorm.DB
}

// NewGormUnitOfWork 创建 GORM 工作单元实例
func NewGormUnitOfWork(db *gorm.DB) UnitOfWork {
	return &gormUnitOfWork{db: db}
}

// Begin 开启一个新的事务
func (uow *gormUnitOfWork) Begin(ctx context.Context) (Transaction, error) {
	tx := uow.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	return NewTransaction(tx), nil
}

// WithTransaction 在事务中执行操作
func (uow *gormUnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := uow.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// panic 时回滚
			_ = tx.Rollback()
			panic(p) // 重新抛出 panic
		}
	}()

	// 将事务添加到 context
	ctxWithTx := ContextWithTx(ctx, tx.GetDB())

	// 执行操作
	err = fn(ctxWithTx)
	if err != nil {
		// 有错误时回滚
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %w, original error: %v", rbErr, err)
		}
		return err
	}

	// 提交事务
	if cmErr := tx.Commit(); cmErr != nil {
		return fmt.Errorf("commit error: %w", cmErr)
	}

	return nil
}
