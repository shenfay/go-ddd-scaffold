// Package transaction 提供事务管理相关的基础设施
package transaction

import (
	"context"

	"gorm.io/gorm"
)

// Transaction 事务接口
type Transaction interface {
	// Commit 提交事务
	Commit() error
	// Rollback 回滚事务
	Rollback() error
	// GetDB 获取底层的 gorm.DB 实例
	GetDB() *gorm.DB
}

// UnitOfWork 工作单元接口
type UnitOfWork interface {
	// Begin 开启一个新的事务
	Begin(ctx context.Context) (Transaction, error)
}

// gormTransaction GORM 事务实现
type gormTransaction struct {
	tx *gorm.DB
}

// NewTransaction 创建新的事务实例
func NewTransaction(tx *gorm.DB) Transaction {
	return &gormTransaction{tx: tx}
}

// Commit 提交事务
func (t *gormTransaction) Commit() error {
	return t.tx.Commit().Error
}

// Rollback 回滚事务
func (t *gormTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

// GetDB 获取底层的 gorm.DB 实例
func (t *gormTransaction) GetDB() *gorm.DB {
	return t.tx
}
