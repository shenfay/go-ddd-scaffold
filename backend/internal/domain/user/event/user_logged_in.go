package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
)

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	*kernel.BaseEvent
	UserID            vo.UserID `json:"user_id"`
	LoginAt           time.Time `json:"login_at"`
	IPAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	Location          string    `json:"location"`           // 地理位置
	DeviceType        string    `json:"device_type"`        // 设备类型
	DeviceFingerprint string    `json:"device_fingerprint"` // 设备指纹
	LoginMethod       string    `json:"login_method"`       // 登录方式：password/sms/email
	Success           bool      `json:"success"`            // 是否成功
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(userID vo.UserID, ipAddress, userAgent, location, deviceType, deviceFingerprint, loginMethod string, success bool) *UserLoggedInEvent {
	event := &UserLoggedInEvent{
		BaseEvent:         kernel.NewBaseEvent("UserLoggedIn", userID, 1),
		UserID:            userID,
		LoginAt:           time.Now(),
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		Location:          location,
		DeviceType:        deviceType,
		DeviceFingerprint: deviceFingerprint,
		LoginMethod:       loginMethod,
		Success:           success,
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}
