package model

import "context"

// ActivityLogRepository 活动日志仓储接口
type ActivityLogRepository interface {
	// Save 保存活动日志
	Save(ctx context.Context, log *ActivityLog) error

	// FindByUserID 按用户 ID 查询
	FindByUserID(ctx context.Context, userID int64, limit int) ([]*ActivityLog, error)

	// FindByAction 按操作类型查询
	FindByAction(ctx context.Context, action ActivityType, limit int) ([]*ActivityLog, error)

	// FindFailed 查询失败的活动
	FindFailed(ctx context.Context, limit int) ([]*ActivityLog, error)
}
