package auth

// ==================== 认证模块 Request 类型 ====================

// LoginRequest 登录请求
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	CurrentToken string `json:"current_token"`
}

// LogoutRequest 登出请求（可选，支持主动传递 access_token）
type LogoutRequest struct {
	AccessToken string `json:"access_token"`
}
