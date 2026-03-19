package audit

import (
	"context"
	"encoding/json"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// EventHandler 审计日志领域事件处理器
type EventHandler struct {
	repo        AuditLogRepository
	idGenerator IDGenerator
}

// IDGenerator ID生成器接口
type IDGenerator interface {
	Generate() (int64, error)
}

// NewEventHandler 创建审计日志事件处理器
func NewEventHandler(repo AuditLogRepository, idGenerator IDGenerator) *EventHandler {
	return &EventHandler{
		repo:        repo,
		idGenerator: idGenerator,
	}
}

// Handle 处理领域事件
func (h *EventHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
	switch e := event.(type) {
	case *user.UserRegisteredEvent:
		return h.handleUserRegistered(ctx, e)
	case *user.UserLoggedInEvent:
		return h.handleUserLoggedIn(ctx, e)
	default:
		return nil // 忽略不关心的事件
	}
}

func (h *EventHandler) handleUserRegistered(ctx context.Context, event *user.UserRegisteredEvent) error {
	metadata, _ := json.Marshal(map[string]interface{}{
		"username": event.Username,
		"email":    event.Email,
	})

	log := &AuditLog{
		ID:           h.generateID(),
		UserID:       event.UserID.Int64(),
		Action:       "USER_REGISTERED",
		ResourceType: "User",
		ResourceID:   int64Ptr(event.UserID.Int64()),
		Metadata:     parseMetadata(metadata),
		Status:       StatusSuccess,
		OccurredAt:   event.RegisteredAt,
	}

	return h.repo.Save(ctx, log)
}

func (h *EventHandler) handleUserLoggedIn(ctx context.Context, event *user.UserLoggedInEvent) error {
	metadata, _ := json.Marshal(map[string]interface{}{
		"ip_address": event.IPAddress,
		"user_agent": event.UserAgent,
	})

	log := &AuditLog{
		ID:           h.generateID(),
		UserID:       event.UserID.Int64(),
		Action:       "USER_LOGIN",
		ResourceType: "User",
		ResourceID:   int64Ptr(event.UserID.Int64()),
		IPAddress:    event.IPAddress,
		UserAgent:    event.UserAgent,
		Metadata:     parseMetadata(metadata),
		Status:       StatusSuccess,
		OccurredAt:   event.LoginAt,
	}

	return h.repo.Save(ctx, log)
}

func (h *EventHandler) generateID() int64 {
	if h.idGenerator != nil {
		id, _ := h.idGenerator.Generate()
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
