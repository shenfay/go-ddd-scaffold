package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// UserCacheAdapter UserCache 端口适配器
type UserCacheAdapter struct {
	client *redis.Client
	prefix string
}

// NewUserCacheAdapter 创建用户缓存适配器
func NewUserCacheAdapter(client *redis.Client, prefix string) *UserCacheAdapter {
	return &UserCacheAdapter{
		client: client,
		prefix: prefix,
	}
}

func (a *UserCacheAdapter) key(userID int64) string {
	return a.prefix + ":user:" + string(rune(userID))
}

// Get 获取用户缓存
func (a *UserCacheAdapter) Get(ctx context.Context, userID int64) ([]byte, error) {
	key := a.key(userID)
	data, err := a.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // 缓存不存在
	}
	return data, err
}

// Set 设置用户缓存
func (a *UserCacheAdapter) Set(ctx context.Context, userID int64, data []byte, expireSeconds int64) error {
	key := a.key(userID)
	return a.client.Set(ctx, key, data, time.Duration(expireSeconds)*time.Second).Err()
}

// Delete 删除用户缓存
func (a *UserCacheAdapter) Delete(ctx context.Context, userID int64) error {
	key := a.key(userID)
	return a.client.Del(ctx, key).Err()
}

// Exists 检查缓存是否存在
func (a *UserCacheAdapter) Exists(ctx context.Context, userID int64) (bool, error) {
	key := a.key(userID)
	n, err := a.client.Exists(ctx, key).Result()
	return n > 0, err
}
