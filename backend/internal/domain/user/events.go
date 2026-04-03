package events

import "time"

// UserLoggedIn 用户登录事件
type UserLoggedIn struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Device    string    `json:"device"`
	Timestamp time.Time `json:"timestamp"`
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
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	IP        string `json:"ip"`
	Reason    string `json:"reason"`
	Timestamp time.Time
}

func (e *LoginFailed) GetPayload() interface{} {
	return e
}

func (e *LoginFailed) GetType() string {
	return "AUTH.LOGIN.FAILED"
}

// AccountLocked 账户锁定事件
type AccountLocked struct {
	UserID        string    `json:"user_id"`
	Email         string    `json:"email"`
	FailedAttempts int       `json:"failed_attempts"`
	LockedUntil   time.Time `json:"locked_until"`
	Timestamp     time.Time `json:"timestamp"`
}

func (e *AccountLocked) GetPayload() interface{} {
	return e
}

func (e *AccountLocked) GetType() string {
	return "SECURITY.ACCOUNT.LOCKED"
}
