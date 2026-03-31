package usecase

import (
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// ============================================================================
// User UseCase Commands (用户用例命令对象)
// ============================================================================

// RegisterUserCommand 用户注册命令
type RegisterUserCommand struct {
	Username string
	Email    string
	Password string
}

// LoginUserCommand 用户登录命令
type LoginUserCommand struct {
	Username  string
	Password  string
	IPAddress string
	UserAgent string
}

// UpdateProfileCommand 更新用户资料命令
type UpdateProfileCommand struct {
	UserID      vo.UserID
	DisplayName *string
	FirstName   *string
	LastName    *string
	Gender      *vo.UserGender // 指针类型支持可选参数
	PhoneNumber *string
}

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID      vo.UserID
	OldPassword string
	NewPassword string
	IPAddress   string
}