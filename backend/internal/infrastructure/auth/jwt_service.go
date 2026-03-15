package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
)

// JWTClaims JWT 声明结构
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// JWTService JWT 服务实现 - 实现 domain/user.TokenService 接口
type JWTService struct {
	secretKey  []byte
	accessExp  time.Duration
	refreshExp time.Duration
	issuer     string
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
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(exp)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
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
