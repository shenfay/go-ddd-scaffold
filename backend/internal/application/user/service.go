package user

import (
	"context"
	"errors"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// UserService 用户应用服务接口（核心流程）
type UserService interface {
	// === 核心流程 ===
	RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*user.User, error)
	AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticationResult, error)
	GetUserByID(ctx context.Context, userID user.UserID) (*user.User, error)

	// === 辅助功能（可选）===
	UpdateUserProfile(ctx context.Context, cmd *UpdateUserProfileCommand) error
	ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error
}

// UserServiceImpl 用户应用服务实现
type UserServiceImpl struct {
	userRepo       user.UserRepository
	eventPublisher ddd.EventPublisher
	passwordHasher user.PasswordHasher
	tokenService   auth.TokenService
}

// NewUserService 创建用户应用服务
func NewUserService(
	userRepo user.UserRepository,
	eventPublisher ddd.EventPublisher,
	passwordHasher user.PasswordHasher,
	tokenService auth.TokenService,
) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

// ============================================================================
// Commands & DTOs
// ============================================================================

// RegisterUserCommand 用户注册命令
type RegisterUserCommand struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// AuthenticateUserCommand 用户认证命令
type AuthenticateUserCommand struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// UpdateUserProfileCommand 更新用户资料命令（可选）
type UpdateUserProfileCommand struct {
	UserID      user.UserID      `json:"user_id" validate:"required"`
	DisplayName *string          `json:"display_name,omitempty"`
	FirstName   *string          `json:"first_name,omitempty"`
	LastName    *string          `json:"last_name,omitempty"`
	Gender      *user.UserGender `json:"gender,omitempty"`
	PhoneNumber *string          `json:"phone_number,omitempty"`
}

// ChangePasswordCommand 修改密码命令（可选）
type ChangePasswordCommand struct {
	UserID      user.UserID `json:"user_id" validate:"required"`
	OldPassword string      `json:"old_password" validate:"required"`
	NewPassword string      `json:"new_password" validate:"required,min=8"`
	IPAddress   string      `json:"ip_address,omitempty"`
}

// AuthenticationResult 认证结果
type AuthenticationResult struct {
	UserID       user.UserID `json:"user_id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    time.Time   `json:"expires_at"`
}

// ============================================================================
// Service Methods - 核心流程实现
// ============================================================================

// GetUserByID 根据 ID 获取用户
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID user.UserID) (*user.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// AuthenticateUser 认证用户
func (s *UserServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticationResult, error) {
	// 1. 查找用户
	u, err := s.userRepo.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, ddd.ErrAggregateNotFound
	}

	// 2. 验证密码
	if !s.passwordHasher.Verify(cmd.Password, u.Password().Value()) {
		// 记录失败登录（可选）
		u.RecordFailedLogin(cmd.IPAddress, cmd.UserAgent, "invalid_password")
		_ = s.userRepo.Save(ctx, u)
		return nil, ddd.NewBusinessError("INVALID_PASSWORD", "密码错误")
	}

	// 3. 检查用户是否可以登录
	if !u.CanLogin() {
		return nil, ddd.NewBusinessError("USER_CANNOT_LOGIN", "用户无法登录")
	}

	// 4. 生成 JWT Token
	tokenPair, err := s.tokenService.GenerateTokenPair(
		u.ID().(user.UserID).Int64(),
		u.Username().Value(),
		u.Email().Value(),
	)
	if err != nil {
		return nil, err
	}

	// 5. 记录成功登录
	u.RecordLogin(cmd.IPAddress, cmd.UserAgent)
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	// 6. 发布领域事件
	events := u.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	u.ClearUncommittedEvents()

	return &AuthenticationResult{
		UserID:       u.ID().(user.UserID),
		Username:     u.Username().Value(),
		Email:        u.Email().Value(),
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// UpdateUserProfile 更新用户资料
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, cmd *UpdateUserProfileCommand) error {
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return ddd.ErrAggregateNotFound
	}

	// 更新用户信息
	if cmd.DisplayName != nil {
		u.SetDisplayName(*cmd.DisplayName)
	}
	if cmd.FirstName != nil {
		u.SetFirstName(*cmd.FirstName)
	}
	if cmd.LastName != nil {
		u.SetLastName(*cmd.LastName)
	}
	if cmd.Gender != nil {
		u.SetGender(*cmd.Gender)
	}
	if cmd.PhoneNumber != nil {
		u.SetPhoneNumber(*cmd.PhoneNumber)
	}

	// 保存用户（会自动发布事件）
	return s.userRepo.Save(ctx, u)
}

// ChangePassword 修改密码
func (s *UserServiceImpl) ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error {
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return ddd.ErrAggregateNotFound
	}

	// 修改密码（内部会验证旧密码）
	if err := u.ChangePassword(cmd.OldPassword, cmd.NewPassword, cmd.IPAddress); err != nil {
		return err
	}

	// 保存用户（会自动发布事件）
	return s.userRepo.Save(ctx, u)
}

// RegisterUser 注册用户
func (s *UserServiceImpl) RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*user.User, error) {
	// 1. 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, cmd.Username)
	if err != nil && !errors.Is(err, ddd.ErrAggregateNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ddd.NewBusinessError("USERNAME_EXISTS", "用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	existingUser, err = s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, ddd.ErrAggregateNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ddd.NewBusinessError("EMAIL_EXISTS", "邮箱已被注册")
	}

	// 3. 哈希密码
	hashedPassword, err := s.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	// 4. 创建用户实体（需要 ID 生成器，这里使用简单实现）
	newUser, err := user.NewUser(cmd.Username, cmd.Email, hashedPassword, func() int64 {
		return time.Now().UnixNano() // 临时实现，实际应该使用 Snowflake
	})
	if err != nil {
		return nil, err
	}

	// 5. 保存用户
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	// 6. 发布领域事件
	events := newUser.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	newUser.ClearUncommittedEvents()

	// 7. 生成 JWT 令牌对（注册后自动登录）
	_, err = s.tokenService.GenerateTokenPair(
		newUser.ID().(user.UserID).Int64(),
		newUser.Username().Value(),
		newUser.Email().Value(),
	)
	if err != nil {
		// 令牌生成失败不影响注册流程，仅记录日志
		// TODO: 添加日志记录
	}

	// TODO: 将令牌信息附加到返回结果中
	// 目前返回用户对象，令牌信息可以通过 AuthenticationResult 获取

	return newUser, nil
}
