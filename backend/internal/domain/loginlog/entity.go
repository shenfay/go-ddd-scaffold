package loginlog

import (
	"context"
	"time"
)

// LoginLog 登录日志实体
type LoginLog struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	TenantID      *int64    `json:"tenant_id"`
	LoginType     string    `json:"login_type"`
	LoginStatus   string    `json:"login_status"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	DeviceType    string    `json:"device_type"`
	OSInfo        string    `json:"os_info"`
	BrowserInfo   string    `json:"browser_info"`
	Country       string    `json:"country"`
	City          string    `json:"city"`
	FailureReason string    `json:"failure_reason"`
	IsSuspicious  bool      `json:"is_suspicious"`
	RiskScore     int       `json:"risk_score"`
	SessionID     string    `json:"session_id"`
	AccessTokenID string    `json:"access_token_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	CreatedAt     time.Time `json:"created_at"`
}

const (
	LoginStatusSuccess = "success"
	LoginStatusFailed  = "failed"
	LoginStatusBlocked = "blocked"
)

// LoginLogRepository 登录日志仓储接口
type LoginLogRepository interface {
	Save(ctx context.Context, log *LoginLog) error
	FindByUserID(ctx context.Context, userID int64, limit int) ([]*LoginLog, error)
	FindSuspiciousLogins(ctx context.Context, limit int) ([]*LoginLog, error)
}
