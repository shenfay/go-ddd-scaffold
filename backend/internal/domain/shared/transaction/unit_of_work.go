// Package transaction 事务管理接口定义
package transaction

import "context"

// Transaction 事务接口
type Transaction interface {
	// Commit 提交事务
	Commit() error
	// Rollback 回滚事务
	Rollback() error
}

// UnitOfWork 工作单元接口
// 用于保证跨仓储操作的原子性和一致性
type UnitOfWork interface {
	// Begin 开始新事务
	Begin(ctx context.Context) (Transaction, error)
	
	// WithTransaction 在事务中执行操作
	// 自动处理 commit/rollback，简化使用
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
