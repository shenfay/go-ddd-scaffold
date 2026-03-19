package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserEmailChangedEvent 用户邮箱变更事件
type UserEmailChangedEvent struct {
	*kernel.BaseEvent
	UserID    valueobject.UserID `json:"user_id"`
	OldEmail  string             `json:"old_email"`
	NewEmail  string             `json:"new_email"`
	ChangedAt time.Time          `json:"changed_at"`
}

// NewUserEmailChangedEvent 创建用户邮箱变更事件
func NewUserEmailChangedEvent(userID valueobject.UserID, oldEmail, newEmail string) *UserEmailChangedEvent {
	event := &UserEmailChangedEvent{
		BaseEvent: kernel.NewBaseEvent("UserEmailChanged", userID, 1),
		UserID:    userID,
		OldEmail:  oldEmail,
		NewEmail:  newEmail,
		ChangedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}
