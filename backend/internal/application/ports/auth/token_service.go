package auth

import "context"

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// TokenService Token 服务端口
type TokenService interface {
	// GenerateToken 生成访问令牌和刷新令牌
	GenerateToken(ctx context.Context, userID int64, email string) (accessToken, refreshToken string, err error)

	// ValidateToken 验证访问令牌
	ValidateToken(ctx context.Context, token string) (int64, error)

	// RefreshToken 刷新访问令牌
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)

	// RevokeToken 吊销令牌
	RevokeToken(ctx context.Context, userID int64) error
}
