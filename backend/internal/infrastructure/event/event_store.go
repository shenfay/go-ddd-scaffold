// Package event 事件存储接口定义
package event

import (
	"context"
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
