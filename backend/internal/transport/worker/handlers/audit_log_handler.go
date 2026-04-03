package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
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
	var data map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &data); err != nil {
		return err
	}

	log := &repository.AuditLog{
		UserID:   data["user_id"].(string),
		Action:   data["action"].(string),
		Status:   data["status"].(string),
		Metadata: data,
	}

	return h.repo.Save(ctx, log)
}
