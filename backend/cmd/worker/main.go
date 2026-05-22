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

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/events"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/listener"
	"github.com/shenfay/go-ddd-scaffold/pkg/logger"
)

// eventRegistry 领域事件反序列化注册表
// 将事件名称映射到对应的空实例构造函数
var eventRegistry = map[string]func() events.DomainEvent{
	"user.registered":      func() events.DomainEvent { return &user.UserRegistered{} },
	"user.logged_in":       func() events.DomainEvent { return &user.UserLoggedIn{} },
	"user.login_failed":    func() events.DomainEvent { return &user.LoginFailed{} },
	"user.account_locked":  func() events.DomainEvent { return &user.AccountLocked{} },
	"user.logged_out":      func() events.DomainEvent { return &user.UserLoggedOut{} },
	"user.token_refreshed": func() events.DomainEvent { return &user.TokenRefreshed{} },
	"user.profile_updated": func() events.DomainEvent { return &user.UserProfileUpdated{} },
}

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

	logger.Info("Starting Asynq Worker...")

	// 3. 初始化 Redis 客户端和数据库
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis: ", err)
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	logger.Info("Redis connection established")

	// 初始化数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database: ", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Database connection established")

	// 4. 初始化仓储
	auditLogRepo := repository.NewAuditLogRepository(db)
	activityLogRepo := repository.NewActivityLogRepository(db)

	// 5. 创建领域事件监听器（直接写 DB）
	auditLogListener := listener.NewAuditLogListener(auditLogRepo)
	activityLogListener := listener.NewActivityLogListener(activityLogRepo)

	// 6. 注册 Asynq 任务处理器
	mux := asynq.NewServeMux()

	// 为每个领域事件注册通用处理器
	mux.HandleFunc("user.registered", createEventHandler(activityLogListener.HandleUserRegistered))
	mux.HandleFunc("user.logged_in", createEventHandler(auditLogListener.HandleUserLoggedIn))
	mux.HandleFunc("user.login_failed", createEventHandler(auditLogListener.HandleLoginFailed))
	mux.HandleFunc("user.account_locked", createEventHandler(auditLogListener.HandleAccountLocked))
	mux.HandleFunc("user.logged_out", createEventHandler(activityLogListener.HandleUserLoggedOut))
	mux.HandleFunc("user.token_refreshed", createEventHandler(activityLogListener.HandleTokenRefreshed))
	mux.HandleFunc("user.profile_updated", createEventHandler(auditLogListener.HandleUserLoggedIn))

	// 7. 创建 Asynq 服务器
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
			StrictPriority: true,
		},
	)

	logger.Info("Asynq server created with concurrency=", cfg.Asynq.Concurrency)

	// 8. 启动 Worker
	go func() {
		logger.Info("Starting Asynq Worker processor...")
		if err := srv.Run(mux); err != nil {
			logger.Error("Failed to run Asynq server: ", err)
			log.Fatalf("Failed to run Asynq server: %v", err)
		}
	}()

	// 9. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")

	// 10. 优雅关闭
	srv.Shutdown()
	logger.Info("Worker stopped gracefully")
}

// createEventHandler 创建 Asynq 任务处理器
// 使用 registry 反序列化事件载荷，然后调用 listener 处理
func createEventHandler(handler events.Handler) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		eventName := task.Type()

		constructor, ok := eventRegistry[eventName]
		if !ok {
			logger.Warn("Unknown event type: ", eventName)
			return nil
		}

		evt := constructor()
		if err := json.Unmarshal(task.Payload(), evt); err != nil {
			logger.Error("Failed to unmarshal event: ", eventName, ", error: ", err)
			return err
		}

		return handler(ctx, evt)
	}
}
