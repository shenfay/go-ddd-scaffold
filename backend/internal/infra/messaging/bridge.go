package messaging

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/events"
	"github.com/shenfay/go-ddd-scaffold/pkg/constants"
)

// DomainToIntegrationBridge 领域事件到集成事件的桥接器
// 订阅 InProcessBus 中的领域事件，转换为 Asynq 集成任务入队
type DomainToIntegrationBridge struct {
	client   *asynq.Client
	queueMap map[constants.EventName]constants.QueueName
}

// NewBridge 创建桥接器
func NewBridge(client *asynq.Client) *DomainToIntegrationBridge {
	return &DomainToIntegrationBridge{
		client: client,
		queueMap: map[constants.EventName]constants.QueueName{
			constants.EventUserRegistered:     constants.QueueDefault,
			constants.EventUserLoggedIn:       constants.QueueDefault,
			constants.EventUserLoginFailed:    constants.QueueCritical,
			constants.EventUserAccountLocked:  constants.QueueCritical,
			constants.EventUserLoggedOut:      constants.QueueDefault,
			constants.EventUserTokenRefreshed: constants.QueueDefault,
			constants.EventUserProfileUpdated: constants.QueueDefault,
		},
	}
}

// SubscribeTo 订阅 InProcessBus 中的所有领域事件
func (b *DomainToIntegrationBridge) SubscribeTo(bus events.Bus) {
	for eventName := range b.queueMap {
		name := eventName // capture for closure
		bus.Subscribe(string(name), func(ctx context.Context, evt events.DomainEvent) error {
			return b.enqueue(ctx, name, evt)
		})
	}
}

// enqueue 将领域事件入队到 Asynq
func (b *DomainToIntegrationBridge) enqueue(ctx context.Context, eventName constants.EventName, evt events.DomainEvent) error {
	payload, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	queue := b.queueMap[eventName]
	maxRetry := 3
	if queue == constants.QueueCritical {
		maxRetry = 5
	}

	_, err = b.client.EnqueueContext(ctx,
		asynq.NewTask(string(eventName), payload),
		asynq.Queue(string(queue)),
		asynq.MaxRetry(maxRetry),
	)
	return err
}
