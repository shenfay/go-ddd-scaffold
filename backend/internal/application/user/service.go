package user

import (
	"context"
	"errors"
	"time"

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
}

// NewUserService 创建用户应用服务
func NewUserService(
	userRepo user.UserRepository,
	eventPublisher ddd.EventPublisher,
	passwordHasher user.PasswordHasher,
) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
		passwordHasher: passwordHasher,
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
	u, err := s.userRepo.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, ddd.ErrAggregateNotFound
	}

	// 验证密码
	if !s.passwordHasher.Verify(cmd.Password, u.Password().Value()) {
		return nil, ddd.NewBusinessError("INVALID_PASSWORD", "密码错误")
	}

	// TODO: 生成 JWT Token
	// token, refreshToken, expiresAt, err := s.jwtService.GenerateToken(u)
	// if err != nil {
	//     return nil, err
	// }

	return &AuthenticationResult{
		UserID:   u.ID().(user.UserID),
		Username: u.Username().Value(),
		Email:    u.Email().Value(),
		// Token:        token,
		// RefreshToken: refreshToken,
		// ExpiresAt:    expiresAt,
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
	// TODO: 实现用户注册逻辑
	return nil, errors.New("REGISTER_USER_NOT_IMPLEMENTED")
}
