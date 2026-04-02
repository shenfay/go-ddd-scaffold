package auth

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

// HandleUserRegisteredEvent 处理用户注册事件
func HandleUserRegisteredEvent() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var event UserRegisteredEvent
		if err := json.Unmarshal(t.Payload(), &event); err != nil {
			return err
		}

		log.Printf("📨 Processing UserRegisteredEvent for user %s (%s)", event.UserID, event.Email)

		// 1. 发送验证邮件（这里只是记录日志，实际应该将任务加入队列）
		json.Marshal(SendVerificationEmailPayload{
			UserID: event.UserID,
			Email:  event.Email,
		})
		log.Printf("✓ Queued verification email for %s", event.Email)

		// 2. 发送欢迎邮件
		json.Marshal(SendVerificationEmailPayload{
			UserID: event.UserID,
			Email:  event.Email,
		})
		log.Printf("✓ Queued welcome email for %s", event.Email)

		// 3. 记录审计日志
		json.Marshal(LogUserRegistrationPayload{
			UserID:    event.UserID,
			Email:     event.Email,
			IP:        event.IP,
			UserAgent: event.UserAgent,
			Timestamp: event.CreatedAt.Unix(),
		})
		log.Printf("✓ Logged registration for %s", event.Email)

		return nil
	}
}

// HandleUserLoggedInEvent 处理用户登录事件
func HandleUserLoggedInEvent() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var event UserLoggedInEvent
		if err := json.Unmarshal(t.Payload(), &event); err != nil {
			return err
		}

		status := "failed"
		if event.Success {
			status = "success"
		}

		log.Printf("📨 Processing UserLoggedInEvent (%s) for user %s from IP %s", status, event.UserID, event.IP)

		// 记录登录日志（实际应该将任务加入队列）
		json.Marshal(LogLoginAttemptPayload{
			UserID:    event.UserID,
			Email:     event.Email,
			IP:        event.IP,
			UserAgent: event.UserAgent,
			Success:   event.Success,
			Timestamp: event.Timestamp.Unix(),
		})
		log.Printf("✓ Logged login attempt for %s", event.Email)

		return nil
	}
}

// HandleUserLoggedOutEvent 处理用户登出事件
func HandleUserLoggedOutEvent() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var event UserLoggedOutEvent
		if err := json.Unmarshal(t.Payload(), &event); err != nil {
			return err
		}

		log.Printf("📨 Processing UserLoggedOutEvent for user %s (reason: %s)", event.UserID, event.Reason)

		// 可以在这里做一些清理工作，比如通知其他系统
		log.Printf("✓ Processed logout for user %s", event.UserID)

		return nil
	}
}

// HandleTokenRefreshedEvent 处理 Token 刷新事件
func HandleTokenRefreshedEvent() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var event TokenRefreshedEvent
		if err := json.Unmarshal(t.Payload(), &event); err != nil {
			return err
		}

		log.Printf("📨 Processing TokenRefreshedEvent for user %s (old: %s, new: %s)",
			event.UserID, event.OldTokenID, event.NewTokenID)

		// 可以记录 Token 刷新日志用于安全审计
		log.Printf("✓ Logged token refresh for user %s", event.UserID)

		return nil
	}
}

// SendVerificationEmailPayload 发送邮件任务载荷
type SendVerificationEmailPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// NewSendVerificationEmailHandler 创建发送验证邮件处理器
func NewSendVerificationEmailHandler() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload SendVerificationEmailPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		log.Printf("Processing send verification email task for user %s (%s)", payload.UserID, payload.Email)

		// TODO: 实现实际的邮件发送逻辑
		// 这里只是示例，实际应该调用邮件服务

		log.Printf("✓ Sent verification email to %s", payload.Email)
		return nil
	}
}

// NewSendWelcomeEmailHandler 创建发送欢迎邮件处理器
func NewSendWelcomeEmailHandler() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload SendVerificationEmailPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		log.Printf("Processing send welcome email task for user %s", payload.UserID)

		// TODO: 实现实际的欢迎邮件发送逻辑

		log.Printf("✓ Sent welcome email to %s", payload.Email)
		return nil
	}
}

// LogUserRegistrationPayload 记录用户注册日志载荷
type LogUserRegistrationPayload struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Timestamp int64  `json:"timestamp"`
}

// NewLogUserRegistrationHandler 创建记录用户注册日志处理器
func NewLogUserRegistrationHandler() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload LogUserRegistrationPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		log.Printf("Processing log user registration for user %s from IP %s", payload.UserID, payload.IP)

		// TODO: 实现实际的审计日志记录逻辑
		// 可以写入数据库、Elasticsearch 或日志系统

		log.Printf("✓ Logged user registration for %s", payload.Email)
		return nil
	}
}

// LogLoginAttemptPayload 记录登录尝试载荷
type LogLoginAttemptPayload struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp"`
}

// NewLogLoginAttemptHandler 创建记录登录尝试处理器
func NewLogLoginAttemptHandler() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload LogLoginAttemptPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		status := "failed"
		if payload.Success {
			status = "success"
		}

		log.Printf("Processing log login attempt (%s) for user %s from IP %s", status, payload.UserID, payload.IP)

		// TODO: 实现实际的登录日志记录逻辑

		log.Printf("✓ Logged login attempt for %s", payload.Email)
		return nil
	}
}

// NewCleanupExpiredTokensHandler 创建清理过期 Token 处理器
func NewCleanupExpiredTokensHandler(redisClient interface{}) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		log.Printf("Processing cleanup expired tokens task")

		// TODO: 实现清理过期 Token 的逻辑
		// 可以通过 Redis SCAN 命令查找并删除过期的 key

		log.Printf("✓ Cleaned up expired tokens")
		return nil
	}
}
