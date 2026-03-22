package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/logging"
	task_queue "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/taskqueue"
	"go.uber.org/zap"
)

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
	srv := task_queue.NewServer(task_queue.Config{
		RedisAddr:     appConfig.Redis.Addr,
		RedisPassword: appConfig.Redis.Password,
		RedisDB:       appConfig.Redis.DB,
	})

	// 5. 创建任务处理器并注册
	// 创建 Processor（可以在此处添加 Handler 来处理特定的领域事件）
	processor := task_queue.NewProcessor(logger)

	// 创建 ServeMux 并注册处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc(task_queue.TaskTypeDomainEvent, processor.ProcessTask)

	logger.Info("Registered task handlers",
		zap.String("task_type", task_queue.TaskTypeDomainEvent))

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
