package usecase

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user"
)

// ============================================================================
// User UseCase Results (用户用例结果对象)
// ============================================================================

// RegisterUserResult 用户注册结果
type RegisterUserResult struct {
	UserID   int64
	Username string
	Email    string
}

// LoginUserResult 用户登录结果
type LoginUserResult struct {
	UserID       int64
	Username     string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// GetUserResult 获取用户结果
type GetUserResult struct {
	UserDTO *user.UserDTO
}

// UpdateProfileResult 更新用户资料结果
type UpdateProfileResult struct {
	Success bool
}

// ChangePasswordResult 修改密码结果
type ChangePasswordResult struct {
	Success bool
}
