package asynq

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// Processor Asynq 任务处理器
type Processor struct {
	srv  *asynq.Server
	mux  *asynq.ServeMux
	task *Router
}

// NewProcessor 创建 Asynq 处理器
func NewProcessor(redisAddr string, concurrency int, queuesConfig map[string]int, taskRouter *Router, logger *zap.Logger) (*Processor, error) {
	if logger == nil {
		logger = zap.L().Named("asynq_processor")
	}

	// 创建 Asynq Server（使用默认配置，暂不实现自定义 Logger 和 ErrorHandler）
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency:    concurrency,
			Queues:         queuesConfig,
			RetryDelayFunc: createRetryDelayFunc(),
		},
	)

	// 创建路由
	mux := asynq.NewServeMux()

	processor := &Processor{
		srv:  srv,
		mux:  mux,
		task: taskRouter,
	}

	// 注册通用处理器
	processor.setupHandlers()

	return processor, nil
}

// setupHandlers 注册通用处理器
func (p *Processor) setupHandlers() {
	// 所有任务都通过 TaskRouter 分发
	p.mux.HandleFunc("*", func(ctx context.Context, t *asynq.Task) error {
		return p.task.ProcessTask(ctx, t)
	})
}

// Start 启动处理器
func (p *Processor) Start() error {
	p.task.logger.Info("Starting Asynq processor")
	return p.srv.Run(p.mux)
}

// Stop 停止处理器
func (p *Processor) Stop() {
	p.task.logger.Info("Stopping Asynq processor")
	p.srv.Shutdown()
}

// createRetryDelayFunc 创建重试延迟函数
func createRetryDelayFunc() asynq.RetryDelayFunc {
	return func(n int, err error, t *asynq.Task) time.Duration {
		// 根据任务类型设置不同的重试策略
		taskType := t.Type()

		if isNotificationTask(taskType) {
			// Notification：快速重试（1 分钟、5 分钟、15 分钟）
			delays := []time.Duration{
				1 * time.Minute,
				5 * time.Minute,
				15 * time.Minute,
			}
			if n < len(delays) {
				return delays[n]
			}
			return 30 * time.Minute
		}

		// Job：慢速重试（5 分钟、30 分钟、1 小时）
		delays := []time.Duration{
			5 * time.Minute,
			30 * time.Minute,
			1 * time.Hour,
		}
		if n < len(delays) {
			return delays[n]
		}
		return 2 * time.Hour
	}
}

// isNotificationTask 判断是否为通知类任务
func isNotificationTask(taskType string) bool {
	return len(taskType) > 12 && taskType[:12] == "notification"
}
