package event

import (
	"time"

	"go-ddd-scaffold/internal/domain/user/entity"

	"github.com/google/uuid"
)

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	UserID    uuid.UUID       `json:"userId"`
	Email     string          `json:"email"`
	Role      entity.UserRole `json:"role"`
	TenantID  *uuid.UUID      `json:"tenantId,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// GetEventType 获取事件类型
func (e *UserRegisteredEvent) GetEventType() string {
	return "UserRegistered"
}

// GetEventID 获取事件 ID（使用 UserID 作为事件 ID）
func (e *UserRegisteredEvent) GetEventID() string {
	return e.UserID.String()
}

// GetAggregateID 获取聚合根 ID
func (e *UserRegisteredEvent) GetAggregateID() uuid.UUID {
	return e.UserID
}

// GetOccurredAt 获取事件发生时间
func (e *UserRegisteredEvent) GetOccurredAt() time.Time {
	return e.Timestamp
}

// GetVersion 获取事件版本
func (e *UserRegisteredEvent) GetVersion() int {
	return 1
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	UserID        uuid.UUID `json:"userId"`
	Timestamp     time.Time `json:"timestamp"`
	IP            string    `json:"ip,omitempty"`
	UserAgent     string    `json:"userAgent,omitempty"`
	DeviceType    string    `json:"deviceType,omitempty"`
	OSInfo        *string   `json:"osInfo,omitempty"`
	BrowserInfo   *string   `json:"browserInfo,omitempty"`
	LoginStatus   string    `json:"loginStatus"`
	FailureReason *string   `json:"failureReason,omitempty"`
}

// GetEventType 获取事件类型
func (e *UserLoggedInEvent) GetEventType() string {
	return "UserLoggedIn"
}

// GetEventID 获取事件 ID（使用 UserID + 时间戳生成）
func (e *UserLoggedInEvent) GetEventID() string {
	return uuid.New().String()
}

// GetAggregateID 获取聚合根 ID
func (e *UserLoggedInEvent) GetAggregateID() uuid.UUID {
	return e.UserID
}

// GetOccurredAt 获取事件发生时间
func (e *UserLoggedInEvent) GetOccurredAt() time.Time {
	return e.Timestamp
}

// GetVersion 获取事件版本
func (e *UserLoggedInEvent) GetVersion() int {
	return 1
}
