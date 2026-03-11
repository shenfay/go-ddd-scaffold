package user

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// User 用户聚合根
type User struct {
	ddd.BaseEntity

	username       *UserName
	email          *Email
	password       *HashedPassword
	status         UserStatus
	displayName    string
	firstName      string
	lastName       string
	gender         UserGender
	phoneNumber    string
	avatarURL      string
	lastLoginAt    *time.Time
	loginCount     int
	lockedUntil    *time.Time
	failedAttempts int
	createdAt      time.Time
	updatedAt      time.Time
}

// NewUser 创建新用户
func NewUser(username, email, password string) (*User, error) {
	user := &User{
		status:         UserStatusPending,
		gender:         UserGenderUnknown,
		loginCount:     0,
		failedAttempts: 0,
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
	}

	// 设置初始ID（在实际应用中应该由ID生成器分配）
	user.SetID(NewUserID(1))

	// 验证和设置用户名
	un, err := NewUserName(username)
	if err != nil {
		return nil, err
	}
	user.username = un

	// 验证和设置邮箱
	em, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	user.email = em

	// 设置密码（这里应该进行加密处理）
	user.password = NewHashedPassword(password) // 实际应用中应该使用bcrypt等加密

	return user, nil
}

// Username 获取用户名
func (u *User) Username() *UserName {
	return u.username
}

// Email 获取邮箱
func (u *User) Email() *Email {
	return u.email
}

// Password 获取密码
func (u *User) Password() *HashedPassword {
	return u.password
}

// Status 获取用户状态
func (u *User) Status() UserStatus {
	return u.status
}

// DisplayName 获取显示名称
func (u *User) DisplayName() string {
	return u.displayName
}

// FirstName 获取名字
func (u *User) FirstName() string {
	return u.firstName
}

// LastName 获取姓氏
func (u *User) LastName() string {
	return u.lastName
}

// Gender 获取性别
func (u *User) Gender() UserGender {
	return u.gender
}

// PhoneNumber 获取电话号码
func (u *User) PhoneNumber() string {
	return u.phoneNumber
}

// AvatarURL 获取头像URL
func (u *User) AvatarURL() string {
	return u.avatarURL
}

// LastLoginAt 获取最后登录时间
func (u *User) LastLoginAt() *time.Time {
	return u.lastLoginAt
}

// LoginCount 获取登录次数
func (u *User) LoginCount() int {
	return u.loginCount
}

// LockedUntil 获取锁定截止时间
func (u *User) LockedUntil() *time.Time {
	return u.lockedUntil
}

// FailedAttempts 获取失败尝试次数
func (u *User) FailedAttempts() int {
	return u.failedAttempts
}

// SetDisplayName 设置显示名称
func (u *User) SetDisplayName(displayName string) {
	u.displayName = displayName
	u.updatedAt = time.Now()
}

// SetFirstName 设置名字
func (u *User) SetFirstName(firstName string) {
	u.firstName = firstName
	u.updatedAt = time.Now()
}

// SetLastName 设置姓氏
func (u *User) SetLastName(lastName string) {
	u.lastName = lastName
	u.updatedAt = time.Now()
}

// SetGender 设置性别
func (u *User) SetGender(gender UserGender) {
	u.gender = gender
	u.updatedAt = time.Now()
}

// SetPhoneNumber 设置电话号码
func (u *User) SetPhoneNumber(phoneNumber string) {
	u.phoneNumber = phoneNumber
	u.updatedAt = time.Now()
}

// SetAvatarURL 设置头像URL
func (u *User) SetAvatarURL(avatarURL string) {
	u.avatarURL = avatarURL
	u.updatedAt = time.Now()
}

// Activate 激活用户
func (u *User) Activate() error {
	if u.status != UserStatusPending {
		return ddd.NewBusinessError("USER_NOT_PENDING", "user is not in pending status")
	}

	u.status = UserStatusActive
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// Deactivate 禁用用户
func (u *User) Deactivate() error {
	if u.status == UserStatusInactive {
		return ddd.NewBusinessError("USER_ALREADY_INACTIVE", "user is already inactive")
	}

	u.status = UserStatusInactive
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// Lock 锁定用户
func (u *User) Lock(duration time.Duration) error {
	if u.status == UserStatusLocked {
		return ddd.NewBusinessError("USER_ALREADY_LOCKED", "user is already locked")
	}

	lockUntil := time.Now().Add(duration)
	u.status = UserStatusLocked
	u.lockedUntil = &lockUntil
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// Unlock 解锁用户
func (u *User) Unlock() error {
	if u.status != UserStatusLocked {
		return ddd.NewBusinessError("USER_NOT_LOCKED", "user is not locked")
	}

	u.status = UserStatusActive
	u.lockedUntil = nil
	u.failedAttempts = 0
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// RecordLogin 记录登录
func (u *User) RecordLogin() {
	now := time.Now()
	u.lastLoginAt = &now
	u.loginCount++
	u.failedAttempts = 0
	u.updatedAt = now
	u.IncrementVersion()
}

// RecordFailedLogin 记录失败登录
func (u *User) RecordFailedLogin() {
	u.failedAttempts++
	u.updatedAt = time.Now()
	u.IncrementVersion()
}

// ResetFailedAttempts 重置失败尝试次数
func (u *User) ResetFailedAttempts() {
	u.failedAttempts = 0
	u.updatedAt = time.Now()
	u.IncrementVersion()
}

// ChangePassword 修改密码
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if !u.password.Matches(oldPassword) {
		return ddd.NewBusinessError("INVALID_OLD_PASSWORD", "invalid old password")
	}

	// 这里应该验证新密码强度
	u.password = NewHashedPassword(newPassword) // 实际应用中应该加密
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(newEmail string) error {
	email, err := NewEmail(newEmail)
	if err != nil {
		return err
	}

	u.email = email
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// IsLocked 检查用户是否被锁定
func (u *User) IsLocked() bool {
	if u.status != UserStatusLocked {
		return false
	}

	if u.lockedUntil != nil && time.Now().After(*u.lockedUntil) {
		// 锁定时间已过，自动解锁
		u.Unlock()
		return false
	}

	return true
}

// CanLogin 检查用户是否可以登录
func (u *User) CanLogin() bool {
	return u.status == UserStatusActive && !u.IsLocked()
}

// FullName 获取完整姓名
func (u *User) FullName() string {
	if u.firstName != "" && u.lastName != "" {
		return u.firstName + " " + u.lastName
	}
	if u.firstName != "" {
		return u.firstName
	}
	if u.lastName != "" {
		return u.lastName
	}
	return ""
}

// GetFullName 获取完整姓名（带默认值）
func (u *User) GetFullName(defaultName string) string {
	fullName := u.FullName()
	if fullName == "" {
		return defaultName
	}
	return fullName
}
