package repository

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
)

// ActivityLogRepository 活动日志仓储端口
type ActivityLogRepository interface {
	// Save 保存活动日志
	Save(ctx context.Context, log *model.ActivityLog) error

	// FindByUserID 根据用户 ID 查找
	FindByUserID(ctx context.Context, userID int64, limit int) ([]*model.ActivityLog, error)
}
