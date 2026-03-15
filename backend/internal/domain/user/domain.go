package user

import (
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
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

// Events
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

// Event constructors
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
