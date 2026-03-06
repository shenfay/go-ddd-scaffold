// Package handlers 领域事件处理器实现
package handlers

import (
	"context"
	"time"

	"go-ddd-scaffold/internal/infrastructure/event"
	"go-ddd-scaffold/internal/pkg/errors"
	"go.uber.org/zap"
)

// UserRegisteredEventHandler 用户注册事件处理器
type UserRegisteredEventHandler struct {
	logger *zap.Logger
}

// NewUserRegisteredEventHandler 创建用户注册事件处理器
func NewUserRegisteredEventHandler(logger *zap.Logger) *UserRegisteredEventHandler {
	return &UserRegisteredEventHandler{
		logger: logger,
	}
}

// Handle 处理用户注册事件
func (h *UserRegisteredEventHandler) Handle(ctx context.Context, event event.DomainEvent) error {
	h.logger.Info("处理用户注册事件",
		zap.String("eventId", event.GetEventID()),
		zap.String("userId", event.GetAggregateID().String()),
		zap.Time("occurredAt", event.GetOccurredAt()),
	)

	// TODO: 实现具体的业务逻辑
	// 例如：
	// 1. 发送欢迎邮件
	// 2. 初始化用户数据
	// 3. 触发其他下游操作

	return nil
}

// RetryPolicy 重试策略
func (h *UserRegisteredEventHandler) RetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:    3,
		InitialDelay:  time.Second,
		Multiplier:    2,
		MaxDelay:      time.Minute,
	}
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
	if attempt >= p.MaxRetries {
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
	if attempt >= p.MaxRetries {
		return false
	}

	// 某些错误不应该重试（如参数错误）
	if appErr, ok := err.(*errors.AppError); ok {
		switch appErr.GetCode() {
		case "InvalidParameter", "ValidationFailed":
			return false
		}
	}

	return true
}
