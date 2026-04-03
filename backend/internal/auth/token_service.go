package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"github.com/shenfay/go-ddd-scaffold/pkg/constants"
	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils/ulid"
)

// DeviceInfo 设备信息
type DeviceInfo struct {
	UserID     string `json:"user_id"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	DeviceType string `json:"device_type"` // desktop, mobile, tablet
	CreatedAt  string `json:"created_at"`
}

// TokenPair Token 对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// JWTClaims JWT 自定义声明
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// TokenService Token 服务
type TokenService struct {
	redisClient   *redis.Client
	jwtSecret     []byte
	issuer        string
	accessExpire  time.Duration
	refreshExpire time.Duration
}

// NewTokenService 创建 Token 服务
func NewTokenService(redisClient *redis.Client, jwtSecret string, issuer string, accessExpire, refreshExpire time.Duration) *TokenService {
	return &TokenService{
		redisClient:   redisClient,
		jwtSecret:     []byte(jwtSecret),
		issuer:        issuer,
		accessExpire:  accessExpire,
		refreshExpire: refreshExpire,
	}
}

// GenerateTokens 生成 Token 对
func (s *TokenService) GenerateTokens(ctx context.Context, userID, email string) (*TokenPair, error) {
	now := time.Now()

	// 生成 Access Token
	accessToken, err := s.generateAccessToken(userID, email, now)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token
	refreshTokenID := ulid.GenerateTokenID()
	refreshToken, err := s.generateRefreshToken(refreshTokenID, now)
	if err != nil {
		return nil, err
	}

	// 存储 Refresh Token 到 Redis
	if err := s.storeRefreshToken(ctx, refreshTokenID, userID); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.accessExpire,
	}, nil
}

// generateAccessToken 生成 Access Token
func (s *TokenService) generateAccessToken(userID, email string, issuedAt time.Time) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(s.accessExpire)),
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// generateRefreshToken 生成 Refresh Token（简单的 UUID 格式）
func (s *TokenService) generateRefreshToken(tokenID string, issuedAt time.Time) (string, error) {
	// Refresh Token 就是一个普通的字符串，存储在 Redis 中
	return tokenID, nil
}

// storeRefreshToken 存储 Refresh Token 到 Redis
func (s *TokenService) storeRefreshToken(ctx context.Context, tokenID, userID string) error {
	key := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, tokenID)
	return s.redisClient.Set(ctx, key, userID, s.refreshExpire).Err()
}

// ValidateRefreshToken 验证 Refresh Token
func (s *TokenService) ValidateRefreshToken(ctx context.Context, token string) (*JWTClaims, error) {
	// 1. 从 Redis 查找 Refresh Token
	key := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, token)
	userID, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrInvalidToken
	}

	// 2. 返回用户信息（用于生成新的 Token 对）
	return &JWTClaims{
		UserID: userID,
	}, nil
}

// ValidateAccessToken 验证 Access Token
func (s *TokenService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	// 检查 Token 类型
	if claims.TokenType != "access" {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}

// RevokeToken 撤销 Token（退出登录时调用）
func (s *TokenService) RevokeToken(ctx context.Context, tokenID string) error {
	// 删除 Refresh Token
	key := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, tokenID)
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		return err
	}

	// TODO: 如果需要支持 Access Token 黑名单机制
	// 可以在这里将未过期的 Access Token 加入黑名单
	// 但由于 Access Token 是短生命周期的，通常不需要立即撤销
	return nil
}

// RevokeAllUserTokens 撤销用户的所有 Token（修改密码等场景）
func (s *TokenService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	// TODO: 如果需要支持多设备登录，需要维护一个用户到 Token 的映射
	// 这里简化实现，暂时不处理
	return nil
}

// ValidateRefreshTokenWithDevice 验证 Refresh Token 并返回设备信息（支持多设备）
func (s *TokenService) ValidateRefreshTokenWithDevice(ctx context.Context, token string) (*DeviceInfo, error) {
	// 1. 从 Redis 查找 Refresh Token（旧格式兼容）
	key := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, token)
	userID, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrInvalidToken
	}

	// 2. 尝试获取设备信息（新格式）
	deviceKey := fmt.Sprintf("%sdevice:%s", constants.RedisKeyPrefix, token)
	deviceData, err := s.redisClient.Get(ctx, deviceKey).Result()
	if err == nil {
		// 有设备信息，解析返回
		var deviceInfo DeviceInfo
		if err := json.Unmarshal([]byte(deviceData), &deviceInfo); err == nil {
			return &deviceInfo, nil
		}
	}

	// 3. 没有设备信息，返回基本用户信息（兼容旧数据）
	return &DeviceInfo{
		UserID: userID,
	}, nil
}

// StoreDeviceInfo 存储设备信息到 Redis（支持多设备）
func (s *TokenService) StoreDeviceInfo(ctx context.Context, token string, deviceInfo DeviceInfo) error {
	now := time.Now()
	deviceInfo.CreatedAt = now.Format(time.RFC3339)

	// 1. 序列化设备信息
	data, err := json.Marshal(deviceInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal device info: %w", err)
	}

	// 2. 存储设备信息到 Redis
	deviceKey := fmt.Sprintf("%sdevice:%s", constants.RedisKeyPrefix, token)
	if err := s.redisClient.Set(ctx, deviceKey, string(data), s.refreshExpire).Err(); err != nil {
		return err
	}

	// 3. 更新用户设备列表（Set 类型，自动去重）
	userDevicesKey := fmt.Sprintf("%suser_devices:%s", constants.RedisKeyPrefix, deviceInfo.UserID)
	if err := s.redisClient.SAdd(ctx, userDevicesKey, token).Err(); err != nil {
		return err
	}

	// 4. 设置用户设备列表的过期时间（与 Refresh Token 相同）
	s.redisClient.Expire(ctx, userDevicesKey, s.refreshExpire)

	return nil
}

// RevokeDevice 撤销特定设备（支持多设备）
func (s *TokenService) RevokeDevice(ctx context.Context, token string) error {
	// 1. 删除设备信息
	deviceKey := fmt.Sprintf("%sdevice:%s", constants.RedisKeyPrefix, token)
	s.redisClient.Del(ctx, deviceKey)

	// 2. 删除 Refresh Token（旧格式兼容）
	tokenKey := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, token)
	s.redisClient.Del(ctx, tokenKey)

	// 3. 从用户设备列表中移除
	// 注意：这里需要先从 device:{token} 获取 userID，再从 user_devices:{userID} 移除
	// 但由于我们已经删除了 device:{token}，所以需要调用方传入 userID
	// 或者在调用 RevokeDevice 之前先获取 deviceInfo
	return nil
}

// GetUserDevices 获取用户的所有设备列表
func (s *TokenService) GetUserDevices(ctx context.Context, userID string) ([]DeviceInfo, error) {
	// 1. 获取用户的所有 Token
	userDevicesKey := fmt.Sprintf("%suser_devices:%s", constants.RedisKeyPrefix, userID)
	tokens, err := s.redisClient.SMembers(ctx, userDevicesKey).Result()
	if err != nil {
		return nil, err
	}

	// 2. 获取每个 Token 的设备信息
	var devices []DeviceInfo
	for _, token := range tokens {
		deviceKey := fmt.Sprintf("%sdevice:%s", constants.RedisKeyPrefix, token)
		deviceData, err := s.redisClient.Get(ctx, deviceKey).Result()
		if err == nil {
			var deviceInfo DeviceInfo
			if err := json.Unmarshal([]byte(deviceData), &deviceInfo); err == nil {
				devices = append(devices, deviceInfo)
			}
		}
	}

	return devices, nil
}

// RevokeDeviceByToken 根据 Token 撤销特定设备（完整的撤销流程）
func (s *TokenService) RevokeDeviceByToken(ctx context.Context, token string) error {
	// 1. 先获取设备信息（用于从用户设备列表中移除）
	deviceInfo, err := s.ValidateRefreshTokenWithDevice(ctx, token)
	if err != nil {
		return err // Token 已过期或无效
	}

	// 2. 删除设备信息
	deviceKey := fmt.Sprintf("%sdevice:%s", constants.RedisKeyPrefix, token)
	s.redisClient.Del(ctx, deviceKey)

	// 3. 删除 Refresh Token（旧格式兼容）
	tokenKey := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, token)
	s.redisClient.Del(ctx, tokenKey)

	// 4. 从用户设备列表中移除
	if deviceInfo.UserID != "" {
		userDevicesKey := fmt.Sprintf("%suser_devices:%s", constants.RedisKeyPrefix, deviceInfo.UserID)
		s.redisClient.SRem(ctx, userDevicesKey, token)
	}

	return nil
}

// RevokeAllDevices 撤销用户的所有设备
func (s *TokenService) RevokeAllDevices(ctx context.Context, userID string) error {
	// 1. 获取所有 Token
	userDevicesKey := fmt.Sprintf("%suser_devices:%s", constants.RedisKeyPrefix, userID)
	tokens, err := s.redisClient.SMembers(ctx, userDevicesKey).Result()
	if err != nil {
		return err
	}

	// 2. 逐个撤销
	for _, token := range tokens {
		if err := s.RevokeDeviceByToken(ctx, token); err != nil {
			// 记录错误但继续处理其他 Token
			continue
		}
	}

	// 3. 删除用户设备列表
	s.redisClient.Del(ctx, userDevicesKey)

	return nil
}
