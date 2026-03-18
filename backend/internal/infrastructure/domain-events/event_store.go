package domainevents

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shenfay/go-ddd-scaffold/shared/kernel"
)

// EventStore 领域事件存储接口
type EventStore interface {
	// AppendEvents 追加事件到存储
	AppendEvents(ctx context.Context, aggregateID interface{}, events []kernel.DomainEvent) error
	// GetEventsForAggregate 获取聚合的事件
	GetEventsForAggregate(ctx context.Context, aggregateID interface{}, afterVersion int) ([]kernel.DomainEvent, error)
	// GetAllEvents 获取所有事件（用于事件回放）
	GetAllEvents(ctx context.Context, afterID int64, limit int) ([]StoredEvent, error)
}

// StoredEvent 存储的事件
type StoredEvent struct {
	ID            int64
	AggregateID   string
	AggregateType string
	EventType     string
	EventVersion  int
	EventData     string
	OccurredOn    time.Time
	CreatedAt     time.Time
}

// EventStoreImpl 事件存储实现
type EventStoreImpl struct {
	db DB
}

// DB 数据库接口
type DB interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (Result, error)
	BeginTx(ctx context.Context) (Tx, error)
}

// Row 单行结果接口
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows 多行结果接口
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

// Result 执行结果接口
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// Tx 事务接口
type Tx interface {
	DB
	Commit() error
	Rollback() error
}

// NewEventStore 创建事件存储
func NewEventStore(db DB) EventStore {
	return &EventStoreImpl{db: db}
}

// AppendEvents 追加事件
func (s *EventStoreImpl) AppendEvents(ctx context.Context, aggregateID interface{}, events []kernel.DomainEvent) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO domain_events (aggregate_id, aggregate_type, event_type, event_version, event_data, occurred_on, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			fmt.Sprintf("%v", aggregateID),
			s.getAggregateType(event),
			event.EventName(),
			event.Version(),
			string(eventData),
			event.OccurredOn(),
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to store event: %w", err)
		}
	}

	return tx.Commit()
}

// GetEventsForAggregate 获取聚合的事件
func (s *EventStoreImpl) GetEventsForAggregate(ctx context.Context, aggregateID interface{}, afterVersion int) ([]kernel.DomainEvent, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, aggregate_id, aggregate_type, event_type, event_version, event_data, occurred_on, created_at 
		FROM domain_events WHERE aggregate_id = ? AND event_version > ? ORDER BY event_version ASC`,
		fmt.Sprintf("%v", aggregateID),
		afterVersion,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanEvents(rows)
}

// GetAllEvents 获取所有事件
func (s *EventStoreImpl) GetAllEvents(ctx context.Context, afterID int64, limit int) ([]StoredEvent, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, aggregate_id, aggregate_type, event_type, event_version, event_data, occurred_on, created_at 
		FROM domain_events WHERE id > ? ORDER BY id ASC LIMIT ?`,
		afterID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []StoredEvent
	for rows.Next() {
		var event StoredEvent
		err := rows.Scan(
			&event.ID, &event.AggregateID, &event.AggregateType, &event.EventType,
			&event.EventVersion, &event.EventData, &event.OccurredOn, &event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// scanEvents 扫描事件列表
func (s *EventStoreImpl) scanEvents(rows Rows) ([]kernel.DomainEvent, error) {
	var events []kernel.DomainEvent
	for rows.Next() {
		var model StoredEvent
		err := rows.Scan(
			&model.ID, &model.AggregateID, &model.AggregateType, &model.EventType,
			&model.EventVersion, &model.EventData, &model.OccurredOn, &model.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		event, err := s.deserializeEvent(model.EventType, model.EventData)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// deserializeEvent 反序列化事件
func (s *EventStoreImpl) deserializeEvent(eventType string, eventData string) (kernel.DomainEvent, error) {
	// 这里返回一个通用的领域事件实现
	// 实际项目中应该根据事件类型注册对应的反序列化器
	var baseEvent struct {
		EventName   string                 `json:"event_name"`
		AggregateID interface{}            `json:"aggregate_id"`
		Version     int                    `json:"version"`
		OccurredOn  time.Time              `json:"occurred_on"`
		Metadata    map[string]interface{} `json:"metadata"`
	}

	if err := json.Unmarshal([]byte(eventData), &baseEvent); err != nil {
		return nil, err
	}

	return &GenericDomainEvent{
		eventName:   eventType,
		aggregateID: baseEvent.AggregateID,
		version:     baseEvent.Version,
		occurredOn:  baseEvent.OccurredOn,
		metadata:    baseEvent.Metadata,
		data:        eventData,
	}, nil
}

// getAggregateType 获取聚合类型
func (s *EventStoreImpl) getAggregateType(event kernel.DomainEvent) string {
	if metadata := event.Metadata(); metadata != nil {
		if aggregateType, ok := metadata["aggregate_type"].(string); ok {
			return aggregateType
		}
	}
	return "unknown"
}

// GenericDomainEvent 通用领域事件实现
type GenericDomainEvent struct {
	eventName   string
	aggregateID interface{}
	version     int
	occurredOn  time.Time
	metadata    map[string]interface{}
	data        string
}

// EventName 返回事件名称
func (e *GenericDomainEvent) EventName() string {
	return e.eventName
}

// AggregateID 返回聚合ID
func (e *GenericDomainEvent) AggregateID() interface{} {
	return e.aggregateID
}

// Version 返回版本
func (e *GenericDomainEvent) Version() int {
	return e.version
}

// OccurredOn 返回发生时间
func (e *GenericDomainEvent) OccurredOn() time.Time {
	return e.occurredOn
}

// Metadata 返回元数据
func (e *GenericDomainEvent) Metadata() map[string]interface{} {
	return e.metadata
}

// Data 返回原始数据
func (e *GenericDomainEvent) Data() string {
	return e.data
}
