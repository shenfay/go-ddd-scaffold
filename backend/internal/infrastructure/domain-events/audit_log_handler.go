package domainevents

import (
	"context"
	"encoding/json"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/audit"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// AuditLogHandler 审计日志事件处理器
type AuditLogHandler struct {
	repo      audit.AuditLogRepository
	snowflake *snowflake.Node
}

func NewAuditLogHandler(repo audit.AuditLogRepository, snowflake *snowflake.Node) *AuditLogHandler {
	return &AuditLogHandler{repo: repo, snowflake: snowflake}
}

// Handle 处理领域事件
func (h *AuditLogHandler) Handle(ctx context.Context, evt ddd.DomainEvent) error {
	switch e := evt.(type) {
	case *user.UserRegisteredEvent:
		return h.handleUserRegistered(ctx, e)
	case *user.UserLoggedInEvent:
		return h.handleUserLoggedIn(ctx, e)
	default:
		return nil // 忽略不关心的事件
	}
}

func (h *AuditLogHandler) handleUserRegistered(ctx context.Context, event *user.UserRegisteredEvent) error {
	metadata, _ := json.Marshal(map[string]interface{}{
		"username": event.Username,
		"email":    event.Email,
	})

	log := &audit.AuditLog{
		ID:           h.generateID(),
		UserID:       event.UserID.Int64(),
		Action:       "USER_REGISTERED",
		ResourceType: "User",
		ResourceID:   int64Ptr(event.UserID.Int64()),
		Metadata:     h.parseMetadata(metadata),
		Status:       audit.StatusSuccess,
		OccurredAt:   event.RegisteredAt,
	}

	return h.repo.Save(ctx, log)
}

func (h *AuditLogHandler) handleUserLoggedIn(ctx context.Context, event *user.UserLoggedInEvent) error {
	metadata, _ := json.Marshal(map[string]interface{}{
		"ip_address": event.IPAddress,
		"user_agent": event.UserAgent,
	})

	log := &audit.AuditLog{
		ID:           h.generateID(),
		UserID:       event.UserID.Int64(),
		Action:       "USER_LOGIN",
		ResourceType: "User",
		ResourceID:   int64Ptr(event.UserID.Int64()),
		IPAddress:    event.IPAddress,
		UserAgent:    event.UserAgent,
		Metadata:     h.parseMetadata(metadata),
		Status:       audit.StatusSuccess,
		OccurredAt:   event.LoginAt,
	}

	return h.repo.Save(ctx, log)
}

func (h *AuditLogHandler) generateID() int64 {
	if h.snowflake != nil {
		id, _ := h.snowflake.Generate()
		return id
	}
	return 0
}

func int64Ptr(i int64) *int64 {
	return &i
}

func (h *AuditLogHandler) parseMetadata(data []byte) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	return m
}
