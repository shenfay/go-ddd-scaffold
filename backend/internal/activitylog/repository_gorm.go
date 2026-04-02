package activitylog

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/shenfay/go-ddd-scaffold/pkg/utils/ulid"
)

// activityLogRepository GORM 实现
type activityLogRepository struct {
	db *gorm.DB
}

// NewActivityLogRepository 创建活动日志仓储
func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

// Create 创建活动日志
func (r *activityLogRepository) Create(ctx context.Context, log *ActivityLog) error {
	if log.ID == "" {
		log.ID = ulid.GenerateSessionID() // 使用 session ID 作为日志 ID
	}
	
	fmt.Printf("[DEBUG] Inserting activity log: ID=%s, UserID=%s, Action=%s\n", 
		log.ID, log.UserID, log.Action)
	
	result := r.db.WithContext(ctx).Create(log)
	
	if result.Error != nil {
		fmt.Printf("[ERROR] Failed to insert activity log: %v\n", result.Error)
		return result.Error
	}
	
	fmt.Printf("[DEBUG] Inserted %d rows\n", result.RowsAffected)
	
	return result.Error
}

// FindByUserID 根据用户 ID 查找活动日志
func (r *activityLogRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*ActivityLog, error) {
	var logs []*ActivityLog
	
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	
	return logs, err
}

// FindByTimeRange 根据时间范围查找活动日志
func (r *activityLogRepository) FindByTimeRange(ctx context.Context, start, end time.Time, limit, offset int) ([]*ActivityLog, error) {
	var logs []*ActivityLog
	
	err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	
	return logs, err
}

// FindRecent 查找最近的活动日志
func (r *activityLogRepository) FindRecent(ctx context.Context, userID string, limit int) ([]*ActivityLog, error) {
	var logs []*ActivityLog
	
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	
	return logs, err
}

// CountByUserID 统计用户的活动日志数量
func (r *activityLogRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	
	err := r.db.WithContext(ctx).
		Model(&ActivityLog{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	
	return count, err
}
