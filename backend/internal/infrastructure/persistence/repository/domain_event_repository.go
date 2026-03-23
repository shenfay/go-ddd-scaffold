package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
)

// DomainEventRepository 领域事件仓储实现
type DomainEventRepository struct {
	query *dao.Query
}

// NewDomainEventRepository 创建领域事件仓储
func NewDomainEventRepository(query *dao.Query) *DomainEventRepository {
	return &DomainEventRepository{query: query}
}

// SaveEvents 保存领域事件到数据库（纯溯源，不设置状态）
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
		}

		err = r.query.DomainEvent.WithContext(ctx).Create(domainEvent)
		if err != nil {
			return fmt.Errorf("failed to save domain event: %w", err)
		}
	}

	return nil
}

// GetEvents 获取聚合根的所有历史事件
func (r *DomainEventRepository) GetEvents(ctx context.Context, aggregateID string) ([]*kernel.EventRecord, error) {
	events, err := r.query.DomainEvent.WithContext(ctx).
		Where(r.query.DomainEvent.AggregateID.Eq(aggregateID)).
		Order(r.query.DomainEvent.OccurredOn.Asc()).
		Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	// 转换为 EventRecord
	records := make([]*kernel.EventRecord, len(events))
	for i, e := range events {
		records[i] = &kernel.EventRecord{
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

// GetEventsByType 按类型获取事件
func (r *DomainEventRepository) GetEventsByType(ctx context.Context, eventType string, limit int) ([]*kernel.EventRecord, error) {
	events, err := r.query.DomainEvent.WithContext(ctx).
		Where(r.query.DomainEvent.EventType.Eq(eventType)).
		Order(r.query.DomainEvent.OccurredOn.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type: %w", err)
	}

	// 转换为 EventRecord
	records := make([]*kernel.EventRecord, len(events))
	for i, e := range events {
		records[i] = &kernel.EventRecord{
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

// 确保 DomainEventRepository 实现了 kernel.DomainEventRepository 接口
var _ kernel.DomainEventRepository = (*DomainEventRepository)(nil)
