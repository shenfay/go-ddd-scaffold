// Package event 事件存储接口定义
package event

import (
	"context"
	"sync"
	"time"
)

// EventStore 事件存储接口（用于持久化）
type EventStore interface {
	// Store 保存事件到存储
	Store(ctx context.Context, event DomainEvent) error
	// GetPendingEvents 获取待处理的未消费事件
	GetPendingEvents(ctx context.Context, limit int) ([]DomainEvent, error)
	// MarkAsProcessed 标记事件已处理
	MarkAsProcessed(ctx context.Context, eventID string) error
	// MarkAsFailed 标记事件处理失败
	MarkAsFailed(ctx context.Context, eventID string, errorMsg string) error
	// DeleteOldEvents 删除旧的已处理事件（清理）
	DeleteOldEvents(ctx context.Context, before time.Time) error
}

// EventRetryPolicy 事件重试策略
type EventRetryPolicy interface {
	// ShouldRetry 判断是否应该重试
	ShouldRetry(attempt int, maxRetries int) bool
	// GetDelay 获取下次重试的延迟时间
	GetDelay(attempt int) time.Duration
}

// exponentialBackoffRetryPolicy 指数退避重试策略
type exponentialBackoffRetryPolicy struct {
	baseDelay time.Duration
	maxDelay  time.Duration
}

// NewExponentialBackoffRetryPolicy 创建指数退避重试策略
func NewExponentialBackoffRetryPolicy(baseDelay, maxDelay time.Duration) EventRetryPolicy {
	return &exponentialBackoffRetryPolicy{
		baseDelay: baseDelay,
		maxDelay:  maxDelay,
	}
}

// ShouldRetry 判断是否应该重试
func (p *exponentialBackoffRetryPolicy) ShouldRetry(attempt int, maxRetries int) bool {
	return attempt < maxRetries
}

// GetDelay 获取下次重试的延迟时间（指数增长）
func (p *exponentialBackoffRetryPolicy) GetDelay(attempt int) time.Duration {
	delay := p.baseDelay * time.Duration(1<<uint(attempt)) // 2^n * baseDelay
	if delay > p.maxDelay {
		return p.maxDelay
	}
	return delay
}

// InMemoryEventStore 内存事件存储实现（用于测试）
type InMemoryEventStore struct {
	mu        sync.RWMutex
	events    map[string][]DomainEvent // 按 aggregateID 字符串存储
	allEvents []DomainEvent            // 所有事件（用于时间范围查询）
}

// NewInMemoryEventStore 创建内存事件存储
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events:    make(map[string][]DomainEvent),
		allEvents: make([]DomainEvent, 0),
	}
}

// SaveEvent 保存事件
func (s *InMemoryEventStore) SaveEvent(ctx context.Context, event DomainEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	aggregateIDStr := event.GetAggregateID().String()
	s.events[aggregateIDStr] = append(s.events[aggregateIDStr], event)
	s.allEvents = append(s.allEvents, event)
	return nil
}

// LoadEvents 加载指定聚合根的事件
func (s *InMemoryEventStore) LoadEvents(ctx context.Context, aggregateID string) ([]DomainEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, exists := s.events[aggregateID]
	if !exists {
		return []DomainEvent{}, nil
	}

	result := make([]DomainEvent, len(events))
	copy(result, events)
	return result, nil
}

// GetEventCount 获取事件总数
func (s *InMemoryEventStore) GetEventCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.allEvents)
}

// LoadEventsInRange 按时间范围加载事件
func (s *InMemoryEventStore) LoadEventsInRange(ctx context.Context, start, end time.Time) ([]DomainEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]DomainEvent, 0)
	for _, event := range s.allEvents {
		occurredAt := event.GetOccurredAt()
		if (occurredAt.Equal(start) || occurredAt.After(start)) &&
			(occurredAt.Equal(end) || occurredAt.Before(end)) {
			result = append(result, event)
		}
	}

	return result, nil
}

// Store 保存事件到存储（实现 EventStore 接口）
func (s *InMemoryEventStore) Store(ctx context.Context, event DomainEvent) error {
	return s.SaveEvent(ctx, event)
}

// GetPendingEvents 获取待处理的未消费事件
func (s *InMemoryEventStore) GetPendingEvents(ctx context.Context, limit int) ([]DomainEvent, error) {
	// 简化实现：返回所有事件
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit > len(s.allEvents) {
		limit = len(s.allEvents)
	}
	result := make([]DomainEvent, limit)
	copy(result, s.allEvents[:limit])
	return result, nil
}

// MarkAsProcessed 标记事件已处理
func (s *InMemoryEventStore) MarkAsProcessed(ctx context.Context, eventID string) error {
	// 内存实现中无需额外操作
	return nil
}

// MarkAsFailed 标记事件处理失败
func (s *InMemoryEventStore) MarkAsFailed(ctx context.Context, eventID string, errorMsg string) error {
	// 内存实现中无需额外操作
	return nil
}

// DeleteOldEvents 删除旧的已处理事件
func (s *InMemoryEventStore) DeleteOldEvents(ctx context.Context, before time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 过滤掉旧的事件
	filtered := make([]DomainEvent, 0)
	for _, event := range s.allEvents {
		if event.GetOccurredAt().After(before) {
			filtered = append(filtered, event)
		}
	}
	s.allEvents = filtered

	// 同时清理 events map 中的旧事件
	for aggregateID, events := range s.events {
		filteredEvents := make([]DomainEvent, 0)
		for _, event := range events {
			if event.GetOccurredAt().After(before) {
				filteredEvents = append(filteredEvents, event)
			}
		}
		if len(filteredEvents) == 0 {
			delete(s.events, aggregateID)
		} else {
			s.events[aggregateID] = filteredEvents
		}
	}

	return nil
}
