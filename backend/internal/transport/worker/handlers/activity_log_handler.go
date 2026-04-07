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

	log := &repository.ActivityLog{
		UserID:      utils.ToString(payload["user_id"]),
		Email:       utils.ToString(payload["email"]),
		Action:      utils.ToString(payload["action"]),
		Status:      utils.ToString(payload["status"]),
		IP:          utils.ToString(payload["ip"]),
		UserAgent:   utils.ToString(payload["user_agent"]),
		Device:      utils.ToString(payload["device"]),
		Browser:     utils.ToString(payload["browser"]),
		OS:          utils.ToString(payload["os"]),
		Description: utils.ToString(payload["description"]),
		Metadata:    payload,
	}

	if err := h.repo.Create(ctx, log); err != nil {
		logger.Error("❌ Failed to create activity log: ", err)
		return err
	}

	logger.Info("✅ Activity log created: user_id=", log.UserID, " action=", log.Action)
	return nil
}
