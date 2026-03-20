package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
)

// UserProfileUpdatedEvent 用户资料更新事件
type UserProfileUpdatedEvent struct {
	*kernel.BaseEvent
	UserID        vo.UserID `json:"user_id"`
	UpdatedFields []string  `json:"updated_fields"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewUserProfileUpdatedEvent 创建用户资料更新事件
func NewUserProfileUpdatedEvent(userID vo.UserID, updatedFields []string) *UserProfileUpdatedEvent {
	event := &UserProfileUpdatedEvent{
		BaseEvent:     kernel.NewBaseEvent("UserProfileUpdated", userID, 1),
		UserID:        userID,
		UpdatedFields: updatedFields,
		UpdatedAt:     time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}
