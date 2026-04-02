package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// AsynqEventBus 基于 Asynq 的事件总线实现
type AsynqEventBus struct {
	client   *asynq.Client
	handlers map[string][]EventHandler
}

// NewAsynqEventBus 创建 Asynq 事件总线
func NewAsynqEventBus(client *asynq.Client) *AsynqEventBus {
	return &AsynqEventBus{
		client:   client,
		handlers: make(map[string][]EventHandler),
	}
}

// Publish 发布事件到 Asynq 队列
func (b *AsynqEventBus) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event.GetPayload())
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	task := asynq.NewTask(event.GetType(), payload)

	// 将事件入队（异步）
	info, err := b.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue event: %w", err)
	}

	log.Printf("✓ Published event %s to queue (ID: %s)", event.GetType(), info.ID)
	return nil
}

// Subscribe 订阅事件（内存中，用于同步处理场景）
func (b *AsynqEventBus) Subscribe(eventType string, handler EventHandler) {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	log.Printf("✓ Subscribed handler for event type: %s", eventType)
}

// GetHandlers 获取事件类型的处理器列表
func (b *AsynqEventBus) GetHandlers(eventType string) []EventHandler {
	return b.handlers[eventType]
}

// ProcessEvent 处理事件（在 Worker 中调用）
func ProcessEvent(ctx context.Context, eventType string, payload []byte, handlers []EventHandler) error {
	for _, handler := range handlers {
		// 创建一个包装事件
		event := &BaseEvent{
			Type:    eventType,
			Payload: payload,
		}

		if err := handler(ctx, event); err != nil {
			log.Printf("✗ Event handler failed for %s: %v", eventType, err)
			return err
		}
	}

	return nil
}
