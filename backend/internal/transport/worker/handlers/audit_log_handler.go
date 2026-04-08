package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// AuditLogHandler 审计日志处理器
type AuditLogHandler struct {
	repo repository.AuditLogRepository
}

// NewAuditLogHandler 创建审计日志处理器
func NewAuditLogHandler(repo repository.AuditLogRepository) *AuditLogHandler {
	return &AuditLogHandler{repo: repo}
}

// ProcessTask 处理审计日志任务
func (h *AuditLogHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	log := &repository.AuditLog{
		UserID:    utils.ToString(payload["user_id"]),
		Email:     utils.ToString(payload["email"]),
		Action:    utils.ToString(payload["action"]),
		Status:    utils.ToString(payload["status"]),
		IP:        utils.ToString(payload["ip"]),
		UserAgent: utils.ToString(payload["user_agent"]),
		Device:    utils.ToString(payload["device"]),
		Browser:   utils.ToString(payload["browser"]),
		OS:        utils.ToString(payload["os"]),
		Metadata:  payload,
	}

	return h.repo.Save(ctx, log)
}
