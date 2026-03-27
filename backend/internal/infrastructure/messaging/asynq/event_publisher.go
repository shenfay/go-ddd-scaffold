package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
	"go.uber.org/zap"
)

// EventPublisherAdapter 事件发布器适配器
// 兼容旧的 common.EventPublisher 接口，内部实现新的三重写逻辑
type EventPublisherAdapter struct {
	query         *dao.Query
	taskPublisher *EventPublisher
	logger        *zap.Logger
}

// NewEventPublisherAdapter 创建事件发布器适配器
func NewEventPublisherAdapter(
	query *dao.Query,
	taskPublisher *EventPublisher,
	logger *zap.Logger,
) *EventPublisherAdapter {
	if logger == nil {
		logger = zap.L().Named("event_publisher")
	}
	return &EventPublisherAdapter{
		query:         query,
		taskPublisher: taskPublisher,
		logger:        logger,
	}
}

// Publish 发布领域事件（实现 common.EventPublisher 接口）
// 记录 DomainEvent 和 Outbox（支持事务）
// 注意：ActivityLog 应由应用层直接写入，不通过事件发布器
func (a *EventPublisherAdapter) Publish(ctx context.Context, event common.DomainEvent) error {
	a.logger.Debug("Publishing event",
		zap.String("event_type", event.EventName()),
		zap.Any("aggregate_id", event.AggregateID()),
	)

	// 1. 记录事件日志（DomainEvent）- 用于事件溯源
	if err := a.saveEventLog(ctx, event); err != nil {
		a.logger.Error("Failed to save event log", zap.Error(err))
		return err // 事务中失败会回滚
	}

	// 2. 记录到 Outbox 表（Outbox Pattern，保证事务一致性）
	if err := a.saveToOutbox(ctx, event); err != nil {
		a.logger.Error("Failed to save to outbox", zap.Error(err))
		return err // 事务中失败会回滚
	}

	return nil
}

// saveEventLog 保存事件日志（轻量级，用于事件溯源）
func (a *EventPublisherAdapter) saveEventLog(ctx context.Context, event common.DomainEvent) error {
	// 将事件转换为 map[string]any
	eventData, err := a.eventToMap(event)
	if err != nil {
		return err
	}

	// 确定聚合类型
	aggregateType := a.inferAggregateType(event.EventName())

	// 序列化事件数据
	eventDataJSON, _ := json.Marshal(eventData)

	now := time.Now()
	daoModel := &model.DomainEvent{
		ID:            idgen.Generate(), // 生成雪花 ID
		AggregateID:   a.aggregateIDToString(event.AggregateID()),
		AggregateType: aggregateType,
		EventType:     event.EventName(),
		EventData:     string(eventDataJSON),
		OccurredAt:    &now,
		CreatedAt:     &now,
	}

	return a.query.DomainEvent.WithContext(ctx).Create(daoModel)
}

// inferAggregateType 从事件类型推断聚合类型
func (a *EventPublisherAdapter) inferAggregateType(eventType string) string {
	// 简单规则：去掉 "Event" 后缀（如果有）
	// UserCreatedEvent -> User
	// OrderPaidEvent -> Order
	if len(eventType) > 5 && eventType[len(eventType)-5:] == "Event" {
		return eventType[:len(eventType)-5]
	}

	// 如果没有 Event 后缀，尝试其他常见模式
	// UserRegistered -> User (去掉 Registered)
	// UserLoggedIn -> User (去掉 LoggedIn/LoggedOut)
	// UserCreated -> User (去掉 Created)
	suffixes := []string{
		"Registered", "LoggedIn", "LoggedOut", "Activated", "Deactivated",
		"Created", "Updated", "Deleted", "Changed", "Verified",
	}

	for _, suffix := range suffixes {
		if len(eventType) > len(suffix) && eventType[len(eventType)-len(suffix):] == suffix {
			return eventType[:len(eventType)-len(suffix)]
		}
	}

	// 如果都不匹配，返回原始事件类型作为备用
	return eventType
}

// aggregateIDToString 将聚合根 ID 转换为字符串
func (a *EventPublisherAdapter) aggregateIDToString(id interface{}) string {
	if id == nil {
		return ""
	}

	switch v := id.(type) {
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	default:
		// 其他类型转为字符串
		return fmt.Sprintf("%v", v)
	}
}

// eventToMap 将事件转换为 map（用于存储）
func (a *EventPublisherAdapter) eventToMap(event common.DomainEvent) (map[string]any, error) {
	// 使用 JSON 序列化再反序列化的方式
	data, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// saveToOutbox 保存事件到 Outbox 表（用于 Outbox Pattern，支持事务）
func (a *EventPublisherAdapter) saveToOutbox(ctx context.Context, event common.DomainEvent) error {
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
	daoModel := &model.Outbox{
		ID:            idgen.Generate(), // 生成雪花 ID
		EventType:     event.EventName(),
		AggregateType: a.inferAggregateType(event.EventName()),
		AggregateID:   a.aggregateIDToString(event.AggregateID()),
		Payload:       string(payload),
		Metadata:      metadataJSON,
		OccurredAt:    &now,
		Processed:     false,
		RetryCount:    0,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	return a.query.Outbox.WithContext(ctx).Create(daoModel)
}
