package main

import (
	"os"
	"os/signal"
	"syscall"

	asynq "github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/cmd/shared/bootstrap"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/email"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/worker"
	eventHandler "github.com/shenfay/go-ddd-scaffold/internal/interfaces/event"
	"github.com/shenfay/go-ddd-scaffold/pkg/useragent"
	"go.uber.org/zap"
)

// Processor Worker 任务处理器
type Processor struct {
	infra  *bootstrap.Infrastructure
	logger *zap.Logger
	config *asynq.Config
}

// NewProcessor 创建任务处理器
func NewProcessor(infra *bootstrap.Infrastructure, logger *zap.Logger) *Processor {
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
	daoQuery := dao.Use(p.infra.DB)

	// 邮件服务
	var emailService event.EmailService
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

	// 领域事件处理器
	sideEffectHandler := event.NewSideEffectHandler(p.logger, emailService)

	// 审计日志
	activityLogRepo := repository.NewActivityLogRepository(daoQuery)
	idGeneratorAdapter := idgen.NewGeneratorAdapter()
	auditSubscriber := eventHandler.NewAuditSubscriber(activityLogRepo, idGeneratorAdapter)
	auditHandler := worker.NewAuditLogHandlerAdapter(auditSubscriber)

	// 登录日志
	uaParser := useragent.NewParser()
	loginLogSubscriber := eventHandler.NewLoginLogSubscriber(activityLogRepo, idGeneratorAdapter, &userAgentParserAdapter{parser: uaParser})
	loginLogHandler := worker.NewLoginLogHandlerAdapter(loginLogSubscriber)

	// 创建 Processor
	processor := worker.NewProcessor(p.logger, sideEffectHandler, auditHandler, loginLogHandler)

	// 注册 Handler
	mux := asynq.NewServeMux()
	mux.HandleFunc("domain_event", processor.ProcessTask)

	p.logger.Info("Registered task handlers",
		zap.String("task_type", "domain_event"),
		zap.Int("handler_count", 3),
		zap.Strings("handlers", []string{"SideEffectHandler", "AuditSubscriber", "LoginLogSubscriber"}))

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

// userAgentParserAdapter 适配 useragent.Parser 到 event.UserAgentParser 接口
type userAgentParserAdapter struct {
	parser *useragent.Parser
}

func (a *userAgentParserAdapter) Parse(ua string) eventHandler.DeviceInfo {
	info := a.parser.Parse(ua)
	return eventHandler.DeviceInfo{
		DeviceType: info.DeviceType,
		OS:         info.OS,
		Browser:    info.Browser,
	}
}
