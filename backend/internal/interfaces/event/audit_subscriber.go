package event

import (
	"context"
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
)

// AuditSubscriber 审计日志事件订阅者
// 负责监听领域事件并记录活动日志
type AuditSubscriber struct {
	repo        aggregate.ActivityLogRepository
	idGenerator IDGenerator
}

// IDGenerator ID生成器接口
type IDGenerator interface {
	Generate() (int64, error)
}

// NewAuditSubscriber 创建审计日志事件订阅者
func NewAuditSubscriber(repo aggregate.ActivityLogRepository, idGenerator IDGenerator) *AuditSubscriber {
	return &AuditSubscriber{
		repo:        repo,
		idGenerator: idGenerator,
	}
}

// Handle 处理领域事件
func (s *AuditSubscriber) Handle(ctx context.Context, event kernel.DomainEvent) error {
	// TODO: 临时调试日志
	fmt.Printf("[AuditSubscriber] Handling event: %s, type: %T\n", event.EventName(), event)
	switch e := event.(type) {
	case *userEvent.UserRegisteredEvent:
		return s.handleUserRegistered(ctx, e)
	case *userEvent.UserLoggedInEvent:
		return s.handleUserLoggedIn(ctx, e)
	default:
		fmt.Printf("[AuditSubscriber] Unknown event type: %T\n", event)
		return nil // 忽略不关心的事件
	}
}

func (s *AuditSubscriber) handleUserRegistered(ctx context.Context, event *userEvent.UserRegisteredEvent) error {
	activity := aggregate.NewActivityLog(
		event.UserID.Int64(),
		aggregate.ActivityUserRegistered,
		aggregate.ActivityStatusSuccess,
	)
	activity.OccurredAt = event.RegisteredAt

	if id, err := s.idGenerator.Generate(); err == nil {
		activity.ID = id
	}

	activity.WithMetadata("username", event.Username)
	activity.WithMetadata("email", event.Email)

	return s.repo.Save(ctx, activity)
}

func (s *AuditSubscriber) handleUserLoggedIn(ctx context.Context, event *userEvent.UserLoggedInEvent) error {
	fmt.Printf("[AuditSubscriber] handleUserLoggedIn called for user: %s\n", event.UserID.String())

	activity := aggregate.NewActivityLog(
		event.UserID.Int64(),
		aggregate.ActivityUserLoggedIn,
		aggregate.ActivityStatusSuccess,
	)
	activity.OccurredAt = event.LoginAt

	if id, err := s.idGenerator.Generate(); err == nil {
		activity.ID = id
	}

	activity.WithMetadata("ip_address", event.IPAddress)
	activity.WithMetadata("user_agent", event.UserAgent)

	err := s.repo.Save(ctx, activity)
	if err != nil {
		fmt.Printf("[AuditSubscriber] Failed to save activity log: %v\n", err)
	} else {
		fmt.Printf("[AuditSubscriber] Activity log saved successfully\n")
	}
	return err
}

func (s *AuditSubscriber) generateID() int64 {
	if s.idGenerator != nil {
		id, _ := s.idGenerator.Generate()
		return id
	}
	return 0
}
