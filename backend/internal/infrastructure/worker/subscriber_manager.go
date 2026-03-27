package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"go.uber.org/zap"
)

// EventHandler 事件处理器接口
type EventHandler interface {
	// HandleEvent 处理事件
	HandleEvent(ctx context.Context, event common.DomainEvent) error
	// EventType 返回处理器能处理的事件类型
	EventType() string
	// Priority 返回处理器优先级（数字越小优先级越高）
	Priority() int
}

// SubscriberManager 事件订阅者管理器
// 提供统一的订阅者注册、管理和生命周期控制
type SubscriberManager struct {
	handlers   map[string][]EventHandler // eventType -> handlers
	handlerMap map[string]EventHandler   // handlerName -> handler (用于快速查找)
	mutex      sync.RWMutex
	logger     *zap.Logger
	config     SubscriberConfig
	metrics    *SubscriberMetrics
}

// SubscriberConfig 订阅者管理器配置
type SubscriberConfig struct {
	// 并发处理配置
	MaxConcurrency int
	BufferSize     int

	// 超时配置
	HandlerTimeout time.Duration

	// 错误处理配置
	MaxRetries      int
	RetryDelay      time.Duration
	DeadLetterQueue string

	// 监控配置
	EnableMetrics bool
	MetricsPrefix string
}

// SubscriberMetrics 订阅者指标
type SubscriberMetrics struct {
	TotalProcessed uint64
	TotalErrors    uint64
	AverageLatency time.Duration
}

// NewSubscriberManager 创建订阅者管理器
func NewSubscriberManager(logger *zap.Logger, config SubscriberConfig) *SubscriberManager {
	if logger == nil {
		logger = zap.L().Named("subscriber_manager")
	}

	// 设置默认配置
	if config.MaxConcurrency <= 0 {
		config.MaxConcurrency = 10
	}
	if config.BufferSize <= 0 {
		config.BufferSize = 100
	}
	if config.HandlerTimeout == 0 {
		config.HandlerTimeout = 30 * time.Second
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = time.Second
	}
	if config.MetricsPrefix == "" {
		config.MetricsPrefix = "events"
	}

	return &SubscriberManager{
		handlers:   make(map[string][]EventHandler),
		handlerMap: make(map[string]EventHandler),
		logger:     logger,
		config:     config,
		metrics:    &SubscriberMetrics{},
	}
}

// RegisterHandler 注册事件处理器
func (sm *SubscriberManager) RegisterHandler(handler EventHandler) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	handlerName := fmt.Sprintf("%s_%s", handler.EventType(), getHandlerTypeName(handler))

	// 检查是否已注册
	if _, exists := sm.handlerMap[handlerName]; exists {
		return fmt.Errorf("handler %s already registered", handlerName)
	}

	// 添加到处理器列表
	eventType := handler.EventType()
	sm.handlers[eventType] = append(sm.handlers[eventType], handler)
	sm.handlerMap[handlerName] = handler

	// 按优先级排序
	sm.sortHandlersByPriority(eventType)

	sm.logger.Info("Handler registered",
		zap.String("event_type", eventType),
		zap.String("handler_name", handlerName),
		zap.Int("priority", handler.Priority()))

	return nil
}

// UnregisterHandler 注销事件处理器
func (sm *SubscriberManager) UnregisterHandler(handler EventHandler) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	handlerName := fmt.Sprintf("%s_%s", handler.EventType(), getHandlerTypeName(handler))

	// 从映射中删除
	if _, exists := sm.handlerMap[handlerName]; !exists {
		return fmt.Errorf("handler %s not found", handlerName)
	}

	delete(sm.handlerMap, handlerName)

	// 从列表中删除
	eventType := handler.EventType()
	handlers := sm.handlers[eventType]
	for i, h := range handlers {
		if h == handler {
			sm.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	sm.logger.Info("Handler unregistered",
		zap.String("event_type", eventType),
		zap.String("handler_name", handlerName))

	return nil
}

// ProcessEvent 处理事件
func (sm *SubscriberManager) ProcessEvent(ctx context.Context, event common.DomainEvent) error {
	eventType := event.EventName()

	sm.mutex.RLock()
	handlers := sm.handlers[eventType]
	sm.mutex.RUnlock()

	if len(handlers) == 0 {
		sm.logger.Debug("No handlers found for event",
			zap.String("event_type", eventType))
		return nil
	}

	// 创建带超时的上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, sm.config.HandlerTimeout)
	defer cancel()

	// 并发处理所有处理器
	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()

			startTime := time.Now()

			// 执行处理器
			if err := sm.executeHandlerWithRetry(timeoutCtx, h, event); err != nil {
				errChan <- fmt.Errorf("handler %s failed: %w", h.EventType(), err)
				sm.metrics.TotalErrors++
			} else {
				sm.metrics.TotalProcessed++
				latency := time.Since(startTime)
				// 简单的平均延迟计算（实际应使用更精确的方法）
				sm.metrics.AverageLatency = (sm.metrics.AverageLatency + latency) / 2
			}
		}(handler)
	}

	// 等待所有处理器完成
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 所有处理器完成
		close(errChan)
		return sm.collectErrors(errChan)
	case <-timeoutCtx.Done():
		// 超时
		sm.logger.Warn("Event processing timeout",
			zap.String("event_type", eventType),
			zap.Duration("timeout", sm.config.HandlerTimeout))
		return fmt.Errorf("event processing timeout for %s", eventType)
	}
}

