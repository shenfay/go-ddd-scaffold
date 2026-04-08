package user

import (
	"time"
)

// UserRegistered 用户注册事件
type UserRegistered struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *UserRegistered) GetType() string {
	return "USER.REGISTERED"
}

func (e *UserRegistered) GetPayload() interface{} {
	return map[string]interface{}{
		"user_id":   e.UserID,
		"email":     e.Email,
		"timestamp": e.Timestamp,
	}
}

// UserProfileUpdated 用户资料更新事件
type UserProfileUpdated struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *UserProfileUpdated) GetType() string {
	return "USER.PROFILE.UPDATED"
}

func (e *UserProfileUpdated) GetPayload() interface{} {
	return map[string]interface{}{
		"user_id":   e.UserID,
		"email":     e.Email,
		"timestamp": e.Timestamp,
	}
}
