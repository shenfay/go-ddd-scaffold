package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/listener"
	workerHandlers "github.com/shenfay/go-ddd-scaffold/internal/transport/worker/handlers"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
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

	// 创建活动日志仓储和处理器
	activityLogRepo := repository.NewActivityLogRepository(db)
	activityLogHandler := workerHandlers.NewActivityLogWorkerHandler(activityLogRepo)
	logger.Info("✓ ActivityLogRepository and ActivityLogHandler initialized")

	// 创建 Asynq Client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Asynq.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer asynqClient.Close()
	logger.Info("✓ Asynq client initialized (Worker)")

	// 创建 EventBus（用于 Worker 端订阅领域事件并发布日志任务）
	eventBus := messaging.NewEventBus(asynqClient)
	logger.Info("✓ EventBus initialized (Worker)")

	// 创建监听器并订阅领域事件
	auditLogListener := listener.NewAuditLogListener(eventBus)
	auditLogListener.SubscribeEvents(eventBus)
	logger.Info("✓ AuditLogListener registered")

	activityLogListener := listener.NewActivityLogListener(eventBus)
	activityLogListener.SubscribeEvents(eventBus)
	logger.Info("✓ ActivityLogListener registered")

	// 转换为 AsynqEventBus 以获取 handlers
	asynqEventBus := eventBus.(*messaging.AsynqEventBus)

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

	// 注册领域事件处理器（由 Listener 处理）
	mux.HandleFunc("USER.REGISTERED", func(ctx context.Context, task *asynq.Task) error {
		return processDomainEvent(ctx, "USER.REGISTERED", task, asynqEventBus)
	})
	logger.Info("✓ Registered handler for: USER.REGISTERED")

	mux.HandleFunc("AUTH.LOGIN.SUCCESS", func(ctx context.Context, task *asynq.Task) error {
		return processDomainEvent(ctx, "AUTH.LOGIN.SUCCESS", task, asynqEventBus)
	})
	logger.Info("✓ Registered handler for: AUTH.LOGIN.SUCCESS")

	mux.HandleFunc("AUTH.LOGIN.FAILED", func(ctx context.Context, task *asynq.Task) error {
		return processDomainEvent(ctx, "AUTH.LOGIN.FAILED", task, asynqEventBus)
	})
	logger.Info("✓ Registered handler for: AUTH.LOGIN.FAILED")

	mux.HandleFunc("SECURITY.ACCOUNT.LOCKED", func(ctx context.Context, task *asynq.Task) error {
		return processDomainEvent(ctx, "SECURITY.ACCOUNT.LOCKED", task, asynqEventBus)
	})
	logger.Info("✓ Registered handler for: SECURITY.ACCOUNT.LOCKED")

	mux.HandleFunc("AUTH.LOGOUT", func(ctx context.Context, task *asynq.Task) error {
		return processDomainEvent(ctx, "AUTH.LOGOUT", task, asynqEventBus)
	})
	logger.Info("✓ Registered handler for: AUTH.LOGOUT")

	mux.HandleFunc("AUTH.TOKEN.REFRESHED", func(ctx context.Context, task *asynq.Task) error {
		return processDomainEvent(ctx, "AUTH.TOKEN.REFRESHED", task, asynqEventBus)
	})
	logger.Info("✓ Registered handler for: AUTH.TOKEN.REFRESHED")

	// 注册日志任务处理器
	mux.HandleFunc("audit.log", auditLogHandler.ProcessTask)
	logger.Info("✓ Registered audit log handler for type: audit.log")

	mux.HandleFunc("activity.log", activityLogHandler.ProcessActivityLog)
	logger.Info("✓ Registered activity log handler for type: activity.log")

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

// processDomainEvent 处理领域事件（调用订阅的 Listener）
func processDomainEvent(ctx context.Context, eventType string, task *asynq.Task, eventBus *messaging.AsynqEventBus) error {
	handlers := eventBus.GetHandlers(eventType)
	if len(handlers) == 0 {
		logger.Warn("No handlers registered for event type: ", eventType)
		return nil
	}

	// 根据事件类型反序列化为具体领域事件
	var evt event.Event
	switch eventType {
	case "USER.REGISTERED":
		evt = &user.UserRegistered{}
	case "AUTH.LOGIN.SUCCESS":
		evt = &user.UserLoggedIn{}
	case "AUTH.LOGIN.FAILED":
		evt = &user.LoginFailed{}
	case "SECURITY.ACCOUNT.LOCKED":
		evt = &user.AccountLocked{}
	case "AUTH.LOGOUT":
		evt = &user.UserLoggedOut{}
	case "AUTH.TOKEN.REFRESHED":
		evt = &user.TokenRefreshed{}
	case "USER.PROFILE.UPDATED":
		evt = &user.UserProfileUpdated{}
	default:
		logger.Warn("Unknown event type: ", eventType)
		return nil
	}

	// 反序列化 payload
	if err := json.Unmarshal(task.Payload(), evt); err != nil {
		logger.Error("Failed to unmarshal event: ", eventType, ", error: ", err)
		return err
	}

	// 调用所有订阅的 handler
	for _, handler := range handlers {
		if err := handler(ctx, evt); err != nil {
			logger.Error("Handler failed for ", eventType, ": ", err)
			return err
		}
	}

	return nil
}
