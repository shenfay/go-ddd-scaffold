package domain_event

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// EventStore 领域事件存储接口（纯事件溯源）
type EventStore interface {
	// SaveEvents 保存聚合根的未提交事件（仅用于溯源，不跟踪状态）
	// aggregateID: 聚合根 ID
	// aggregateType: 聚合类型（如 user, tenant）
	// events: 待保存的事件列表
	SaveEvents(ctx context.Context, aggregateID string, aggregateType string, events []kernel.DomainEvent) error

	// GetEvents 获取聚合根的所有历史事件（用于事件回放）
	// aggregateID: 聚合根 ID
	GetEvents(ctx context.Context, aggregateID string) ([]*EventRecord, error)

	// GetEventsByType 按类型获取事件（用于分析和审计）
	// eventType: 事件类型
	// limit: 限制数量
	GetEventsByType(ctx context.Context, eventType string, limit int) ([]*EventRecord, error)
}

// EventRecord 事件记录（从数据库读取）
type EventRecord struct {
	ID            int64
	AggregateID   string
	AggregateType string
	EventType     string
	EventVersion  int32
	EventData     string
	OccurredOn    string
	Metadata      *string
}
