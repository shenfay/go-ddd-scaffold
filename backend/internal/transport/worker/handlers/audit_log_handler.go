package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

// AuditLogHandler 审计日志处理器
type AuditLogHandler struct {
	// TODO: 注入 repository
}

// NewAuditLogHandler 创建审计日志处理器
func NewAuditLogHandler() *AuditLogHandler {
	return &AuditLogHandler{}
}

// ProcessTask 处理审计日志任务
func (h *AuditLogHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var data map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &data); err != nil {
		return err
	}

	// TODO: 调用 repository 保存审计日志
	// log := &AuditLog{...}
	// return h.repo.Save(ctx, log)

	return nil
}
