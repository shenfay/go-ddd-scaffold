package asynq

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

// TaskType 定义任务类型常量
const (
	TaskTypeDomainEvent = "domain:event" // 领域事件任务类型
)

// DomainEventPayload 领域事件任务负载
type DomainEventPayload struct {
	AggregateID   string            `json:"aggregate_id"`
	AggregateType string            `json:"aggregate_type"`
	EventType     string            `json:"event_type"`
	EventVersion  int32             `json:"event_version"`
	EventData     json.RawMessage   `json:"event_data"`
	OccurredOn    string            `json:"occurred_on"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// Publisher asynq 任务发布器
type Publisher struct {
	client *asynq.Client
}

// NewPublisher 创建任务发布器
func NewPublisher(client *asynq.Client) *Publisher {
	return &Publisher{client: client}
}

// PublishDomainEvent 发布领域事件任务
func (p *Publisher) PublishDomainEvent(ctx context.Context, payload DomainEventPayload, queue string) error {
	task, err := NewDomainEventTask(payload)
	if err != nil {
		return err
	}

	_, err = p.client.EnqueueContext(ctx, task, asynq.Queue(queue))
	return err
}

// NewDomainEventTask 创建领域事件任务
func NewDomainEventTask(payload DomainEventPayload) (*asynq.Task, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskTypeDomainEvent, b), nil
}

// ExtractDomainEventPayload 从任务中提取领域事件负载
func ExtractDomainEventPayload(task *asynq.Task) (*DomainEventPayload, error) {
	var payload DomainEventPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
