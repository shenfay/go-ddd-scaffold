package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
)

// JWTClaims JWT 声明结构
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	JTI      string `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

// JWTService JWT 服务实现 - 实现 domain/user.TokenService 接口
type JWTService struct {
	secretKey  []byte
	accessExp  time.Duration
	refreshExp time.Duration
	issuer     string
	redis      *redis.Client // Redis 客户端用于令牌黑名单
}

// NewJWTService 创建 JWT 服务
func NewJWTService(secretKey string, accessExp, refreshExp time.Duration, issuer string) *JWTService {
	return &JWTService{
		secretKey:  []byte(secretKey),
		accessExp:  accessExp,
		refreshExp: refreshExp,
		issuer:     issuer,
	}
}

// SetRedisClient 设置 Redis 客户端（用于令牌黑名单）
func (s *JWTService) SetRedisClient(client *redis.Client) {
	s.redis = client
}

// GenerateTokenPair 生成令牌对 - 实现 TokenService 接口
func (s *JWTService) GenerateTokenPair(userID int64, username, email string) (*auth.TokenPair, error) {
	now := time.Now()

	// 生成 Access Token
	accessToken, err := s.generateToken(userID, username, email, now, s.accessExp)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token
	refreshToken, err := s.generateToken(userID, username, email, now, s.refreshExp)
	if err != nil {
		return nil, err
	}

	return &auth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(s.accessExp),
	}, nil
}

// generateToken 生成单个令牌
func (s *JWTService) generateToken(userID int64, username, email string, issuedAt time.Time, exp time.Duration) (string, error) {
	// 生成随机 JTI（Token 唯一标识），确保每次生成的 token 都不同
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		JTI:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(exp)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// generateJTI 生成唯一的 Token ID
func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ParseAccessToken 解析访问令牌 - 实现 TokenService 接口
func (s *JWTService) ParseAccessToken(tokenString string) (*auth.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return &auth.TokenClaims{
			UserID:    claims.UserID,
			Username:  claims.Username,
			Email:     claims.Email,
			JTI:       claims.JTI,
			IssuedAt:  claims.IssuedAt.Time,
			ExpiresAt: claims.ExpiresAt.Time,
		}, nil
	}

	return nil, errors.New("invalid token claims")
}

// ParseRefreshToken 解析刷新令牌 - 实现 TokenService 接口
func (s *JWTService) ParseRefreshToken(tokenString string) (*auth.TokenClaims, error) {
	// Refresh Token 使用相同的解析逻辑
	return s.ParseAccessToken(tokenString)
}

// ValidateToken 验证令牌 - 实现 TokenService 接口
func (s *JWTService) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	return s.ParseAccessToken(tokenString)
}

// BlacklistToken 将令牌加入黑名单 - 实现 TokenService 接口
func (s *JWTService) BlacklistToken(token string, expiresAt time.Time) error {
	if s.redis == nil {
		// Redis 未初始化，跳过黑名单（不返回错误）
		return nil
	}

	// 计算剩余有效期
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// 令牌已过期，不需要加入黑名单
		return nil
	}

	// 使用 token 的 jti 或 hash 作为 key
	// 这里简单使用 token 本身作为 key（实际可以使用 jti 或 hash 节省空间）
	key := fmt.Sprintf("blacklist:%s", token)

	// 存入 Redis，设置 TTL
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.redis.SetEx(ctx, key, "1", ttl).Err()
}

// IsTokenBlacklisted 检查令牌是否已在黑名单中 - 实现 TokenService 接口
func (s *JWTService) IsTokenBlacklisted(token string) (bool, error) {
	if s.redis == nil {
		// Redis 未初始化，认为不在黑名单中
		return false, nil
	}

	key := fmt.Sprintf("blacklist:%s", token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查 key 是否存在
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
