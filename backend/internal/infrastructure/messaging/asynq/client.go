package asynq

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

// Config asynq 配置
type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// NewClient 创建 asynq 客户端
func NewClient(cfg Config) *asynq.Client {
	r := &asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	return asynq.NewClient(r)
}

// NewServer 创建 asynq 服务器
func NewServer(cfg Config) *asynq.Server {
	r := &asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	return asynq.NewServer(
		r,
		asynq.Config{
			Concurrency: 10, // 并发处理数
			Queues: map[string]int{
				"critical": 6, // 高优先级队列
				"default":  3, // 默认队列
				"low":      1, // 低优先级队列
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				// 这里可以集成日志系统
				fmt.Printf("Error processing task %s: %v\n", task.Type(), err)
			}),
		},
	)
}
