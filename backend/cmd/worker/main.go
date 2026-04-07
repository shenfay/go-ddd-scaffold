package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shenfay/go-ddd-scaffold/internal/infra/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	workerHandlers "github.com/shenfay/go-ddd-scaffold/internal/transport/worker/handlers"
	"github.com/shenfay/go-ddd-scaffold/pkg/logger"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("development")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. 初始化日志系统
	if err := logger.Init(cfg.Logger.Level); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("🚀 Starting Asynq Worker...")

	// 2. 初始化 Redis 客户端和数据库
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("❌ Failed to connect to Redis: ", err)
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	logger.Info("✓ Redis connection established")

	// 初始化数据库（用于活动日志仓储）
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Error("❌ Failed to connect to database: ", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("✓ Database connection established")

	// 创建审计日志仓储和处理器
	auditLogRepo := repository.NewAuditLogRepository(db)
	auditLogHandler := workerHandlers.NewAuditLogHandler(auditLogRepo)
	logger.Info("✓ AuditLogRepository and AuditLogHandler initialized")

	// 3. 创建 Asynq 服务器
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Asynq.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
		asynq.Config{
			Concurrency: cfg.Asynq.Concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			StrictPriority: true, // 严格按优先级处理
		},
	)

	logger.Info("✓ Asynq server created with concurrency=", cfg.Asynq.Concurrency)

	// 4. 注册任务处理器
	mux := asynq.NewServeMux()

	// 审计日志任务
	mux.HandleFunc("audit.log.task", auditLogHandler.ProcessTask)
	logger.Info("✓ Registered audit log handler for type: audit.log.task")

	// 5. 启动 Worker
	go func() {
		logger.Info("🎯 Starting Asynq Worker processor...")
		if err := srv.Run(mux); err != nil {
			logger.Error("❌ Failed to run Asynq server: ", err)
			log.Fatalf("Failed to run Asynq server: %v", err)
		}
	}()

	// 6. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("⏹ Shutting down worker...")

	// 7. 优雅关闭
	srv.Shutdown()
	logger.Info("✅ Worker stopped gracefully")
}
