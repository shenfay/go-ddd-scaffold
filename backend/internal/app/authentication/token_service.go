package authentication

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

// TokenServiceImpl Token服务实现
type TokenServiceImpl struct {
	redisClient   *redis.Client
	jwtSecret     []byte
	issuer        string
	accessExpire  time.Duration
	refreshExpire time.Duration
}

// NewTokenServiceImpl 创建Token服务实现
func NewTokenServiceImpl(redisClient *redis.Client, jwtSecret string, issuer string, accessExpire, refreshExpire time.Duration) TokenService {
	return &TokenServiceImpl{
		redisClient:   redisClient,
		jwtSecret:     []byte(jwtSecret),
		issuer:        issuer,
		accessExpire:  accessExpire,
		refreshExpire: refreshExpire,
	}
}

// GenerateTokens 生成 Token 对
func (s *TokenServiceImpl) GenerateTokens(ctx context.Context, userID, email string) (*TokenPair, error) {
	now := time.Now()

	accessToken, err := s.generateAccessToken(userID, email, now)
	if err != nil {
		return nil, err
	}

	refreshTokenID := ulid.GenerateTokenID()
	refreshToken := refreshTokenID

	if err := s.storeRefreshToken(ctx, refreshTokenID, userID); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.accessExpire,
	}, nil
}

func (s *TokenServiceImpl) generateAccessToken(userID, email string, issuedAt time.Time) (string, error) {
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

func (s *TokenServiceImpl) storeRefreshToken(ctx context.Context, tokenID, userID string) error {
	key := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, tokenID)
	return s.redisClient.Set(ctx, key, userID, s.refreshExpire).Err()
}

// ValidateAccessToken 验证 Access Token
func (s *TokenServiceImpl) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
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

	if claims.TokenType != "access" {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshTokenWithDevice 验证 Refresh Token 并返回设备信息
func (s *TokenServiceImpl) ValidateRefreshTokenWithDevice(ctx context.Context, token string) (*DeviceInfo, error) {
	key := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, token)
	userID, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrInvalidToken
	}

	deviceKey := fmt.Sprintf("%sdevice:%s", constants.RedisKeyPrefix, token)
	deviceData, err := s.redisClient.Get(ctx, deviceKey).Result()
	if err == nil {
		var deviceInfo DeviceInfo
		if err := json.Unmarshal([]byte(deviceData), &deviceInfo); err == nil {
			return &deviceInfo, nil
		}
	}

	return &DeviceInfo{
		UserID: userID,
	}, nil
}

// StoreDeviceInfo 存储设备信息到 Redis
func (s *TokenServiceImpl) StoreDeviceInfo(ctx context.Context, token string, deviceInfo DeviceInfo) error {
	now := time.Now()
	deviceInfo.CreatedAt = now.Format(time.RFC3339)

	data, err := json.Marshal(deviceInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal device info: %w", err)
	}

	deviceKey := fmt.Sprintf("%sdevice:%s", constants.RedisKeyPrefix, token)
	if err := s.redisClient.Set(ctx, deviceKey, string(data), s.refreshExpire).Err(); err != nil {
		return err
	}

	userDevicesKey := fmt.Sprintf("%suser_devices:%s", constants.RedisKeyPrefix, deviceInfo.UserID)
	if err := s.redisClient.SAdd(ctx, userDevicesKey, token).Err(); err != nil {
		return err
	}

	s.redisClient.Expire(ctx, userDevicesKey, s.refreshExpire)

	return nil
}

// RevokeToken 撤销 Token
func (s *TokenServiceImpl) RevokeToken(ctx context.Context, tokenID string) error {
	key := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, tokenID)
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

// RevokeDeviceByToken 根据 Token 撤销特定设备
func (s *TokenServiceImpl) RevokeDeviceByToken(ctx context.Context, token string) error {
	deviceInfo, err := s.ValidateRefreshTokenWithDevice(ctx, token)
	if err != nil {
		return err
	}

	deviceKey := fmt.Sprintf("%sdevice:%s", constants.RedisKeyPrefix, token)
	s.redisClient.Del(ctx, deviceKey)

	tokenKey := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, token)
	s.redisClient.Del(ctx, tokenKey)

	if deviceInfo.UserID != "" {
		userDevicesKey := fmt.Sprintf("%suser_devices:%s", constants.RedisKeyPrefix, deviceInfo.UserID)
		s.redisClient.SRem(ctx, userDevicesKey, token)
	}

	return nil
}

// RevokeAllDevices 撤销用户的所有设备
func (s *TokenServiceImpl) RevokeAllDevices(ctx context.Context, userID string) error {
	userDevicesKey := fmt.Sprintf("%suser_devices:%s", constants.RedisKeyPrefix, userID)
	tokens, err := s.redisClient.SMembers(ctx, userDevicesKey).Result()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if err := s.RevokeDeviceByToken(ctx, token); err != nil {
			continue
		}
	}

	s.redisClient.Del(ctx, userDevicesKey)

	return nil
}

// GetUserDevices 获取用户的所有设备列表
func (s *TokenServiceImpl) GetUserDevices(ctx context.Context, userID string) ([]DeviceInfo, error) {
	userDevicesKey := fmt.Sprintf("%suser_devices:%s", constants.RedisKeyPrefix, userID)
	tokens, err := s.redisClient.SMembers(ctx, userDevicesKey).Result()
	if err != nil {
		return nil, err
	}

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
