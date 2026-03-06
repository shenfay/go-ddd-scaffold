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
	EventID   string          `json:"eventId"`
	EventType string          `json:"eventType"`
	AggregateID uuid.UUID     `json:"aggregateId"`
	OccurredAt  time.Time     `json:"occurredAt"`
	Version     int           `json:"version"`
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID uuid.UUID, email string, role entity.UserRole, tenantID *uuid.UUID) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		UserID:      userID,
		Email:       email,
		Role:        role,
		TenantID:    tenantID,
		EventID:     uuid.New().String(),
		EventType:   "UserRegistered",
		AggregateID: userID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

// GetEventType 获取事件类型
func (e *UserRegisteredEvent) GetEventType() string {
	return e.EventType
}

// GetEventID 获取事件 ID
func (e *UserRegisteredEvent) GetEventID() string {
	return e.EventID
}

// GetAggregateID 获取聚合根 ID
func (e *UserRegisteredEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

// GetOccurredAt 获取事件发生时间
func (e *UserRegisteredEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetVersion 获取事件版本
func (e *UserRegisteredEvent) GetVersion() int {
	return e.Version
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	UserID        uuid.UUID `json:"userId"`
	IP            string    `json:"ip,omitempty"`
	UserAgent     string    `json:"userAgent,omitempty"`
	DeviceType    string    `json:"deviceType,omitempty"`
	OSInfo        *string   `json:"osInfo,omitempty"`
	BrowserInfo   *string   `json:"browserInfo,omitempty"`
	LoginStatus   string    `json:"loginStatus"`
	FailureReason *string   `json:"failureReason,omitempty"`
	EventID       string    `json:"eventId"`
	EventType     string    `json:"eventType"`
	AggregateID   uuid.UUID `json:"aggregateId"`
	OccurredAt    time.Time `json:"occurredAt"`
	Version       int       `json:"version"`
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(
	userID uuid.UUID,
	ip string,
	userAgent string,
	deviceType string,
	loginStatus string,
	failureReason *string,
) *UserLoggedInEvent {
	return &UserLoggedInEvent{
		UserID:        userID,
		IP:            ip,
		UserAgent:     userAgent,
		DeviceType:    deviceType,
		LoginStatus:   loginStatus,
		FailureReason: failureReason,
		EventID:       uuid.New().String(),
		EventType:     "UserLoggedIn",
		AggregateID:   userID,
		OccurredAt:    time.Now(),
		Version:       1,
	}
}

// GetEventType 获取事件类型
func (e *UserLoggedInEvent) GetEventType() string {
	return e.EventType
}

// GetEventID 获取事件 ID
func (e *UserLoggedInEvent) GetEventID() string {
	return e.EventID
}

// GetAggregateID 获取聚合根 ID
func (e *UserLoggedInEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

// GetOccurredAt 获取事件发生时间
func (e *UserLoggedInEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetVersion 获取事件版本
func (e *UserLoggedInEvent) GetVersion() int {
	return e.Version
}
