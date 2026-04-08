package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/pkg/logger"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// ActivityLogWorkerHandler 活动日志Worker处理器
type ActivityLogWorkerHandler struct {
	repo repository.ActivityLogRepository
}

// NewActivityLogWorkerHandler 创建活动日志Worker处理器
func NewActivityLogWorkerHandler(repo repository.ActivityLogRepository) *ActivityLogWorkerHandler {
	return &ActivityLogWorkerHandler{
		repo: repo,
	}
}

// ProcessActivityLog 处理活动日志任务
func (h *ActivityLogWorkerHandler) ProcessActivityLog(ctx context.Context, task *asynq.Task) error {
	logger.Info("📝 Processing activity log task")

	var payload map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error("❌ Failed to unmarshal activity log payload: ", err)
		return err
	}

	// 轻量级设计：仅保留核心字段，其他信息存入 Metadata
	log := &repository.ActivityLog{
		UserID:   utils.ToString(payload["user_id"]),
		Action:   utils.ToString(payload["action"]),
		Metadata: payload, // 所有信息（IP、设备、描述等）统一存储在 Metadata 中
	}

	if err := h.repo.Create(ctx, log); err != nil {
		logger.Error("❌ Failed to create activity log: ", err)
		return err
	}

	logger.Info("✅ Activity log created: user_id=", log.UserID, " action=", log.Action)
	return nil
}
