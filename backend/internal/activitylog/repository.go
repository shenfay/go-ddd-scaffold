package activitylog

import (
	"context"
	"time"
)

// ActivityLogRepository 活动日志仓储接口
type ActivityLogRepository interface {
	// Create 创建活动日志
	Create(ctx context.Context, log *ActivityLog) error
	
	// FindByUserID 根据用户 ID 查找活动日志
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*ActivityLog, error)
	
	// FindByTimeRange 根据时间范围查找活动日志
	FindByTimeRange(ctx context.Context, start, end time.Time, limit, offset int) ([]*ActivityLog, error)
	
	// FindRecent 查找最近的活动日志
	FindRecent(ctx context.Context, userID string, limit int) ([]*ActivityLog, error)
	
	// CountByUserID 统计用户的活动日志数量
	CountByUserID(ctx context.Context, userID string) (int64, error)
}
