package listener

import (
	"context"
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventBus 模拟事件总线
type MockEventBus struct {
	mock.Mock
	PublishedEvents []event.Event
}

func (m *MockEventBus) Publish(ctx context.Context, evt event.Event) error {
	m.PublishedEvents = append(m.PublishedEvents, evt)
	args := m.Called(ctx, evt)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler event.EventHandler) {
	m.Called(eventType, handler)
}

func TestActivityLogListener_HandleUserRegistered(t *testing.T) {
	t.Run("should publish activity log task on user registered", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewActivityLogListener(mockBus)

		mockBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

		evt := &user.UserRegistered{
			UserID: "user-123",
			Email:  "test@example.com",
		}

		err := listener.HandleUserRegistered(context.Background(), evt)

		assert.NoError(t, err)
		assert.Len(t, mockBus.PublishedEvents, 1)

		task := mockBus.PublishedEvents[0].(*ActivityLogTask)
		assert.Equal(t, "activity.log", task.Type)
		assert.Equal(t, "USER.REGISTERED", task.Action)
		assert.Equal(t, "user-123", task.UserID)
		assert.Equal(t, "test@example.com", task.Email)
		assert.Equal(t, "用户注册", task.Description)
		assert.Equal(t, "SUCCESS", task.Status)
	})
}

func TestActivityLogListener_HandleUserLoggedIn(t *testing.T) {
	t.Run("should publish activity log with device info on login", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewActivityLogListener(mockBus)

		mockBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

		evt := &user.UserLoggedIn{
			UserID:    "user-456",
			Email:     "user@example.com",
			IP:        "192.168.1.1",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		}

		err := listener.HandleUserLoggedIn(context.Background(), evt)

		assert.NoError(t, err)
		assert.Len(t, mockBus.PublishedEvents, 1)

		task := mockBus.PublishedEvents[0].(*ActivityLogTask)
		assert.Equal(t, "USER.LOGIN", task.Action)
		assert.Equal(t, "user-456", task.UserID)
		assert.Equal(t, "192.168.1.1", task.IP)
		assert.Contains(t, task.Metadata, "ip")
		assert.Contains(t, task.Metadata, "user_agent")
	})
}

func TestActivityLogListener_HandleUserLoggedOut(t *testing.T) {
	t.Run("should publish activity log on logout", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewActivityLogListener(mockBus)

		mockBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

		evt := &user.UserLoggedOut{
			UserID: "user-789",
			Email:  "logout@example.com",
		}

		err := listener.HandleUserLoggedOut(context.Background(), evt)

		assert.NoError(t, err)
		assert.Len(t, mockBus.PublishedEvents, 1)

		task := mockBus.PublishedEvents[0].(*ActivityLogTask)
		assert.Equal(t, "USER.LOGOUT", task.Action)
		assert.Equal(t, "user-789", task.UserID)
		assert.Equal(t, "用户登出", task.Description)
	})
}

func TestActivityLogListener_HandleTokenRefreshed(t *testing.T) {
	t.Run("should publish activity log on token refresh", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewActivityLogListener(mockBus)

		mockBus.On("Publish", mock.Anything, mock.Anything).Return(nil)

		evt := &user.TokenRefreshed{
			UserID:   "user-999",
			OldToken: "old-token",
			NewToken: "new-token",
		}

		err := listener.HandleTokenRefreshed(context.Background(), evt)

		assert.NoError(t, err)
		assert.Len(t, mockBus.PublishedEvents, 1)

		task := mockBus.PublishedEvents[0].(*ActivityLogTask)
		assert.Equal(t, "USER.TOKEN_REFRESHED", task.Action)
		assert.Equal(t, "user-999", task.UserID)
		assert.Equal(t, "Token刷新", task.Description)
		assert.Equal(t, "old-token", task.Metadata["old_token"])
		assert.Equal(t, "new-token", task.Metadata["new_token"])
	})
}

func TestActivityLogListener_SubscribeEvents(t *testing.T) {
	t.Run("should subscribe to all events", func(t *testing.T) {
		mockBus := new(MockEventBus)
		listener := NewActivityLogListener(mockBus)

		// Subscribe should be called for each event type
		mockBus.On("Subscribe", "USER.REGISTERED", mock.Anything)
		mockBus.On("Subscribe", "AUTH.LOGIN.SUCCESS", mock.Anything)
		mockBus.On("Subscribe", "AUTH.LOGOUT", mock.Anything)
		mockBus.On("Subscribe", "AUTH.TOKEN.REFRESHED", mock.Anything)

		listener.SubscribeEvents(mockBus)

		mockBus.AssertExpectations(t)
	})
}
