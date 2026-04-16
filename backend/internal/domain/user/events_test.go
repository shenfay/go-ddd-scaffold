package user_test

import (
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestUserRegistered_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"

	event := user.NewUserRegisteredEvent(userID, email)

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, "USER.REGISTERED", event.GetType())
	assert.NotNil(t, event.GetPayload())
	assert.Equal(t, event, event.GetPayload())
}

func TestUserLoggedIn_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"
	ip := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	device := "web"

	event := user.NewUserLoggedInEvent(userID, email, ip, userAgent, device)

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, ip, event.IP)
	assert.Equal(t, userAgent, event.UserAgent)
	assert.Equal(t, device, event.Device)
	assert.Equal(t, "AUTH.LOGIN.SUCCESS", event.GetType())
}

func TestLoginFailed_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"
	ip := "192.168.1.1"
	reason := "invalid_password"

	event := user.NewLoginFailedEvent(userID, email, ip, reason)

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, ip, event.IP)
	assert.Equal(t, reason, event.Reason)
	assert.Equal(t, "AUTH.LOGIN.FAILED", event.GetType())
}

func TestAccountLocked_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"
	failedAttempts := 5
	lockedUntil := time.Now().Add(30 * time.Minute)

	event := user.NewAccountLockedEvent(userID, email, failedAttempts, lockedUntil)

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, failedAttempts, event.FailedAttempts)
	assert.Equal(t, lockedUntil, event.LockedUntil)
	assert.Equal(t, "SECURITY.ACCOUNT.LOCKED", event.GetType())
}

func TestUserLoggedOut_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"

	event := user.NewUserLoggedOutEvent(userID, email)

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, "AUTH.LOGOUT", event.GetType())
}

func TestEmailVerified_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"

	event := user.NewEmailVerifiedEvent(userID, email)

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, "USER.EMAIL.VERIFIED", event.GetType())
}

func TestEmailVerificationRequested_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"
	token := "verification-token-abc123"

	event := user.NewEmailVerificationRequestedEvent(userID, email, token)

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, token, event.Token)
	assert.Equal(t, "USER.EMAIL.VERIFICATION.REQUESTED", event.GetType())
}
