package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/activitylog"
	"github.com/shenfay/go-ddd-scaffold/internal/asynq/tasks"
	"github.com/shenfay/go-ddd-scaffold/pkg/logger"
)

// ActivityLogHandler 活动日志任务处理器
type ActivityLogHandler struct {
	repo activitylog.ActivityLogRepository
}

// NewActivityLogHandler 创建活动日志处理器
func NewActivityLogHandler(repo activitylog.ActivityLogRepository) *ActivityLogHandler {
	return &ActivityLogHandler{repo: repo}
}

// HandleActivityLogRecord 处理活动日志记录任务
func (h *ActivityLogHandler) HandleActivityLogRecord(ctx context.Context, t *asynq.Task) error {
	var p tasks.ActivityLogPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		logger.Error("Failed to unmarshal activity log payload: ", err)
		return err
	}

	log := &activitylog.ActivityLog{
		ID:          "", // 由仓储层自动生成
		UserID:      p.UserID,
		Email:       p.Email,
		Action:      activitylog.ActivityType(p.Action),
		Status:      activitylog.ActivityStatus(p.Status),
		IP:          p.IP,
		UserAgent:   p.UserAgent,
		Device:      p.Device,
		Browser:     p.Browser,
		OS:          p.OS,
		Description: p.Description,
		Metadata:    "{}",
		CreatedAt:   time.Now(),
	}

	if err := h.repo.Create(ctx, log); err != nil {
		logger.Error("Failed to create activity log: ", err)
		return err
	}

	logger.Info("✓ Activity log created: user_id=", p.UserID, " action=", p.Action)
	return nil
}
