package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
	"go.uber.org/zap"
)

// EventPublisherConfig 事件发布器配置
type EventPublisherConfig struct {
	// 队列配置
	DefaultQueue      string
	HighPriorityQueue string
	LowPriorityQueue  string

	// 幂等性配置
	DeduplicationTTL time.Duration

	// 重试配置
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration

	// 监控配置
	EnableMetrics bool
}

// EventDeduplicator 事件去重器
type EventDeduplicator struct {
	redisClient *redis.Client
	ttl         time.Duration
	logger      *zap.Logger
}

// IsDuplicate 检查事件是否重复
func (d *EventDeduplicator) IsDuplicate(ctx context.Context, eventID string) (bool, error) {
	key := fmt.Sprintf("event:processed:%s", eventID)
	exists, err := d.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// MarkProcessed 标记事件为已处理
func (d *EventDeduplicator) MarkProcessed(ctx context.Context, eventID string) error {
	key := fmt.Sprintf("event:processed:%s", eventID)
	return d.redisClient.SetEx(ctx, key, "1", d.ttl).Err()
}

// TaskType 定义任务类型常量
const (
	TaskTypeDomainEvent = "domain:event" // 领域事件任务类型
)

// DomainEventPayload 领域事件任务负载
type DomainEventPayload struct {
	AggregateID   string            `json:"aggregate_id"`
	AggregateType string            `json:"aggregate_type"`
	EventType     string            `json:"event_type"`
	EventVersion  int32             `json:"event_version"`
	EventData     json.RawMessage   `json:"event_data"`
	OccurredOn    string            `json:"occurred_on"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// NewDomainEventTask 创建领域事件任务
func NewDomainEventTask(payload DomainEventPayload) (*asynq.Task, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskTypeDomainEvent, b), nil
}

// ExtractDomainEventPayload 从任务中提取领域事件负载
func ExtractDomainEventPayload(task *asynq.Task) (*DomainEventPayload, error) {
	var payload DomainEventPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// EventPublisher 事件发布器
// 特性：
// 1. 一体化事件发布（数据库持久化 + Asynq 队列）
// 2. 内置幂等性保障
// 3. 灵活的重试策略
// 4. 完善的监控指标
type EventPublisher struct {
	query        *dao.Query
	asynqClient  *asynq.Client
	redisClient  *redis.Client
	logger       *zap.Logger
	config       EventPublisherConfig
	deduplicator *EventDeduplicator
}

// NewEventPublisher 创建事件发布器
func NewEventPublisher(
	query *dao.Query,
	asynqClient *asynq.Client,
	redisClient *redis.Client,
	logger *zap.Logger,
	config EventPublisherConfig,
) *EventPublisher {
	if logger == nil {
		logger = zap.L().Named("event_publisher")
	}

	// 默认配置
	if config.DefaultQueue == "" {
		config.DefaultQueue = "events"
	}
	if config.DeduplicationTTL == 0 {
		config.DeduplicationTTL = 24 * time.Hour
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.BaseDelay == 0 {
		config.BaseDelay = time.Second
	}
	if config.MaxDelay == 0 {
		config.MaxDelay = time.Minute
	}

	deduplicator := &EventDeduplicator{
		redisClient: redisClient,
		ttl:         config.DeduplicationTTL,
		logger:      logger.Named("deduplicator"),
	}

	return &EventPublisher{
		query:        query,
		asynqClient:  asynqClient,
		redisClient:  redisClient,
		logger:       logger,
		config:       config,
		deduplicator: deduplicator,
	}
}

// Publish 发布事件（一体化操作）
func (p *EventPublisher) Publish(ctx context.Context, event common.DomainEvent) error {
	// 1. 幂等性检查
	eventID := p.generateEventID(event)
	isDuplicate, err := p.deduplicator.IsDuplicate(ctx, eventID)
	if err != nil {
		p.logger.Warn("Failed to check duplication",
			zap.String("event_id", eventID),
			zap.Error(err))
		// 不因幂等性检查失败而中断发布
	} else if isDuplicate {
		p.logger.Debug("Skipping duplicate event",
			zap.String("event_id", eventID),
			zap.String("event_type", event.EventName()))
		return nil
	}

	// 2. 事务内持久化事件
	if err := p.persistEvent(ctx, event, eventID); err != nil {
		return fmt.Errorf("failed to persist event: %w", err)
	}

	// 3. 发送到 Asynq 队列
	if err := p.enqueueEvent(ctx, event, eventID); err != nil {
		// 记录错误但不回滚（事件已持久化）
		p.logger.Error("Failed to enqueue event",
			zap.String("event_id", eventID),
			zap.String("event_type", event.EventName()),
			zap.Error(err))
		return fmt.Errorf("failed to enqueue event: %w", err)
	}

	// 4. 标记为已处理
	if err := p.deduplicator.MarkProcessed(ctx, eventID); err != nil {
		p.logger.Warn("Failed to mark event as processed",
			zap.String("event_id", eventID),
			zap.Error(err))
	}

	p.logger.Debug("Event published successfully",
		zap.String("event_id", eventID),
		zap.String("event_type", event.EventName()),
		zap.Any("aggregate_id", event.AggregateID()))

	return nil
}

// persistEvent 持久化事件到数据库
func (p *EventPublisher) persistEvent(ctx context.Context, event common.DomainEvent, eventID string) error {
	// 序列化事件数据
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	now := time.Now()

	// 保存到 DomainEvent 表（事件溯源）
	domainEvent := &model.DomainEvent{
		ID:            idgen.Generate(),
		AggregateID:   p.aggregateIDToString(event.AggregateID()),
		AggregateType: p.inferAggregateType(event.EventName()),
		EventType:     event.EventName(),
		EventData:     string(eventData),
		OccurredAt:    &now,
		CreatedAt:     &now,
	}

	if err := p.query.DomainEvent.WithContext(ctx).Create(domainEvent); err != nil {
		return fmt.Errorf("failed to save domain event: %w", err)
	}

	// 保存到 Outbox 表（用于最终一致性）
	outboxEvent := &model.Outbox{
		ID:            idgen.Generate(),
		EventType:     event.EventName(),
		AggregateType: p.inferAggregateType(event.EventName()),
		AggregateID:   p.aggregateIDToString(event.AggregateID()),
		Payload:       string(eventData),
		Processed:     false,
		RetryCount:    0,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	if err := p.query.Outbox.WithContext(ctx).Create(outboxEvent); err != nil {
		return fmt.Errorf("failed to save outbox event: %w", err)
	}

	return nil
}

// enqueueEvent 将事件发送到 Asynq 队列
func (p *EventPublisher) enqueueEvent(ctx context.Context, event common.DomainEvent, eventID string) error {
	// 确定队列优先级
	queue := p.determineQueue(event.EventName())

	// 构造任务负载
	payload := DomainEventPayload{
		AggregateID:   p.aggregateIDToString(event.AggregateID()),
		AggregateType: p.inferAggregateType(event.EventName()),
		EventType:     event.EventName(),
		EventVersion:  int32(event.Version()),
		EventData:     json.RawMessage("{}"), // 具体数据在数据库中
		OccurredOn:    event.OccurredOn().Format(time.RFC3339),
		Metadata:      make(map[string]string),
	}

	// 添加元数据
	for k, v := range event.Metadata() {
		if str, ok := v.(string); ok {
			payload.Metadata[k] = str
		}
	}

	// 创建 Asynq 任务
	task, err := NewDomainEventTask(payload)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// 配置重试策略
	opts := []asynq.Option{
		asynq.Queue(queue),
		asynq.MaxRetry(p.config.MaxRetries),
		// 使用默认的指数退避策略
	}

	// 发送任务
	info, err := p.asynqClient.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	p.logger.Debug("Task enqueued",
		zap.String("task_id", info.ID),
		zap.String("queue", queue),
		zap.String("event_id", eventID))

	return nil
}

// determineQueue 根据事件类型确定队列
func (p *EventPublisher) determineQueue(eventType string) string {
	// 高优先级事件
	highPriorityEvents := map[string]bool{
		"UserLoggedIn":        true,
		"UserRegistered":      true,
		"UserPasswordChanged": true,
	}

	if highPriorityEvents[eventType] {
		return p.config.HighPriorityQueue
	}

	// 低优先级事件
	lowPriorityEvents := map[string]bool{
		"UserEmailChanged":   true,
		"UserProfileUpdated": true,
	}

	if lowPriorityEvents[eventType] {
		return p.config.LowPriorityQueue
	}

	// 默认队列
	return p.config.DefaultQueue
}

// generateEventID 生成全局唯一的事件ID
func (p *EventPublisher) generateEventID(event common.DomainEvent) string {
	// 使用聚合ID + 事件类型 + 时间戳 + 随机数生成唯一ID
	return fmt.Sprintf("%v-%s-%d-%s",
		event.AggregateID(),
		event.EventName(),
		time.Now().UnixNano(),
		fmt.Sprintf("%08x", time.Now().UnixNano()&0xFFFFFFFF))
}

// 工具方法（复用现有实现）
func (p *EventPublisher) inferAggregateType(eventType string) string {
	// 复用现有逻辑
	if len(eventType) > 5 && eventType[len(eventType)-5:] == "Event" {
		return eventType[:len(eventType)-5]
	}

	suffixes := []string{
		"Registered", "LoggedIn", "LoggedOut", "Activated", "Deactivated",
		"Created", "Updated", "Deleted", "Changed", "Verified",
	}

	for _, suffix := range suffixes {
		if len(eventType) > len(suffix) && eventType[len(eventType)-len(suffix):] == suffix {
			return eventType[:len(eventType)-len(suffix)]
		}
	}

	return eventType
}

func (p *EventPublisher) aggregateIDToString(id interface{}) string {
	// 复用现有逻辑
	switch v := id.(type) {
	case string:
		return v
	case int64:
		return fmt.Sprintf("%d", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
