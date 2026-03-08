// Package transaction 事务管理辅助工具
package transaction

import (
	"context"

	"gorm.io/gorm"
)

type txKeyType struct{}

var txKey = txKeyType{}

// ContextWithTx 将事务添加到 context
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// GetTxFromContext 从 context 获取事务
func GetTxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// getTxFromContext 简化版本（内部使用）
func getTxFromContext(ctx context.Context) *gorm.DB {
	return GetTxFromContext(ctx)
}
