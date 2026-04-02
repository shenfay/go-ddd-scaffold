package event

import (
	"context"
	"encoding/json"
)

// Event 领域事件基接口
type Event interface {
	// GetType 获取事件类型
	GetType() string
	// GetPayload 获取事件载荷
	GetPayload() interface{}
}

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event Event) error
	// Subscribe 订阅事件
	Subscribe(eventType string, handler EventHandler)
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event Event) error

// BaseEvent 基础事件实现
type BaseEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// GetType 获取事件类型
func (e *BaseEvent) GetType() string {
	return e.Type
}

// GetPayload 获取事件载荷
func (e *BaseEvent) GetPayload() interface{} {
	return e.Payload
}

// Marshal 序列化事件
func Marshal(event Event) ([]byte, error) {
	return json.Marshal(event)
}

// Unmarshal 反序列化事件
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
