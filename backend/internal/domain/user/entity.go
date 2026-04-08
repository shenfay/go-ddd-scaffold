package user

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// User 用户聚合根
// 封装用户身份、凭据和生命周期状态
type User struct {
	ID             string     `json:"id"`                       // 用户唯一标识
	Email          string     `json:"email"`                    // 用户邮箱
	Password       string     `json:"-"`                        // 密码哈希（不序列化）
	EmailVerified  bool       `json:"email_verified"`           // 邮箱验证状态
	Locked         bool       `json:"locked"`                   // 账户锁定状态
	FailedAttempts int        `json:"failed_attempts"`          // 连续登录失败次数
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`  // 最后登录时间
	CreatedAt      time.Time  `json:"created_at"`               // 创建时间
	UpdatedAt      time.Time  `json:"updated_at"`               // 更新时间
}

// NewUser 创建新用户
// 邮箱格式无效或密码不符合要求时返回错误
func NewUser(email, password string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := utils.Now()
	return &User{
		ID:             utils.GenerateID(),
		Email:          email,
		Password:       hashedPassword,
		EmailVerified:  false,
		Locked:         false,
		FailedAttempts: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// VerifyPassword 验证密码是否匹配
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsLocked 检查账户是否因登录失败次数过多而被锁定
func (u *User) IsLocked() bool {
	return u.Locked
}

// IncrementFailedAttempts 增加失败尝试次数
func (u *User) IncrementFailedAttempts(maxAttempts int) {
	u.FailedAttempts++
	u.UpdatedAt = utils.Now()

	if u.FailedAttempts >= maxAttempts {
		u.Locked = true
	}
}

// ResetFailedAttempts 重置失败尝试次数
func (u *User) ResetFailedAttempts() {
	u.FailedAttempts = 0
	u.UpdatedAt = utils.Now()
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := utils.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// VerifyEmail 验证邮箱
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = utils.Now()
}

// ChangePassword 修改密码
func (u *User) ChangePassword(newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	u.UpdatedAt = utils.Now()
	return nil
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(newEmail string) error {
	u.Email = newEmail
	u.UpdatedAt = utils.Now()
	return nil
}
