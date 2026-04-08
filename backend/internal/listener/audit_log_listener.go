package listener

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// AuditLogListener 审计日志事件监听器
type AuditLogListener struct {
	eventBus messaging.EventBus
}

// NewAuditLogListener 创建审计日志监听器实例
func NewAuditLogListener(eventBus messaging.EventBus) *AuditLogListener {
	return &AuditLogListener{eventBus: eventBus}
}

// SubscribeEvents 订阅事件（在 API 初始化时调用）
func (l *AuditLogListener) SubscribeEvents(eventBus messaging.EventBus) {
	eventBus.Subscribe("AUTH.LOGIN.SUCCESS", l.HandleUserLoggedIn)
	eventBus.Subscribe("AUTH.LOGIN.FAILED", l.HandleLoginFailed)
	eventBus.Subscribe("SECURITY.ACCOUNT.LOCKED", l.HandleAccountLocked)
}

// HandleUserLoggedIn 处理用户登录成功事件
func (l *AuditLogListener) HandleUserLoggedIn(ctx context.Context, evt event.Event) error {
	e := evt.(*user.UserLoggedIn)

	// 解析 User-Agent
	uaInfo := utils.ParseUserAgent(e.UserAgent)

	// 转换为审计日志任务并发布到 EventBus（异步）
	task := &AuditLogTask{
		Type:   "audit.log",
		Action: "AUTH.LOGIN.SUCCESS",
		Status: "SUCCESS",
		Data: map[string]interface{}{
			"user_id":    e.UserID,
			"email":      e.Email,
			"ip":         e.IP,
			"user_agent": e.UserAgent,
			"device":     uaInfo.Device,
			"browser":    uaInfo.Browser,
			"os":         uaInfo.OS,
		},
	}

	return l.eventBus.Publish(ctx, task)
}

// HandleLoginFailed 处理用户登录失败事件
func (l *AuditLogListener) HandleLoginFailed(ctx context.Context, evt event.Event) error {
	e := evt.(*user.LoginFailed)

	task := &AuditLogTask{
		Type:   "audit.log",
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
		Type:   "audit.log",
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
