package event

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// StoredEvent 已存储的事件
type StoredEvent struct {
	ID          string                 `json:"id"`
	EventType   string                 `json:"eventType"`
	AggregateID string                 `json:"aggregateId"`
	EventData   DomainEvent            `json:"eventData"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// EventStore 事件存储接口
type EventStore interface {
	// SaveEvent 保存事件
	SaveEvent(ctx context.Context, event DomainEvent) error

	// LoadEvents 加载指定聚合的所有事件
	LoadEvents(ctx context.Context, aggregateID string) ([]*StoredEvent, error)

	// LoadEventsByType 加载指定类型的事件
	LoadEventsByType(ctx context.Context, eventType string) ([]*StoredEvent, error)

	// LoadEventsInRange 加载指定时间范围内的事件
	LoadEventsInRange(ctx context.Context, start, end time.Time) ([]*StoredEvent, error)
}

// InMemoryEventStore 内存事件存储（用于开发和测试）
type InMemoryEventStore struct {
	events []StoredEvent
	mutex  sync.RWMutex
}

// NewInMemoryEventStore 创建内存事件存储
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make([]StoredEvent, 0),
	}
}

// SaveEvent 保存事件到内存
func (s *InMemoryEventStore) SaveEvent(ctx context.Context, event DomainEvent) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if event == nil {
		return fmt.Errorf("事件不能为空")
	}

	stored := StoredEvent{
		ID:          event.GetEventID(),
		EventType:   event.GetEventType(),
		AggregateID: event.GetAggregateID().String(),
		EventData:   event,
		Timestamp:   event.GetOccurredAt(),
		Metadata:    make(map[string]interface{}),
	}

	s.events = append(s.events, stored)
	return nil
}

// LoadEvents 加载指定聚合的所有事件
func (s *InMemoryEventStore) LoadEvents(ctx context.Context, aggregateID string) ([]*StoredEvent, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*StoredEvent
	for i := range s.events {
		if s.events[i].AggregateID == aggregateID {
			result = append(result, &s.events[i])
		}
	}

	return result, nil
}

// LoadEventsByType 加载指定类型的事件
func (s *InMemoryEventStore) LoadEventsByType(ctx context.Context, eventType string) ([]*StoredEvent, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*StoredEvent
	for i := range s.events {
		if s.events[i].EventType == eventType {
			result = append(result, &s.events[i])
		}
	}

	return result, nil
}

// LoadEventsInRange 加载指定时间范围内的事件
func (s *InMemoryEventStore) LoadEventsInRange(ctx context.Context, start, end time.Time) ([]*StoredEvent, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*StoredEvent
	for i := range s.events {
		if s.events[i].Timestamp.After(start) && s.events[i].Timestamp.Before(end) {
			result = append(result, &s.events[i])
		}
	}

	return result, nil
}

// GetEventCount 获取存储的事件总数
func (s *InMemoryEventStore) GetEventCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.events)
}

// Clear 清空所有事件（仅用于测试）
func (s *InMemoryEventStore) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.events = make([]StoredEvent, 0)
}
