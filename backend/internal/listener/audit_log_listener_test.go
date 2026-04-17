package listener

import (
	"context"
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuditLogListener_HandleUserLoggedIn(t *testing.T) {
	t.Run("should publish audit log on successful login", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewAuditLogListener(mockBus)

		mockBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

		evt := &user.UserLoggedIn{
			UserID:    "user-123",
			Email:     "test@example.com",
			IP:        "10.0.0.1",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		}

		err := listener.HandleUserLoggedIn(context.Background(), evt)

		assert.NoError(t, err)
		assert.Len(t, mockBus.PublishedEvents, 1)

		task := mockBus.PublishedEvents[0].(*AuditLogTask)
		assert.Equal(t, "audit.log", task.Type)
		assert.Equal(t, "AUTH.LOGIN.SUCCESS", task.Action)
		assert.Equal(t, "SUCCESS", task.Status)
		assert.Equal(t, "user-123", task.Data["user_id"])
		assert.Equal(t, "test@example.com", task.Data["email"])
		assert.Equal(t, "10.0.0.1", task.Data["ip"])
	})
}

func TestAuditLogListener_HandleLoginFailed(t *testing.T) {
	t.Run("should publish audit log on login failure", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewAuditLogListener(mockBus)

		mockBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

		evt := &user.LoginFailed{
			UserID: "user-456",
			Email:  "fail@example.com",
			IP:     "192.168.1.100",
			Reason: "invalid_password",
		}

		err := listener.HandleLoginFailed(context.Background(), evt)

		assert.NoError(t, err)
		assert.Len(t, mockBus.PublishedEvents, 1)

		task := mockBus.PublishedEvents[0].(*AuditLogTask)
		assert.Equal(t, "audit.log", task.Type)
		assert.Equal(t, "AUTH.LOGIN.FAILED", task.Action)
		assert.Equal(t, "FAILED", task.Status)
		assert.Equal(t, "user-456", task.Data["user_id"])
		assert.Equal(t, "invalid_password", task.Data["reason"])
	})
}

func TestAuditLogListener_HandleAccountLocked(t *testing.T) {
	t.Run("should publish audit log on account lock", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewAuditLogListener(mockBus)

		mockBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

		evt := &user.AccountLocked{
			UserID:         "user-789",
			Email:          "locked@example.com",
			FailedAttempts: 5,
			LockedUntil:    time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
		}

		err := listener.HandleAccountLocked(context.Background(), evt)

		assert.NoError(t, err)
		assert.Len(t, mockBus.PublishedEvents, 1)

		task := mockBus.PublishedEvents[0].(*AuditLogTask)
		assert.Equal(t, "audit.log", task.Type)
		assert.Equal(t, "SECURITY.ACCOUNT.LOCKED", task.Action)
		assert.Equal(t, "FAILED", task.Status)
		assert.Equal(t, "user-789", task.Data["user_id"])
		assert.Equal(t, 5, task.Data["failed_attempts"])
		assert.NotNil(t, task.Data["locked_until"])
	})
}

func TestAuditLogListener_SubscribeEvents(t *testing.T) {
	t.Run("should subscribe to authentication events", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewAuditLogListener(mockBus)

		mockBus.On("Subscribe", "AUTH.LOGIN.SUCCESS", mock.Anything)
		mockBus.On("Subscribe", "AUTH.LOGIN.FAILED", mock.Anything)
		mockBus.On("Subscribe", "SECURITY.ACCOUNT.LOCKED", mock.Anything)

		listener.SubscribeEvents(mockBus)

		mockBus.AssertExpectations(t)
	})
}
