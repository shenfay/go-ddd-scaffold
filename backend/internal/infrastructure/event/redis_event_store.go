// Package event Redis Stream 事件存储实现
package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisEventStore Redis 事件存储实现
type RedisEventStore struct {
	client   *redis.Client
	streamKey string
}

// RedisEventStoreConfig Redis 事件存储配置
type RedisEventStoreConfig struct {
	StreamKey     string // Stream 名称，默认 "domain_events"
	MaxRetries    int    // 最大重试次数，默认 3
	RetryBaseDelay time.Duration // 重试基础延迟，默认 1 秒
}

// NewRedisEventStore 创建 Redis 事件存储
func NewRedisEventStore(client *redis.Client, config RedisEventStoreConfig) *RedisEventStore {
	if config.StreamKey == "" {
		config.StreamKey = "domain_events"
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryBaseDelay == 0 {
		config.RetryBaseDelay = time.Second
	}

	return &RedisEventStore{
		client:    client,
		streamKey: config.StreamKey,
	}
}

// StoredEvent 已存储的事件（带处理状态）
type StoredEvent struct {
	EventID      string                 `json:"eventId"`
	EventType    string                 `json:"eventType"`
	AggregateID  string                 `json:"aggregateId"`
	EventData    map[string]interface{} `json:"eventData"`
	Timestamp    time.Time              `json:"timestamp"`
	Status       string                 `json:"status"` // pending, processed, failed
	Attempt      int                    `json:"attempt"`
	LastError    string                 `json:"lastError,omitempty"`
	NextRetryAt  *time.Time             `json:"nextRetryAt,omitempty"`
}

const (
	EventStatusPending   = "pending"
	EventStatusProcessed = "processed"
	EventStatusFailed    = "failed"
)

// Store 保存事件到 Redis Stream
func (s *RedisEventStore) Store(ctx context.Context, event DomainEvent) error {
	if event == nil {
		return fmt.Errorf("事件不能为空")
	}

	// 序列化事件数据
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("序列化事件失败：%w", err)
	}

	stored := StoredEvent{
		EventID:      event.GetEventID(),
		EventType:    event.GetEventType(),
		AggregateID:  event.GetAggregateID().String(),
		EventData:    make(map[string]interface{}),
		Timestamp:    event.GetOccurredAt(),
		Status:       EventStatusPending,
		Attempt:      0,
		NextRetryAt:  nil,
	}

	// 将事件数据反序列化为 map（便于查询）
	if err := json.Unmarshal(eventData, &stored.EventData); err != nil {
		return fmt.Errorf("转换事件数据失败：%w", err)
	}

	// 序列化为 JSON 存入 Stream
	value, err := json.Marshal(stored)
	if err != nil {
		return fmt.Errorf("序列化存储事件失败：%w", err)
	}

	// 添加到 Redis Stream
	err = s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: s.streamKey,
		ID:     "*", // 让 Redis 自动生成 ID
		Values: map[string]interface{}{
			"data": string(value),
		},
	}).Err()

	if err != nil {
		return fmt.Errorf("写入 Redis Stream 失败：%w", err)
	}

	return nil
}

// GetPendingEvents 获取待处理的未消费事件
func (s *RedisEventStore) GetPendingEvents(ctx context.Context, limit int) ([]DomainEvent, error) {
	// 从 Stream 读取消息
	msgs, err := s.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{s.streamKey, "0"}, // 从头开始读
		Count:   int64(limit),
		Block:   0, // 不阻塞
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("读取 Redis Stream 失败：%w", err)
	}

	if len(msgs) == 0 || len(msgs[0].Messages) == 0 {
		return []DomainEvent{}, nil
	}

	var events []DomainEvent
	for _, msg := range msgs[0].Messages {
		dataStr := msg.Values["data"].(string)
		
		var stored StoredEvent
		if err := json.Unmarshal([]byte(dataStr), &stored); err != nil {
			continue // 跳过解析失败的消息
		}

		// 只返回 pending 状态的事件
		if stored.Status != EventStatusPending {
			continue
		}

		// 检查是否需要重试
		if stored.Attempt > 0 && stored.NextRetryAt != nil {
			if time.Now().Before(*stored.NextRetryAt) {
				continue // 还没到重试时间
			}
		}

		// 反序列化回具体事件类型（这里需要根据 EventType 动态创建）
		// 简化处理：返回通用事件包装器
		events = append(events, s.deserializeEvent(stored))
	}

	return events, nil
}

