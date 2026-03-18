package messaging

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/shared/kernel"
)

// EventStore 领域事件存储接口
type EventStore interface {
	// SaveEvents 保存聚合根的未提交事件
	// aggregateID: 聚合根 ID
	// aggregateType: 聚合类型（如 user, tenant）
	// events: 待保存的事件列表
	SaveEvents(ctx context.Context, aggregateID string, aggregateType string, events []kernel.DomainEvent) error

	// GetPendingEvents 获取待处理的领域事件（用于后台处理器）
	// limit: 每次获取的最大事件数
	GetPendingEvents(ctx context.Context, limit int) ([]*EventRecord, error)

	// MarkAsProcessed 标记事件为已处理
	MarkAsProcessed(ctx context.Context, eventID int64) error

	// MarkAsFailed 标记事件处理失败
	MarkAsFailed(ctx context.Context, eventID int64, errorMsg string) error
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
