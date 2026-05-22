package messaging

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

// TaskPublisher 任务发布器
// 将任务入队到 Asynq 队列，供 Worker 消费
type TaskPublisher struct {
	client *asynq.Client
}

// NewTaskPublisher 创建任务发布器
func NewTaskPublisher(client *asynq.Client) *TaskPublisher {
	return &TaskPublisher{client: client}
}

// Publish 发布任务到 Asynq 队列
func (p *TaskPublisher) Publish(ctx context.Context, taskType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = p.client.EnqueueContext(ctx,
		asynq.NewTask(taskType, data),
		asynq.MaxRetry(3),
	)
	return err
}