// MarkAsProcessed 标记事件已处理
func (s *RedisEventStore) MarkAsProcessed(ctx context.Context, eventID string) error {
	// 在 Stream 中查找并更新状态
	msgs, err := s.client.XRange(ctx, s.streamKey, "-", "+").Result()
	if err != nil {
		return fmt.Errorf("查询事件失败：%w", err)
	}

	for _, msg := range msgs {
		dataStr := msg.Values["data"].(string)
		var stored StoredEvent
		if err := json.Unmarshal([]byte(dataStr), &stored); err != nil {
			continue
		}

		if stored.EventID == eventID {
			stored.Status = EventStatusProcessed
			updatedData, _ := json.Marshal(stored)
			
			// 更新消息
			err := s.client.XAdd(ctx, &redis.XAddArgs{
				Stream: s.streamKey + "_processed",
				Values: map[string]interface{}{
					"data": string(updatedData),
				},
			}).Err()
			
			if err != nil {
				return fmt.Errorf("标记已处理失败：%w", err)
			}
			
			// 从原 Stream 删除（可选：也可以保留用于审计）
			s.client.XDel(ctx, s.streamKey, msg.ID).Err()
			return nil
		}
	}

	return fmt.Errorf("事件未找到：%s", eventID)
}

// MarkAsFailed 标记事件处理失败
func (s *RedisEventStore) MarkAsFailed(ctx context.Context, eventID string, errorMsg string) error {
	msgs, err := s.client.XRange(ctx, s.streamKey, "-", "+").Result()
	if err != nil {
		return fmt.Errorf("查询事件失败：%w", err)
	}

	for _, msg := range msgs {
		dataStr := msg.Values["data"].(string)
		var stored StoredEvent
		if err := json.Unmarshal([]byte(dataStr), &stored); err != nil {
			continue
		}

		if stored.EventID == eventID {
			stored.Attempt++
			stored.LastError = errorMsg
			stored.Status = EventStatusFailed
			
			// 计算下次重试时间（指数退避）
			retryPolicy := NewExponentialBackoffRetryPolicy(time.Second, time.Minute*5)
			if retryPolicy.ShouldRetry(stored.Attempt, 3) {
				delay := retryPolicy.GetDelay(stored.Attempt)
				nextRetry := time.Now().Add(delay)
				stored.NextRetryAt = &nextRetry
				stored.Status = EventStatusPending // 待重试
			}
			
			updatedData, _ := json.Marshal(stored)
			
			// 更新消息
			err := s.client.XAdd(ctx, &redis.XAddArgs{
				Stream: s.streamKey + "_failed",
				Values: map[string]interface{}{
					"data": string(updatedData),
				},
			}).Err()
			
			if err != nil {
				return fmt.Errorf("标记失败失败：%w", err)
			}
			
			// 从原 Stream 删除
			s.client.XDel(ctx, s.streamKey, msg.ID).Err()
			return nil
		}
	}

	return fmt.Errorf("事件未找到：%s", eventID)
}

// DeleteOldEvents 删除旧的已处理事件
func (s *RedisEventStore) DeleteOldEvents(ctx context.Context, before time.Time) error {
	// 清理 processed stream
	err := s.client.XTrim(ctx, &redis.XTrimArgs{
		Stream: s.streamKey + "_processed",
		MinID:  "",
		MaxLen: 1000, // 保留最近 1000 条
	}).Err()
	
	if err != nil {
		return fmt.Errorf("清理已处理事件失败：%w", err)
	}
	
	return nil
}

// deserializeEvent 反序列化事件（简化版本）
func (s *RedisEventStore) deserializeEvent(stored StoredEvent) DomainEvent {
	// 实际项目中需要根据 EventType 动态创建具体事件类型
	// 这里返回一个通用的事件包装器
	return &GenericDomainEvent{
		ID:          stored.EventID,
		EventType:   stored.EventType,
		AggregateID: uuid.MustParse(stored.AggregateID),
		Data:        stored.EventData,
		Timestamp:   stored.Timestamp,
	}
}

// GenericDomainEvent 通用领域事件（用于反序列化）
type GenericDomainEvent struct {
	ID          string                 `json:"id"`
	EventType   string                 `json:"eventType"`
	AggregateID uuid.UUID              `json:"aggregateId"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     int                    `json:"version"`
}

func (e *GenericDomainEvent) GetEventID() string           { return e.ID }
func (e *GenericDomainEvent) GetEventType() string         { return e.EventType }
func (e *GenericDomainEvent) GetAggregateID() uuid.UUID    { return e.AggregateID }
func (e *GenericDomainEvent) GetOccurredAt() time.Time     { return e.Timestamp }
func (e *GenericDomainEvent) GetData() map[string]interface{} { return e.Data }
func (e *GenericDomainEvent) GetVersion() int              { return e.Version }
