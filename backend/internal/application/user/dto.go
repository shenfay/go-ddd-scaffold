package user

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// ============================================================================
// User DTOs (用于 API 层和应用层之间的数据传输)
// ============================================================================

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	DisplayName *string        `json:"display_name,omitempty"`
	FirstName   *string        `json:"first_name,omitempty"`
	LastName    *string        `json:"last_name,omitempty"`
	Gender      *vo.UserGender `json:"gender,omitempty"`
	PhoneNumber *string        `json:"phone_number,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
	IPAddress   string `json:"ip_address,omitempty"`
}

// ConvertUserToDTO 将领域对象转换为 DTO
func ConvertUserToDTO(user *aggregate.User) *UserDTO {
	return &UserDTO{
		ID:          user.ID().(vo.UserID).Int64(),
		Username:    user.Username().Value(),
		Email:       user.Email().Value(),
		DisplayName: user.DisplayName(),
		Status:      user.Status().String(),
		CreatedAt:   user.CreatedAt(),
	}
}

// AuthenticateUserResult 认证结果（保持兼容）
type AuthenticateUserResult struct {
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
