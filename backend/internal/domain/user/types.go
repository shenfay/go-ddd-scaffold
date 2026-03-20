package user

import (
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
)

// ============================================================================
// Type aliases for backward compatibility - 保持向后兼容
// ============================================================================

// Core types
type UserID = vo.UserID
type UserName = vo.UserName
type Email = vo.Email
type HashedPassword = vo.HashedPassword
type UserStatus = vo.UserStatus
type UserGender = vo.UserGender
type User = aggregate.User
type UserBuilder = aggregate.UserBuilder

// Services
type PasswordHasher = service.PasswordHasher
type PasswordPolicy = service.PasswordPolicy
type PasswordPolicyConfig = service.PasswordPolicyConfig
type UserRepository = repository.UserRepository

// Events - 领域事件（统一从 event 包导出）
type UserRegisteredEvent = event.UserRegisteredEvent
type UserActivatedEvent = event.UserActivatedEvent
type UserDeactivatedEvent = event.UserDeactivatedEvent
type UserLoggedInEvent = event.UserLoggedInEvent
type UserPasswordChangedEvent = event.UserPasswordChangedEvent
type UserProfileUpdatedEvent = event.UserProfileUpdatedEvent
type UserEmailChangedEvent = event.UserEmailChangedEvent
type UserLockedEvent = event.UserLockedEvent
type UserUnlockedEvent = event.UserUnlockedEvent
type UserFailedLoginAttemptEvent = event.UserFailedLoginAttemptEvent

// Event constructors - 事件构造函数
var (
	NewUserRegisteredEvent         = event.NewUserRegisteredEvent
	NewUserActivatedEvent          = event.NewUserActivatedEvent
	NewUserDeactivatedEvent        = event.NewUserDeactivatedEvent
	NewUserLoggedInEvent           = event.NewUserLoggedInEvent
	NewUserPasswordChangedEvent    = event.NewUserPasswordChangedEvent
	NewUserProfileUpdatedEvent     = event.NewUserProfileUpdatedEvent
	NewUserEmailChangedEvent       = event.NewUserEmailChangedEvent
	NewUserLockedEvent             = event.NewUserLockedEvent
	NewUserUnlockedEvent           = event.NewUserUnlockedEvent
	NewUserFailedLoginAttemptEvent = event.NewUserFailedLoginAttemptEvent
)

// Constants
const (
	UserStatusPending  = vo.UserStatusPending
	UserStatusActive   = vo.UserStatusActive
	UserStatusInactive = vo.UserStatusInactive
	UserStatusLocked   = vo.UserStatusLocked

	UserGenderUnknown = vo.UserGenderUnknown
	UserGenderMale    = vo.UserGenderMale
	UserGenderFemale  = vo.UserGenderFemale
	UserGenderOther   = vo.UserGenderOther
)

// Constructor functions
var (
	NewUserID                   = vo.NewUserID
	NewUserName                 = vo.NewUserName
	NewEmail                    = vo.NewEmail
	NewHashedPassword           = vo.NewHashedPassword
	NewUser                     = aggregate.NewUser
	NewUserBuilder              = aggregate.NewUserBuilder
	NewBcryptPasswordHasher     = service.NewBcryptPasswordHasher
	DefaultPasswordPolicyConfig = service.DefaultPasswordPolicyConfig
)
