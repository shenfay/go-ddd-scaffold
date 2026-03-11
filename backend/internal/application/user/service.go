package user

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// UserService 用户应用服务接口
type UserService interface {
	// 用户生命周期管理
	RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*user.User, error)
	ActivateUser(ctx context.Context, cmd *ActivateUserCommand) error
	DeactivateUser(ctx context.Context, cmd *DeactivateUserCommand) error

	// 认证相关
	AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticationResult, error)
	ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error
	UpdateEmail(ctx context.Context, cmd *UpdateEmailCommand) error

	// 用户资料管理
	UpdateUserProfile(ctx context.Context, cmd *UpdateUserProfileCommand) error
	UpdateUserAvatar(ctx context.Context, cmd *UpdateUserAvatarCommand) error

	// 安全相关
	LockUser(ctx context.Context, cmd *LockUserCommand) error
	UnlockUser(ctx context.Context, cmd *UnlockUserCommand) error
	ResetUserPassword(ctx context.Context, cmd *ResetUserPasswordCommand) error

	// 查询服务
	GetUserByID(ctx context.Context, userID user.UserID) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)

	// 批量操作
	BatchActivateUsers(ctx context.Context, userIDs []user.UserID) error
	BatchDeactivateUsers(ctx context.Context, userIDs []user.UserID) error
}

// RegisterUserCommand 用户注册命令
type RegisterUserCommand struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// ActivateUserCommand 用户激活命令
type ActivateUserCommand struct {
	UserID user.UserID `json:"user_id" validate:"required"`
}

// DeactivateUserCommand 用户禁用命令
type DeactivateUserCommand struct {
	UserID user.UserID `json:"user_id" validate:"required"`
	Reason string      `json:"reason,omitempty"`
}

// AuthenticateUserCommand 用户认证命令
type AuthenticateUserCommand struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID      user.UserID `json:"user_id" validate:"required"`
	OldPassword string      `json:"old_password" validate:"required"`
	NewPassword string      `json:"new_password" validate:"required,min=8"`
	IPAddress   string      `json:"ip_address,omitempty"`
}

// UpdateEmailCommand 更新邮箱命令
type UpdateEmailCommand struct {
	UserID    user.UserID `json:"user_id" validate:"required"`
	NewEmail  string      `json:"new_email" validate:"required,email"`
	IPAddress string      `json:"ip_address,omitempty"`
}

// UpdateUserProfileCommand 更新用户资料命令
type UpdateUserProfileCommand struct {
	UserID      user.UserID      `json:"user_id" validate:"required"`
	DisplayName *string          `json:"display_name,omitempty"`
	FirstName   *string          `json:"first_name,omitempty"`
	LastName    *string          `json:"last_name,omitempty"`
	Gender      *user.UserGender `json:"gender,omitempty"`
	PhoneNumber *string          `json:"phone_number,omitempty"`
}

// UpdateUserAvatarCommand 更新用户头像命令
type UpdateUserAvatarCommand struct {
	UserID    user.UserID `json:"user_id" validate:"required"`
	AvatarURL string      `json:"avatar_url" validate:"required,url"`
}

// LockUserCommand 锁定用户命令
type LockUserCommand struct {
	UserID   user.UserID     `json:"user_id" validate:"required"`
	Duration time.Duration   `json:"duration" validate:"required"`
	Reason   string          `json:"reason" validate:"required"`
}

// UnlockUserCommand 解锁用户命令
type UnlockUserCommand struct {
	UserID user.UserID `json:"user_id" validate:"required"`
}

// ResetUserPasswordCommand 重置用户密码命令
type ResetUserPasswordCommand struct {
	UserID        user.UserID `json:"user_id" validate:"required"`
	NewPassword   string      `json:"new_password" validate:"required,min=8"`
	Administrator string      `json:"administrator" validate:"required"`
	Reason        string      `json:"reason" validate:"required"`
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
