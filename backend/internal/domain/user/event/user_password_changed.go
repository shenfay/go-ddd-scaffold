package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserPasswordChangedEvent 用户密码修改事件
type UserPasswordChangedEvent struct {
	*kernel.BaseEvent
	UserID    valueobject.UserID `json:"user_id"`
	ChangedAt time.Time          `json:"changed_at"`
	IPAddress string             `json:"ip_address"`
}

// NewUserPasswordChangedEvent 创建用户密码修改事件
func NewUserPasswordChangedEvent(userID valueobject.UserID, ipAddress string) *UserPasswordChangedEvent {
	event := &UserPasswordChangedEvent{
		BaseEvent: kernel.NewBaseEvent("UserPasswordChanged", userID, 1),
		UserID:    userID,
		ChangedAt: time.Now(),
		IPAddress: ipAddress,
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}
