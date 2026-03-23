package aggregate

import (
	"context"
	"time"
)

// ActivityType 活动类型
type ActivityType string

const (
	// 用户相关
	ActivityUserRegistered      ActivityType = "USER_REGISTERED"
	ActivityUserLoggedIn        ActivityType = "USER_LOGIN"
	ActivityUserLoggedOut       ActivityType = "USER_LOGOUT"
	ActivityUserActivated       ActivityType = "USER_ACTIVATED"
	ActivityUserDeactivated     ActivityType = "USER_DEACTIVATED"
	ActivityUserLocked          ActivityType = "USER_LOCKED"
	ActivityUserUnlocked        ActivityType = "USER_UNLOCKED"
	ActivityUserPasswordChanged ActivityType = "USER_PASSWORD_CHANGED"
	ActivityUserEmailChanged    ActivityType = "USER_EMAIL_CHANGED"
	ActivityUserProfileUpdated  ActivityType = "USER_PROFILE_UPDATED"

	// 订单相关（示例）
	ActivityOrderCreated   ActivityType = "ORDER_CREATED"
	ActivityOrderPaid      ActivityType = "ORDER_PAID"
	ActivityOrderShipped   ActivityType = "ORDER_SHIPPED"
	ActivityOrderCancelled ActivityType = "ORDER_CANCELLED"

	// 系统相关
	ActivitySystemError   ActivityType = "SYSTEM_ERROR"
	ActivitySecurityAlert ActivityType = "SECURITY_ALERT"
)

// ActivityStatus 活动状态
type ActivityStatus int16

const (
	ActivityStatusSuccess ActivityStatus = 0 // 成功
	ActivityStatusFailed  ActivityStatus = 1 // 失败
)

// ActivityLog 活动日志实体（统一记录所有活动）
type ActivityLog struct {
	ID         int64          `json:"id"`
	TenantID   *int64         `json:"tenant_id"`   // 租户 ID: NULL 表示系统级操作
	UserID     int64          `json:"user_id"`     // 用户 ID
	Action     ActivityType   `json:"action"`      // 操作类型
	Status     ActivityStatus `json:"status"`      // 状态
	IPAddress  string         `json:"ip_address"`  // IP 地址
	UserAgent  string         `json:"user_agent"`  // User-Agent
	Metadata   map[string]any `json:"metadata"`    // 元数据
	OccurredAt time.Time      `json:"occurred_at"` // 发生时间
	CreatedAt  time.Time      `json:"created_at"`  // 创建时间
}

// NewActivityLog 创建活动日志
func NewActivityLog(userID int64, action ActivityType, status ActivityStatus) *ActivityLog {
	return &ActivityLog{
		UserID:     userID,
		Action:     action,
		Status:     status,
		Metadata:   make(map[string]any),
		OccurredAt: time.Now(),
		CreatedAt:  time.Now(),
	}
}

// WithTenantID 设置租户 ID（链式调用）
func (l *ActivityLog) WithTenantID(tenantID int64) *ActivityLog {
	l.TenantID = &tenantID
	return l
}

// WithIPAddress 设置 IP 地址
func (l *ActivityLog) WithIPAddress(ip string) *ActivityLog {
	l.IPAddress = ip
	return l
}

// WithUserAgent 设置 User-Agent
func (l *ActivityLog) WithUserAgent(ua string) *ActivityLog {
	l.UserAgent = ua
	return l
}

// WithMetadata 添加元数据
func (l *ActivityLog) WithMetadata(key string, value any) *ActivityLog {
	if l.Metadata == nil {
		l.Metadata = make(map[string]any)
	}
	l.Metadata[key] = value
	return l
}

// WithOccurrenceTime 设置发生时间
func (l *ActivityLog) WithOccurrenceTime(t time.Time) *ActivityLog {
	l.OccurredAt = t
	return l
}

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
