// Package auth_test Token 黑名单服务性能基准测试
package auth_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"go-ddd-scaffold/internal/config"
	"go-ddd-scaffold/internal/infrastructure/auth"
)

// setupTestRedis 初始化测试用 Redis
func setupTestRedis(b *testing.B) (*redis.Client, func()) {
	// 加载配置
	cfg, err := config.LoadConfig("../../config/config.yaml")
	if err != nil {
		b.Skipf("跳过基准测试：无法加载配置文件 - %v", err)
	}

	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		b.Skipf("跳过基准测试：无法连接 Redis - %v", err)
	}

	// 清理函数
	cleanup := func() {
		iter := rdb.Scan(ctx, 0, "bench:blacklist:*", 0).Iterator()
		for iter.Next(ctx) {
			rdb.Del(ctx, iter.Val()).Result()
		}
		rdb.Close()
	}

	return rdb, cleanup
}

// BenchmarkIsBlacklisted_Single 基准测试：单次检查（无 Pipeline）
func BenchmarkIsBlacklisted_Single(b *testing.B) {
	rdb, cleanup := setupTestRedis(b)
	defer cleanup()

	service := auth.NewRedisTokenBlacklistService(rdb, "bench:blacklist:")
	ctx := context.Background()

	// 准备测试数据：预先加入一些 token 到黑名单
	tokenCount := 100
	tokens := make([]string, tokenCount)
	for i := 0; i < tokenCount; i++ {
		tokens[i] = fmt.Sprintf("token_%d", i)
		expireAt := time.Now().Add(24 * time.Hour)
		service.AddToBlacklist(ctx, tokens[i], expireAt)
	}

	b.ResetTimer()

	// 测试：逐个检查（N 次网络往返）
	for i := 0; i < b.N; i++ {
		for _, token := range tokens[:10] { // 每次检查 10 个
			_, _ = service.IsBlacklisted(ctx, token)
		}
	}
}

// BenchmarkIsBlacklisted_Batch 基准测试：批量检查（使用 Pipeline）
func BenchmarkIsBlacklisted_Batch(b *testing.B) {
	rdb, cleanup := setupTestRedis(b)
	defer cleanup()

	service := auth.NewRedisTokenBlacklistService(rdb, "bench:blacklist:")
	ctx := context.Background()

	// 准备测试数据
	tokenCount := 100
	tokens := make([]string, tokenCount)
	for i := 0; i < tokenCount; i++ {
		tokens[i] = fmt.Sprintf("token_%d", i)
		expireAt := time.Now().Add(24 * time.Hour)
		service.AddToBlacklist(ctx, tokens[i], expireAt)
	}

	b.ResetTimer()

	// 测试：批量检查（1 次网络往返）
	for i := 0; i < b.N; i++ {
		_, _ = service.IsBlacklistedBatch(ctx, tokens[:10]) // 每次检查 10 个
	}
}

// BenchmarkIsBlacklisted_Batch_Large 基准测试：大批量检查
func BenchmarkIsBlacklisted_Batch_Large(b *testing.B) {
	rdb, cleanup := setupTestRedis(b)
	defer cleanup()

	service := auth.NewRedisTokenBlacklistService(rdb, "bench:blacklist:")
	ctx := context.Background()

	// 准备测试数据：1000 个 token
	tokenCount := 1000
	tokens := make([]string, tokenCount)
	for i := 0; i < tokenCount; i++ {
		tokens[i] = fmt.Sprintf("token_%d", i)
		if i%2 == 0 { // 一半在黑名单中
			expireAt := time.Now().Add(24 * time.Hour)
			service.AddToBlacklist(ctx, tokens[i], expireAt)
		}
	}

	b.ResetTimer()

	// 测试：批量检查 100 个 token
	for i := 0; i < b.N; i++ {
		start := (i * 100) % tokenCount
		end := start + 100
		if end > tokenCount {
			end = tokenCount
		}
		_, _ = service.IsBlacklistedBatch(ctx, tokens[start:end])
	}
}

// TestIsBlacklistedBatch_Correctness 功能测试：验证批量检查正确性
func TestIsBlacklistedBatch_Correctness(t *testing.T) {
	// 检查是否运行在 CI 环境
	if os.Getenv("CI") != "" {
		t.Skip("跳过测试：CI 环境未配置 Redis")
	}

	rdb, cleanup := setupTestRedisForTest(t)
	defer cleanup()

	service := auth.NewRedisTokenBlacklistService(rdb, "test:blacklist:")
	ctx := context.Background()

	// 准备测试数据
	tokens := []string{
		"token_in_blacklist_1",
		"token_in_blacklist_2",
		"token_not_in_blacklist_1",
		"token_not_in_blacklist_2",
	}

	// 将部分 token 加入黑名单
	expireAt := time.Now().Add(24 * time.Hour)
	_ = service.AddToBlacklist(ctx, tokens[0], expireAt)
	_ = service.AddToBlacklist(ctx, tokens[1], expireAt)

	// 执行批量检查
	results, err := service.IsBlacklistedBatch(ctx, tokens)
	assert.NoError(t, err)
	assert.Len(t, results, len(tokens))

	// 验证结果
	assert.True(t, results[tokens[0]], "token_0 应该在黑名单中")
	assert.True(t, results[tokens[1]], "token_1 应该在黑名单中")
	assert.False(t, results[tokens[2]], "token_2 不应该在黑名单中")
	assert.False(t, results[tokens[3]], "token_3 不应该在黑名单中")
}

// TestIsBlacklistedBatch_Empty 功能测试：空列表
func TestIsBlacklistedBatch_Empty(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("跳过测试：CI 环境未配置 Redis")
	}

	rdb, cleanup := setupTestRedisForTest(t)
	defer cleanup()

	service := auth.NewRedisTokenBlacklistService(rdb, "test:blacklist:")
	ctx := context.Background()

	// 测试空列表
	results, err := service.IsBlacklistedBatch(ctx, []string{})
	assert.NoError(t, err)
	assert.Empty(t, results)
}

// setupTestRedisForTest 初始化为测试用 Redis
func setupTestRedisForTest(t *testing.T) (*redis.Client, func()) {
	cfg, err := config.LoadConfig("../../config/config.yaml")
	if err != nil {
		t.Skipf("跳过测试：无法加载配置文件 - %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("跳过测试：无法连接 Redis - %v", err)
	}

	cleanup := func() {
		iter := rdb.Scan(ctx, 0, "test:blacklist:*", 0).Iterator()
		for iter.Next(ctx) {
			rdb.Del(ctx, iter.Val()).Result()
		}
	}

	return rdb, cleanup
}
