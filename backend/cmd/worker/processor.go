package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/cmd/app/factory"
	"github.com/shenfay/go-ddd-scaffold/internal/application/user/subscriber"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/email"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/worker"
	"go.uber.org/zap"
)

// Processor Worker 任务处理器
type Processor struct {
	infra  *factory.Infrastructure
	logger *zap.Logger
	config *asynq.Config
}

// NewProcessor 创建任务处理器
func NewProcessor(infra *factory.Infrastructure, logger *zap.Logger) *Processor {
	return &Processor{
		infra:  infra,
		logger: logger,
		config: createAsynqConfig(),
	}
}

// Run 运行 Worker（包含完整的启动和关闭流程）
func (p *Processor) Run() {
	redisOpt := asynq.RedisClientOpt{
		Addr:     p.infra.Config.Redis.Addr,
		Password: p.infra.Config.Redis.Password,
		DB:       p.infra.Config.Redis.DB,
	}

	srv := asynq.NewServer(redisOpt, *p.config)
	mux := p.createTaskHandlers()

	p.logger.Info("Worker started, waiting for tasks...")

	// 启动 Worker
	go func() {
		if err := srv.Run(mux); err != nil {
			p.logger.Fatal("Failed to run asynq server", zap.Error(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	p.logger.Info("Shutting down worker...")
	srv.Shutdown()
	p.logger.Info("Worker stopped")
}

// createTaskHandlers 创建任务处理器
func (p *Processor) createTaskHandlers() *asynq.ServeMux {

	// 邮件服务
	var emailService subscriber.EmailService
	if p.infra.Config.Email.SMTPHost != "" {
		emailService = email.NewSMTPService(p.infra.Config.Email, p.logger)
		p.logger.Info("邮件服务已配置",
			zap.String("smtp_host", p.infra.Config.Email.SMTPHost),
			zap.String("from", p.infra.Config.Email.From),
		)
	} else {
		emailService = email.NewNoOpService(p.logger)
		p.logger.Info("邮件服务未配置，使用空实现")
	}

	// 领域事件订阅器（✅ 新版在应用层，实现 Handler 接口）
	eventSubscriber := subscriber.NewUserEventSubscriber(
		p.logger,
		emailService,
		nil, // TODO: 需要注入 StatisticsRepository
	)

	// 创建 Processor（传入实现了 Handler 接口的 UserEventSubscriber）
	processor := worker.NewProcessor(p.logger, eventSubscriber)

	// 注册 Handler
	mux := asynq.NewServeMux()
	mux.HandleFunc("domain_event", processor.ProcessTask)

	p.logger.Info("Registered task handlers",
		zap.String("task_type", "domain_event"),
		zap.Int("handler_count", 1),
		zap.Strings("handlers", []string{"UserEventSubscriber"}))

	return mux
}

// createAsynqConfig 创建 Asynq 配置
func createAsynqConfig() *asynq.Config {
	return &asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	}
}
