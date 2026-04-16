package authentication_test

import (
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"github.com/stretchr/testify/assert"
)

func TestPasswordResetRequested_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"
	timestamp := time.Now()

	event := &authentication.PasswordResetRequested{
		UserID:    userID,
		Email:     email,
		Timestamp: timestamp,
	}

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, "AUTH.PASSWORD_RESET.REQUESTED", event.GetType())
	assert.NotNil(t, event.GetPayload())
}

func TestPasswordResetCompleted_Event(t *testing.T) {
	userID := "user-123"
	email := "test@example.com"
	timestamp := time.Now()

	event := &authentication.PasswordResetCompleted{
		UserID:    userID,
		Email:     email,
		Timestamp: timestamp,
	}

	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, email, event.Email)
	assert.Equal(t, "AUTH.PASSWORD_RESET.COMPLETED", event.GetType())
	assert.NotNil(t, event.GetPayload())
}
