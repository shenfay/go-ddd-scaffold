package auth

import (
	"time"

	ports_auth "github.com/shenfay/go-ddd-scaffold/internal/application/ports/auth"
)

// TokenServiceAdapter TokenService 端口适配器
type TokenServiceAdapter struct {
	service TokenService
}

// NewTokenServiceAdapter 创建 TokenService 适配器
func NewTokenServiceAdapter(service TokenService) *TokenServiceAdapter {
	return &TokenServiceAdapter{
		service: service,
	}
}

// GenerateTokenPair 生成令牌对
func (a *TokenServiceAdapter) GenerateTokenPair(userID int64, username, email string) (*ports_auth.TokenPair, error) {
	pair, err := a.service.GenerateTokenPair(userID, username, email)
	if err != nil {
		return nil, err
	}
	// 转换为 Port 的 TokenPair
	return &ports_auth.TokenPair{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    pair.ExpiresAt,
	}, nil
}

// ParseAccessToken 解析访问令牌
func (a *TokenServiceAdapter) ParseAccessToken(token string) (*ports_auth.TokenClaims, error) {
	claims, err := a.service.ParseAccessToken(token)
	if err != nil {
		return nil, err
	}
	// 转换为 Port 的 TokenClaims
	return &ports_auth.TokenClaims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Email:     claims.Email,
		JTI:       claims.JTI,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}

// ParseRefreshToken 解析刷新令牌
func (a *TokenServiceAdapter) ParseRefreshToken(token string) (*ports_auth.TokenClaims, error) {
	claims, err := a.service.ParseRefreshToken(token)
	if err != nil {
		return nil, err
	}
	return &ports_auth.TokenClaims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Email:     claims.Email,
		JTI:       claims.JTI,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}

// ValidateToken 验证令牌
func (a *TokenServiceAdapter) ValidateToken(token string) (*ports_auth.TokenClaims, error) {
	claims, err := a.service.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return &ports_auth.TokenClaims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Email:     claims.Email,
		JTI:       claims.JTI,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}

// BlacklistToken 将令牌加入黑名单
func (a *TokenServiceAdapter) BlacklistToken(token string, expiresAt time.Time) error {
	return a.service.BlacklistToken(token, expiresAt)
}

// IsTokenBlacklisted 检查令牌是否已在黑名单中
func (a *TokenServiceAdapter) IsTokenBlacklisted(token string) (bool, error) {
	return a.service.IsTokenBlacklisted(token)
}

var _ ports_auth.TokenService = (*TokenServiceAdapter)(nil)
