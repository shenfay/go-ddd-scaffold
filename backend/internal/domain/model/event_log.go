package model

import (
	"context"
	"time"
)

// EventLog 事件日志（轻量级事件存储）
type EventLog struct {
	ID            int64          `json:"id"`
	AggregateID   string         `json:"aggregate_id"`   // 聚合根 ID
	AggregateType string         `json:"aggregate_type"` // 聚合类型：User, Tenant, Order
	EventType     string         `json:"event_type"`     // 事件类型
	EventData     map[string]any `json:"event_data"`     // 事件数据
	OccurredAt    time.Time      `json:"occurred_at"`    // 发生时间
	CreatedAt     time.Time      `json:"created_at"`     // 创建时间
}

// NewEventLog 创建事件日志
func NewEventLog(aggregateID, aggregateType, eventType string, eventData map[string]any) *EventLog {
	return &EventLog{
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventType:     eventType,
		EventData:     eventData,
		OccurredAt:    time.Now(),
		CreatedAt:     time.Now(),
	}
}

// EventLogRepository 事件日志仓储接口
type EventLogRepository interface {
	// Save 保存事件日志
	Save(ctx context.Context, log *EventLog) error

	// GetByAggregateID 按聚合根 ID 查询（用于事件回放）
	GetByAggregateID(ctx context.Context, aggregateID string, limit int) ([]*EventLog, error)

	// GetByEventType 按事件类型查询
	GetByEventType(ctx context.Context, eventType string, limit int) ([]*EventLog, error)
}
