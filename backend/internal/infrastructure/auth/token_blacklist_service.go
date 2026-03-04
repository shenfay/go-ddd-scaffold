// Package auth JWT Token 黑名单服务
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklistService Token 黑名单服务接口
type TokenBlacklistService interface {
	// AddToBlacklist 将 token 加入黑名单
	AddToBlacklist(ctx context.Context, token string, expireAt time.Time) error
	// IsBlacklisted 检查 token 是否在黑名单中
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	// RemoveFromBlacklist 从黑名单移除 token（可选操作）
	RemoveFromBlacklist(ctx context.Context, token string) error
}

// redisTokenBlacklistService Redis 实现的 Token 黑名单服务
type redisTokenBlacklistService struct {
	client *redis.Client
	prefix string // key 前缀，默认 "token:blacklist:"
}

// NewRedisTokenBlacklistService 创建 Redis Token 黑名单服务
func NewRedisTokenBlacklistService(client *redis.Client, prefix string) TokenBlacklistService {
	if prefix == "" {
		prefix = "token:blacklist:"
	}
	return &redisTokenBlacklistService{
		client: client,
		prefix: prefix,
	}
}

// AddToBlacklist 将 token 加入黑名单（使用 SETEX 设置过期时间）
func (s *redisTokenBlacklistService) AddToBlacklist(ctx context.Context, token string, expireAt time.Time) error {
	key := s.prefix + token
	
	// 计算剩余有效期
	ttl := time.Until(expireAt)
	if ttl <= 0 {
		return fmt.Errorf("token 已过期")
	}

	// 使用 SETEX 设置带过期时间的值
	err := s.client.SetEx(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("添加 token 到黑名单失败：%w", err)
	}

	return nil
}

// IsBlacklisted 检查 token 是否在黑名单中
func (s *redisTokenBlacklistService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := s.prefix + token
	
	// 检查 key 是否存在
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查 token 黑名单失败：%w", err)
	}

	return exists > 0, nil
}

// RemoveFromBlacklist 从黑名单移除 token
func (s *redisTokenBlacklistService) RemoveFromBlacklist(ctx context.Context, token string) error {
	key := s.prefix + token
	
	_, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("从黑名单移除 token 失败：%w", err)
	}

	return nil
}
