package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户聚合根
type User struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	Password       string     `json:"-"` // 不序列化到 JSON
	EmailVerified  bool       `json:"email_verified"`
	Locked         bool       `json:"locked"`
	FailedAttempts int        `json:"failed_attempts"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewUser 创建新用户
func NewUser(email, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:             generateUserID(),
		Email:          email,
		Password:       hashedPassword,
		EmailVerified:  false,
		Locked:         false,
		FailedAttempts: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsLocked 检查账户是否被锁定
func (u *User) IsLocked() bool {
	return u.Locked
}

// IncrementFailedAttempts 增加失败尝试次数
func (u *User) IncrementFailedAttempts(maxAttempts int) {
	u.FailedAttempts++
	u.UpdatedAt = time.Now()

	if u.FailedAttempts >= maxAttempts {
		u.Locked = true
	}
}

// ResetFailedAttempts 重置失败尝试次数
func (u *User) ResetFailedAttempts() {
	u.FailedAttempts = 0
	u.UpdatedAt = time.Now()
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// VerifyEmail 验证邮箱
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}

// ChangePassword 修改密码
func (u *User) ChangePassword(newPassword string) error {
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// hashPassword 哈希密码
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// generateUserID 生成用户 ID（实际会注入）
var generateUserID = func() string {
	// 这里会在 utils/ulid 中实现
	// 暂时使用简单实现
	return "user_" + time.Now().Format("20060102150405")
}
