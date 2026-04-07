package listener

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// ActivityLogListener 活动日志监听器
type ActivityLogListener struct {
	asynqClient *asynq.Client
}

// NewActivityLogListener 创建活动日志监听器
func NewActivityLogListener(asynqClient *asynq.Client) *ActivityLogListener {
	return &ActivityLogListener{asynqClient: asynqClient}
}

// SubscribeEvents 订阅事件（在 API 初始化时调用）
func (l *ActivityLogListener) SubscribeEvents(eventBus messaging.EventBus) {
	eventBus.Subscribe("USER.REGISTERED", l.HandleUserRegistered)
	eventBus.Subscribe("AUTH.LOGOUT", l.HandleUserLoggedOut)
	eventBus.Subscribe("AUTH.TOKEN.REFRESHED", l.HandleTokenRefreshed)
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

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":     task.UserID,
		"email":       task.Email,
		"action":      task.Action,
		"description": task.Description,
		"metadata":    task.Metadata,
	})
	_, err := l.asynqClient.EnqueueContext(ctx,
		asynq.NewTask("activity:record", payload),
		asynq.Queue("default"),
	)
	return err
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

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":     task.UserID,
		"email":       task.Email,
		"action":      task.Action,
		"description": task.Description,
		"metadata":    task.Metadata,
	})
	_, err := l.asynqClient.EnqueueContext(ctx,
		asynq.NewTask("activity:record", payload),
		asynq.Queue("default"),
	)
	return err
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

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":     task.UserID,
		"email":       task.Email,
		"action":      task.Action,
		"description": task.Description,
		"metadata":    task.Metadata,
	})
	_, err := l.asynqClient.EnqueueContext(ctx,
		asynq.NewTask("activity:record", payload),
		asynq.Queue("default"),
	)
	return err
}
