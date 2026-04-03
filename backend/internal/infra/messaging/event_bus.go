package messaging

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// EventBus 事件总线接口
type EventBus interface {
	Publish(ctx context.Context, evt event.Event) error
	Subscribe(eventType string, handler event.EventHandler)
}

// QueueConfig 队列配置
type QueueConfig struct {
	Critical string // 高优先级队列（审计日志）
	Default  string // 普通优先级队列（活动日志）
}

// asynqEventBus EventBus 的 Asynq 实现
type asynqEventBus struct {
	client *asynq.Client
	config QueueConfig
}

// NewEventBus 工厂方法（统一入口）
func NewEventBus(redisAddr string, config QueueConfig) EventBus {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &asynqEventBus{client: client, config: config}
}

// Publish 发布事件到 Asynq 队列
func (b *asynqEventBus) Publish(ctx context.Context, evt event.Event) error {
	payload, _ := json.Marshal(evt.GetPayload())

	_, err := b.client.EnqueueContext(ctx,
		asynq.NewTask(evt.GetType(), payload),
		asynq.Queue(b.getQueueForEvent(evt.GetType())),
	)

	return err
}

// getQueueForEvent 根据事件类型选择队列
func (b *asynqEventBus) getQueueForEvent(eventType string) string {
	// 审计日志 → critical 队列（高优先级）
	if strings.HasPrefix(eventType, "AUTH.") || strings.HasPrefix(eventType, "SECURITY.") {
		return b.config.Critical
	}
	// 活动日志 → default 队列
	return b.config.Default
}

// Subscribe 订阅事件（由 Listener 调用）
func (b *asynqEventBus) Subscribe(eventType string, handler event.EventHandler) {
	// Listener 负责注册到 Worker 的 ServeMux
	// 这里只是标记已订阅
}
