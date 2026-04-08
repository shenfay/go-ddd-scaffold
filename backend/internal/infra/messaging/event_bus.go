package messaging

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// EventBus 进程内事件总线实现接口
type EventBus interface {
	Publish(ctx context.Context, evt event.Event) error
	Subscribe(eventType string, handler event.EventHandler)
}

// AsynqEventBus 基于 Asynq 的事件总线实现
type AsynqEventBus struct {
	client   *asynq.Client
	handlers map[string][]event.EventHandler // Worker 端订阅使用
}

// NewEventBus 创建 Asynq 事件总线（工厂方法）
func NewEventBus(client *asynq.Client) EventBus {
	return &AsynqEventBus{
		client:   client,
		handlers: make(map[string][]event.EventHandler),
	}
}

// Publish 发布领域事件到所有订阅者到 Asynq 队列（异步）
func (b *AsynqEventBus) Publish(ctx context.Context, evt event.Event) error {
	payload, err := json.Marshal(evt.GetPayload())
	if err != nil {
		return err
	}

	// 根据事件类型自动选择队列
	queue := b.getQueueForEvent(evt.GetType())

	// 配置重试策略：根据队列优先级设置不同重试次数
	maxRetry := 3
	if queue == "critical" {
		maxRetry = 5 // 高优先级队列更多重试
	}

	_, err = b.client.EnqueueContext(ctx,
		asynq.NewTask(evt.GetType(), payload),
		asynq.Queue(queue),
		asynq.MaxRetry(maxRetry),
	)

	return err
}

// getQueueForEvent 根据事件类型选择队列
func (b *AsynqEventBus) getQueueForEvent(eventType string) string {
	// 认证/安全事件 → critical 队列（高优先级）
	if strings.HasPrefix(eventType, "AUTH.") || strings.HasPrefix(eventType, "SECURITY.") {
		return "critical"
	}
	// 用户事件 → default 队列
	if strings.HasPrefix(eventType, "USER.") {
		return "default"
	}
	// 日志任务 → 根据类型
	if eventType == "audit.log" {
		return "critical"
	}
	if eventType == "activity.log" {
		return "default"
	}
	return "default"
}

// Subscribe 订阅指定类型的领域事件（由 Worker 端调用）
func (b *AsynqEventBus) Subscribe(eventType string, handler event.EventHandler) {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// GetHandlers 获取事件类型的处理器列表（Worker 端使用）
func (b *AsynqEventBus) GetHandlers(eventType string) []event.EventHandler {
	return b.handlers[eventType]
}

// GetSubscriptions 获取所有订阅的事件类型（Worker 端使用）
func (b *AsynqEventBus) GetSubscriptions() map[string][]event.EventHandler {
	return b.handlers
}

// DispatchEvent 分发事件到所有订阅的处理器（Worker 端使用）
func (b *AsynqEventBus) DispatchEvent(ctx context.Context, evt event.Event) error {
	handlers := b.handlers[evt.GetType()]
	if len(handlers) == 0 {
		return nil // 没有订阅者，不报错
	}

	// 调用所有处理器
	for _, handler := range handlers {
		if err := handler(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}
