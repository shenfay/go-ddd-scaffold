package domain_event

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	task_queue "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/taskqueue"
	"go.uber.org/zap"
)

// EventPublisherAdapter 事件发布器适配器
// 兼容旧的 kernel.EventPublisher 接口，内部实现新的三重写逻辑
type EventPublisherAdapter struct {
	query         *dao.Query
	taskPublisher *task_queue.Publisher
	logger        *zap.Logger
}

// NewEventPublisherAdapter 创建事件发布器适配器
func NewEventPublisherAdapter(
	query *dao.Query,
	taskPublisher *task_queue.Publisher,
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

// Publish 发布领域事件（实现 kernel.EventPublisher 接口）
// 同时记录 ActivityLog、EventLog 并发布到 Asynq 队列
func (a *EventPublisherAdapter) Publish(ctx context.Context, event kernel.DomainEvent) error {
	a.logger.Debug("Publishing event",
		zap.String("event_type", event.EventName()),
		zap.Any("aggregate_id", event.AggregateID()),
	)

	// 1. 记录活动日志（ActivityLog）
	if err := a.saveActivityLog(ctx, event); err != nil {
		a.logger.Error("Failed to save activity log", zap.Error(err))
		// 不阻断后续流程
	}

	// 2. 记录事件日志（EventLog）
	if err := a.saveEventLog(ctx, event); err != nil {
		a.logger.Error("Failed to save event log", zap.Error(err))
		// 不阻断后续流程
	}

	// 3. 发布到 Asynq 队列
	if err := a.publishToQueue(ctx, event); err != nil {
		a.logger.Error("Failed to publish to queue", zap.Error(err))
		return err
	}

	return nil
}

// saveActivityLog 保存活动日志
func (a *EventPublisherAdapter) saveActivityLog(ctx context.Context, event kernel.DomainEvent) error {
	// 获取用户 ID（如果有的话）
	var userID int64
	if id, ok := event.AggregateID().(int64); ok {
		userID = id
	}

	// 根据事件类型创建不同的活动日志
	action := a.eventTypeToAction(event.EventName())
	if action == "" {
		// 忽略不需要记录为活动的系统事件
		return nil
	}

	now := time.Now()
	metadata := event.Metadata()
	metadataJSON, _ := json.Marshal(metadata)
	status := int16(0) // Success
	metadataStr := string(metadataJSON)

	daoModel := &model.ActivityLog{
		ID:         now.UnixNano(), // 使用时间戳作为临时 ID
		UserID:     userID,
		Action:     string(action),
		Status:     &status,
		Metadata:   &metadataStr,
		OccurredAt: event.OccurredOn(),
		CreatedAt:  &now,
	}

	return a.query.ActivityLog.WithContext(ctx).Create(daoModel)
}

// saveEventLog 保存事件日志（轻量级，用于事件溯源）
func (a *EventPublisherAdapter) saveEventLog(ctx context.Context, event kernel.DomainEvent) error {
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
	daoModel := &model.EventLog{
		AggregateID:   a.aggregateIDToString(event.AggregateID()),
		AggregateType: aggregateType,
		EventType:     event.EventName(),
		EventData:     string(eventDataJSON),
		OccurredAt:    event.OccurredOn(),
		CreatedAt:     &now,
	}

	return a.query.EventLog.WithContext(ctx).Create(daoModel)
}

// publishToQueue 发布到 Asynq 队列
func (a *EventPublisherAdapter) publishToQueue(ctx context.Context, event kernel.DomainEvent) error {
	// 序列化事件数据
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// 创建任务负载
	payload := task_queue.DomainEventPayload{
		AggregateID:   a.aggregateIDToString(event.AggregateID()),
		AggregateType: a.inferAggregateType(event.EventName()),
		EventType:     event.EventName(),
		EventVersion:  int32(event.Version()),
		EventData:     eventData,
		OccurredOn:    event.OccurredOn().Format(time.RFC3339),
		Metadata:      a.convertMetadata(event.Metadata()),
	}

	// 根据事件类型选择队列
	queue := a.getQueueForEvent(event.EventName())

	// 发布到队列
	return a.taskPublisher.PublishDomainEvent(ctx, payload, queue)
}

// eventTypeToAction 将事件类型转换为活动类型
func (a *EventPublisherAdapter) eventTypeToAction(eventType string) string {
	switch eventType {
	case "UserRegistered":
		return "USER_REGISTERED"
	case "UserLoggedIn":
		return "USER_LOGIN"
	case "UserLoggedOut":
		return "USER_LOGOUT"
	case "UserActivated":
		return "USER_ACTIVATED"
	case "UserDeactivated":
		return "USER_DEACTIVATED"
	case "UserLocked":
		return "USER_LOCKED"
	case "UserUnlocked":
		return "USER_UNLOCKED"
	case "UserPasswordChanged":
		return "USER_PASSWORD_CHANGED"
	case "UserEmailChanged":
		return "USER_EMAIL_CHANGED"
	case "UserProfileUpdated":
		return "USER_PROFILE_UPDATED"
	default:
		return "" // 不需要记录为活动的事件
	}
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

// convertMetadata 转换元数据类型
func (a *EventPublisherAdapter) convertMetadata(meta map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range meta {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}
	return result
}

// eventToMap 将事件转换为 map（用于存储）
func (a *EventPublisherAdapter) eventToMap(event kernel.DomainEvent) (map[string]any, error) {
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

// getQueueForEvent 根据事件类型返回相应的队列
func (a *EventPublisherAdapter) getQueueForEvent(eventType string) string {
	// 高优先级事件
	criticalEvents := map[string]bool{
		"UserRegistered":   true, // 用户注册
		"PaymentCompleted": true, // 支付完成
	}

	if criticalEvents[eventType] {
		return "critical"
	}

	// 低优先级事件
	lowEvents := map[string]bool{
		"UserLoggedIn":     true, // 用户登录（日志类）
		"NotificationSent": true, // 通知发送
	}

	if lowEvents[eventType] {
		return "low"
	}

	return "default"
}
