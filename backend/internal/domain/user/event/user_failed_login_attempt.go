package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
)

// UserFailedLoginAttemptEvent 用户失败登录尝试事件
type UserFailedLoginAttemptEvent struct {
	*kernel.BaseEvent
	UserID    vo.UserID `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Reason    string    `json:"reason"`
	AttemptAt time.Time `json:"attempt_at"`
}

// NewUserFailedLoginAttemptEvent 创建用户失败登录尝试事件
func NewUserFailedLoginAttemptEvent(userID vo.UserID, ipAddress, userAgent, reason string) *UserFailedLoginAttemptEvent {
	event := &UserFailedLoginAttemptEvent{
		BaseEvent: kernel.NewBaseEvent("UserFailedLoginAttempt", userID, 1),
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Reason:    reason,
		AttemptAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}
