package listener

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// AuditLogListener 审计日志监听器
type AuditLogListener struct {
	eventBus messaging.EventBus
}

// NewAuditLogListener 创建审计日志监听器
func NewAuditLogListener(eventBus messaging.EventBus) *AuditLogListener {
	l := &AuditLogListener{eventBus: eventBus}

	// 订阅认证相关事件
	eventBus.Subscribe("AUTH.LOGIN.SUCCESS", l.HandleUserLoggedIn)
	eventBus.Subscribe("AUTH.LOGIN.FAILED", l.HandleLoginFailed)
	eventBus.Subscribe("SECURITY.ACCOUNT.LOCKED", l.HandleAccountLocked)

	return l
}

// HandleUserLoggedIn 处理用户登录成功事件
func (l *AuditLogListener) HandleUserLoggedIn(ctx context.Context, evt event.Event) error {
	e := evt.(*user.UserLoggedIn)

	// 转换为审计日志任务并发布到 Worker 队列
	task := &AuditLogTask{
		Type:   "audit.log.task",
		Action: "AUTH.LOGIN.SUCCESS",
		Status: "SUCCESS",
		Data: map[string]interface{}{
			"user_id":    e.UserID,
			"email":      e.Email,
			"ip":         e.IP,
			"user_agent": e.UserAgent,
			"device":     e.Device,
		},
	}

	return l.eventBus.Publish(ctx, task)
}

// HandleLoginFailed 处理用户登录失败事件
func (l *AuditLogListener) HandleLoginFailed(ctx context.Context, evt event.Event) error {
	e := evt.(*user.LoginFailed)

	task := &AuditLogTask{
		Type:   "audit.log.task",
		Action: "AUTH.LOGIN.FAILED",
		Status: "FAILED",
		Data: map[string]interface{}{
			"user_id": e.UserID,
			"email":   e.Email,
			"ip":      e.IP,
			"reason":  e.Reason,
		},
	}

	return l.eventBus.Publish(ctx, task)
}

// HandleAccountLocked 处理账户锁定事件
func (l *AuditLogListener) HandleAccountLocked(ctx context.Context, evt event.Event) error {
	e := evt.(*user.AccountLocked)

	task := &AuditLogTask{
		Type:   "audit.log.task",
		Action: "SECURITY.ACCOUNT.LOCKED",
		Status: "FAILED",
		Data: map[string]interface{}{
			"user_id":         e.UserID,
			"email":           e.Email,
			"failed_attempts": e.FailedAttempts,
			"locked_until":    e.LockedUntil,
		},
	}

	return l.eventBus.Publish(ctx, task)
}
