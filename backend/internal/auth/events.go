package auth

import (
	"context"
	"log"
	"time"

	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// ==================== 用户认证事件 ====================

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IP        string    `json:"ip,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

// GetType 获取事件类型
func (e *UserRegisteredEvent) GetType() string {
	return "user.registered"
}

// GetPayload 获取事件载荷
func (e *UserRegisteredEvent) GetPayload() interface{} {
	return e
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID, email, ip, userAgent string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		UserID:    userID,
		Email:     email,
		CreatedAt: time.Now(),
		IP:        ip,
		UserAgent: userAgent,
	}
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

// GetType 获取事件类型
func (e *UserLoggedInEvent) GetType() string {
	return "user.logged_in"
}

// GetPayload 获取事件载荷
func (e *UserLoggedInEvent) GetPayload() interface{} {
	return e
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(userID, email, ip, userAgent string, success bool) *UserLoggedInEvent {
	return &UserLoggedInEvent{
		UserID:    userID,
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		Success:   success,
		Timestamp: time.Now(),
	}
}

// UserLoggedOutEvent 用户退出事件
type UserLoggedOutEvent struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"` // logout, token_expired, kicked_out
}

// GetType 获取事件类型
func (e *UserLoggedOutEvent) GetType() string {
	return "user.logged_out"
}

// GetPayload 获取事件载荷
func (e *UserLoggedOutEvent) GetPayload() interface{} {
	return e
}

// NewUserLoggedOutEvent 创建用户退出事件
func NewUserLoggedOutEvent(userID, reason string) *UserLoggedOutEvent {
	return &UserLoggedOutEvent{
		UserID:    userID,
		Timestamp: time.Now(),
		Reason:    reason,
	}
}

// TokenRefreshedEvent Token 刷新事件
type TokenRefreshedEvent struct {
	UserID     string    `json:"user_id"`
	OldTokenID string    `json:"old_token_id"`
	NewTokenID string    `json:"new_token_id"`
	Timestamp  time.Time `json:"timestamp"`
}

// GetType 获取事件类型
func (e *TokenRefreshedEvent) GetType() string {
	return "token.refreshed"
}

// GetPayload 获取事件载荷
func (e *TokenRefreshedEvent) GetPayload() interface{} {
	return e
}

// NewTokenRefreshedEvent 创建 Token 刷新事件
func NewTokenRefreshedEvent(userID, oldTokenID, newTokenID string) *TokenRefreshedEvent {
	return &TokenRefreshedEvent{
		UserID:     userID,
		OldTokenID: oldTokenID,
		NewTokenID: newTokenID,
		Timestamp:  time.Now(),
	}
}

// PublishEvent 发布事件（辅助函数）
func PublishEvent(eventBus event.EventBus, ctx context.Context, event event.Event) error {
	if event == nil {
		// 事件为 nil，直接返回
		return nil
	}

	if eventBus == nil {
		// 如果没有事件总线，记录日志但不处理
		log.Printf("⚠ Event bus not configured, skipping event: %s", event.GetType())
		return nil
	}

	return eventBus.Publish(ctx, event)
}
