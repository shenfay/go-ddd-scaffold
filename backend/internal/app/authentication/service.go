package authentication

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	authErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/auth"
	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/pkg/metrics"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// JWTClaims JWT 自定义声明
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenService Token 服务接口
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

// TokenPair 访问令牌和刷新令牌对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// DeviceInfo 设备会话信息
type DeviceInfo struct {
	UserID     string `json:"user_id"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	DeviceType string `json:"device_type"`
	CreatedAt  string `json:"created_at"`
}

// Service 认证应用服务
type Service struct {
	userRepo       user.UserRepository
	resetTokenRepo authentication.PasswordResetTokenRepository
	emailTokenRepo authentication.EmailVerificationTokenRepository
	tokenService   TokenService
	publisher      *event.Publisher
	metrics        *metrics.Metrics
	maxAttempts    int
}

// NewService 创建认证服务实例
func NewService(
	userRepo user.UserRepository,
	resetTokenRepo authentication.PasswordResetTokenRepository,
	emailTokenRepo authentication.EmailVerificationTokenRepository,
	tokenService TokenService,
	publisher *event.Publisher,
	m *metrics.Metrics,
) *Service {
	return &Service{
		userRepo:       userRepo,
		resetTokenRepo: resetTokenRepo,
		emailTokenRepo: emailTokenRepo,
		tokenService:   tokenService,
		publisher:      publisher,
		metrics:        m,
		maxAttempts:    5,
	}
}

// Register 创建用户账户并返回认证令牌
//
// 注册流程：
// 1. 验证邮箱唯一性
// 2. 创建用户实体（密码已加密）
// 3. 生成访问令牌和刷新令牌
// 4. 发布 UserRegistered 领域事件
//
// 参数：
//   - ctx: 请求上下文
//   - cmd: 注册命令（包含邮箱和密码）
//
// 返回：
//   - *ServiceAuthResponse: 用户数据和认证令牌
//   - error: 注册失败时返回错误（邮箱已存在、验证错误等）
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (*ServiceAuthResponse, error) {
	// 1. 检查邮箱是否已存在
	if s.userRepo.ExistsByEmail(ctx, cmd.Email) {
		return nil, userErr.ErrEmailAlreadyExists
	}

	// 2. 创建用户实体
	u, err := user.NewUser(cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	// 记录用户注册指标
	if s.metrics != nil {
		s.metrics.IncUserRegistration()
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

	// 5. 触发邮箱验证流程
	s.SendVerificationEmail(ctx, u.ID)

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

		// 记录认证失败指标
		if s.metrics != nil {
			if u.IsLocked() {
				s.metrics.IncAuthFailure("password", "account_locked")
			} else {
				s.metrics.IncAuthFailure("password", "invalid_credentials")
			}
		}

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

	// 记录认证成功指标
	if s.metrics != nil {
		s.metrics.IncAuthSuccess("password")
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

// RequestPasswordResetCommand 请求密码重置命令
type RequestPasswordResetCommand struct {
	Email string
}

// ResetPasswordCommand 重置密码命令
type ResetPasswordCommand struct {
	Token       string
	NewPassword string
}

// RequestPasswordReset 处理密码重置请求
//
// 业务流程:
// 1. 查找用户(不暴露用户是否存在)
// 2. 检查账户状态
// 3. 生成重置令牌
// 4. 保存令牌到数据库
// 5. 发布 PasswordResetRequested 事件(触发邮件发送)
//
// 安全考虑:
// - 即使邮箱不存在也返回成功,防止邮箱枚举攻击
// - 令牌有效期 1 小时
//
// 参数:
//   - ctx: 请求上下文
//   - cmd: 密码重置请求命令
//
// 返回:
//   - error: 处理失败时返回错误
func (s *Service) RequestPasswordReset(ctx context.Context, cmd RequestPasswordResetCommand) error {
	// 1. 查找用户(不暴露用户是否存在)
	u, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		// 安全考虑:即使邮箱不存在也返回成功
		return nil
	}

	// 2. 检查邮箱是否已验证
	if !u.EmailVerified {
		return userErr.ErrEmailNotVerified
	}

	// 3. 检查账户状态
	if u.IsLocked() {
		return nil // 安全考虑:不暴露账户锁定状态
	}

	// 4. 生成重置令牌
	resetToken, err := authentication.NewPasswordResetToken(u.ID)
	if err != nil {
		return err
	}

	// 5. 保存令牌到数据库
	if err := s.resetTokenRepo.Create(ctx, resetToken); err != nil {
		return err
	}

	// 6. 记录指标
	if s.metrics != nil {
		s.metrics.IncPasswordResetRequested()
	}

	// 7. 发布事件(异步触发邮件发送)
	s.publisher.Publish(ctx, &authentication.PasswordResetRequested{
		UserID: u.ID,
		Email:  u.Email,
	})

	return nil
}

// ResetPassword 执行密码重置
//
// 业务流程:
// 1. 验证令牌有效性
// 2. 检查令牌是否过期(双重验证)
// 3. 查找用户
// 4. 调用领域行为重置密码
// 5. 标记令牌已使用
// 6. 保存用户
// 7. 发布 PasswordResetCompleted 事件
//
// 参数:
//   - ctx: 请求上下文
//   - cmd: 密码重置命令(包含令牌和新密码)
//
// 返回:
//   - error: 重置失败时返回错误(令牌无效、密码强度不足等)
func (s *Service) ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error {
	// 1. 验证令牌
	resetToken, err := s.resetTokenRepo.FindByToken(ctx, cmd.Token)
	if err != nil {
		if s.metrics != nil {
			s.metrics.IncPasswordResetFailed("invalid_token")
		}
		return authErr.ErrInvalidResetToken
	}

	// 2. 检查令牌是否过期(双重验证)
	if resetToken.IsExpired() {
		if s.metrics != nil {
			s.metrics.IncPasswordResetFailed("expired_token")
		}
		return authErr.ErrInvalidResetToken
	}

	// 3. 查找用户
	u, err := s.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		if s.metrics != nil {
			s.metrics.IncPasswordResetFailed("user_not_found")
		}
		return userErr.ErrNotFound
	}

	// 4. 重置密码(领域行为)
	if err := u.ChangePassword(cmd.NewPassword); err != nil {
		if s.metrics != nil {
			s.metrics.IncPasswordResetFailed("password_change_failed")
		}
		return err
	}

	// 5. 标记令牌已使用
	resetToken.MarkAsUsed()
	if err := s.resetTokenRepo.MarkAsUsed(ctx, resetToken.ID); err != nil {
		return err
	}

	// 6. 保存用户
	if err := s.userRepo.Update(ctx, u); err != nil {
		return err
	}

	// 7. 记录指标
	if s.metrics != nil {
		s.metrics.IncPasswordResetCompleted()
	}

	// 8. 发布事件
	s.publisher.Publish(ctx, &authentication.PasswordResetCompleted{
		UserID: u.ID,
		Email:  u.Email,
	})

	return nil
}
