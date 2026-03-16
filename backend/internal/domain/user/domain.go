package user

import (
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/model"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
)

// ============================================================================
// Type aliases for backward compatibility - 保持向后兼容
// ============================================================================

// Core types
type UserID = model.UserID
type UserName = model.UserName
type Email = model.Email
type HashedPassword = model.HashedPassword
type UserStatus = model.UserStatus
type UserGender = model.UserGender
type User = model.User
type UserBuilder = model.UserBuilder

// Services
type PasswordHasher = service.PasswordHasher
type UserRepository = repository.UserRepository

// Events - 领域事件（统一从 model 包导出）
type UserRegisteredEvent = model.UserRegisteredEvent
type UserActivatedEvent = model.UserActivatedEvent
type UserDeactivatedEvent = model.UserDeactivatedEvent
type UserLoggedInEvent = model.UserLoggedInEvent
type UserPasswordChangedEvent = model.UserPasswordChangedEvent
type UserProfileUpdatedEvent = model.UserProfileUpdatedEvent
type UserEmailChangedEvent = model.UserEmailChangedEvent
type UserLockedEvent = model.UserLockedEvent
type UserUnlockedEvent = model.UserUnlockedEvent
type UserFailedLoginAttemptEvent = model.UserFailedLoginAttemptEvent

// Event constructors - 事件构造函数
var (
	NewUserRegisteredEvent         = model.NewUserRegisteredEvent
	NewUserActivatedEvent          = model.NewUserActivatedEvent
	NewUserDeactivatedEvent        = model.NewUserDeactivatedEvent
	NewUserLoggedInEvent           = model.NewUserLoggedInEvent
	NewUserPasswordChangedEvent    = model.NewUserPasswordChangedEvent
	NewUserProfileUpdatedEvent     = model.NewUserProfileUpdatedEvent
	NewUserEmailChangedEvent       = model.NewUserEmailChangedEvent
	NewUserLockedEvent             = model.NewUserLockedEvent
	NewUserUnlockedEvent           = model.NewUserUnlockedEvent
	NewUserFailedLoginAttemptEvent = model.NewUserFailedLoginAttemptEvent
)

// Constants
const (
	UserStatusPending  = model.UserStatusPending
	UserStatusActive   = model.UserStatusActive
	UserStatusInactive = model.UserStatusInactive
	UserStatusLocked   = model.UserStatusLocked

	UserGenderUnknown = model.UserGenderUnknown
	UserGenderMale    = model.UserGenderMale
	UserGenderFemale  = model.UserGenderFemale
	UserGenderOther   = model.UserGenderOther
)

// Constructor functions
var (
	NewUserID               = model.NewUserID
	NewUserName             = model.NewUserName
	NewEmail                = model.NewEmail
	NewHashedPassword       = model.NewHashedPassword
	NewUser                 = model.NewUser
	NewUserBuilder          = model.NewUserBuilder
	NewBcryptPasswordHasher = service.NewBcryptPasswordHasher
)
