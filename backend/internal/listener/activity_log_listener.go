package listener

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// ActivityLogListener 活动日志事件监听器
type ActivityLogListener struct {
	eventBus messaging.EventBus
}

// NewActivityLogListener 创建活动日志监听器实例
func NewActivityLogListener(eventBus messaging.EventBus) *ActivityLogListener {
	return &ActivityLogListener{eventBus: eventBus}
}

// SubscribeEvents 订阅事件（在 API 初始化时调用）
func (l *ActivityLogListener) SubscribeEvents(eventBus messaging.EventBus) {
	eventBus.Subscribe("USER.REGISTERED", l.HandleUserRegistered)
	eventBus.Subscribe("AUTH.LOGIN.SUCCESS", l.HandleUserLoggedIn)
	eventBus.Subscribe("AUTH.LOGOUT", l.HandleUserLoggedOut)
	eventBus.Subscribe("AUTH.TOKEN.REFRESHED", l.HandleTokenRefreshed)
}

// HandleUserRegistered 处理用户注册事件
func (l *ActivityLogListener) HandleUserRegistered(ctx context.Context, evt event.Event) error {
	e := evt.(*user.UserRegistered)

	task := &ActivityLogTask{
		Type:        "activity.log",
		Action:      "USER.REGISTERED",
		UserID:      e.UserID,
		Email:       e.Email,
		Description: "用户注册",
		Status:      "SUCCESS",
		Metadata:    nil,
	}

	return l.eventBus.Publish(ctx, task)
}

// HandleUserLoggedIn 处理用户登录事件
func (l *ActivityLogListener) HandleUserLoggedIn(ctx context.Context, evt event.Event) error {
	e := evt.(*user.UserLoggedIn)

	// 解析 User-Agent
	uaInfo := utils.ParseUserAgent(e.UserAgent)

	task := &ActivityLogTask{
		Type:        "activity.log",
		Action:      "USER.LOGIN",
		UserID:      e.UserID,
		Email:       e.Email,
		Description: "用户登录",
		Status:      "SUCCESS",
		IP:          e.IP,
		UserAgent:   e.UserAgent,
		Device:      uaInfo.Device,
		Browser:     uaInfo.Browser,
		OS:          uaInfo.OS,
		Metadata: map[string]interface{}{
			"ip":         e.IP,
			"user_agent": e.UserAgent,
			"device":     uaInfo.Device,
			"browser":    uaInfo.Browser,
			"os":         uaInfo.OS,
		},
	}

	return l.eventBus.Publish(ctx, task)
}

// HandleUserLoggedOut 处理用户登出事件
func (l *ActivityLogListener) HandleUserLoggedOut(ctx context.Context, evt event.Event) error {
	e := evt.(*user.UserLoggedOut)

	task := &ActivityLogTask{
		Type:        "activity.log",
		Action:      "USER.LOGOUT",
		UserID:      e.UserID,
		Email:       e.Email,
		Description: "用户登出",
		Status:      "SUCCESS",
		Metadata:    nil,
	}

	return l.eventBus.Publish(ctx, task)
}

// HandleTokenRefreshed 处理Token刷新事件
func (l *ActivityLogListener) HandleTokenRefreshed(ctx context.Context, evt event.Event) error {
	e := evt.(*user.TokenRefreshed)

	task := &ActivityLogTask{
		Type:        "activity.log",
		Action:      "USER.TOKEN_REFRESHED",
		UserID:      e.UserID,
		Email:       "",
		Description: "Token刷新",
		Status:      "SUCCESS",
		Metadata: map[string]interface{}{
			"old_token": e.OldToken,
			"new_token": e.NewToken,
		},
	}

	return l.eventBus.Publish(ctx, task)
}
