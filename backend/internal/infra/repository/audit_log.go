package repository

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/pkg/utils/ulid"
	"gorm.io/gorm"
)

// AuditLog 审计日志实体
type AuditLog struct {
	ID          string                 `json:"id" gorm:"primaryKey"`
	UserID      string                 `json:"user_id" gorm:"not null;index:idx_user_id"`
	Email       string                 `json:"email" gorm:"index:idx_email"`
	Action      string                 `json:"action" gorm:"not null;index:idx_action"`
	Status      string                 `json:"status" gorm:"not null;index:idx_status"`
	IP          string                 `json:"ip" gorm:"size:45"`
	UserAgent   string                 `json:"user_agent" gorm:"size:500"`
	Device      string                 `json:"device" gorm:"size:100"`
	Browser     string                 `json:"browser" gorm:"size:50"`
	OS          string                 `json:"os" gorm:"size:50"`
	Description string                 `json:"description" gorm:"type:text"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb;default:'{}'::jsonb;serializer:json"`
	CreatedAt   time.Time              `json:"created_at" gorm:"not null;index:idx_created_at"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogRepository 审计日志仓储接口
type AuditLogRepository interface {
	// Save 保存审计日志
	Save(ctx context.Context, log *AuditLog) error

	// FindByUserID 根据用户 ID 查找日志
	FindByUserID(ctx context.Context, userID string, limit int, offset int) ([]*AuditLog, error)

	// FindByAction 根据操作类型查找
	FindByAction(ctx context.Context, action string, limit int, offset int) ([]*AuditLog, error)
}

// auditLogRepository GORM 实现
type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository 创建审计日志仓储
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Save 保存审计日志
func (r *auditLogRepository) Save(ctx context.Context, log *AuditLog) error {
	if log.ID == "" {
		log.ID = ulid.GenerateAuditLogID()
	}

	return r.db.WithContext(ctx).Create(log).Error
}

// FindByUserID 根据用户 ID 查找日志
func (r *auditLogRepository) FindByUserID(ctx context.Context, userID string, limit int, offset int) ([]*AuditLog, error) {
	var logs []*AuditLog

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, err
}

// FindByAction 根据操作类型查找
func (r *auditLogRepository) FindByAction(ctx context.Context, action string, limit int, offset int) ([]*AuditLog, error) {
	var logs []*AuditLog

	err := r.db.WithContext(ctx).
		Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, err
}
