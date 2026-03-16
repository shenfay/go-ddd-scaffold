package user

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// ============================================================================
// Input DTOs (Commands)
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

// UpdateUserProfileCommand 更新用户资料命令
type UpdateUserProfileCommand struct {
	UserID      user.UserID      `json:"user_id" validate:"required"`
	DisplayName *string          `json:"display_name,omitempty"`
	FirstName   *string          `json:"first_name,omitempty"`
	LastName    *string          `json:"last_name,omitempty"`
	Gender      *user.UserGender `json:"gender,omitempty"`
	PhoneNumber *string          `json:"phone_number,omitempty"`
}

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID      user.UserID `json:"user_id" validate:"required"`
	OldPassword string      `json:"old_password" validate:"required"`
	NewPassword string      `json:"new_password" validate:"required,min=8"`
	IPAddress   string      `json:"ip_address,omitempty"`
}

// ============================================================================
// Output DTOs (Results)
// ============================================================================

// RegisterUserResult 用户注册结果
type RegisterUserResult struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// AuthenticateUserResult 认证结果
type AuthenticateUserResult struct {
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GetUserResult 获取用户结果
type GetUserResult struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Gender      string    `json:"gender"`
	PhoneNumber string    `json:"phone_number"`
	AvatarURL   string    `json:"avatar_url"`
	Status      int32     `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ============================================================================
// Auxiliary DTOs
// ============================================================================

// UserProfileUpdate 用户资料更新数据（用于内部传输）
type UserProfileUpdate struct {
	DisplayName *string
	FirstName   *string
	LastName    *string
	Gender      *user.UserGender
	PhoneNumber *string
	AvatarURL   *string
}
