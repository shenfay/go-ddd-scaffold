package auth

import (
	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
)

// Mapper DTO 转换器
// 负责 HTTP 请求/响应与 Application DTO 之间的转换
type Mapper struct{}

// NewMapper 创建转换器
func NewMapper() *Mapper {
	return &Mapper{}
}

// ToAuthenticateCommand 转换为认证命令
func (m *Mapper) ToAuthenticateCommand(req *LoginRequest, ipAddress, userAgent string) *authApp.AuthenticateCommand {
	return &authApp.AuthenticateCommand{
		Identifier: req.UsernameOrEmail,
		Password:   req.Password,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}
}

// ToRegisterCommand 转换为注册命令
func (m *Mapper) ToRegisterCommand(req *RegisterRequest) *authApp.RegisterCommand {
	return &authApp.RegisterCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
}

// ToRefreshTokenCommand 转换为刷新令牌命令
func (m *Mapper) ToRefreshTokenCommand(req *RefreshTokenRequest, ipAddress, userAgent string) *authApp.RefreshTokenCommand {
	return &authApp.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
		CurrentToken: req.CurrentToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}
}

// ToLogoutCommand 转换为登出命令
func (m *Mapper) ToLogoutCommand(userID int64, accessToken, ipAddress, userAgent string) *authApp.LogoutCommand {
	return &authApp.LogoutCommand{
		UserID:      userID,
		AccessToken: accessToken,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}
}

// ToLoginResponse 转换为登录响应
func (m *Mapper) ToLoginResponse(result *authApp.AuthenticateResult) *LoginResponse {
	return &LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
		User: UserInfo{
			ID:       result.UserID,
			Username: result.Username,
			Email:    result.Email,
		},
	}
}

// ToRegisterResponse 转换为注册响应
func (m *Mapper) ToRegisterResponse(result *authApp.RegisterResult) *RegisterResponse {
	return &RegisterResponse{
		UserID:   result.UserID,
		Username: result.Username,
		Email:    result.Email,
	}
}

// ToRefreshTokenResponse 转换为刷新令牌响应
func (m *Mapper) ToRefreshTokenResponse(result *authApp.RefreshTokenResult) *RefreshTokenResponse {
	return &RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
	}
}

// ToCurrentUserResponse 转换为当前用户响应
func (m *Mapper) ToCurrentUserResponse(user *authApp.UserInfoResult) *CurrentUserResponse {
	return &CurrentUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Status:      user.Status,
	}
}
