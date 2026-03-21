package event

import (
	"context"
	"encoding/json"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/aggregate"
)

// AuditSubscriber 审计日志事件订阅者
// 负责监听领域事件并记录审计日志
type AuditSubscriber struct {
	repo        aggregate.AuditLogRepository
	idGenerator IDGenerator
}

// IDGenerator ID生成器接口
type IDGenerator interface {
	Generate() (int64, error)
}

// NewAuditSubscriber 创建审计日志事件订阅者
func NewAuditSubscriber(repo aggregate.AuditLogRepository, idGenerator IDGenerator) *AuditSubscriber {
	return &AuditSubscriber{
		repo:        repo,
		idGenerator: idGenerator,
	}
}

// Handle 处理领域事件
func (s *AuditSubscriber) Handle(ctx context.Context, event kernel.DomainEvent) error {
	switch e := event.(type) {
	case *userEvent.UserRegisteredEvent:
		return s.handleUserRegistered(ctx, e)
	case *userEvent.UserLoggedInEvent:
		return s.handleUserLoggedIn(ctx, e)
	default:
		return nil // 忽略不关心的事件
	}
}

func (s *AuditSubscriber) handleUserRegistered(ctx context.Context, event *userEvent.UserRegisteredEvent) error {
	metadata, _ := json.Marshal(map[string]interface{}{
		"username": event.Username,
		"email":    event.Email,
	})

	log := &aggregate.AuditLog{
		ID:           s.generateID(),
		UserID:       event.UserID.Int64(),
		Action:       "USER_REGISTERED",
		ResourceType: "User",
		ResourceID:   int64Ptr(event.UserID.Int64()),
		Metadata:     parseMetadata(metadata),
		Status:       aggregate.StatusSuccess,
		OccurredAt:   event.RegisteredAt,
	}

	return s.repo.Save(ctx, log)
}

func (s *AuditSubscriber) handleUserLoggedIn(ctx context.Context, event *userEvent.UserLoggedInEvent) error {
	metadata, _ := json.Marshal(map[string]interface{}{
		"ip_address": event.IPAddress,
		"user_agent": event.UserAgent,
	})

	log := &aggregate.AuditLog{
		ID:           s.generateID(),
		UserID:       event.UserID.Int64(),
		Action:       "USER_LOGIN",
		ResourceType: "User",
		ResourceID:   int64Ptr(event.UserID.Int64()),
		IPAddress:    event.IPAddress,
		UserAgent:    event.UserAgent,
		Metadata:     parseMetadata(metadata),
		Status:       aggregate.StatusSuccess,
		OccurredAt:   event.LoginAt,
	}

	return s.repo.Save(ctx, log)
}

func (s *AuditSubscriber) generateID() int64 {
	if s.idGenerator != nil {
		id, _ := s.idGenerator.Generate()
		return id
	}
	return 0
}

func int64Ptr(i int64) *int64 {
	return &i
}

func parseMetadata(data []byte) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	return m
}
