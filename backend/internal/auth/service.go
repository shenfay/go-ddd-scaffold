package auth

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// RegisterCommand 注册命令
type RegisterCommand struct {
	Email    string
	Password string
}

// LoginCommand 登录命令
type LoginCommand struct {
	Email    string
	Password string
	IP       string
	UserAgent string
}

// RefreshTokenCommand 刷新 Token 命令
type RefreshTokenCommand struct {
	RefreshToken string
}

// LogoutCommand 退出登录命令
type LogoutCommand struct {
	UserID string
}

// ServiceAuthResponse 服务层认证响应（内部使用）
type ServiceAuthResponse struct {
	User         *User
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// Service 认证应用服务
type Service struct {
	userRepo     UserRepository
	tokenService *TokenService
	maxAttempts  int // 最大登录尝试次数
}

// NewService 创建认证服务
func NewService(userRepo UserRepository, tokenService *TokenService) *Service {
	return &Service{
		userRepo:     userRepo,
		tokenService: tokenService,
		maxAttempts:  5,
	}
}

// Register 处理用户注册
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (*ServiceAuthResponse, error) {
	// 1. 检查邮箱是否已存在
	if s.userRepo.ExistsByEmail(ctx, cmd.Email) {
		return nil, errors.ErrEmailAlreadyExists
	}
	
	// 2. 创建用户
	user, err := NewUser(cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}
	
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	
	// 3. 生成 Token
	tokens, err := s.tokenService.GenerateTokens(ctx, user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	
	// 4. TODO: 发布领域事件（异步）
	// s.eventBus.Publish(ctx, &UserRegisteredEvent{...})
	
	return &ServiceAuthResponse{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

// Login 处理用户登录
func (s *Service) Login(ctx context.Context, cmd LoginCommand) (*ServiceAuthResponse, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}
	
	// 2. 检查账户是否被锁定
	if user.IsLocked() {
		return nil, errors.ErrAccountLocked
	}
	
	// 3. 验证密码
	if !user.VerifyPassword(cmd.Password) {
		// 增加失败尝试次数
		user.IncrementFailedAttempts(s.maxAttempts)
		s.userRepo.Update(ctx, user)
		
		if user.IsLocked() {
			return nil, errors.ErrAccountLocked
		}
		
		return nil, errors.ErrInvalidCredentials
	}
	
	// 4. 重置失败次数，更新最后登录时间
	user.ResetFailedAttempts()
	user.UpdateLastLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	
	// 5. 生成 Token
	tokens, err := s.tokenService.GenerateTokens(ctx, user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	
	// 6. TODO: 发布领域事件
	// s.eventBus.Publish(ctx, &UserLoggedInEvent{...})
	
	return &ServiceAuthResponse{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

// Logout 处理用户退出
func (s *Service) Logout(ctx context.Context, cmd LogoutCommand) error {
	// 1. 撤销 Refresh Token
	if err := s.tokenService.RevokeToken(ctx, cmd.UserID); err != nil {
		return err
	}
	
	// 2. TODO: 发布领域事件
	// s.eventBus.Publish(ctx, &UserLoggedOutEvent{...})
	
	return nil
}

// RefreshToken 刷新 Access Token
func (s *Service) RefreshToken(ctx context.Context, cmd RefreshTokenCommand) (*ServiceAuthResponse, error) {
	// 1. 验证并解析 Refresh Token
	claims, err := s.tokenService.ValidateRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}
	
	// 2. 查找用户
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	
	// 3. 生成新的 Token 对
	tokens, err := s.tokenService.GenerateTokens(ctx, user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	
	// 4. 更新最后登录时间
	user.UpdateLastLogin()
	s.userRepo.Update(ctx, user)
	
	// 5. TODO: 发布领域事件
	// s.eventBus.Publish(ctx, &TokenRefreshedEvent{...})
	
	return &ServiceAuthResponse{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}
