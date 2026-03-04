// Package auth JWT Token 黑名单服务
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"go-ddd-scaffold/internal/pkg/metrics"
	"go-ddd-scaffold/internal/pkg/ratelimit"
)

// TokenBlacklistService Token 黑名单服务接口
type TokenBlacklistService interface {
	// AddToBlacklist 将 token 加入黑名单
	AddToBlacklist(ctx context.Context, token string, expireAt time.Time) error
	// IsBlacklisted 检查 token 是否在黑名单中
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	// IsBlacklistedBatch 批量检查多个 token 是否在黑名单中（使用 Pipeline 优化性能）
	IsBlacklistedBatch(ctx context.Context, tokens []string) (map[string]bool, error)
	// RemoveFromBlacklist 从黑名单移除 token（可选操作）
	RemoveFromBlacklist(ctx context.Context, token string) error
}

// redisTokenBlacklistService Redis 实现的 Token 黑名单服务
type redisTokenBlacklistService struct {
	client         *redis.Client
	prefix         string // key 前缀，默认 "token:blacklist:"
	rateLimiter    *ratelimit.RateLimiter     // 限流器
	circuitBreaker *ratelimit.CircuitBreaker  // 熔断器
	metrics        *metrics.Metrics           // 监控指标
}

// NewRedisTokenBlacklistService 创建 Redis Token 黑名单服务（带限流熔断）
func NewRedisTokenBlacklistService(
	client *redis.Client,
	prefix string,
	rateLimiter *ratelimit.RateLimiter,
	circuitBreaker *ratelimit.CircuitBreaker,
	metrics *metrics.Metrics,
) TokenBlacklistService {
	if prefix == "" {
		prefix = "token:blacklist:"
	}
	return &redisTokenBlacklistService{
		client:         client,
		prefix:         prefix,
		rateLimiter:    rateLimiter,
		circuitBreaker: circuitBreaker,
		metrics:        metrics,
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

// IsBlacklisted 检查 token 是否在黑名单中（带监控和限流熔断）
func (s *redisTokenBlacklistService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	startTime := time.Now()
	
	// 1. 限流检查
	if s.rateLimiter != nil && !s.rateLimiter.Allow() {
		s.metrics.RecordTokenBlacklistCheck("single", false, time.Since(startTime))
		return false, ratelimit.ErrRateLimited
	}
	
	var result bool
	var err error
	
	// 2. 熔断器保护
	if s.circuitBreaker != nil {
		err = s.circuitBreaker.Execute(ctx, func() error {
			result, err = s.checkBlacklist(ctx, token)
			return err
		})
		if err != nil {
			s.metrics.RecordTokenBlacklistCheck("single", false, time.Since(startTime))
			return false, err
		}
	} else {
		result, err = s.checkBlacklist(ctx, token)
	}
	
	// 3. 记录监控指标
	s.metrics.RecordTokenBlacklistCheck("single", result, time.Since(startTime))
	
	return result, err
}

// checkBlacklist 实际检查黑名单的内部方法
func (s *redisTokenBlacklistService) checkBlacklist(ctx context.Context, token string) (bool, error) {
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

// IsBlacklistedBatch 批量检查多个 token 是否在黑名单中（使用 Pipeline 优化）
func (s *redisTokenBlacklistService) IsBlacklistedBatch(ctx context.Context, tokens []string) (map[string]bool, error) {
	startTime := time.Now()
	
	// 1. 限流检查
	if s.rateLimiter != nil && !s.rateLimiter.Allow() {
		s.metrics.RecordTokenBlacklistCheck("batch", false, time.Since(startTime))
		return nil, ratelimit.ErrRateLimited
	}
	
	var result map[string]bool
	var err error
	
	// 2. 熔断器保护
	if s.circuitBreaker != nil {
		err = s.circuitBreaker.Execute(ctx, func() error {
			result, err = s.checkBlacklistBatch(ctx, tokens)
			return err
		})
		if err != nil {
			s.metrics.RecordTokenBlacklistCheck("batch", false, time.Since(startTime))
			return nil, err
		}
	} else {
		result, err = s.checkBlacklistBatch(ctx, tokens)
	}
	
	// 3. 记录监控指标
	duration := time.Since(startTime)
	s.metrics.RecordRedisPipeline("blacklist_batch", len(tokens))
	
	// 统计命中情况
	hits := 0
	for _, inBlacklist := range result {
		if inBlacklist {
			hits++
		}
	}
	s.metrics.RecordTokenBlacklistCheck("batch", hits > 0, duration)
	
	return result, err
}

// checkBlacklistBatch 实际批量检查黑名单的内部方法
func (s *redisTokenBlacklistService) checkBlacklistBatch(ctx context.Context, tokens []string) (map[string]bool, error) {
	if len(tokens) == 0 {
		return make(map[string]bool), nil
	}

	// 创建 Pipeline
	pipe := s.client.Pipeline()
	
	// 构建所有 EXISTS 命令
	cmds := make([]*redis.IntCmd, len(tokens))
	for i, token := range tokens {
		key := s.prefix + token
		cmds[i] = pipe.Exists(ctx, key)
	}

	// 执行 Pipeline（一次网络往返）
	_, execErr := pipe.Exec(ctx)
	if execErr != nil && execErr != redis.Nil {
		return nil, fmt.Errorf("批量检查 token 黑名单失败：%w", execErr)
	}

	// 收集结果
	result := make(map[string]bool, len(tokens))
	for i, cmd := range cmds {
		exists, _ := cmd.Result()
		result[tokens[i]] = exists > 0
	}

	return result, nil
}
