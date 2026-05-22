package listener

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/events"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// AuditLogListener 审计日志事件监听器
// 直接注入 Repository，将领域事件同步写入审计日志表
type AuditLogListener struct {
	repo repository.AuditLogRepository
}

// NewAuditLogListener 创建审计日志监听器实例
func NewAuditLogListener(repo repository.AuditLogRepository) *AuditLogListener {
	return &AuditLogListener{repo: repo}
}

// HandleUserLoggedIn 处理用户登录成功事件
func (l *AuditLogListener) HandleUserLoggedIn(ctx context.Context, evt events.DomainEvent) error {
	e := evt.(*user.UserLoggedIn)
	uaInfo := utils.ParseUserAgent(e.UserAgent)

	return l.repo.Save(ctx, &repository.AuditLog{
		UserID:    e.UserID,
		Email:     e.Email,
		Action:    e.EventName(),
		Status:    "SUCCESS",
		IP:        e.IP,
		UserAgent: e.UserAgent,
		Device:    uaInfo.Device,
		Browser:   uaInfo.Browser,
		OS:        uaInfo.OS,
		Metadata: map[string]interface{}{
			"ip":         e.IP,
			"user_agent": e.UserAgent,
			"device":     uaInfo.Device,
			"browser":    uaInfo.Browser,
			"os":         uaInfo.OS,
		},
	})
}

// HandleLoginFailed 处理用户登录失败事件
func (l *AuditLogListener) HandleLoginFailed(ctx context.Context, evt events.DomainEvent) error {
	e := evt.(*user.LoginFailed)

	return l.repo.Save(ctx, &repository.AuditLog{
		UserID: e.UserID,
		Email:  e.Email,
		Action: e.EventName(),
		Status: "FAILED",
		IP:     e.IP,
		Metadata: map[string]interface{}{
			"user_id": e.UserID,
			"email":   e.Email,
			"ip":      e.IP,
			"reason":  e.Reason,
		},
	})
}

// HandleAccountLocked 处理账户锁定事件
func (l *AuditLogListener) HandleAccountLocked(ctx context.Context, evt events.DomainEvent) error {
	e := evt.(*user.AccountLocked)

	return l.repo.Save(ctx, &repository.AuditLog{
		UserID: e.UserID,
		Email:  e.Email,
		Action: e.EventName(),
		Status: "FAILED",
		Metadata: map[string]interface{}{
			"user_id":         e.UserID,
			"email":           e.Email,
			"failed_attempts": e.FailedAttempts,
			"locked_until":    e.LockedUntil,
		},
	})
}
