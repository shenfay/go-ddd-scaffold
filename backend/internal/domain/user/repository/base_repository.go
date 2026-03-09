package repository

import (
	"context"

	"gorm.io/gorm"
)

// BaseRepository 通用仓储接口（泛型）
type BaseRepository[T any, ID any] interface {
	// GetByID 根据 ID 获取实体
	GetByID(ctx context.Context, id ID) (*T, error)
	// Create 创建实体
	Create(ctx context.Context, entity *T) error
	// Update 更新实体
	Update(ctx context.Context, entity *T) error
	// Delete 删除实体
	Delete(ctx context.Context, id ID) error
	// WithTx 切换到事务上下文
	WithTx(tx *gorm.DB) BaseRepository[T, ID]
}
