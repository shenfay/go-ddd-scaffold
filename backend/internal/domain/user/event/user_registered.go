package event

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	*kernel.BaseEvent
	UserID         valueobject.UserID `json:"user_id"`
	Username       string             `json:"username"`
	Email          string             `json:"email"`
	Status         string             `json:"status"`          // 用户状态
	DisplayName    string             `json:"display_name"`    // 显示名称
	RegistrationIP string             `json:"registration_ip"` // 注册 IP
	TenantID       int64              `json:"tenant_id"`       // 租户 ID
	RegisteredAt   time.Time          `json:"registered_at"`
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID valueobject.UserID, username, email, status, displayName, registrationIP string, tenantID int64) *UserRegisteredEvent {
	event := &UserRegisteredEvent{
		BaseEvent:      kernel.NewBaseEvent("UserRegistered", userID, 1),
		UserID:         userID,
		Username:       username,
		Email:          email,
		Status:         status,
		DisplayName:    displayName,
		RegistrationIP: registrationIP,
		TenantID:       tenantID,
		RegisteredAt:   time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}
