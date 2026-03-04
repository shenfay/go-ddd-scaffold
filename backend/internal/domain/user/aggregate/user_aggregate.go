package aggregate

import (
	"time"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserAggregate 用户聚合根
type UserAggregate struct {
	*entity.User
}

// NewUserAggregate 创建用户聚合根
func NewUserAggregate(user *entity.User) *UserAggregate {
	return &UserAggregate{User: user}
}

// RegisterUser 注册新用户（工厂方法）
func RegisterUser(email valueobject.Email, password valueobject.Password, nickname valueobject.Nickname, role entity.UserRole) (*UserAggregate, error) {
	if !email.IsValid() {
		return nil, ErrInvalidEmail
	}

	if !password.IsValid() {
		return nil, ErrInvalidPassword
	}

	if !nickname.IsValid() {
		return nil, ErrInvalidNickname
	}

	user := &entity.User{
		ID:       uuid.New(),
		Email:    email.String(),
		Password: password.String(),
		Nickname: nickname.String(),
		Role:     role,
		Status:   entity.StatusActive,
	}

	return NewUserAggregate(user), nil
}

// UpdateProfile 更新用户资料
func (ua *UserAggregate) UpdateProfile(nickname *valueobject.Nickname, avatar *string) error {
	if nickname != nil {
		if !nickname.IsValid() {
			return ErrInvalidNickname
		}
		ua.Nickname = nickname.String()
	}

	if avatar != nil {
		ua.Avatar = avatar
	}

	ua.UpdatedAt = time.Now()
	return nil
}

// ChangePassword 修改密码
func (ua *UserAggregate) ChangePassword(oldPassword, newPassword valueobject.Password) error {
	if !newPassword.IsValid() {
		return ErrInvalidPassword
	}

	// 这里应该验证旧密码，但在聚合根中不直接处理密码验证逻辑
	// 密码验证应该在应用层或领域服务中处理

	ua.Password = newPassword.String()
	ua.UpdatedAt = time.Now()
	return nil
}

// Lock 锁定用户
func (ua *UserAggregate) Lock() {
	ua.Status = entity.StatusLocked
	ua.UpdatedAt = time.Now()
}

// Activate 激活用户
func (ua *UserAggregate) Activate() {
	ua.Status = entity.StatusActive
	ua.UpdatedAt = time.Now()
}

// 领域错误定义
var (
	ErrInvalidEmail    = &DomainError{"invalid_email", "邮箱格式不正确"}
	ErrInvalidPassword = &DomainError{"invalid_password", "密码不符合要求"}
	ErrInvalidNickname = &DomainError{"invalid_nickname", "昵称格式不正确"}
)

// DomainError 领域错误
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}
