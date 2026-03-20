package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
)

// UserLockedEvent 用户锁定事件
type UserLockedEvent struct {
	*kernel.BaseEvent
	UserID      vo.UserID `json:"user_id"`
	Reason      string    `json:"reason"`
	LockedUntil time.Time `json:"locked_until"`
}

// NewUserLockedEvent 创建用户锁定事件
func NewUserLockedEvent(userID vo.UserID, reason string, lockedUntil time.Time) *UserLockedEvent {
	event := &UserLockedEvent{
		BaseEvent:   kernel.NewBaseEvent("UserLocked", userID, 1),
		UserID:      userID,
		Reason:      reason,
		LockedUntil: lockedUntil,
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}
