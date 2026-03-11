package user

import (
	"context"
	"errors"
	"time"

	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// AuthenticationService 认证服务
type AuthenticationService struct {
	userRepo       UserRepository
	tokenService   TokenService
	passwordPolicy PasswordPolicyService
}

// NewAuthenticationService 创建认证服务
func NewAuthenticationService(
	userRepo UserRepository,
	tokenService TokenService,
	passwordPolicy PasswordPolicyService,
) *AuthenticationService {
	return &AuthenticationService{
		userRepo:       userRepo,
		tokenService:   tokenService,
		passwordPolicy: passwordPolicy,
	}
}

// AuthenticateResult 认证结果
type AuthenticateResult struct {
	UserID       UserID
	Username     string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// Authenticate 用户认证
func (s *AuthenticationService) Authenticate(ctx context.Context, usernameOrEmail, password string, ipAddress, userAgent string) (*AuthenticateResult, error) {
	// 1. 查找用户（支持用户名或邮箱登录）
	var u *User
	var err error

	// 尝试作为邮箱查找
	u, err = s.userRepo.FindByEmail(ctx, usernameOrEmail)
	if err != nil {
		// 尝试作为用户名查找
		u, err = s.userRepo.FindByUsername(ctx, usernameOrEmail)
		if err != nil {
			return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "invalid username or password")
		}
	}

	// 2. 验证密码
	if !u.Password().Matches(password) {
		// 记录失败登录
		u.RecordFailedLogin(ipAddress, userAgent, "invalid_password")
		_ = s.userRepo.Save(ctx, u)
		return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "invalid username or password")
	}

	// 3. 检查账户状态
	if !u.CanLogin() {
		if u.IsLocked() {
			return nil, ddd.NewBusinessError("ACCOUNT_LOCKED", "account is locked")
		}
		if u.Status() == UserStatusInactive {
			return nil, ddd.NewBusinessError("ACCOUNT_INACTIVE", "account is inactive")
		}
		if u.Status() == UserStatusPending {
			return nil, ddd.NewBusinessError("ACCOUNT_PENDING", "account is pending activation")
		}
		return nil, ddd.NewBusinessError("ACCOUNT_DISABLED", "account cannot login")
	}

	// 4. 记录成功登录
	u.RecordLogin(ipAddress, userAgent)
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	// 5. 生成令牌
	tokenPair, err := s.tokenService.GenerateTokenPair(u.ID().(UserID))
	if err != nil {
		return nil, ddd.NewBusinessErrorWithDetails("TOKEN_GENERATION_FAILED", "failed to generate tokens", err.Error())
	}

	return &AuthenticateResult{
		UserID:       u.ID().(UserID),
		Username:     u.Username().Value(),
		Email:        u.Email().Value(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthenticationService) RefreshToken(ctx context.Context, refreshToken string) (*AuthenticateResult, error) {
	// 1. 验证刷新令牌
	claims, err := s.tokenService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, ddd.NewBusinessError("INVALID_REFRESH_TOKEN", "invalid or expired refresh token")
	}

	// 2. 查找用户
	userID := NewUserID(claims.UserID)
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ddd.NewBusinessError("USER_NOT_FOUND", "user not found")
	}

	// 3. 检查用户状态
	if !u.CanLogin() {
		return nil, ddd.NewBusinessError("ACCOUNT_CANNOT_LOGIN", "account cannot login")
	}

	// 4. 生成新令牌
	tokenPair, err := s.tokenService.GenerateTokenPair(userID)
	if err != nil {
		return nil, ddd.NewBusinessErrorWithDetails("TOKEN_GENERATION_FAILED", "failed to generate tokens", err.Error())
	}

	return &AuthenticateResult{
		UserID:       userID,
		Username:     u.Username().Value(),
		Email:        u.Email().Value(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// Logout 用户登出
func (s *AuthenticationService) Logout(ctx context.Context, userID UserID) error {
	// 这里可以实现令牌黑名单等逻辑
	return nil
}

// ChangePassword 修改密码
func (s *AuthenticationService) ChangePassword(ctx context.Context, userID UserID, oldPassword, newPassword string, ipAddress string) error {
	// 1. 查找用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ddd.ErrAggregateNotFound
	}

	// 2. 验证新密码强度
	if err := s.passwordPolicy.Validate(newPassword); err != nil {
		return ddd.NewBusinessErrorWithDetails("WEAK_PASSWORD", "password does not meet policy requirements", err.Error())
	}

	// 3. 修改密码
	if err := u.ChangePassword(oldPassword, newPassword, ipAddress); err != nil {
		return err
	}

	// 4. 保存用户
	return s.userRepo.Save(ctx, u)
}

// ResetPassword 重置密码（管理员操作）
func (s *AuthenticationService) ResetPassword(ctx context.Context, userID UserID, newPassword string, ipAddress string) error {
	// 1. 查找用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ddd.ErrAggregateNotFound
	}

	// 2. 验证新密码强度
	if err := s.passwordPolicy.Validate(newPassword); err != nil {
		return ddd.NewBusinessErrorWithDetails("WEAK_PASSWORD", "password does not meet policy requirements", err.Error())
	}

	// 3. 直接设置新密码（跳过旧密码验证）
	// 这里使用一个特殊标记来跳过旧密码验证
	hashedPassword := NewHashedPassword(newPassword)
	u.password = hashedPassword
	u.updatedAt = time.Now()
	u.IncrementVersion()

	// 发布密码修改事件
	event := NewUserPasswordChangedEvent(userID, ipAddress)
	u.ApplyEvent(event)

	// 4. 保存用户
	return s.userRepo.Save(ctx, u)
}

// TokenService 令牌服务接口
type TokenService interface {
	// GenerateTokenPair 生成令牌对
	GenerateTokenPair(userID UserID) (*TokenPair, error)
	// ParseAccessToken 解析访问令牌
	ParseAccessToken(token string) (*TokenClaims, error)
	// ParseRefreshToken 解析刷新令牌
	ParseRefreshToken(token string) (*TokenClaims, error)
	// ValidateToken 验证令牌
	ValidateToken(token string) (*TokenClaims, error)
}

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
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// PasswordPolicyService 密码策略服务接口
type PasswordPolicyService interface {
	// Validate 验证密码是否符合策略
	Validate(password string) error
	// GetPolicy 获取密码策略
	GetPolicy() PasswordPolicy
}

// PasswordPolicy 密码策略
type PasswordPolicy struct {
	MinLength           int
	MaxLength           int
	RequireUppercase    bool
	RequireLowercase    bool
	RequireDigit        bool
	RequireSpecialChar  bool
	SpecialChars        string
	DisallowUsername    bool
	MaxRepeatedChars    int
	PasswordHistorySize int
}

// DefaultPasswordPolicyService 默认密码策略服务实现
type DefaultPasswordPolicyService struct {
	policy PasswordPolicy
}

// NewDefaultPasswordPolicyService 创建默认密码策略服务
func NewDefaultPasswordPolicyService() PasswordPolicyService {
	return &DefaultPasswordPolicyService{
		policy: PasswordPolicy{
			MinLength:           8,
			MaxLength:           128,
			RequireUppercase:    true,
			RequireLowercase:    true,
			RequireDigit:        true,
			RequireSpecialChar:  true,
			SpecialChars:        "!@#$%^&*()_+-=[]{}|;:,.<>?",
			DisallowUsername:    true,
			MaxRepeatedChars:    3,
			PasswordHistorySize: 5,
		},
	}
}

// Validate 验证密码
func (s *DefaultPasswordPolicyService) Validate(password string) error {
	if len(password) < s.policy.MinLength {
		return errors.New("password is too short")
	}

	if len(password) > s.policy.MaxLength {
		return errors.New("password is too long")
	}

	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUppercase = true
		}
		if char >= 'a' && char <= 'z' {
			hasLowercase = true
		}
		if char >= '0' && char <= '9' {
			hasDigit = true
		}
		for _, special := range s.policy.SpecialChars {
			if char == special {
				hasSpecial = true
				break
			}
		}
	}

	if s.policy.RequireUppercase && !hasUppercase {
		return errors.New("password must contain at least one uppercase letter")
	}

	if s.policy.RequireLowercase && !hasLowercase {
		return errors.New("password must contain at least one lowercase letter")
	}

	if s.policy.RequireDigit && !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	if s.policy.RequireSpecialChar && !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	// 检查连续重复字符
	if s.policy.MaxRepeatedChars > 0 {
		maxRepeat := 1
		currentRepeat := 1
		for i := 1; i < len(password); i++ {
			if password[i] == password[i-1] {
				currentRepeat++
				if currentRepeat > maxRepeat {
					maxRepeat = currentRepeat
				}
			} else {
				currentRepeat = 1
			}
		}
		if maxRepeat > s.policy.MaxRepeatedChars {
			return errors.New("password contains too many repeated characters")
		}
	}

	return nil
}

// GetPolicy 获取密码策略
func (s *DefaultPasswordPolicyService) GetPolicy() PasswordPolicy {
	return s.policy
}

// PasswordHasher 密码哈希接口
type PasswordHasher interface {
	// Hash 哈希密码
	Hash(password string) (string, error)
	// Verify 验证密码
	Verify(password, hash string) bool
}

// SimplePasswordHasher 简单密码哈希实现（仅用于演示）
type SimplePasswordHasher struct{}

// Hash 哈希密码
func (h *SimplePasswordHasher) Hash(password string) (string, error) {
	// 实际应用中应该使用 bcrypt 等安全哈希算法
	return password, nil
}

// Verify 验证密码
func (h *SimplePasswordHasher) Verify(password, hash string) bool {
	// 实际应用中应该使用 bcrypt 等安全哈希算法
	return password == hash
}
