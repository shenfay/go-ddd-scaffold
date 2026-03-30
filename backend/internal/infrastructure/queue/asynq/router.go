package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/queue"
)

// QueueConfig 队列配置
type QueueConfig struct {
	Name     string
	Priority int
	MaxRetry int
	Timeout  time.Duration
}

// Router 任务路由器（Asynq 实现）
type Router struct {
	jobs          map[string]queue.JobHandler
	notifications map[string]queue.NotificationHandler
	scheduler     *asynq.Scheduler
	logger        *zap.Logger
}

// NewRouter 创建任务路由器
func NewRouter(redisAddr string, logger *zap.Logger) (*Router, error) {
	if logger == nil {
		logger = zap.L().Named("router")
	}

	// 创建调度器
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisAddr},
		nil,
	)

	return &Router{
		jobs:          make(map[string]queue.JobHandler),
		notifications: make(map[string]queue.NotificationHandler),
		scheduler:     scheduler,
		logger:        logger,
	}, nil
}

// RegisterJob 注册 Job
func (r *Router) RegisterJob(name string, handler queue.JobHandler) {
	r.jobs[name] = handler
	r.logger.Info("Job registered", zap.String("name", name), zap.String("queue", handler.Queue()))

	// 如果是定时任务，自动注册到调度器
	if scheduled, ok := handler.(queue.ScheduledJob); ok {
		r.registerScheduledJob(name, scheduled)
	}
}

// registerScheduledJob 注册定时任务
func (r *Router) registerScheduledJob(name string, job queue.ScheduledJob) {
	task := asynq.NewTask(name, nil, asynq.Queue(job.Queue()))

	entryID, err := r.scheduler.Register(job.Schedule(), task)
	if err != nil {
		r.logger.Error("Failed to register scheduled job",
			zap.String("name", name),
			zap.String("schedule", job.Schedule()),
			zap.Error(err))
		return
	}

	r.logger.Info("Scheduled job registered",
		zap.String("name", name),
		zap.String("entry_id", entryID),
		zap.String("schedule", job.Schedule()))
}

// RegisterNotification 注册 Notification
func (r *Router) RegisterNotification(name string, handler queue.NotificationHandler) {
	r.notifications[name] = handler
	r.logger.Info("Notification registered", zap.String("name", name), zap.String("queue", handler.Queue()))
}

// ProcessTask 处理任务（供 Asynq Handler 调用）
func (r *Router) ProcessTask(ctx context.Context, task *asynq.Task) error {
	taskType := task.Type()
	payload, err := extractPayload(task)
	if err != nil {
		return fmt.Errorf("failed to extract payload: %w", err)
	}

	// 1. 尝试作为 Job 处理
	if handler, ok := r.jobs[taskType]; ok {
		r.logger.Debug("Processing job",
			zap.String("type", taskType),
			zap.Any("payload", payload))

		// 设置超时上下文
		timeoutCtx, cancel := context.WithTimeout(ctx, handler.Timeout())
		defer cancel()

		return handler.Execute(timeoutCtx, payload)
	}

	// 2. 尝试作为 Notification 处理
	if handler, ok := r.notifications[taskType]; ok {
		r.logger.Debug("Processing notification",
			zap.String("type", taskType),
			zap.Any("payload", payload))

		return handler.Handle(ctx, payload)
	}

	return fmt.Errorf("unknown task type: %s", taskType)
}

// GetQueueConfig 获取队列配置
func (r *Router) GetQueueConfig(taskType string) QueueConfig {
	// 1. 检查是否为 Notification
	if handler, ok := r.notifications[taskType]; ok {
		return QueueConfig{
			Name:     handler.Queue(),
			Priority: getNotificationPriority(handler.Queue()),
			MaxRetry: 3,
			Timeout:  5 * time.Minute,
		}
	}

	// 2. 检查是否为 Job
	if handler, ok := r.jobs[taskType]; ok {
		return QueueConfig{
			Name:     handler.Queue(),
			Priority: getJobPriority(handler.Queue()),
			MaxRetry: handler.MaxRetry(),
			Timeout:  handler.Timeout(),
		}
	}

	// 默认配置
	return QueueConfig{
		Name:     "default",
		Priority: 1,
		MaxRetry: 3,
		Timeout:  5 * time.Minute,
	}
}

// StartScheduler 启动调度器
func (r *Router) StartScheduler() error {
	r.logger.Info("Starting task scheduler")
	return r.scheduler.Start()
}

// StopScheduler 停止调度器
func (r *Router) StopScheduler() {
	r.logger.Info("Stopping task scheduler")
	r.scheduler.Shutdown()
}

// extractPayload 从 Asynq Task 中提取 payload
func extractPayload(task *asynq.Task) (map[string]interface{}, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// getNotificationPriority 获取通知队列优先级
func getNotificationPriority(queueName string) int {
	switch queueName {
	case "notifications_critical":
		return 10
	case "notifications_high":
		return 8
	default:
		return 6 // notifications_default
	}
}

// getJobPriority 获取作业队列优先级
func getJobPriority(queueName string) int {
	switch queueName {
	case "jobs_realtime":
		return 7
	case "jobs_scheduled":
		return 4
	case "jobs_batch":
		return 2
	default:
		return 4 // jobs_default
	}
}
