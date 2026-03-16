package model

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// ============================================================================
// Domain Events - 领域事件
// ============================================================================
// 领域事件是聚合根状态变更的忠实记录，用于：
// 1. 发布到事件总线，触发后续业务流程
// 2. 持久化到 domain_events 表，用于审计和事件溯源
// 3. 在 CQRS 中作为 Read Model 的更新源（如实施）
// ============================================================================

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	*ddd.BaseEvent
	UserID         UserID    `json:"user_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Status         string    `json:"status"`          // 用户状态
	DisplayName    string    `json:"display_name"`    // 显示名称
	RegistrationIP string    `json:"registration_ip"` // 注册 IP
	TenantID       int64     `json:"tenant_id"`       // 租户 ID
	RegisteredAt   time.Time `json:"registered_at"`
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID UserID, username, email, status, displayName, registrationIP string, tenantID int64) *UserRegisteredEvent {
	event := &UserRegisteredEvent{
		BaseEvent:      ddd.NewBaseEvent("UserRegistered", userID, 1),
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

// UserActivatedEvent 用户激活事件
type UserActivatedEvent struct {
	*ddd.BaseEvent
	UserID      UserID    `json:"user_id"`
	ActivatedAt time.Time `json:"activated_at"`
}

// NewUserActivatedEvent 创建用户激活事件
func NewUserActivatedEvent(userID UserID) *UserActivatedEvent {
	event := &UserActivatedEvent{
		BaseEvent:   ddd.NewBaseEvent("UserActivated", userID, 1),
		UserID:      userID,
		ActivatedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}

// UserDeactivatedEvent 用户禁用事件
type UserDeactivatedEvent struct {
	*ddd.BaseEvent
	UserID        UserID    `json:"user_id"`
	Reason        string    `json:"reason"`
	DeactivatedAt time.Time `json:"deactivated_at"`
}

// NewUserDeactivatedEvent 创建用户禁用事件
func NewUserDeactivatedEvent(userID UserID, reason string) *UserDeactivatedEvent {
	event := &UserDeactivatedEvent{
		BaseEvent:     ddd.NewBaseEvent("UserDeactivated", userID, 1),
		UserID:        userID,
		Reason:        reason,
		DeactivatedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	*ddd.BaseEvent
	UserID            UserID    `json:"user_id"`
	LoginAt           time.Time `json:"login_at"`
	IPAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	Location          string    `json:"location"`           // 地理位置
	DeviceType        string    `json:"device_type"`        // 设备类型
	DeviceFingerprint string    `json:"device_fingerprint"` // 设备指纹
	LoginMethod       string    `json:"login_method"`       // 登录方式：password/sms/email
	Success           bool      `json:"success"`            // 是否成功
}

// NewUserLoggedInEvent 创建用户登录事件
func NewUserLoggedInEvent(userID UserID, ipAddress, userAgent, location, deviceType, deviceFingerprint, loginMethod string, success bool) *UserLoggedInEvent {
	event := &UserLoggedInEvent{
		BaseEvent:         ddd.NewBaseEvent("UserLoggedIn", userID, 1),
		UserID:            userID,
		LoginAt:           time.Now(),
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		Location:          location,
		DeviceType:        deviceType,
		DeviceFingerprint: deviceFingerprint,
		LoginMethod:       loginMethod,
		Success:           success,
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}

// UserPasswordChangedEvent 用户密码修改事件
type UserPasswordChangedEvent struct {
	*ddd.BaseEvent
	UserID    UserID    `json:"user_id"`
	ChangedAt time.Time `json:"changed_at"`
	IPAddress string    `json:"ip_address"`
}

// NewUserPasswordChangedEvent 创建用户密码修改事件
func NewUserPasswordChangedEvent(userID UserID, ipAddress string) *UserPasswordChangedEvent {
	event := &UserPasswordChangedEvent{
		BaseEvent: ddd.NewBaseEvent("UserPasswordChanged", userID, 1),
		UserID:    userID,
		ChangedAt: time.Now(),
		IPAddress: ipAddress,
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}

// UserProfileUpdatedEvent 用户资料更新事件
type UserProfileUpdatedEvent struct {
	*ddd.BaseEvent
	UserID        UserID    `json:"user_id"`
	UpdatedFields []string  `json:"updated_fields"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewUserProfileUpdatedEvent 创建用户资料更新事件
func NewUserProfileUpdatedEvent(userID UserID, updatedFields []string) *UserProfileUpdatedEvent {
	event := &UserProfileUpdatedEvent{
		BaseEvent:     ddd.NewBaseEvent("UserProfileUpdated", userID, 1),
		UserID:        userID,
		UpdatedFields: updatedFields,
		UpdatedAt:     time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}

// UserEmailChangedEvent 用户邮箱变更事件
type UserEmailChangedEvent struct {
	*ddd.BaseEvent
	UserID    UserID    `json:"user_id"`
	OldEmail  string    `json:"old_email"`
	NewEmail  string    `json:"new_email"`
	ChangedAt time.Time `json:"changed_at"`
}

// NewUserEmailChangedEvent 创建用户邮箱变更事件
func NewUserEmailChangedEvent(userID UserID, oldEmail, newEmail string) *UserEmailChangedEvent {
	event := &UserEmailChangedEvent{
		BaseEvent: ddd.NewBaseEvent("UserEmailChanged", userID, 1),
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

// UserLockedEvent 用户锁定事件
type UserLockedEvent struct {
	*ddd.BaseEvent
	UserID      UserID    `json:"user_id"`
	Reason      string    `json:"reason"`
	LockedUntil time.Time `json:"locked_until"`
}

// NewUserLockedEvent 创建用户锁定事件
func NewUserLockedEvent(userID UserID, reason string, lockedUntil time.Time) *UserLockedEvent {
	event := &UserLockedEvent{
		BaseEvent:   ddd.NewBaseEvent("UserLocked", userID, 1),
		UserID:      userID,
		Reason:      reason,
		LockedUntil: lockedUntil,
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}

// UserUnlockedEvent 用户解锁事件
type UserUnlockedEvent struct {
	*ddd.BaseEvent
	UserID     UserID    `json:"user_id"`
	UnlockedAt time.Time `json:"unlocked_at"`
}

// NewUserUnlockedEvent 创建用户解锁事件
func NewUserUnlockedEvent(userID UserID) *UserUnlockedEvent {
	event := &UserUnlockedEvent{
		BaseEvent:  ddd.NewBaseEvent("UserUnlocked", userID, 1),
		UserID:     userID,
		UnlockedAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	return event
}

// UserFailedLoginAttemptEvent 用户失败登录尝试事件
type UserFailedLoginAttemptEvent struct {
	*ddd.BaseEvent
	UserID    UserID    `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Reason    string    `json:"reason"`
	AttemptAt time.Time `json:"attempt_at"`
}

// NewUserFailedLoginAttemptEvent 创建用户失败登录尝试事件
func NewUserFailedLoginAttemptEvent(userID UserID, ipAddress, userAgent, reason string) *UserFailedLoginAttemptEvent {
	event := &UserFailedLoginAttemptEvent{
		BaseEvent: ddd.NewBaseEvent("UserFailedLoginAttempt", userID, 1),
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Reason:    reason,
		AttemptAt: time.Now(),
	}
	event.SetMetadata("event_type", "domain_event")
	event.SetMetadata("aggregate_type", "user")
	event.SetMetadata("security_event", true)
	return event
}
