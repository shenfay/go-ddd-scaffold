package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	asynq_pkg "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/messaging/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/email"
	logging "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/log"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/worker"
	eventHandler "github.com/shenfay/go-ddd-scaffold/internal/interfaces/event"
	"github.com/shenfay/go-ddd-scaffold/pkg/useragent"
	"go.uber.org/zap"
)

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

func main() {
	// 1. 加载环境变量和配置（与 API 入口相同的方式）
	env := os.Getenv("ENV_MODE")
	if env == "" {
		env = "development"
	}

	configLoader := config.NewConfigLoader(nil)
	appConfig, err := configLoader.Load(env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 创建 Logger
	logConfig := &config.LoggingConfig{
		Level:      appConfig.Logging.Level,
		Format:     appConfig.Logging.Format,
		File:       appConfig.Logging.File,
		MaxSize:    appConfig.Logging.MaxSize,
		MaxBackups: appConfig.Logging.MaxBackups,
		MaxAge:     appConfig.Logging.MaxAge,
	}
	appLogger, err := logging.New(logConfig)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer appLogger.Sync()

	logger := appLogger.Logger.Named("worker")

	logger.Info("Starting Asynq Worker...")

	logger.Info("Configuration loaded",
		zap.String("env", env),
		zap.String("redis_addr", appConfig.Redis.Addr))

	// 3. 创建基础设施
	infra, cleanup, err := bootstrap.NewInfra(appConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create infrastructure", zap.Error(err))
	}
	defer cleanup()

	// 4. 创建 Asynq Server
	srv := asynq_pkg.NewServer(asynq_pkg.Config{
		RedisAddr:     appConfig.Redis.Addr,
		RedisPassword: appConfig.Redis.Password,
		RedisDB:       appConfig.Redis.DB,
	})

	// 5. 创建任务处理器并注册
	// 创建 DAO Query
	daoQuery := dao.Use(infra.DB)

	// 创建邮件服务
	var emailService event.EmailService
	if appConfig.Email.SMTPHost != "" {
		emailService = email.NewSMTPService(appConfig.Email, logger)
		logger.Info("邮件服务已配置",
			zap.String("smtp_host", appConfig.Email.SMTPHost),
			zap.String("from", appConfig.Email.From),
		)
	} else {
		emailService = email.NewNoOpService(logger)
		logger.Info("邮件服务未配置，使用空实现")
	}

	// 创建领域事件处理器（副作用处理：发邮件等）
	sideEffectHandler := event.NewSideEffectHandler(logger, emailService)

	// 创建活动日志订阅器（使用新的 ActivityLogRepository）
	activityLogRepo := repository.NewActivityLogRepository(daoQuery)
	auditSubscriber := eventHandler.NewAuditSubscriber(activityLogRepo, infra.Snowflake)
	auditHandler := worker.NewAuditLogHandlerAdapter(auditSubscriber)

	// 创建活动日志订阅器（登录日志也使用同一个 ActivityLogRepository）
	uaParser := useragent.NewParser()
	loginLogSubscriber := eventHandler.NewLoginLogSubscriber(activityLogRepo, infra.Snowflake, &userAgentParserAdapter{parser: uaParser})
	loginLogHandler := worker.NewLoginLogHandlerAdapter(loginLogSubscriber)

	// 创建 Processor 并注册所有 Handler
	processor := worker.NewProcessor(logger, sideEffectHandler, auditHandler, loginLogHandler)

	// 创建 ServeMux 并注册处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc(asynq_pkg.TaskTypeDomainEvent, processor.ProcessTask)

	logger.Info("Registered task handlers",
		zap.String("task_type", asynq_pkg.TaskTypeDomainEvent),
		zap.Int("handler_count", 3),
		zap.Strings("handlers", []string{"SideEffectHandler", "AuditSubscriber", "LoginLogSubscriber"}))

	// 6. 启动 Worker（在 goroutine 中运行）
	go func() {
		logger.Info("Worker started, waiting for tasks...")
		if err := srv.Run(mux); err != nil {
			logger.Fatal("Failed to run asynq server", zap.Error(err))
		}
	}()

	// 7. 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")

	// 8. 优雅关闭
	srv.Shutdown()

	// 使用 context 确保关闭完成
	_ = context.Background()
	_ = infra // infra 会通过 defer cleanup() 自动清理

	logger.Info("Worker stopped")
}
