package authentication

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	apperrors "github.com/shenfay/go-ddd-scaffold/pkg/errors"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// JWTClaims JWT声明
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenService Token服务接口
type TokenService interface {
	GenerateTokens(ctx context.Context, userID, email string) (*TokenPair, error)
	RevokeToken(ctx context.Context, tokenID string) error
	ValidateRefreshTokenWithDevice(ctx context.Context, token string) (*DeviceInfo, error)
	ValidateAccessToken(tokenString string) (*JWTClaims, error)
	StoreDeviceInfo(ctx context.Context, token string, deviceInfo DeviceInfo) error
	RevokeDeviceByToken(ctx context.Context, token string) error
	RevokeAllDevices(ctx context.Context, userID string) error
	GetUserDevices(ctx context.Context, userID string) ([]DeviceInfo, error)
}

// TokenPair Token对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	UserID     string `json:"user_id"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	DeviceType string `json:"device_type"`
	CreatedAt  string `json:"created_at"`
}

// Service 认证应用服务
type Service struct {
	userRepo     user.UserRepository
	tokenService TokenService
	eventBus     messaging.EventBus
	maxAttempts  int
}

// NewService 创建认证服务
func NewService(userRepo user.UserRepository, tokenService TokenService) *Service {
	return &Service{
		userRepo:     userRepo,
		tokenService: tokenService,
		eventBus:     nil,
		maxAttempts:  5,
	}
}

// SetEventBus 设置事件总线（可选）
func (s *Service) SetEventBus(eventBus messaging.EventBus) {
	s.eventBus = eventBus
}

// Register 处理用户注册
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (*ServiceAuthResponse, error) {
	// 1. 检查邮箱是否已存在
	if s.userRepo.ExistsByEmail(ctx, cmd.Email) {
		return nil, apperrors.ErrEmailAlreadyExists
	}

	// 2. 创建用户
	u, err := user.NewUser(cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	// 3. 生成 Token
	tokens, err := s.tokenService.GenerateTokens(ctx, u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	// 4. 发布领域事件（异步）
	if s.eventBus != nil {
		evt := &user.UserRegistered{
			UserID:    u.ID,
			Email:     u.Email,
			Timestamp: utils.Now(),
		}
		if err := s.eventBus.Publish(ctx, evt); err != nil {
			log.Printf("Failed to publish UserRegistered event: %v", err)
		}
	}

	return &ServiceAuthResponse{
		User:         u,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

// Login 处理用户登录
func (s *Service) Login(ctx context.Context, cmd LoginCommand) (*ServiceAuthResponse, error) {
	// 1. 查找用户
	u, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// 2. 检查账户是否被锁定
	if u.IsLocked() {
		return nil, apperrors.ErrAccountLocked
	}

	// 3. 验证密码
	if !u.VerifyPassword(cmd.Password) {
		u.IncrementFailedAttempts(s.maxAttempts)
		s.userRepo.Update(ctx, u)

		if u.IsLocked() {
			return nil, apperrors.ErrAccountLocked
		}

		return nil, apperrors.ErrInvalidCredentials
	}

	// 4. 重置失败次数，更新最后登录时间
	u.ResetFailedAttempts()
	u.UpdateLastLogin()
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	// 5. 生成 Token
	tokens, err := s.tokenService.GenerateTokens(ctx, u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	// 6. 存储设备信息到 Redis
	if err := s.tokenService.StoreDeviceInfo(ctx, tokens.RefreshToken, DeviceInfo{
		UserID:     u.ID,
		IP:         cmd.IP,
		UserAgent:  cmd.UserAgent,
		DeviceType: cmd.DeviceType,
	}); err != nil {
		log.Printf("Failed to store device info: %v", err)
	}

	// 7. 发布领域事件（异步）
	if s.eventBus != nil {
		evt := &user.UserLoggedIn{
			UserID:    u.ID,
			Email:     u.Email,
			IP:        cmd.IP,
			UserAgent: cmd.UserAgent,
			Device:    cmd.DeviceType,
			Timestamp: utils.Now(),
		}
		if err := s.eventBus.Publish(ctx, evt); err != nil {
			log.Printf("Failed to publish UserLoggedIn event: %v", err)
		}
	}

	return &ServiceAuthResponse{
		User:         u,
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

	// 2. 发布领域事件（异步）
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err == nil && s.eventBus != nil {
		evt := &user.UserLoggedOut{
			UserID:    u.ID,
			Email:     u.Email,
			Timestamp: utils.Now(),
		}
		if err := s.eventBus.Publish(ctx, evt); err != nil {
			log.Printf("Failed to publish UserLoggedOut event: %v", err)
		}
	}

	return nil
}

// RefreshToken 刷新 Access Token
func (s *Service) RefreshToken(ctx context.Context, cmd RefreshTokenCommand) (*ServiceAuthResponse, error) {
	// 1. 验证并解析 Refresh Token
	deviceInfo, err := s.tokenService.ValidateRefreshTokenWithDevice(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 2. 查找用户
	u, err := s.userRepo.FindByID(ctx, deviceInfo.UserID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	// 3. 撤销旧的 Refresh Token
	if err := s.tokenService.RevokeDeviceByToken(ctx, cmd.RefreshToken); err != nil {
		return nil, err
	}

	// 4. 生成新的 Token 对
	tokens, err := s.tokenService.GenerateTokens(ctx, u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	// 5. 存储新设备信息
	if err := s.tokenService.StoreDeviceInfo(ctx, tokens.RefreshToken, DeviceInfo{
		UserID:     u.ID,
		IP:         deviceInfo.IP,
		UserAgent:  deviceInfo.UserAgent,
		DeviceType: deviceInfo.DeviceType,
	}); err != nil {
		return nil, err
	}

	// 6. 更新最后登录时间
	u.UpdateLastLogin()
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	// 7. 发布领域事件
	if s.eventBus != nil {
		evt := &user.TokenRefreshed{
			UserID:    u.ID,
			OldToken:  cmd.RefreshToken,
			NewToken:  tokens.RefreshToken,
			Timestamp: utils.Now(),
		}
		if err := s.eventBus.Publish(ctx, evt); err != nil {
			log.Printf("Failed to publish TokenRefreshed event: %v", err)
		}
	}

	return &ServiceAuthResponse{
		User:         u,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

// GetUserByID 根据 ID 获取用户
func (s *Service) GetUserByID(ctx context.Context, userID string) (*user.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
