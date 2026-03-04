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
