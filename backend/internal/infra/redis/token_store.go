package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/shenfay/go-ddd-scaffold/pkg/metrics"
)

// TokenStore Redis Token 存储
type TokenStore struct {
	client  *redis.Client
	metrics *metrics.Metrics
}

// TokenData Token 数据结构
type TokenData struct {
	UserID       string    `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	DeviceID     string    `json:"device_id"`
}

// NewTokenStore 创建 Token 存储实例
func NewTokenStore(client *redis.Client, m *metrics.Metrics) *TokenStore {
	return &TokenStore{client: client, metrics: m}
}

// observeRedisCommand 记录 Redis 命令指标
func (s *TokenStore) observeRedisCommand(command string, start time.Time) {
	if s.metrics != nil {
		duration := time.Since(start).Seconds()
		s.metrics.IncRedisCommand(command)
		s.metrics.ObserveRedisCommandDuration(command, duration)
	}
}

// Store 存储 Token（7 天有效期）
func (s *TokenStore) Store(ctx context.Context, refreshToken string, data *TokenData) error {
	start := time.Now()
	key := s.buildKey(refreshToken)
	value, _ := json.Marshal(data)

	err := s.client.Set(ctx, key, value, 7*24*time.Hour).Err()
	s.observeRedisCommand("SET", start)
	return err
}

// Get 获取 Token 信息
func (s *TokenStore) Get(ctx context.Context, refreshToken string) (*TokenData, error) {
	start := time.Now()
	key := s.buildKey(refreshToken)
	value, err := s.client.Get(ctx, key).Bytes()
	s.observeRedisCommand("GET", start)

	if err == redis.Nil {
		return nil, ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	var data TokenData
	if err := json.Unmarshal(value, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Delete 删除 Token（登出时使用）
func (s *TokenStore) Delete(ctx context.Context, refreshToken string) error {
	start := time.Now()
	key := s.buildKey(refreshToken)
	err := s.client.Del(ctx, key).Err()
	s.observeRedisCommand("DEL", start)
	return err
}

// IsBlacklisted 检查 Token 是否在黑名单中
func (s *TokenStore) IsBlacklisted(ctx context.Context, refreshToken string) bool {
	start := time.Now()
	key := "auth:blacklist:" + refreshToken
	exists, _ := s.client.Exists(ctx, key).Result()
	s.observeRedisCommand("EXISTS", start)
	return exists > 0
}

// AddToBlacklist 将 Token 加入黑名单
func (s *TokenStore) AddToBlacklist(ctx context.Context, refreshToken string, expiresAt time.Time) error {
	start := time.Now()
	key := "auth:blacklist:" + refreshToken
	ttl := time.Until(expiresAt)
	err := s.client.Set(ctx, key, "1", ttl).Err()
	s.observeRedisCommand("SET", start)
	return err
}

func (s *TokenStore) buildKey(refreshToken string) string {
	return "auth:token:" + refreshToken
}

// ErrTokenNotFound Token 未找到错误
var ErrTokenNotFound = redis.Nil
