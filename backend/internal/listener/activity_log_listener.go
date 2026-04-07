package listener

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// ActivityLogListener 活动日志监听器
type ActivityLogListener struct {
	eventBus messaging.EventBus
}

// NewActivityLogListener 创建活动日志监听器
func NewActivityLogListener(eventBus messaging.EventBus) *ActivityLogListener {
	l := &ActivityLogListener{eventBus: eventBus}

	// 订阅用户相关事件
	eventBus.Subscribe("USER.REGISTERED", l.HandleUserRegistered)
	eventBus.Subscribe("AUTH.LOGOUT", l.HandleUserLoggedOut)
	eventBus.Subscribe("AUTH.TOKEN.REFRESHED", l.HandleTokenRefreshed)

	return l
}

// HandleUserRegistered 处理用户注册事件
func (l *ActivityLogListener) HandleUserRegistered(ctx context.Context, evt event.Event) error {
	e := evt.(*user.UserRegistered)

	task := &ActivityLogTask{
		Type:        "activity:record",
		Action:      "USER.REGISTERED",
		UserID:      e.UserID,
		Email:       e.Email,
		Description: "用户注册",
		Metadata:    nil,
	}

	return l.eventBus.Publish(ctx, task)
}

// HandleUserLoggedOut 处理用户登出事件
func (l *ActivityLogListener) HandleUserLoggedOut(ctx context.Context, evt event.Event) error {
	e := evt.(*user.UserLoggedOut)

	task := &ActivityLogTask{
		Type:        "activity:record",
		Action:      "USER.LOGOUT",
		UserID:      e.UserID,
		Email:       e.Email,
		Description: "用户登出",
		Metadata:    nil,
	}

	return l.eventBus.Publish(ctx, task)
}

// HandleTokenRefreshed 处理Token刷新事件
func (l *ActivityLogListener) HandleTokenRefreshed(ctx context.Context, evt event.Event) error {
	e := evt.(*user.TokenRefreshed)

	task := &ActivityLogTask{
		Type:        "activity:record",
		Action:      "USER.TOKEN_REFRESHED",
		UserID:      e.UserID,
		Email:       "",
		Description: "Token刷新",
		Metadata: map[string]interface{}{
			"old_token": e.OldToken,
			"new_token": e.NewToken,
		},
	}

	return l.eventBus.Publish(ctx, task)
}
