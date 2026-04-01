package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/cmd/app/factory"
	"github.com/shenfay/go-ddd-scaffold/internal/application/jobs"
	"github.com/shenfay/go-ddd-scaffold/internal/application/notifications"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/eventstore"
	queue_asynq "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/queue"
	"go.uber.org/zap"
)

// Processor Worker 任务处理器
type Processor struct {
	infra           *factory.Infrastructure
	logger          *zap.Logger
	taskRouter      *queue_asynq.Router
	processor       *queue_asynq.Processor
	outboxProcessor *eventstore.OutboxProcessor
}

// NewProcessor 创建任务处理器
func NewProcessor(infra *factory.Infrastructure, logger *zap.Logger) *Processor {
	// 创建 Router
	taskRouter, err := queue_asynq.NewRouter(infra.Config.Redis.Addr, logger)
	if err != nil {
		logger.Fatal("Failed to create task router", zap.Error(err))
	}

	// 注册 Jobs
	dailyReportJob := jobs.NewDailyReportJob(logger.Named("daily_report"))
	taskRouter.RegisterJob("job:daily_report", dailyReportJob)

	publishEventJob := jobs.NewPublishDomainEventJob(nil, logger.Named("publish_event"))
	taskRouter.RegisterJob("job:publish_domain_event", publishEventJob)

	// 注册 Notifications
	emailNotification := notifications.NewEmailNotification(logger.Named("email"))
	taskRouter.RegisterNotification("notification:email", emailNotification)

	// 创建队列配置
	queuesConfig := map[string]int{
		// Notifications - 最高优先级
		"notifications_critical": 10,
		"notifications_high":     8,
		"notifications_default":  6,

		// Jobs - Realtime - 高优先级
		"jobs_realtime": 7,

		// Jobs - Default - 中优先级
		"jobs_default": 4,

		// Jobs - Batch - 低优先级
		"jobs_batch": 2,
		"jobs_low":   1,
	}

	// 创建 Processor
	processor, err := queue_asynq.NewProcessor(
		infra.Config.Redis.Addr,
		20, // 并发度
		queuesConfig,
		taskRouter,
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to create processor", zap.Error(err))
	}

	// 创建 Asynq Client（用于 Outbox Processor）
	asynqClient := asynq.NewClient(
		asynq.RedisClientOpt{Addr: infra.Config.Redis.Addr},
	)

	// 创建 Outbox Processor
	outboxProcessor := eventstore.NewOutboxProcessor(
		infra.DB,
		asynqClient,
		logger.Named("outbox"),
	)

	return &Processor{
		infra:           infra,
		logger:          logger,
		taskRouter:      taskRouter,
		processor:       processor,
		outboxProcessor: outboxProcessor,
	}
}

// Run 运行 Worker（包含完整的启动和关闭流程）
func (p *Processor) Run() {
	p.logger.Info("Starting worker...")

	// 启动调度器（处理定时任务）
	if err := p.taskRouter.StartScheduler(); err != nil {
		p.logger.Error("Failed to start scheduler", zap.Error(err))
	}

	// 启动 Outbox Processor（后台协程）
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := p.outboxProcessor.Start(ctx); err != nil {
			p.logger.Error("Outbox processor failed", zap.Error(err))
		}
	}()

	// 启动 Worker
	go func() {
		if err := p.processor.Start(); err != nil {
			p.logger.Fatal("Failed to run processor", zap.Error(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	p.logger.Info("Shutting down worker...")

	// 优雅关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// 停止 Outbox Processor
	cancel()

	// 停止调度器
	p.taskRouter.StopScheduler()

	// 停止处理器
	p.processor.Stop()

	// 等待所有任务完成
	<-shutdownCtx.Done()

	p.logger.Info("Worker stopped")
}
