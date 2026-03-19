package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserUnlockedEvent 用户解锁事件
type UserUnlockedEvent struct {
	*kernel.BaseEvent
	UserID     valueobject.UserID `json:"user_id"`
	UnlockedAt time.Time          `json:"unlocked_at"`
}

// NewUserUnlockedEvent 创建用户解锁事件
func NewUserUnlockedEvent(userID valueobject.UserID) *UserUnlockedEvent {
	event := &UserUnlockedEvent{
		BaseEvent:  kernel.NewBaseEvent("UserUnlocked", userID, 1),
		UserID:     userID,
		UnlockedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}
