package authentication

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	authErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/auth"
	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// JWTClaims represents custom JWT token claims.
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenService defines the interface for token generation and validation.
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

// TokenPair contains access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// DeviceInfo contains device session information.
type DeviceInfo struct {
	UserID     string `json:"user_id"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	DeviceType string `json:"device_type"`
	CreatedAt  string `json:"created_at"`
}

// Service handles authentication business logic.
type Service struct {
	userRepo     user.UserRepository
	tokenService TokenService
	publisher    *event.Publisher
	maxAttempts  int
}

// NewService creates a new authentication service instance.
func NewService(userRepo user.UserRepository, tokenService TokenService, publisher *event.Publisher) *Service {
	return &Service{
		userRepo:     userRepo,
		tokenService: tokenService,
		publisher:    publisher,
		maxAttempts:  5,
	}
}

// Register creates a new user account and returns authentication tokens.
//
// The registration process:
// 1. Validates email uniqueness
// 2. Creates user entity with encrypted password
// 3. Generates access and refresh tokens
// 4. Publishes UserRegistered domain event
//
// Parameters:
//   - ctx: context for request lifecycle and tracing
//   - cmd: registration command containing email and password
//
// Returns:
//   - *ServiceAuthResponse: user data and authentication tokens
//   - error: when registration fails (email exists, validation error, etc.)
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (*ServiceAuthResponse, error) {
	// 1. Check if email already exists
	if s.userRepo.ExistsByEmail(ctx, cmd.Email) {
		return nil, userErr.ErrEmailAlreadyExists
	}

	// 2. Create user entity
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
	s.publisher.Publish(ctx, &user.UserRegistered{
		UserID:    u.ID,
		Email:     u.Email,
		Timestamp: utils.Now(),
	})

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
		return nil, authErr.ErrInvalidCredentials
	}

	// 2. 检查账户是否被锁定
	if u.IsLocked() {
		return nil, authErr.ErrAccountLocked
	}

	// 3. 验证密码
	if !u.VerifyPassword(cmd.Password) {
		u.IncrementFailedAttempts(s.maxAttempts)
		s.userRepo.Update(ctx, u)

		if u.IsLocked() {
			return nil, authErr.ErrAccountLocked
		}

		return nil, authErr.ErrInvalidCredentials
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
		// 设备信息存储失败不影响登录流程，仅记录警告
		// 日志已在 StoreDeviceInfo 内部处理
	}

	// 7. 发布领域事件（异步）
	s.publisher.Publish(ctx, &user.UserLoggedIn{
		UserID:    u.ID,
		Email:     u.Email,
		IP:        cmd.IP,
		UserAgent: cmd.UserAgent,
		Device:    cmd.DeviceType,
		Timestamp: utils.Now(),
	})

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
	if err == nil {
		s.publisher.Publish(ctx, &user.UserLoggedOut{
			UserID:    u.ID,
			Email:     u.Email,
			Timestamp: utils.Now(),
		})
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
		return nil, userErr.ErrNotFound
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
	s.publisher.Publish(ctx, &user.TokenRefreshed{
		UserID:    u.ID,
		OldToken:  cmd.RefreshToken,
		NewToken:  tokens.RefreshToken,
		Timestamp: utils.Now(),
	})

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
