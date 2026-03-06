// Package event 事件管理器 - 自动发现和注册事件处理器
package event

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// EventProcessor 事件处理器接口（支持重试）
type EventProcessor interface {
	// Handle 处理事件
	Handle(ctx context.Context, event DomainEvent) error
	// RetryPolicy 获取重试策略
	RetryPolicy() *RetryPolicy
}

// EventManager 事件管理器
type EventManager struct {
	eventBus       *EventBus
	logger         *zap.Logger
	processors     map[string][]EventProcessor
	retryQueue     *RetryQueue
	maxRetries     int
	initialDelay   time.Duration
	multiplier     float64
}

// NewEventManager 创建事件管理器
func NewEventManager(eventBus *EventBus, logger *zap.Logger) *EventManager {
	return &EventManager{
		eventBus:     eventBus,
		logger:       logger,
		processors:   make(map[string][]EventProcessor),
		retryQueue:   NewRetryQueue(),
		maxRetries:   3,
		initialDelay: time.Second,
		multiplier:   2.0,
	}
}

// RegisterProcessor 注册事件处理器
func (m *EventManager) RegisterProcessor(eventType string, processor EventProcessor) {
	m.processors[eventType] = append(m.processors[eventType], processor)
	
	// 同时在事件总线中注册
	m.eventBus.RegisterHandler(eventType, func(ctx context.Context, event DomainEvent) error {
		return m.processEvent(ctx, event, processor)
	})
	
	m.logger.Info("注册事件处理器",
		zap.String("eventType", eventType),
		zap.String("processor", fmt.Sprintf("%T", processor)),
	)
}

// RegisterHandler 注册事件处理器（兼容旧 API）
func (m *EventManager) RegisterHandler(eventType string, handler EventHandler) {
	m.eventBus.RegisterHandler(eventType, handler)
	
	m.logger.Info("注册事件处理器",
		zap.String("eventType", eventType),
	)
}

// GetEventBus 获取事件总线
func (m *EventManager) GetEventBus() *EventBus {
	return m.eventBus
}

// processEvent 处理事件（带重试）
func (m *EventManager) processEvent(ctx context.Context, event DomainEvent, processor EventProcessor) error {
	policy := processor.RetryPolicy()
	if policy == nil {
		policy = &RetryPolicy{
			MaxRetries:   m.maxRetries,
			InitialDelay: m.initialDelay,
			Multiplier:   m.multiplier,
		}
	}

	var lastErr error
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		lastErr = processor.Handle(ctx, event)
		if lastErr == nil {
			// 处理成功
			m.logger.Info("事件处理成功",
				zap.String("eventId", event.GetEventID()),
				zap.String("eventType", event.GetEventType()),
				zap.Int("attempt", attempt+1),
			)
			return nil
		}

		// 处理失败，判断是否重试
		if !policy.ShouldRetry(attempt+1, lastErr) {
			m.logger.Error("事件处理失败，不重试",
				zap.String("eventId", event.GetEventID()),
				zap.String("eventType", event.GetEventType()),
				zap.Error(lastErr),
			)
			return lastErr
		}

		// 等待后重试
		delay := policy.CalculateNextRetryDelay(attempt + 1)
		if delay > 0 {
			m.logger.Warn("事件处理失败，准备重试",
				zap.String("eventId", event.GetEventID()),
				zap.String("eventType", event.GetEventType()),
				zap.Error(lastErr),
				zap.Int("attempt", attempt+1),
				zap.Duration("delay", delay),
			)
			time.Sleep(delay)
		}
	}

	m.logger.Error("事件处理失败，已达最大重试次数",
		zap.String("eventId", event.GetEventID()),
		zap.String("eventType", event.GetEventType()),
		zap.Error(lastErr),
	)
	
	return lastErr
}

// RetryPolicy 重试策略配置
type RetryPolicy struct {
	MaxRetries   int           // 最大重试次数
	InitialDelay time.Duration // 初始延迟
	Multiplier   float64       // 延迟倍增系数
	MaxDelay     time.Duration // 最大延迟
}

// CalculateNextRetryDelay 计算下一次重试的延迟
func (p *RetryPolicy) CalculateNextRetryDelay(attempt int) time.Duration {
	if attempt > p.MaxRetries {
		return 0
	}

	delay := p.InitialDelay
	for i := 1; i < attempt && delay < p.MaxDelay; i++ {
		delay = time.Duration(float64(delay) * p.Multiplier)
	}

	if delay > p.MaxDelay {
		return p.MaxDelay
	}

	return delay
}

// ShouldRetry 判断是否应该重试
func (p *RetryPolicy) ShouldRetry(attempt int, err error) bool {
	if attempt > p.MaxRetries {
		return false
	}

	// 某些错误不应该重试（如参数错误）
	// 这里可以根据具体错误类型判断
	return true
}

// RetryQueue 重试队列（简单实现）
type RetryQueue struct {
	items chan *RetryItem
}

// RetryItem 重试项
type RetryItem struct {
	Event      DomainEvent
	Processor  EventProcessor
	Attempt    int
	NextRetry  time.Time
}

// NewRetryQueue 创建重试队列
func NewRetryQueue() *RetryQueue {
	return &RetryQueue{
		items: make(chan *RetryItem, 1000),
	}
}

// Enqueue 加入重试队列
func (q *RetryQueue) Enqueue(item *RetryItem) {
	select {
	case q.items <- item:
	default:
		// 队列已满，丢弃或记录日志
	}
}

// Dequeue 从重试队列取出
func (q *RetryQueue) Dequeue() *RetryItem {
	select {
	case item := <-q.items:
		return item
	default:
		return nil
	}
}
