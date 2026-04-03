package activitylog

import (
	"time"
)

// ActivityType 活动类型
type ActivityType string

const (
	// 认证相关
	ActivityLogin          ActivityType = "LOGIN"           // 用户登录
	ActivityLogout         ActivityType = "LOGOUT"          // 用户登出
	ActivityRegister       ActivityType = "REGISTER"        // 用户注册
	ActivityRefreshToken   ActivityType = "REFRESH_TOKEN"   // 刷新 Token
	ActivityPasswordChange ActivityType = "PASSWORD_CHANGE" // 修改密码

	// 用户信息相关
	ActivityProfileUpdate ActivityType = "PROFILE_UPDATE" // 更新资料
	ActivityEmailVerify   ActivityType = "EMAIL_VERIFY"   // 验证邮箱

	// 系统相关
	ActivityAccountLock   ActivityType = "ACCOUNT_LOCK"   // 账户锁定
	ActivityAccountUnlock ActivityType = "ACCOUNT_UNLOCK" // 账户解锁
)

// ActivityStatus 活动状态
type ActivityStatus string

const (
	ActivitySuccess ActivityStatus = "SUCCESS"
	ActivityFailed  ActivityStatus = "FAILED"
)

// ActivityLog 活动日志
type ActivityLog struct {
	ID          string         `gorm:"primaryKey;type:varchar(50)" json:"id"`
	UserID      string         `gorm:"index;type:varchar(50);not null" json:"user_id"`
	Email       string         `gorm:"type:varchar(255)" json:"email"`
	Action      ActivityType   `gorm:"type:varchar(50);not null" json:"action"`
	Status      ActivityStatus `gorm:"type:varchar(20)" json:"status"`
	IP          string         `gorm:"type:varchar(45)" json:"ip"`          // IPv6 最大长度
	UserAgent   string         `gorm:"type:varchar(500)" json:"user_agent"` // 浏览器/设备信息
	Device      string         `gorm:"type:varchar(100)" json:"device"`     // 设备类型（mobile/desktop）
	Browser     string         `gorm:"type:varchar(50)" json:"browser"`     // 浏览器名称
	OS          string         `gorm:"type:varchar(50)" json:"os"`          // 操作系统
	Description string         `gorm:"type:text" json:"description"`        // 活动描述
	Metadata    string         `gorm:"type:json" json:"metadata,omitempty"` // 额外元数据（JSON 格式）
	CreatedAt   time.Time      `gorm:"index;not null" json:"created_at"`
}

// TableName 指定表名
func (ActivityLog) TableName() string {
	return "activity_logs"
}