// executeHandlerWithRetry 带重试的处理器执行
func (sm *SubscriberManager) executeHandlerWithRetry(ctx context.Context, handler EventHandler, event common.DomainEvent) error {
	var lastErr error

	for attempt := 0; attempt <= sm.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// 重试前等待
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(sm.config.RetryDelay * time.Duration(attempt)):
			}
		}

		err := handler.HandleEvent(ctx, event)
		if err == nil {
			// 成功处理
			if attempt > 0 {
				sm.logger.Info("Handler succeeded after retries",
					zap.String("event_type", event.EventName()),
					zap.Int("attempts", attempt+1))
			}
			return nil
		}

		lastErr = err

		// 检查是否应该跳过重试
		if sm.shouldSkipRetry(err) {
			sm.logger.Warn("Skipping retry for handler",
				zap.String("event_type", event.EventName()),
				zap.Error(err))
			break
		}

		sm.logger.Warn("Handler failed, will retry",
			zap.String("event_type", event.EventName()),
			zap.Int("attempt", attempt+1),
			zap.Int("max_retries", sm.config.MaxRetries),
			zap.Error(err))
	}

	// 所有重试都失败了
	sm.sendToDeadLetterQueue(ctx, event, lastErr)
	return fmt.Errorf("handler failed after %d attempts: %w", sm.config.MaxRetries+1, lastErr)
}

// shouldSkipRetry 判断是否应该跳过重试
func (sm *SubscriberManager) shouldSkipRetry(err error) bool {
	// 某些错误不应该重试（如业务逻辑错误）
	errorMsg := err.Error()
	skipErrorTypes := []string{
		"validation failed",
		"business rule violation",
		"invalid state",
	}

	for _, skipType := range skipErrorTypes {
		if contains(errorMsg, skipType) {
			return true
		}
	}

	return false
}

// sendToDeadLetterQueue 发送到死信队列
func (sm *SubscriberManager) sendToDeadLetterQueue(ctx context.Context, event common.DomainEvent, err error) {
	if sm.config.DeadLetterQueue == "" {
		return
	}

	sm.logger.Error("Event sent to dead letter queue",
		zap.String("event_type", event.EventName()),
		zap.String("dead_letter_queue", sm.config.DeadLetterQueue),
		zap.Error(err))

	// 这里应该实现实际的死信队列发送逻辑
	// 例如：记录到专门的数据库表或发送到另一个 Asynq 队列
}

// sortHandlersByPriority 按优先级对处理器排序
func (sm *SubscriberManager) sortHandlersByPriority(eventType string) {
	handlers := sm.handlers[eventType]

	// 简单的冒泡排序（处理器数量通常很少）
	for i := 0; i < len(handlers)-1; i++ {
		for j := 0; j < len(handlers)-1-i; j++ {
			if handlers[j].Priority() > handlers[j+1].Priority() {
				handlers[j], handlers[j+1] = handlers[j+1], handlers[j]
			}
		}
	}

	sm.handlers[eventType] = handlers
}

// collectErrors 收集所有错误
func (sm *SubscriberManager) collectErrors(errChan chan error) error {
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) == 0 {
		return nil
	}

	// 合并错误信息
	errorMsg := "Multiple handler errors:"
	for _, err := range errors {
		errorMsg += fmt.Sprintf("\n- %v", err)
	}

	return fmt.Errorf(errorMsg)
}

// GetMetrics 获取指标信息
func (sm *SubscriberManager) GetMetrics() *SubscriberMetrics {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// 返回副本以避免竞态条件
	metrics := *sm.metrics
	return &metrics
}

// GetHandlerCount 获取处理器数量统计
func (sm *SubscriberManager) GetHandlerCount() map[string]int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	counts := make(map[string]int)
	for eventType, handlers := range sm.handlers {
		counts[eventType] = len(handlers)
	}
	return counts
}

// Shutdown 优雅关闭
func (sm *SubscriberManager) Shutdown(ctx context.Context) error {
	sm.logger.Info("Shutting down subscriber manager")

	// 这里可以添加清理逻辑
	// 例如：等待正在处理的任务完成

	return nil
}

// 辅助函数
func getHandlerTypeName(handler EventHandler) string {
	// 简化的类型名称获取（实际项目中可能需要反射）
	return fmt.Sprintf("%T", handler)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
