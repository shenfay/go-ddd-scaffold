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

	// 安全的类型转换
	userID := toString(data["user_id"])
	action := toString(data["action"])
	status := toString(data["status"])

	log := &repository.AuditLog{
		UserID:   userID,
		Action:   action,
		Status:   status,
		Metadata: data,
	}

	return h.repo.Save(ctx, log)
}

// toString 安全地将 interface{} 转换为 string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
