package auth

// ============================================================================
// Input DTOs (Commands)
// ============================================================================

// AuthenticateCommand 认证命令
type AuthenticateCommand struct {
	Identifier string `json:"identifier"` // 用户名或邮箱
	Password   string `json:"password"`
	IPAddress  string `json:"ip_address,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
}

// RegisterCommand 注册命令
type RegisterCommand struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RefreshTokenCommand 刷新令牌命令
type RefreshTokenCommand struct {
	RefreshToken string `json:"refresh_token"`
	CurrentToken string `json:"current_token,omitempty"`
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
}

// LogoutCommand 登出命令
type LogoutCommand struct {
	UserID      int64  `json:"user_id"`
	AccessToken string `json:"access_token"`
	IPAddress   string `json:"ip_address,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
}

// ============================================================================
// Output DTOs (Results)
// ============================================================================

// AuthenticateResult 认证结果
type AuthenticateResult struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 过期时间（秒）
}

// RegisterResult 注册结果
type RegisterResult struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RefreshTokenResult 刷新令牌结果
type RefreshTokenResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// LogoutResult 登出结果
type LogoutResult struct {
	Success bool `json:"success"`
}

// UserInfoResult 用户信息结果（用于获取当前用户）
type UserInfoResult struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}
