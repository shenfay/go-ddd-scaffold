package cache

import "context"

// UserCache 用户缓存端口
type UserCache interface {
	// Get 获取用户缓存
	Get(ctx context.Context, userID int64) ([]byte, error)

	// Set 设置用户缓存
	Set(ctx context.Context, userID int64, data []byte, expireSeconds int64) error

	// Delete 删除用户缓存
	Delete(ctx context.Context, userID int64) error

	// Exists 检查缓存是否存在
	Exists(ctx context.Context, userID int64) (bool, error)
}
