package auth

// ============================================================================
// Input DTOs (Commands)
// ============================================================================

// AuthenticateCommand 认证命令
type AuthenticateCommand struct {
	Identifier string // 用户名或邮箱
	Password   string
	IPAddress  string
	UserAgent  string
}

// RegisterCommand 注册命令
type RegisterCommand struct {
	Username string
	Email    string
	Password string
}

// RefreshTokenCommand 刷新令牌命令
type RefreshTokenCommand struct {
	RefreshToken string
	IPAddress    string
	UserAgent    string
}

// LogoutCommand 登出命令
type LogoutCommand struct {
	UserID      int64
	AccessToken string
	IPAddress   string
	UserAgent   string
}

// ============================================================================
// Output DTOs (Results)
// ============================================================================

// AuthenticateResult 认证结果
type AuthenticateResult struct {
	UserID       string
	Username     string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // 过期时间（秒）
}

// RegisterResult 注册结果
type RegisterResult struct {
	UserID   string
	Username string
	Email    string
}

// RefreshTokenResult 刷新令牌结果
type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// LogoutResult 登出结果
type LogoutResult struct {
	Success bool
}
