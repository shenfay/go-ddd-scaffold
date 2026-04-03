package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenStore Redis Token 存储
type TokenStore struct {
	client *redis.Client
}

// TokenData Token 数据结构
type TokenData struct {
	UserID       string    `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	DeviceID     string    `json:"device_id"`
}

// NewTokenStore 创建 Token 存储
func NewTokenStore(client *redis.Client) *TokenStore {
	return &TokenStore{client: client}
}

// Store 存储 Token（7 天有效期）
func (s *TokenStore) Store(ctx context.Context, refreshToken string, data *TokenData) error {
	key := s.buildKey(refreshToken)
	value, _ := json.Marshal(data)

	return s.client.Set(ctx, key, value, 7*24*time.Hour).Err()
}

// Get 获取 Token 信息
func (s *TokenStore) Get(ctx context.Context, refreshToken string) (*TokenData, error) {
	key := s.buildKey(refreshToken)
	value, err := s.client.Get(ctx, key).Bytes()

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
	key := s.buildKey(refreshToken)
	return s.client.Del(ctx, key).Err()
}

// IsBlacklisted 检查 Token 是否在黑名单中
func (s *TokenStore) IsBlacklisted(ctx context.Context, refreshToken string) bool {
	key := "auth:blacklist:" + refreshToken
	exists, _ := s.client.Exists(ctx, key).Result()
	return exists > 0
}

// AddToBlacklist 将 Token 加入黑名单
func (s *TokenStore) AddToBlacklist(ctx context.Context, refreshToken string, expiresAt time.Time) error {
	key := "auth:blacklist:" + refreshToken
	ttl := time.Until(expiresAt)
	return s.client.Set(ctx, key, "1", ttl).Err()
}

func (s *TokenStore) buildKey(refreshToken string) string {
	return "auth:token:" + refreshToken
}

// ErrTokenNotFound Token 未找到错误
var ErrTokenNotFound = redis.Nil
