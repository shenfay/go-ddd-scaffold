package user

import "time"

// UserRegistered 用户注册领域事件
type UserRegistered struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID, email string) *UserRegistered {
	return &UserRegistered{
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
	}
}

// GetPayload 获取事件数据
func (e *UserRegistered) GetPayload() interface{} {
	return e
}

// GetType 获取事件类型
func (e *UserRegistered) GetType() string {
	return "USER.REGISTERED"
}

// UserLoggedIn 用户登录领域事件
type UserLoggedIn struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Device    string    `json:"device"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(userID, email, ip, userAgent, device string) *UserLoggedIn {
	return &UserLoggedIn{
		UserID:    userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Device:    device,
		Timestamp: time.Now(),
	}
}

// GetPayload 获取事件数据
func (e *UserLoggedIn) GetPayload() interface{} {
	return e
}

// GetType 获取事件类型
func (e *UserLoggedIn) GetType() string {
	return "AUTH.LOGIN.SUCCESS"
}

// LoginFailed 登录失败事件
type LoginFailed struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	IP        string    `json:"ip"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// NewLoginFailedEvent 创建登录失败事件
func NewLoginFailedEvent(userID, email, ip, reason string) *LoginFailed {
	return &LoginFailed{
		UserID:    userID,
		Email:     email,
		IP:        ip,
		Reason:    reason,
		Timestamp: time.Now(),
	}
}

func (e *LoginFailed) GetPayload() interface{} {
	return e
}

func (e *LoginFailed) GetType() string {
	return "AUTH.LOGIN.FAILED"
}

// AccountLocked 账户锁定事件
type AccountLocked struct {
	UserID         string    `json:"user_id"`
	Email          string    `json:"email"`
	FailedAttempts int       `json:"failed_attempts"`
	LockedUntil    time.Time `json:"locked_until"`
	Timestamp      time.Time `json:"timestamp"`
}

// NewAccountLockedEvent 创建账户锁定事件
func NewAccountLockedEvent(userID, email string, failedAttempts int, lockedUntil time.Time) *AccountLocked {
	return &AccountLocked{
		UserID:         userID,
		Email:          email,
		FailedAttempts: failedAttempts,
		LockedUntil:    lockedUntil,
		Timestamp:      time.Now(),
	}
}

func (e *AccountLocked) GetPayload() interface{} {
	return e
}

func (e *AccountLocked) GetType() string {
	return "SECURITY.ACCOUNT.LOCKED"
}

// UserLoggedOut 用户登出领域事件
type UserLoggedOut struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUserLoggedOutEvent 创建用户登出事件
func NewUserLoggedOutEvent(userID, email string) *UserLoggedOut {
	return &UserLoggedOut{
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
	}
}

func (e *UserLoggedOut) GetPayload() interface{} {
	return e
}

func (e *UserLoggedOut) GetType() string {
	return "AUTH.LOGOUT"
}

// TokenRefreshed Token刷新事件
type TokenRefreshed struct {
	UserID    string    `json:"user_id"`
	OldToken  string    `json:"old_token"`
	NewToken  string    `json:"new_token"`
	Timestamp time.Time `json:"timestamp"`
}

// NewTokenRefreshedEvent 创建Token刷新事件
func NewTokenRefreshedEvent(userID, oldToken, newToken string) *TokenRefreshed {
	return &TokenRefreshed{
		UserID:    userID,
		OldToken:  oldToken,
		NewToken:  newToken,
		Timestamp: time.Now(),
	}
}

func (e *TokenRefreshed) GetPayload() interface{} {
	return e
}

func (e *TokenRefreshed) GetType() string {
	return "AUTH.TOKEN.REFRESHED"
}

// UserProfileUpdated 用户资料更新事件
type UserProfileUpdated struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUserProfileUpdatedEvent 创建用户资料更新事件
func NewUserProfileUpdatedEvent(userID, email string) *UserProfileUpdated {
	return &UserProfileUpdated{
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
	}
}

func (e *UserProfileUpdated) GetPayload() interface{} {
	return e
}

func (e *UserProfileUpdated) GetType() string {
	return "USER.PROFILE.UPDATED"
}
