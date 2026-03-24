package auth

import (
	"context"
)

// TokenServiceAdapter TokenService 端口适配器
// 将 infrastructure 的 TokenService 适配为 application/ports 的 TokenService 接口
type TokenServiceAdapter struct {
	service TokenService
}

// NewTokenServiceAdapter 创建 TokenService 适配器
func NewTokenServiceAdapter(service TokenService) *TokenServiceAdapter {
	return &TokenServiceAdapter{
		service: service,
	}
}

// GenerateToken 生成访问令牌和刷新令牌
func (a *TokenServiceAdapter) GenerateToken(ctx context.Context, userID int64, email string) (string, string, error) {
	pair, err := a.service.GenerateTokenPair(userID, "", email)
	if err != nil {
		return "", "", err
	}
	return pair.AccessToken, pair.RefreshToken, nil
}

// ValidateToken 验证访问令牌
func (a *TokenServiceAdapter) ValidateToken(ctx context.Context, token string) (int64, error) {
	claims, err := a.service.ValidateToken(token)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// RefreshToken 刷新访问令牌
func (a *TokenServiceAdapter) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// TODO: 实现刷新令牌逻辑
	return "", "", nil
}

// RevokeToken 吊销令牌
func (a *TokenServiceAdapter) RevokeToken(ctx context.Context, userID int64) error {
	// TODO: 实现吊销令牌逻辑
	return nil
}
