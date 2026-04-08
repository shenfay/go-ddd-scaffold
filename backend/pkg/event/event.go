package event

import (
	"context"
	"encoding/json"
)

// Event 领域事件接口
type Event interface {
	// GetType 获取事件类型
	GetType() string
	// GetPayload 获取事件载荷
	GetPayload() interface{}
}

// EventHandler 事件处理器类型
type EventHandler func(ctx context.Context, event Event) error

// Marshal 序列化事件
func Marshal(event Event) ([]byte, error) {
	return json.Marshal(event)
}

// Unmarshal 反序列化事件
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
