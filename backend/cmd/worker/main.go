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

	"github.com/shenfay/go-ddd-scaffold/internal/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"github.com/shenfay/go-ddd-scaffold/pkg/constants"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("development")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. 初始化 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established")

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

	// 4. 注册任务处理器
	mux := asynq.NewServeMux()
	
	// 认证相关任务
	mux.HandleFunc(constants.AsynqTaskSendVerificationEmail, auth.NewSendVerificationEmailHandler())
	mux.HandleFunc(constants.AsynqTaskSendWelcomeEmail, auth.NewSendWelcomeEmailHandler())
	mux.HandleFunc(constants.AsynqTaskLogUserRegistration, auth.NewLogUserRegistrationHandler())
	mux.HandleFunc(constants.AsynqTaskLogLoginAttempt, auth.NewLogLoginAttemptHandler())
	mux.HandleFunc(constants.AsynqTaskCleanupExpiredTokens, auth.NewCleanupExpiredTokensHandler(redisClient))

	// 5. 启动 Worker
	go func() {
		log.Printf("Starting Asynq Worker with concurrency=%d", cfg.Asynq.Concurrency)
		if err := srv.Run(mux); err != nil {
			log.Fatalf("Failed to run Asynq server: %v", err)
		}
	}()

	// 6. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	// 7. 优雅关闭
	srv.Shutdown()
	log.Println("Worker stopped gracefully")
}
