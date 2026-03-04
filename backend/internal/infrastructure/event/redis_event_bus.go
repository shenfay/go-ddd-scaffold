// Package event Redis Stream 事件总线实现（生产化版本）
package event

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisEventBus Redis Stream 事件总线（生产化实现）
type RedisEventBus struct {
	client      *redis.Client
	eventStore  EventStore
	retryPolicy EventRetryPolicy
	handlers    map[string][]EventHandler
	mutex       sync.RWMutex
	running     bool
	ctx         context.Context
	cancel      context.CancelFunc
}

// RedisEventBusConfig Redis 事件总线配置
type RedisEventBusConfig struct {
	MaxRetries     int           // 最大重试次数
	RetryBaseDelay time.Duration // 重试基础延迟
	PollInterval   time.Duration // 轮询间隔
	BatchSize      int           // 批处理大小
}

// NewRedisEventBus 创建 Redis 事件总线
func NewRedisEventBus(
	client *redis.Client,
	config RedisEventBusConfig,
) *RedisEventBus {
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryBaseDelay == 0 {
		config.RetryBaseDelay = time.Second
	}
	if config.PollInterval == 0 {
		config.PollInterval = time.Second * 5
	}
	if config.BatchSize == 0 {
		config.BatchSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	eventStore := NewRedisEventStore(client, RedisEventStoreConfig{
		MaxRetries:     config.MaxRetries,
		RetryBaseDelay: config.RetryBaseDelay,
	})

	return &RedisEventBus{
		client:      client,
		eventStore:  eventStore,
		retryPolicy: NewExponentialBackoffRetryPolicy(config.RetryBaseDelay, time.Minute*5),
		handlers:    make(map[string][]EventHandler),
		running:     false,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// RegisterHandler 注册事件处理器
func (eb *RedisEventBus) RegisterHandler(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// UnregisterHandler 注销事件处理器
func (eb *RedisEventBus) UnregisterHandler(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	delete(eb.handlers, eventType)
}

// Publish 发布领域事件到 Redis Stream（异步持久化）
func (eb *RedisEventBus) Publish(ctx context.Context, event DomainEvent) error {
	if event == nil {
		return fmt.Errorf("事件不能为空")
	}

	// 持久化到 Redis Stream
	err := eb.eventStore.Store(ctx, event)
	if err != nil {
		return fmt.Errorf("存储事件失败：%w", err)
	}

	return nil
}

// PublishSync 同步发布事件（立即执行处理器）
func (eb *RedisEventBus) PublishSync(ctx context.Context, event DomainEvent) error {
	if event == nil {
		return fmt.Errorf("事件不能为空")
	}

	// 先持久化
	err := eb.eventStore.Store(ctx, event)
	if err != nil {
		return fmt.Errorf("存储事件失败：%w", err)
	}

	// 然后同步执行处理器
	eb.mutex.RLock()
	handlers := eb.handlers[event.GetEventType()]
	eb.mutex.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return fmt.Errorf("执行事件处理器失败：%w", err)
		}
	}

	return nil
}

// Start 启动事件总线（开始消费事件）
func (eb *RedisEventBus) Start() {
	if eb.running {
		return
	}

	eb.running = true
	
	// 启动后台协程消费事件
	go eb.consumeEvents()
}

// Stop 停止事件总线
func (eb *RedisEventBus) Stop() {
	if !eb.running {
		return
	}

	eb.running = false
	eb.cancel()
}

// consumeEvents 消费事件循环
func (eb *RedisEventBus) consumeEvents() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-eb.ctx.Done():
			return
		case <-ticker.C:
			eb.processPendingEvents()
		}
	}
}

// processPendingEvents 处理待处理的事件
func (eb *RedisEventBus) processPendingEvents() {
	ctx, cancel := context.WithTimeout(eb.ctx, time.Second*30)
	defer cancel()

	// 获取待处理事件
	events, err := eb.eventStore.GetPendingEvents(ctx, 100)
	if err != nil {
		return
	}

	if len(events) == 0 {
		return
	}

	// 批量处理事件
	for _, event := range events {
		eb.handleEventWithRetry(ctx, event)
	}
}

// handleEventWithRetry 带重试的事件处理
func (eb *RedisEventBus) handleEventWithRetry(ctx context.Context, event DomainEvent) {
	eb.mutex.RLock()
	handlers := eb.handlers[event.GetEventType()]
	eb.mutex.RUnlock()

	if len(handlers) == 0 {
		return // 没有处理器，跳过
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		lastErr = nil
		
		// 执行所有处理器
		for _, handler := range handlers {
			if err := handler(ctx, event); err != nil {
				lastErr = err
				break
			}
		}

		// 成功则标记为已处理
		if lastErr == nil {
			_ = eb.eventStore.MarkAsProcessed(ctx, event.GetEventID())
			return
		}

		// 失败则等待后重试
		if eb.retryPolicy.ShouldRetry(attempt, 3) {
			delay := eb.retryPolicy.GetDelay(attempt)
			time.Sleep(delay)
		}
	}

	// 超过最大重试次数，标记为失败
	_ = eb.eventStore.MarkAsFailed(ctx, event.GetEventID(), lastErr.Error())
}

// GetHandlerCount 获取处理器数量
func (eb *RedisEventBus) GetHandlerCount(eventType string) int {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	return len(eb.handlers[eventType])
}

// ListRegisteredEvents 列出所有已注册事件
func (eb *RedisEventBus) ListRegisteredEvents() []string {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	events := make([]string, 0, len(eb.handlers))
	for eventType := range eb.handlers {
		events = append(events, eventType)
	}

	return events
}
