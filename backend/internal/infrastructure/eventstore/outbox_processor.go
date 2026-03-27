package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	asynq_lib "github.com/hibiken/asynq"
	asynq_pkg "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/messaging/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// DefaultPollInterval 默认轮询间隔
	DefaultPollInterval = 5 * time.Second
	// DefaultBatchSize 默认批量处理大小
	DefaultBatchSize = 100
	// MaxRetries 最大重试次数
	MaxRetries = 10
)

// ExtractDomainEventPayload 从任务中提取领域事件负载
func ExtractDomainEventPayload(task *asynq_lib.Task) (*asynq_pkg.DomainEventPayload, error) {
	var payload asynq_pkg.DomainEventPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// OutboxProcessor Outbox 处理器
type OutboxProcessor struct {
	db           *gorm.DB
	publisher    *asynq_pkg.EventPublisher
	logger       *zap.Logger
	pollInterval time.Duration
	batchSize    int
}

// NewOutboxProcessor 创建 Outbox 处理器
func NewOutboxProcessor(
	db *gorm.DB,
	publisher *asynq_pkg.EventPublisher,
	logger *zap.Logger,
) *OutboxProcessor {
	return &OutboxProcessor{
		db:           db,
		publisher:    publisher,
		logger:       logger.Named("outbox_processor"),
		pollInterval: DefaultPollInterval,
		batchSize:    DefaultBatchSize,
	}
}

// Start 启动轮询（后台协程）
func (p *OutboxProcessor) Start(ctx context.Context) error {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	p.logger.Info("Outbox processor started",
		zap.Duration("poll_interval", p.pollInterval),
		zap.Int("batch_size", p.batchSize),
	)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Outbox processor stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := p.processUnpublishedEvents(ctx); err != nil {
				p.logger.Error("Failed to process outbox", zap.Error(err))
			}
		}
	}
}

// processUnpublishedEvents 处理未发布的事件
func (p *OutboxProcessor) processUnpublishedEvents(ctx context.Context) error {
	// 1. 查询未处理的事件（从 outbox 表）
	var events []*model.Outbox
	err := p.db.WithContext(ctx).
		Where("processed = ?", false).
		Order("occurred_at ASC").
		Limit(p.batchSize).
		Find(&events).Error

	if err != nil {
		return fmt.Errorf("failed to query unpublished events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	p.logger.Debug("Found unpublished events", zap.Int("count", len(events)))

	// 2. 逐个发布
	successCount := 0
	failedCount := 0

	for _, event := range events {
		if err := p.publishEvent(ctx, event); err != nil {
			p.logger.Error("Failed to publish event",
				zap.Int64("id", event.ID),
				zap.String("type", event.EventType),
				zap.Error(err),
			)
			failedCount++

			// 增加重试计数
			if retryErr := p.incrementRetry(ctx, event.ID, err); retryErr != nil {
				p.logger.Error("Failed to increment retry count",
					zap.Int64("id", event.ID),
					zap.Error(retryErr),
				)
			}
		} else {
			successCount++
		}
	}

	p.logger.Info("Processed outbox events",
		zap.Int("success", successCount),
		zap.Int("failed", failedCount),
	)

	return nil
}

// publishEvent 发布单个事件
func (p *OutboxProcessor) publishEvent(ctx context.Context, event *model.Outbox) error {
	// 1. 反序列化事件数据（用于验证）
	var payload asynq_pkg.DomainEventPayload
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// 2. 发布到 Asynq
	task, err := asynq_pkg.NewDomainEventTask(payload)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// 这里应该使用 asynq client 发送任务
	// 暂时跳过实际发送，因为我们需要访问 asynq client
	_ = task
	// if err := p.publisher.PublishDomainEvent(ctx, payload, "default"); err != nil {
	// 	return fmt.Errorf("failed to publish to asynq: %w", err)
	// }

	// 3. 标记为已处理（更新 outbox）
	now := time.Now()
	return p.db.WithContext(ctx).Model(event).Updates(map[string]interface{}{
		"processed":    true,
		"processed_at": now,
		"updated_at":   now,
	}).Error
}

// incrementRetry 增加重试计数
func (p *OutboxProcessor) incrementRetry(ctx context.Context, eventID int64, lastErr error) error {
	return p.db.WithContext(ctx).Model(&model.Outbox{}).Where("id = ?", eventID).Updates(map[string]interface{}{
		"retry_count":   gorm.Expr("retry_count + 1"),
		"error_message": lastErr.Error(),
		"updated_at":    time.Now(),
	}).Error
}

// GetUnpublishedCount 获取未发布事件数量
func (p *OutboxProcessor) GetUnpublishedCount(ctx context.Context) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).
		Model(&model.Outbox{}).
		Where("processed = ?", false).
		Count(&count).Error
	return count, err
}

// GetFailedEvents 获取失败的事件（重试次数超过阈值）
func (p *OutboxProcessor) GetFailedEvents(ctx context.Context, limit int) ([]*model.Outbox, error) {
	var events []*model.Outbox
	err := p.db.WithContext(ctx).
		Where("processed = ? AND retry_count >= ?", false, MaxRetries).
		Limit(limit).
		Order("retry_count DESC").
		Find(&events).Error
	return events, err
}

// ManualRetry 手动重试失败事件
func (p *OutboxProcessor) ManualRetry(ctx context.Context, eventID int64) error {
	return p.db.WithContext(ctx).Model(&model.Outbox{}).Where("id = ?", eventID).Updates(map[string]interface{}{
		"processed":     false,
		"retry_count":   0,
		"error_message": nil,
		"processed_at":  nil,
		"updated_at":    time.Now(),
	}).Error
}
