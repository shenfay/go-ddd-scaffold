package auth

import "time"

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// TokenClaims 令牌声明
type TokenClaims struct {
	UserID    int64
	Username  string
	Email     string
	JTI       string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// TokenService Token 服务端口
type TokenService interface {
	// GenerateTokenPair 生成令牌对
	GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)

	// ParseAccessToken 解析访问令牌
	ParseAccessToken(token string) (*TokenClaims, error)

	// ParseRefreshToken 解析刷新令牌
	ParseRefreshToken(token string) (*TokenClaims, error)

	// ValidateToken 验证令牌
	ValidateToken(token string) (*TokenClaims, error)

	// BlacklistToken 将令牌加入黑名单（登出时使用）
	BlacklistToken(token string, expiresAt time.Time) error

	// IsTokenBlacklisted 检查令牌是否已在黑名单中
	IsTokenBlacklisted(token string) (bool, error)
}
