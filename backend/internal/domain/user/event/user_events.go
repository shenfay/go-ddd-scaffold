package event

import (
	"time"

	"github.com/google/uuid"
)

// UserLockedEvent 用户被锁定事件
type UserLockedEvent struct {
	UserID      uuid.UUID
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewUserLockedEvent 创建用户锁定事件
func NewUserLockedEvent(userID uuid.UUID) *UserLockedEvent {
	return &UserLockedEvent{
		UserID:      userID,
		EventID:     uuid.New().String(),
		EventType:   "UserLocked",
		AggregateID: userID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

func (e *UserLockedEvent) GetEventType() string      { return e.EventType }
func (e *UserLockedEvent) GetEventID() string        { return e.EventID }
func (e *UserLockedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *UserLockedEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *UserLockedEvent) GetVersion() int           { return e.Version }

// UserActivatedEvent 用户被激活事件
type UserActivatedEvent struct {
	UserID      uuid.UUID
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewUserActivatedEvent 创建用户激活事件
func NewUserActivatedEvent(userID uuid.UUID) *UserActivatedEvent {
	return &UserActivatedEvent{
		UserID:      userID,
		EventID:     uuid.New().String(),
		EventType:   "UserActivated",
		AggregateID: userID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

func (e *UserActivatedEvent) GetEventType() string      { return e.EventType }
func (e *UserActivatedEvent) GetEventID() string        { return e.EventID }
func (e *UserActivatedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *UserActivatedEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *UserActivatedEvent) GetVersion() int           { return e.Version }

// UserProfileUpdatedEvent 用户资料更新事件
type UserProfileUpdatedEvent struct {
	UserID      uuid.UUID
	Nickname    string
	Phone       *string
	Bio         *string
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewUserProfileUpdatedEvent 创建用户资料更新事件
func NewUserProfileUpdatedEvent(
	userID uuid.UUID,
	nickname string,
	phone *string,
	bio *string,
) *UserProfileUpdatedEvent {
	return &UserProfileUpdatedEvent{
		UserID:      userID,
		Nickname:    nickname,
		Phone:       phone,
		Bio:         bio,
		EventID:     uuid.New().String(),
		EventType:   "UserProfileUpdated",
		AggregateID: userID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

func (e *UserProfileUpdatedEvent) GetEventType() string      { return e.EventType }
func (e *UserProfileUpdatedEvent) GetEventID() string        { return e.EventID }
func (e *UserProfileUpdatedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *UserProfileUpdatedEvent) GetOccurredAt() time.Time  { return e.OccurredAt }
func (e *UserProfileUpdatedEvent) GetVersion() int           { return e.Version }

// UserEmailChangedEvent 用户邮箱变更事件
type UserEmailChangedEvent struct {
	UserID      uuid.UUID
	OldEmail    string
	NewEmail    string
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewUserEmailChangedEvent 创建用户邮箱变更事件
func NewUserEmailChangedEvent(
	userID uuid.UUID,
	oldEmail, newEmail string,
) *UserEmailChangedEvent {
	return &UserEmailChangedEvent{
		UserID:      userID,
		OldEmail:    oldEmail,
		NewEmail:    newEmail,
		EventID:     uuid.New().String(),
		EventType:   "UserEmailChanged",
		AggregateID: userID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

func (e *UserEmailChangedEvent) GetEventType() string     { return e.EventType }
func (e *UserEmailChangedEvent) GetEventID() string       { return e.EventID }
func (e *UserEmailChangedEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *UserEmailChangedEvent) GetOccurredAt() time.Time { return e.OccurredAt }
func (e *UserEmailChangedEvent) GetVersion() int          { return e.Version }

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	UserID      uuid.UUID
	Email       string
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID uuid.UUID, email string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		UserID:      userID,
		Email:       email,
		EventID:     uuid.New().String(),
		EventType:   "UserRegistered",
		AggregateID: userID,
		OccurredAt:  time.Now(),
		Version:     1,
	}
}

func (e *UserRegisteredEvent) GetEventType() string     { return e.EventType }
func (e *UserRegisteredEvent) GetEventID() string       { return e.EventID }
func (e *UserRegisteredEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *UserRegisteredEvent) GetOccurredAt() time.Time { return e.OccurredAt }
func (e *UserRegisteredEvent) GetVersion() int          { return e.Version }

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	UserID      uuid.UUID
	IP          string
	UserAgent   string
	DeviceType  string
	LoginStatus string
	FailureReason *string
	EventID     string
	EventType   string
	AggregateID uuid.UUID
	OccurredAt  time.Time
	Version     int
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(
	userID uuid.UUID,
	ip, userAgent, deviceType, loginStatus string,
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

func (e *UserLoggedInEvent) GetEventType() string     { return e.EventType }
func (e *UserLoggedInEvent) GetEventID() string       { return e.EventID }
func (e *UserLoggedInEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *UserLoggedInEvent) GetOccurredAt() time.Time { return e.OccurredAt }
func (e *UserLoggedInEvent) GetVersion() int          { return e.Version }
