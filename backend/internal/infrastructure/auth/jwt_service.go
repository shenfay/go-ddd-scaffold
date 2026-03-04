package auth

import (
	"time"

	"go-ddd-scaffold/internal/domain/user/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// jwtService JWT服务实现
type jwtService struct {
	secretKey []byte
	expireIn  time.Duration
}

// NewJWTService 创建JWT服务实例
func NewJWTService(secretKey string, expireIn time.Duration) entity.JWTService {
	return &jwtService{
		secretKey: []byte(secretKey),
		expireIn:  expireIn,
	}
}

// GenerateToken 生成 JWT 令牌（仅包含用户 ID）
func (s *jwtService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID.String(),
		"exp":    time.Now().Add(s.expireIn).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken 验证 JWT 令牌
func (s *jwtService) ValidateToken(tokenString string) (*entity.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenUnverifiable
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	// 解析用户 ID
	userIDStr, ok := claims["userId"].(string)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, jwt.ErrTokenMalformed
	}

	// 简化版 TokenClaims，只包含 UserID
	claimsResult := &entity.TokenClaims{
		UserID: userID,
	}

	return claimsResult, nil
}
