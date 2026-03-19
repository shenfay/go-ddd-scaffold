package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserDeactivatedEvent 用户禁用事件
type UserDeactivatedEvent struct {
	*kernel.BaseEvent
	UserID        valueobject.UserID `json:"user_id"`
	Reason        string             `json:"reason"`
	DeactivatedAt time.Time          `json:"deactivated_at"`
}

// NewUserDeactivatedEvent 创建用户禁用事件
func NewUserDeactivatedEvent(userID valueobject.UserID, reason string) *UserDeactivatedEvent {
	event := &UserDeactivatedEvent{
		BaseEvent:     kernel.NewBaseEvent("UserDeactivated", userID, 1),
		UserID:        userID,
		Reason:        reason,
		DeactivatedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}
