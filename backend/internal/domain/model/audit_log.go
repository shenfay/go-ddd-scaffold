package model

import (
	"context"
	"time"
)

// AuditLog 审计日志实体
type AuditLog struct {
	ID           int64                  `json:"id"`
	TenantID     *int64                 `json:"tenant_id"` // 租户 ID: NULL 表示系统级操作
	UserID       int64                  `json:"user_id"`   // 用户 ID
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   *int64                 `json:"resource_id"`
	RequestID    string                 `json:"request_id"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	Metadata     map[string]interface{} `json:"metadata"`
	Status       int16                  `json:"status"` // 0-成功，1-失败
	ErrorMessage string                 `json:"error_message"`
	OccurredAt   time.Time              `json:"occurred_at"`
	CreatedAt    time.Time              `json:"created_at"`
}

const (
	StatusSuccess = 0
	StatusFailed  = 1
)

// AuditLogRepository 审计日志仓储接口
type AuditLogRepository interface {
	Save(ctx context.Context, log *AuditLog) error
	FindByUserID(ctx context.Context, userID int64, limit int) ([]*AuditLog, error)
	FindByAction(ctx context.Context, action string, limit int) ([]*AuditLog, error)
}
