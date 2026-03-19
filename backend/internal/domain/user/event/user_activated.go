package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserActivatedEvent 用户激活事件
type UserActivatedEvent struct {
	*kernel.BaseEvent
	UserID      valueobject.UserID `json:"user_id"`
	ActivatedAt time.Time          `json:"activated_at"`
}

// NewUserActivatedEvent 创建用户激活事件
func NewUserActivatedEvent(userID valueobject.UserID) *UserActivatedEvent {
	event := &UserActivatedEvent{
		BaseEvent:   kernel.NewBaseEvent("UserActivated", userID, 1),
		UserID:      userID,
		ActivatedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}
