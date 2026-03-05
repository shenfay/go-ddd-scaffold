// Package integration 集成测试 - Token 黑名单机制端到端测试
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"go-ddd-scaffold/internal/config"
	"go-ddd-scaffold/internal/infrastructure/auth"
	"go-ddd-scaffold/internal/infrastructure/wire"
)

// TestTokenBlacklist_EndToEnd Token 黑名单端到端测试
func TestTokenBlacklist_EndToEnd(t *testing.T) {
	// 1. 加载配置
	cfg, err := config.LoadConfig("../../config/config.yaml")
	if err != nil {
		t.Skipf("跳过测试：无法加载配置文件 - %v", err)
		return
	}

	// 2. 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	
	// 测试 Redis 连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("跳过测试：无法连接 Redis - %v", err)
		return
	}
	defer rdb.Close()

	// 3. 初始化 JWT 服务
	jwtService := auth.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.ExpireIn)
	
	// 4. 初始化监控指标、限流器和熔断器
	metrics := wire.InitializeMetrics()
	rateLimiter := wire.InitializeRateLimiter(metrics)
	circuitBreaker := wire.InitializeCircuitBreaker(metrics)
	
	// 5. 初始化 Token 黑名单服务
	tokenBlacklist := auth.NewRedisTokenBlacklistService(rdb, "test:blacklist:", rateLimiter, circuitBreaker, metrics)

	// 6. 清理测试数据
	cleanupKeys(ctx, rdb, "test:blacklist:*")
	defer cleanupKeys(ctx, rdb, "test:blacklist:*")

	// ========== 测试场景 1: 生成 Token 并验证 ==========
	t.Run("生成 Token 并验证成功", func(t *testing.T) {
		userID := uuid.New()
		token, err := jwtService.GenerateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// 验证 Token
		claims, err := jwtService.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})

	// ========== 测试场景 2: 将 Token 加入黑名单 ==========
	t.Run("将 Token 加入黑名单", func(t *testing.T) {
		userID := uuid.New()
		token, _ := jwtService.GenerateToken(userID)

		// 加入黑名单
		expireAt := time.Now().Add(24 * time.Hour)
		err := tokenBlacklist.AddToBlacklist(ctx, token, expireAt)
		assert.NoError(t, err)

		// 验证是否在黑名单中
		isBlacklisted, err := tokenBlacklist.IsBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.True(t, isBlacklisted)
	})

	// ========== 测试场景 3: 检查未加入黑名单的 Token ==========
	t.Run("检查未加入黑名单的 Token", func(t *testing.T) {
		userID := uuid.New()
		token, _ := jwtService.GenerateToken(userID)

		// 未加入黑名单，应该返回 false
		isBlacklisted, err := tokenBlacklist.IsBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.False(t, isBlacklisted)
	})

	// ========== 测试场景 4: Token 自动过期 ==========
	t.Run("Token 自动过期", func(t *testing.T) {
		userID := uuid.New()
		token, _ := jwtService.GenerateToken(userID)

		// 加入黑名单，设置很短的过期时间（1 秒）
		expireAt := time.Now().Add(time.Second)
		err := tokenBlacklist.AddToBlacklist(ctx, token, expireAt)
		assert.NoError(t, err)

		// 立即检查应该在黑名单中
		isBlacklisted, _ := tokenBlacklist.IsBlacklisted(ctx, token)
		assert.True(t, isBlacklisted)

		// 等待 2 秒让 Redis 自动删除
		time.Sleep(2 * time.Second)

		// 再次检查应该不在黑名单中（已过期）
		isBlacklisted, _ = tokenBlacklist.IsBlacklisted(ctx, token)
		assert.False(t, isBlacklisted)
	})

	// ========== 测试场景 5: 带租户上下文的 Token ==========
	t.Run("带租户上下文的 Token 黑名单", func(t *testing.T) {
		userID := uuid.New()
		tenantID := uuid.New()
		
		// 生成带租户上下文的 Token
		token, err := jwtService.GenerateTokenWithTenant(userID, tenantID)
		assert.NoError(t, err)

		// 验证 Token 包含租户信息
		claims, err := jwtService.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims.TenantID)
		assert.Equal(t, tenantID, *claims.TenantID)

		// 加入黑名单
		expireAt := time.Now().Add(24 * time.Hour)
		err = tokenBlacklist.AddToBlacklist(ctx, token, expireAt)
		assert.NoError(t, err)

		// 验证在黑名单中
		isBlacklisted, _ := tokenBlacklist.IsBlacklisted(ctx, token)
		assert.True(t, isBlacklisted)
	})

	// ========== 测试场景 6: 从黑名单移除 Token ==========
	t.Run("从黑名单移除 Token", func(t *testing.T) {
		userID := uuid.New()
		token, _ := jwtService.GenerateToken(userID)

		// 加入黑名单
		expireAt := time.Now().Add(24 * time.Hour)
		err := tokenBlacklist.AddToBlacklist(ctx, token, expireAt)
		assert.NoError(t, err)

		// 验证在黑名单中
		isBlacklisted, _ := tokenBlacklist.IsBlacklisted(ctx, token)
		assert.True(t, isBlacklisted)

		// 从黑名单移除
		err = tokenBlacklist.RemoveFromBlacklist(ctx, token)
		assert.NoError(t, err)

		// 验证不在黑名单中
		isBlacklisted, _ = tokenBlacklist.IsBlacklisted(ctx, token)
		assert.False(t, isBlacklisted)
	})
}

// cleanupKeys 清理指定模式的 Redis 键
func cleanupKeys(ctx context.Context, rdb *redis.Client, pattern string) {
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		rdb.Del(ctx, key).Result()
	}
}
