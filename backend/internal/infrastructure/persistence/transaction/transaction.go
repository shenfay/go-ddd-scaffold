package transaction

import (
	"context"

	"gorm.io/gorm"
)

// Manager 事务管理器
type Manager struct {
	db *gorm.DB
}

// NewManager 创建事务管理器
func NewManager(db *gorm.DB) *Manager {
	return &Manager{db: db}
}

// RunInTransaction 在事务中执行函数
// 如果函数返回错误，则回滚事务；否则提交事务
func (m *Manager) RunInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(fn)
}

// WithTransaction 返回带事务的上下文
// 用于在Repository中识别当前是否在事务中
func WithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// GetTransaction 从上下文中获取事务
// 如果不在事务中，返回nil
func GetTransaction(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txContextKey{}).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

type txContextKey struct{}
