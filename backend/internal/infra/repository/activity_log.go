package repository

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
	"gorm.io/gorm"
)

// ActivityLog 活动日志实体（轻量级设计）
type ActivityLog struct {
	ID        string                 `json:"id" gorm:"primaryKey"`
	UserID    string                 `json:"user_id" gorm:"not null;index:idx_user_id"`
	Action    string                 `json:"action" gorm:"not null;index:idx_action"`
	Metadata  map[string]interface{} `json:"metadata" gorm:"type:jsonb;default:'{}'::jsonb;serializer:json"`
	CreatedAt time.Time              `json:"created_at" gorm:"not null;index:idx_created_at"`
}

// TableName 指定表名
func (ActivityLog) TableName() string {
	return "activity_logs"
}

// ActivityLogRepository 活动日志仓储接口
type ActivityLogRepository interface {
	// Create 创建活动日志
	Create(ctx context.Context, log *ActivityLog) error

	// FindByUserID 根据用户 ID 查找日志
	FindByUserID(ctx context.Context, userID string, limit int, offset int) ([]*ActivityLog, error)

	// FindRecent 查找最近的日志
	FindRecent(ctx context.Context, userID string, limit int) ([]*ActivityLog, error)
}

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
		log.ID = utils.GenerateID()
	}

	return r.db.WithContext(ctx).Create(log).Error
}

// FindByUserID 根据用户 ID 查找日志
func (r *activityLogRepository) FindByUserID(ctx context.Context, userID string, limit int, offset int) ([]*ActivityLog, error) {
	var logs []*ActivityLog

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, err
}

// FindRecent 查找最近的日志
func (r *activityLogRepository) FindRecent(ctx context.Context, userID string, limit int) ([]*ActivityLog, error) {
	var logs []*ActivityLog

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}
