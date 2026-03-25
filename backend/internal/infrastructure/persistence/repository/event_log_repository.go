package repository

import (
	"context"
	"encoding/json"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
)

// EventLogRepositoryImpl 事件日志仓储实现
type EventLogRepositoryImpl struct {
	query *dao.Query
}

// NewEventLogRepository 创建事件日志仓储
func NewEventLogRepository(db *dao.Query) aggregate.EventLogRepository {
	return &EventLogRepositoryImpl{query: db}
}

// Save 保存事件日志
func (r *EventLogRepositoryImpl) Save(ctx context.Context, log *aggregate.EventLog) error {
	daoModel := r.fromDomain(log)
	return r.query.DomainEvent.WithContext(ctx).Create(daoModel)
}

// GetByAggregateID 按聚合根 ID 查询（用于事件回放）
func (r *EventLogRepositoryImpl) GetByAggregateID(ctx context.Context, aggregateID string, limit int) ([]*aggregate.EventLog, error) {
	daoModels, err := r.query.DomainEvent.WithContext(ctx).
		Where(r.query.DomainEvent.AggregateID.Eq(aggregateID)).
		Order(r.query.DomainEvent.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*aggregate.EventLog, len(daoModels))
	for i, m := range daoModels {
		logs[i] = r.toDomain(m)
	}
	return logs, nil
}

// GetByEventType 按事件类型查询
func (r *EventLogRepositoryImpl) GetByEventType(ctx context.Context, eventType string, limit int) ([]*aggregate.EventLog, error) {
	daoModels, err := r.query.DomainEvent.WithContext(ctx).
		Where(r.query.DomainEvent.EventType.Eq(eventType)).
		Order(r.query.DomainEvent.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*aggregate.EventLog, len(daoModels))
	for i, m := range daoModels {
		logs[i] = r.toDomain(m)
	}
	return logs, nil
}

// fromDomain 将领域模型转换为 DAO 模型
func (r *EventLogRepositoryImpl) fromDomain(log *aggregate.EventLog) *model.DomainEvent {
	eventDataJSON, _ := json.Marshal(log.EventData)

	return &model.DomainEvent{
		ID:            log.ID,
		AggregateID:   log.AggregateID,
		AggregateType: log.AggregateType,
		EventType:     log.EventType,
		EventData:     string(eventDataJSON),
		OccurredAt:    &log.OccurredAt,
		CreatedAt:     &log.CreatedAt,
	}
}

// toDomain 将 DAO 模型转换为领域模型
func (r *EventLogRepositoryImpl) toDomain(m *model.DomainEvent) *aggregate.EventLog {
	var eventData map[string]any
	if m.EventData != "" {
		json.Unmarshal([]byte(m.EventData), &eventData)
	}

	return &aggregate.EventLog{
		ID:            m.ID,
		AggregateID:   m.AggregateID,
		AggregateType: m.AggregateType,
		EventType:     m.EventType,
		EventData:     eventData,
		OccurredAt:    *m.OccurredAt,
		CreatedAt:     *m.CreatedAt,
	}
}
