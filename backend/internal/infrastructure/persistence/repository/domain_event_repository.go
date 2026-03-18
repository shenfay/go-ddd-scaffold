package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/messaging"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
	"github.com/shenfay/go-ddd-scaffold/shared/kernel"
)

// DomainEventRepository 领域事件仓储实现
type DomainEventRepository struct {
	query *dao.Query
}

// NewDomainEventRepository 创建领域事件仓储
func NewDomainEventRepository(query *dao.Query) *DomainEventRepository {
	return &DomainEventRepository{query: query}
}

// SaveEvents 保存领域事件到数据库
func (r *DomainEventRepository) SaveEvents(ctx context.Context, aggregateID string, aggregateType string, events []kernel.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		domainEvent := &model.DomainEvent{
			AggregateID:   aggregateID,
			AggregateType: aggregateType,
			EventType:     event.EventName(),
			EventVersion:  int32(event.Version()),
			EventData:     string(eventData),
			OccurredOn:    event.OccurredOn(),
			Processed:     util.Bool(false),
		}

		err = r.query.DomainEvent.WithContext(ctx).Create(domainEvent)
		if err != nil {
			return fmt.Errorf("failed to save domain event: %w", err)
		}
	}

	return nil
}

// GetPendingEvents 获取待处理的事件
func (r *DomainEventRepository) GetPendingEvents(ctx context.Context, limit int) ([]*messaging.EventRecord, error) {
	// 查询未处理的事件，按发生时间排序
	var events []*model.DomainEvent
	err := r.query.DomainEvent.WithContext(ctx).UnderlyingDB().
		Where("processed = ?", false).
		Order("occurred_on ASC").
		Limit(limit).
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pending events: %w", err)
	}

	// 转换为 EventRecord
	records := make([]*messaging.EventRecord, len(events))
	for i, e := range events {
		records[i] = &messaging.EventRecord{
			ID:            e.ID,
			AggregateID:   e.AggregateID,
			AggregateType: e.AggregateType,
			EventType:     e.EventType,
			EventVersion:  e.EventVersion,
			EventData:     e.EventData,
			OccurredOn:    e.OccurredOn.Format("2006-01-02T15:04:05Z07:00"),
			Metadata:      e.Metadata,
		}
	}

	return records, nil
}

// MarkAsProcessed 标记事件为已处理
func (r *DomainEventRepository) MarkAsProcessed(ctx context.Context, eventID int64) error {
	result, err := r.query.DomainEvent.WithContext(ctx).
		Where(r.query.DomainEvent.ID.Eq(eventID)).
		Update(r.query.DomainEvent.Processed, true)
	if err != nil {
		return fmt.Errorf("failed to mark event as processed: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("event not found: %d", eventID)
	}
	return nil
}

// MarkAsFailed 标记事件处理失败（将 processed 设为 true 并记录错误到 metadata）
func (r *DomainEventRepository) MarkAsFailed(ctx context.Context, eventID int64, errorMsg string) error {
	// 获取当前事件
	event, err := r.query.DomainEvent.WithContext(ctx).
		Where(r.query.DomainEvent.ID.Eq(eventID)).
		First()
	if err != nil {
		return fmt.Errorf("failed to find event: %w", err)
	}

	// 更新 metadata 记录错误信息
	metadata := "{}"
	if event.Metadata != nil {
		metadata = *event.Metadata
	}

	// 简单地将错误信息附加到 metadata（实际项目中可能需要更复杂的结构）
	result, err := r.query.DomainEvent.WithContext(ctx).
		Where(r.query.DomainEvent.ID.Eq(eventID)).
		Updates(map[string]interface{}{
			"processed": true,
			"metadata":  fmt.Sprintf(`%s,"error":"%s"`, metadata, errorMsg),
		})
	if err != nil {
		return fmt.Errorf("failed to mark event as failed: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("event not found: %d", eventID)
	}
	return nil
}

// 确保 DomainEventRepository 实现了 EventStore 接口
var _ messaging.EventStore = (*DomainEventRepository)(nil)
