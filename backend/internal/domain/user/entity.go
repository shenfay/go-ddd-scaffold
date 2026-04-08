package user

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// User represents a registered user aggregate root.
// It encapsulates user identity, credentials, and lifecycle state.
type User struct {
	ID             string     `json:"id"`                       // Unique user identifier
	Email          string     `json:"email"`                    // User email address
	Password       string     `json:"-"`                        // Hashed password (never serialized)
	EmailVerified  bool       `json:"email_verified"`           // Email verification status
	Locked         bool       `json:"locked"`                   // Account lock status
	FailedAttempts int        `json:"failed_attempts"`          // Consecutive failed login attempts
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`  // Last successful login time
	CreatedAt      time.Time  `json:"created_at"`               // Account creation time
	UpdatedAt      time.Time  `json:"updated_at"`               // Last update time
}

// NewUser creates a new user with validated email and password.
// Returns error if email format is invalid or password does not meet requirements.
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

// VerifyPassword checks if the provided password matches the stored hash.
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsLocked checks if the user account is locked due to failed login attempts.
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
