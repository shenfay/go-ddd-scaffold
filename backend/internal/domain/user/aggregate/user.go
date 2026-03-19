package aggregate

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// User 用户聚合根
type User struct {
	kernel.BaseEntity

	username       *valueobject.UserName
	email          *valueobject.Email
	password       *valueobject.HashedPassword
	status         valueobject.UserStatus
	displayName    string
	firstName      string
	lastName       string
	gender         valueobject.UserGender
	phoneNumber    string
	avatarURL      string
	lastLoginAt    *time.Time
	loginCount     int
	lockedUntil    *time.Time
	failedAttempts int
	createdAt      time.Time
	updatedAt      time.Time
}

// NewUser 使用已哈希的密码创建新用户
func NewUser(username, email, hashedPassword string, idGenerator func() int64) (*User, error) {
	user := &User{
		status:         valueobject.UserStatusActive,
		gender:         valueobject.UserGenderUnknown,
		loginCount:     0,
		failedAttempts: 0,
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
	}

	// 使用 ID 生成器生成唯一 ID
	userID := idGenerator()
	user.SetID(valueobject.NewUserID(userID))

	// 验证和设置用户名
	un, err := valueobject.NewUserName(username)
	if err != nil {
		return nil, err
	}
	user.username = un

	// 验证和设置邮箱
	em, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}
	user.email = em

	// 设置已哈希的密码
	user.password = valueobject.NewHashedPassword(hashedPassword)

	// 发布用户注册事件（使用默认值）
	event := user.newRegisteredEvent(username, email, user.status.String(), username, "", 0)
	user.ApplyEvent(event)

	return user, nil
}

// Username 获取用户名
func (u *User) Username() *valueobject.UserName {
	return u.username
}

// Email 获取邮箱
func (u *User) Email() *valueobject.Email {
	return u.email
}

// Password 获取密码
func (u *User) Password() *valueobject.HashedPassword {
	return u.password
}

// Status 获取用户状态
func (u *User) Status() valueobject.UserStatus {
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
func (u *User) Gender() valueobject.UserGender {
	return u.gender
}

// PhoneNumber 获取电话号码
func (u *User) PhoneNumber() string {
	return u.phoneNumber
}

// AvatarURL 获取头像 URL
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
func (u *User) SetGender(gender valueobject.UserGender) {
	u.gender = gender
	u.updatedAt = time.Now()
}

// SetPhoneNumber 设置电话号码
func (u *User) SetPhoneNumber(phoneNumber string) {
	u.phoneNumber = phoneNumber
	u.updatedAt = time.Now()
}

// SetAvatarURL 设置头像 URL
func (u *User) SetAvatarURL(avatarURL string) {
	u.avatarURL = avatarURL
	u.updatedAt = time.Now()
}

// Activate 激活用户
func (u *User) Activate() error {
	if u.status != valueobject.UserStatusPending {
		return kernel.NewBusinessError(kernel.CodeUserNotPending, "user is not in pending status")
	}

	u.status = valueobject.UserStatusActive
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// Deactivate 禁用用户
func (u *User) Deactivate(reason string) error {
	if u.status == valueobject.UserStatusInactive {
		return kernel.NewBusinessError(kernel.CodeUserAlreadyInactive, "user is already inactive")
	}

	u.status = valueobject.UserStatusInactive
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// Lock 锁定用户
func (u *User) Lock(duration time.Duration, reason string) error {
	if u.status == valueobject.UserStatusLocked {
		return kernel.NewBusinessError(kernel.CodeUserAlreadyLocked, "user is already locked")
	}

	lockUntil := time.Now().Add(duration)
	u.status = valueobject.UserStatusLocked
	u.lockedUntil = &lockUntil
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// Unlock 解锁用户
func (u *User) Unlock() error {
	if u.status != valueobject.UserStatusLocked {
		return kernel.NewBusinessError(kernel.CodeUserNotLocked, "user is not locked")
	}

	u.status = valueobject.UserStatusActive
	u.lockedUntil = nil
	u.failedAttempts = 0
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// RecordLogin 记录登录
func (u *User) RecordLogin(ipAddress, userAgent string) {
	now := time.Now()
	u.lastLoginAt = &now
	u.loginCount++
	u.failedAttempts = 0
	u.updatedAt = now
	u.IncrementVersion()

	// 发布用户登录事件（使用默认值）
	event := u.newLoggedInEvent(ipAddress, userAgent, "", "", "", "password", true)
	u.ApplyEvent(event)
}

// RecordFailedLogin 记录失败登录
func (u *User) RecordFailedLogin(ipAddress, userAgent, reason string) {
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
func (u *User) ChangePassword(newPassword string, ipAddress string) error {
	// TODO: 这里应该验证新密码强度并加密
	u.password = valueobject.NewHashedPassword(newPassword)
	u.updatedAt = time.Now()
	u.IncrementVersion()

	return nil
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(newEmail string) error {
	email, err := valueobject.NewEmail(newEmail)
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
	if u.status != valueobject.UserStatusLocked {
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
	return u.status == valueobject.UserStatusActive && !u.IsLocked()
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

// newRegisteredEvent 创建用户注册事件
func (u *User) newRegisteredEvent(username, email, status, displayName, registrationIP string, tenantID int64) *event.UserRegisteredEvent {
	ev := &event.UserRegisteredEvent{
		BaseEvent:      kernel.NewBaseEvent("UserRegistered", u.ID(), 1),
		UserID:         u.ID().(valueobject.UserID),
		Username:       username,
		Email:          email,
		Status:         status,
		DisplayName:    displayName,
		RegistrationIP: registrationIP,
		TenantID:       tenantID,
		RegisteredAt:   time.Now(),
	}
	ev.SetMetadata("event_type", "domain_event")
	ev.SetMetadata("aggregate_type", "user")
	return ev
}

// newLoggedInEvent 创建用户登录事件
func (u *User) newLoggedInEvent(ipAddress, userAgent, location, deviceType, deviceFingerprint, loginMethod string, success bool) *event.UserLoggedInEvent {
	ev := &event.UserLoggedInEvent{
		BaseEvent:         kernel.NewBaseEvent("UserLoggedIn", u.ID(), 1),
		UserID:            u.ID().(valueobject.UserID),
		LoginAt:           time.Now(),
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		Location:          location,
		DeviceType:        deviceType,
		DeviceFingerprint: deviceFingerprint,
		LoginMethod:       loginMethod,
		Success:           success,
	}
	ev.SetMetadata("event_type", "domain_event")
	ev.SetMetadata("aggregate_type", "user")
	ev.SetMetadata("security_event", true)
	return ev
}
