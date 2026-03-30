package queue

import (
	"context"
	"time"
)

// JobHandler 后台作业处理器接口
type JobHandler interface {
	// Execute 执行任务
	Execute(ctx context.Context, payload map[string]interface{}) error
	// Queue 返回队列名称
	Queue() string
	// MaxRetry 返回最大重试次数
	MaxRetry() int
	// Timeout 返回超时时间
	Timeout() time.Duration
}

// ScheduledJob 定时任务接口（可选实现）
type ScheduledJob interface {
	JobHandler
	// Schedule 返回 Cron 表达式
	Schedule() string
}

// NotificationHandler 通知处理器接口
type NotificationHandler interface {
	// Handle 处理通知
	Handle(ctx context.Context, payload map[string]interface{}) error
	// Queue 返回队列名称
	Queue() string
}
