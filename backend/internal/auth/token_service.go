package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	
	"github.com/shenfay/go-ddd-scaffold/pkg/constants"
	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils/ulid"
)

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
	redisClient *redis.Client
	jwtSecret   []byte
	issuer      string
	accessExpire time.Duration
	refreshExpire time.Duration
}

// NewTokenService 创建 Token 服务
func NewTokenService(redisClient *redis.Client, jwtSecret string, issuer string, accessExpire, refreshExpire time.Duration) *TokenService {
	return &TokenService{
		redisClient:  redisClient,
		jwtSecret:    []byte(jwtSecret),
		issuer:       issuer,
		accessExpire: accessExpire,
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
	key := fmt.Sprintf("%s%s", constants.RedisKeyRefreshToken, tokenID)
	return s.redisClient.Del(ctx, key).Err()
}

// RevokeAllUserTokens 撤销用户的所有 Token（修改密码等场景）
func (s *TokenService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	// TODO: 如果需要支持多设备登录，需要维护一个用户到 Token 的映射
	// 这里简化实现，暂时不处理
	return nil
}
