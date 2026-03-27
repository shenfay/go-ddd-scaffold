package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

// DAOQuery 定义 GORM Gen DAO 的通用接口
type DAOQuery interface {
	WithContext(ctx context.Context) interface{}
}

// DAOFinder 定义查找接口
type DAOFinder[M any] interface {
	Where(query interface{}, args ...interface{}) DAOFinder[M]
	First() (M, error)
	Count() (int64, error)
}

// DomainConverter 领域对象与数据模型转换接口
type DomainConverter[T any, M any] interface {
	// ToDomain 将数据模型转换为领域对象
	ToDomain(model M) T
	// FromDomain 将领域对象转换为数据模型
	FromDomain(domain T) M
}

// BaseRepository 泛型仓储基类
// T: 领域对象类型
// M: 数据模型类型
type BaseRepository[T any, M any] struct {
	db        *gorm.DB
	converter DomainConverter[T, M]
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any, M any](db *gorm.DB, converter DomainConverter[T, M]) *BaseRepository[T, M] {
	return &BaseRepository[T, M]{
		db:        db,
		converter: converter,
	}
}

// SaveWithOutbox 保存聚合根并记录领域事件到 Outbox（同一事务）
// 用于保证业务操作与事件记录的原子性
func (r *BaseRepository[T, M]) SaveWithOutbox(ctx context.Context, aggregate *T, events []common.DomainEvent) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	// 1. 保存聚合根
	if err := tx.WithContext(ctx).Save(aggregate).Error; err != nil {
		return fmt.Errorf("failed to save aggregate: %w", err)
	}

	// 2. 保存领域事件到 Outbox（同一事务）
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		// 序列化元数据
		var metadataJSON *string
		if meta := event.Metadata(); len(meta) > 0 {
			metadataBytes, _ := json.Marshal(meta)
			metadataStr := string(metadataBytes)
			metadataJSON = &metadataStr
		}

		now := time.Now()
		outbox := &model.Outbox{
			ID:            now.UnixNano(),
			EventType:     event.EventName(),
			AggregateType: getAggregateType(event),
			AggregateID:   fmt.Sprintf("%v", event.AggregateID()),
			Payload:       string(payload),
			Metadata:      metadataJSON,
			OccurredAt:    &now,
			Processed:     false,
			RetryCount:    0,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}

		if err := tx.WithContext(ctx).Create(outbox).Error; err != nil {
			return fmt.Errorf("failed to save to outbox: %w", err)
		}
	}

	return tx.Commit().Error
}

// DB 获取 GORM DB 实例（供子类使用）
func (r *BaseRepository[T, M]) DB() *gorm.DB {
	return r.db
}

// Converter 获取转换器（供子类使用）
func (r *BaseRepository[T, M]) Converter() DomainConverter[T, M] {
	return r.converter
}

// HandleNotFound 统一处理记录不存在错误
func (r *BaseRepository[T, M]) HandleNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.ErrAggregateNotFound
	}
	return err
}

// CheckRowsAffected 检查影响行数
func (r *BaseRepository[T, M]) CheckRowsAffected(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrAggregateNotFound
	}
	return nil
}

// SaveEventsToOutbox 保存领域事件到 Outbox（同一事务中）
func (r *BaseRepository[T, M]) SaveEventsToOutbox(ctx context.Context, events []common.DomainEvent) error {
	for _, event := range events {
		if err := r.saveEventToOutbox(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// saveEventToOutbox 保存单个事件到 Outbox
func (r *BaseRepository[T, M]) saveEventToOutbox(ctx context.Context, event common.DomainEvent) error {
	// 序列化事件数据
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 序列化元数据
	var metadataJSON *string
	if meta := event.Metadata(); len(meta) > 0 {
		metadataBytes, _ := json.Marshal(meta)
		metadataStr := string(metadataBytes)
		metadataJSON = &metadataStr
	}

	now := time.Now()

	// 创建 Outbox 记录
	outboxEvent := &model.Outbox{
		ID:            now.UnixNano(),
		EventType:     event.EventName(),
		AggregateType: getAggregateType(event),
		AggregateID:   fmt.Sprintf("%v", event.AggregateID()),
		Payload:       string(payload),
		Metadata:      metadataJSON,
		OccurredAt:    &now,
		Processed:     false,
		RetryCount:    0,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	// 保存到数据库（在事务中执行）
	return r.db.WithContext(ctx).Create(outboxEvent).Error
}

// SaveEventsToStore 保存领域事件到 Event Store（永久存储，用于审计）
func (r *BaseRepository[T, M]) SaveEventsToStore(ctx context.Context, events []common.DomainEvent) error {
	for _, event := range events {
		if err := r.saveEventToStore(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// saveEventToStore 保存单个事件到 Event Store
func (r *BaseRepository[T, M]) saveEventToStore(ctx context.Context, event common.DomainEvent) error {
	// 序列化事件数据
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 序列化元数据
	var metadataJSON *string
	if meta := event.Metadata(); len(meta) > 0 {
		metadataBytes, _ := json.Marshal(meta)
		metadataStr := string(metadataBytes)
		metadataJSON = &metadataStr
	}

	now := time.Now()

	// 创建 DomainEvent 记录（用于事件溯源和审计，永久保存）
	daoEvent := &model.DomainEvent{
		ID:            now.UnixNano(),
		EventType:     event.EventName(),
		AggregateType: getAggregateType(event),
		AggregateID:   fmt.Sprintf("%v", event.AggregateID()),
		EventData:     string(eventData),
		Metadata:      metadataJSON,
		OccurredAt:    &now,
	}

	// 保存到数据库
	return r.db.WithContext(ctx).Create(daoEvent).Error
}

// getAggregateType 从事件中获取聚合根类型（通过类型推断）
func getAggregateType(event common.DomainEvent) string {
	// 尝试从事件名称推断聚合根类型
	// 例如：user.registered -> user
	eventName := event.EventName()
	if idx := strings.LastIndex(eventName, "."); idx != -1 {
		return eventName[:idx]
	}
	return "unknown"
}
