package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/pkg/logger"
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

	var payload struct {
		UserID      string                 `json:"user_id"`
		Email       string                 `json:"email"`
		Action      string                 `json:"action"`
		Status      string                 `json:"status"`
		IP          string                 `json:"ip"`
		UserAgent   string                 `json:"user_agent"`
		Device      string                 `json:"device"`
		Browser     string                 `json:"browser"`
		OS          string                 `json:"os"`
		Description string                 `json:"description"`
		Metadata    map[string]interface{} `json:"metadata"`
	}

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error("❌ Failed to unmarshal activity log payload: ", err)
		return err
	}

	log := &repository.ActivityLog{
		UserID:      payload.UserID,
		Email:       payload.Email,
		Action:      payload.Action,
		Status:      payload.Status,
		IP:          payload.IP,
		UserAgent:   payload.UserAgent,
		Device:      payload.Device,
		Browser:     payload.Browser,
		OS:          payload.OS,
		Description: payload.Description,
		Metadata:    "{}",
	}

	if err := h.repo.Create(ctx, log); err != nil {
		logger.Error("❌ Failed to create activity log: ", err)
		return err
	}

	logger.Info("✅ Activity log created: user_id=", payload.UserID, " action=", payload.Action)
	return nil
}
