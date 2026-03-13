package auth

// ==================== 认证模块 Swagger Response 类型 ====================

// LoginResponse 登录响应
// @Description 用户登录返回的令牌信息
type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	User         UserInfo `json:"user"`
}

// UserInfo 用户基本信息
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RegisterResponse 注册响应
// @Description 用户注册成功返回的信息
type RegisterResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RefreshTokenResponse 刷新令牌响应
// @Description 刷新令牌后返回的新令牌信息
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// CurrentUserResponse 当前用户响应
// @Description 获取当前登录用户信息
type CurrentUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
