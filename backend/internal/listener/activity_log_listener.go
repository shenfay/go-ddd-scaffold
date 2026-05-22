package listener

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/events"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// ActivityLogListener 活动日志事件监听器
// 直接注入 Repository，将领域事件同步写入活动日志表
type ActivityLogListener struct {
	repo repository.ActivityLogRepository
}

// NewActivityLogListener 创建活动日志监听器实例
func NewActivityLogListener(repo repository.ActivityLogRepository) *ActivityLogListener {
	return &ActivityLogListener{repo: repo}
}

// HandleUserRegistered 处理用户注册事件
func (l *ActivityLogListener) HandleUserRegistered(ctx context.Context, evt events.DomainEvent) error {
	e := evt.(*user.UserRegistered)
	return l.repo.Create(ctx, &repository.ActivityLog{
		UserID: e.UserID,
		Action: e.EventName(),
		Metadata: map[string]interface{}{
			"email":     e.Email,
			"timestamp": e.Timestamp,
		},
	})
}

// HandleUserLoggedIn 处理用户登录事件
func (l *ActivityLogListener) HandleUserLoggedIn(ctx context.Context, evt events.DomainEvent) error {
	e := evt.(*user.UserLoggedIn)
	uaInfo := utils.ParseUserAgent(e.UserAgent)
	return l.repo.Create(ctx, &repository.ActivityLog{
		UserID: e.UserID,
		Action: e.EventName(),
		Metadata: map[string]interface{}{
			"email":      e.Email,
			"ip":         e.IP,
			"user_agent": e.UserAgent,
			"device":     uaInfo.Device,
			"browser":    uaInfo.Browser,
			"os":         uaInfo.OS,
			"timestamp":  e.Timestamp,
		},
	})
}

// HandleUserLoggedOut 处理用户登出事件
func (l *ActivityLogListener) HandleUserLoggedOut(ctx context.Context, evt events.DomainEvent) error {
	e := evt.(*user.UserLoggedOut)
	return l.repo.Create(ctx, &repository.ActivityLog{
		UserID: e.UserID,
		Action: e.EventName(),
		Metadata: map[string]interface{}{
			"email":     e.Email,
			"timestamp": e.Timestamp,
		},
	})
}

// HandleTokenRefreshed 处理Token刷新事件
func (l *ActivityLogListener) HandleTokenRefreshed(ctx context.Context, evt events.DomainEvent) error {
	e := evt.(*user.TokenRefreshed)
	return l.repo.Create(ctx, &repository.ActivityLog{
		UserID: e.UserID,
		Action: e.EventName(),
		Metadata: map[string]interface{}{
			"old_token": e.OldToken,
			"new_token": e.NewToken,
			"timestamp": e.Timestamp,
		},
	})
}
