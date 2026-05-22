package user

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/events"
)

// 编译期断言：所有领域事件实现 events.DomainEvent 接口
var (
	_ events.DomainEvent = (*UserRegistered)(nil)
	_ events.DomainEvent = (*UserLoggedIn)(nil)
	_ events.DomainEvent = (*LoginFailed)(nil)
	_ events.DomainEvent = (*AccountLocked)(nil)
	_ events.DomainEvent = (*UserLoggedOut)(nil)
	_ events.DomainEvent = (*TokenRefreshed)(nil)
	_ events.DomainEvent = (*UserProfileUpdated)(nil)
)

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

// EventName 返回事件名
func (e *UserRegistered) EventName() string { return "user.registered" }

// OccurredAt 返回事件发生时间
func (e *UserRegistered) OccurredAt() time.Time { return e.Timestamp }

// GetPayload 获取事件数据（保留兼容旧接口）
func (e *UserRegistered) GetPayload() interface{} { return e }

// GetType 获取事件类型（保留兼容 Worker 反序列化）
func (e *UserRegistered) GetType() string { return "user.registered" }

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

func (e *UserLoggedIn) EventName() string       { return "user.logged_in" }
func (e *UserLoggedIn) OccurredAt() time.Time   { return e.Timestamp }
func (e *UserLoggedIn) GetPayload() interface{} { return e }
func (e *UserLoggedIn) GetType() string         { return "user.logged_in" }

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

func (e *LoginFailed) EventName() string       { return "user.login_failed" }
func (e *LoginFailed) OccurredAt() time.Time   { return e.Timestamp }
func (e *LoginFailed) GetPayload() interface{} { return e }
func (e *LoginFailed) GetType() string         { return "user.login_failed" }

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

func (e *AccountLocked) EventName() string       { return "user.account_locked" }
func (e *AccountLocked) OccurredAt() time.Time   { return e.Timestamp }
func (e *AccountLocked) GetPayload() interface{} { return e }
func (e *AccountLocked) GetType() string         { return "user.account_locked" }

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

func (e *UserLoggedOut) EventName() string       { return "user.logged_out" }
func (e *UserLoggedOut) OccurredAt() time.Time   { return e.Timestamp }
func (e *UserLoggedOut) GetPayload() interface{} { return e }
func (e *UserLoggedOut) GetType() string         { return "user.logged_out" }

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

func (e *TokenRefreshed) EventName() string       { return "user.token_refreshed" }
func (e *TokenRefreshed) OccurredAt() time.Time   { return e.Timestamp }
func (e *TokenRefreshed) GetPayload() interface{} { return e }
func (e *TokenRefreshed) GetType() string         { return "user.token_refreshed" }

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

func (e *UserProfileUpdated) EventName() string       { return "user.profile_updated" }
func (e *UserProfileUpdated) OccurredAt() time.Time   { return e.Timestamp }
func (e *UserProfileUpdated) GetPayload() interface{} { return e }
func (e *UserProfileUpdated) GetType() string         { return "user.profile_updated" }
