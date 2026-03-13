package commands

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// EventPublisher 领域事件发布器接口
type EventPublisher interface {
	Publish(ctx context.Context, event ddd.DomainEvent) error
}
