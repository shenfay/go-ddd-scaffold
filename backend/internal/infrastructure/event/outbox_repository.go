// Package event 提供事务性发件箱模式支持
package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OutboxEvent 发件箱事件（数据库表模型）
type OutboxEvent struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	EventType   string     `gorm:"size:100;not null;index:idx_status_type" json:"eventType"`
	AggregateID uuid.UUID  `gorm:"type:uuid;not null;index" json:"aggregateId"`
	Payload     string     `gorm:"type:text;not null" json:"payload"` // JSON 序列化的事件数据
	Status      string     `gorm:"size:20;not null;default:'pending';index:idx_status_type" json:"status"`
	RetryCount  int        `gorm:"not null;default:0" json:"retryCount"`
	MaxRetries  int        `gorm:"not null;default:3" json:"maxRetries"`
	NextRetryAt *time.Time `json:"nextRetryAt,omitempty"`
	LastError   *string    `gorm:"type:text" json:"lastError,omitempty"`
	ProcessedAt *time.Time `json:"processedAt,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (OutboxEvent) TableName() string {
	return "domain_events"
}

// OutboxEventStatus 事件状态
type OutboxEventStatus string

const (
	OutboxStatusPending    OutboxEventStatus = "pending"    // 待发布
	OutboxStatusProcessing OutboxEventStatus = "processing" // 处理中
	OutboxStatusPublished  OutboxEventStatus = "published"  // 已发布
	OutboxStatusFailed     OutboxEventStatus = "failed"     // 失败
)

// OutboxRepository 发件箱仓储接口
type OutboxRepository interface {
	// Save 保存事件到发件箱（在事务中调用）
	Save(ctx context.Context, tx *gorm.DB, event DomainEvent) error
	// GetPendingEvents 获取待发布的事件
	GetPendingEvents(ctx context.Context, limit int) ([]*OutboxEvent, error)
	// MarkAsPublished 标记事件已发布
	MarkAsPublished(ctx context.Context, eventID uuid.UUID) error
	// MarkAsFailed 标记事件处理失败
	MarkAsFailed(ctx context.Context, eventID uuid.UUID, errorMsg string) error
	// MarkAsProcessing 标记事件正在处理
	MarkAsProcessing(ctx context.Context, eventID uuid.UUID) error
	// DeleteOldEvents 删除旧事件（用于清理）
	DeleteOldEvents(ctx context.Context, before time.Time, statuses ...OutboxEventStatus) error
}

// gormOutboxRepository GORM 发件箱仓储实现
type gormOutboxRepository struct {
	db *gorm.DB
}

// NewGormOutboxRepository 创建发件箱仓储实例
func NewGormOutboxRepository(db *gorm.DB) OutboxRepository {
	return &gormOutboxRepository{db: db}
}

// Save 保存事件到发件箱（必须在事务中调用）
func (r *gormOutboxRepository) Save(ctx context.Context, tx *gorm.DB, event DomainEvent) error {
	if event == nil {
		return ErrNilEvent
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return ErrSerializationFailed.WithCause(err)
	}

	now := time.Now()
	outboxEvent := &OutboxEvent{
		ID:          uuid.New(),
		EventType:   event.GetEventType(),
		AggregateID: event.GetAggregateID(),
		Payload:     string(payload),
		Status:      string(OutboxStatusPending),
		RetryCount:  0,
		MaxRetries:  3,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 使用传入的事务保存
	result := tx.WithContext(ctx).Create(outboxEvent)
	return result.Error
}

// GetPendingEvents 获取待发布的事件
func (r *gormOutboxRepository) GetPendingEvents(ctx context.Context, limit int) ([]*OutboxEvent, error) {
	var events []*OutboxEvent

	now := time.Now()
	query := r.db.WithContext(ctx).
		Where("status = ?", OutboxStatusPending).
		// 只获取到了重试时间的或者不需要重试的
		Where("next_retry_at IS NULL OR next_retry_at <= ?", now).
		Order("created_at ASC").
		Limit(limit)

	err := query.Find(&events).Error
	if err != nil {
		return nil, ErrDatabaseOperation.WithCause(err)
	}

	return events, nil
}

// MarkAsPublished 标记事件已发布
func (r *gormOutboxRepository) MarkAsPublished(ctx context.Context, eventID uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&OutboxEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"status":       OutboxStatusPublished,
			"processed_at": now,
			"updated_at":   now,
		})

	if result.Error != nil {
		return ErrDatabaseOperation.WithCause(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

// MarkAsProcessing 标记事件正在处理
func (r *gormOutboxRepository) MarkAsProcessing(ctx context.Context, eventID uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&OutboxEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"status":     OutboxStatusProcessing,
			"updated_at": now,
		})

	if result.Error != nil {
		return ErrDatabaseOperation.WithCause(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

// MarkAsFailed 标记事件处理失败
func (r *gormOutboxRepository) MarkAsFailed(ctx context.Context, eventID uuid.UUID, errorMsg string) error {
	event, err := r.getEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	now := time.Now()
	event.RetryCount++
	event.LastError = &errorMsg
	event.UpdatedAt = now

	// 判断是否还需要重试
	if event.RetryCount >= event.MaxRetries {
		event.Status = string(OutboxStatusFailed)
	} else {
		// 计算下次重试时间（指数退避）
		delay := calculateExponentialBackoff(event.RetryCount)
		nextRetryAt := now.Add(delay)
		event.NextRetryAt = &nextRetryAt
		event.Status = string(OutboxStatusPending)
	}

	result := r.db.WithContext(ctx).Save(event)
	if result.Error != nil {
		return ErrDatabaseOperation.WithCause(result.Error)
	}

	return nil
}

// DeleteOldEvents 删除旧事件
func (r *gormOutboxRepository) DeleteOldEvents(ctx context.Context, before time.Time, statuses ...OutboxEventStatus) error {
	query := r.db.WithContext(ctx).
		Where("created_at < ?", before)

	if len(statuses) > 0 {
		statusStrings := make([]string, len(statuses))
		for i, s := range statuses {
			statusStrings[i] = string(s)
		}
		query = query.Where("status IN (?)", statusStrings)
	}

	result := query.Delete(&OutboxEvent{})
	if result.Error != nil {
		return ErrDatabaseOperation.WithCause(result.Error)
	}

	return nil
}

// getEventByID 根据 ID 获取事件
func (r *gormOutboxRepository) getEventByID(ctx context.Context, eventID uuid.UUID) (*OutboxEvent, error) {
	var event OutboxEvent
	err := r.db.WithContext(ctx).First(&event, eventID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEventNotFound
		}
		return nil, ErrDatabaseOperation.WithCause(err)
	}
	return &event, nil
}

// calculateExponentialBackoff 计算指数退避延迟
func calculateExponentialBackoff(retryCount int) time.Duration {
	baseDelay := time.Second
	maxDelay := time.Hour

	// 指数增长：1s, 2s, 4s, 8s, 16s...
	delay := baseDelay * time.Duration(1<<uint(retryCount-1))
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}
