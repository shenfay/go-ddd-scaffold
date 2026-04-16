package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	domainUser "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domainUser.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Save(ctx context.Context, user *domainUser.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domainUser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	args := m.Called(ctx, email)
	return args.Bool(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domainUser.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockEventBus 模拟事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, evt event.Event) error {
	args := m.Called(ctx, evt)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler event.EventHandler) {
	m.Called(eventType, handler)
}

func TestService_CreateUser(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockBus := new(MockEventBus)

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockBus.On("Publish", mock.Anything, mock.AnythingOfType("*user.UserRegistered")).Return(nil)

		service := NewService(mockRepo, mockBus)

		dto, err := service.CreateUser(context.Background(), CreateUserCommand{
			Email:    "test@example.com",
			Password: "ValidPassword123!",
		})

		assert.NoError(t, err)
		assert.NotNil(t, dto)
		assert.Equal(t, "test@example.com", dto.Email)
		assert.False(t, dto.EmailVerified)
		assert.False(t, dto.Locked)

		mockRepo.AssertExpectations(t)
		mockBus.AssertExpectations(t)
	})

	t.Run("should fail when save fails", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockBus := new(MockEventBus)

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("database error"))

		service := NewService(mockRepo, mockBus)

		dto, err := service.CreateUser(context.Background(), CreateUserCommand{
			Email:    "test@example.com",
			Password: "ValidPassword123!",
		})

		assert.Error(t, err)
		assert.Nil(t, dto)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("should succeed when event publish fails", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockBus := new(MockEventBus)

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockBus.On("Publish", mock.Anything, mock.AnythingOfType("*user.UserRegistered")).Return(errors.New("event bus error"))

		service := NewService(mockRepo, mockBus)

		dto, err := service.CreateUser(context.Background(), CreateUserCommand{
			Email:    "test@example.com",
			Password: "ValidPassword123!",
		})

		assert.NoError(t, err)
		assert.NotNil(t, dto)
	})

	t.Run("should succeed with nil event bus", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

		service := NewService(mockRepo, nil)

		dto, err := service.CreateUser(context.Background(), CreateUserCommand{
			Email:    "test@example.com",
			Password: "ValidPassword123!",
		})

		assert.NoError(t, err)
		assert.NotNil(t, dto)
	})
}

func TestService_GetUserByID(t *testing.T) {
	t.Run("should get user by id successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		user := &domainUser.User{
			ID:            "user-123",
			Email:         "test@example.com",
			EmailVerified: true,
			Locked:        false,
			LastLoginAt:   nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		mockRepo.On("FindByID", mock.Anything, "user-123").Return(user, nil)

		service := NewService(mockRepo, nil)

		dto, err := service.GetUserByID(context.Background(), "user-123")

		assert.NoError(t, err)
		assert.NotNil(t, dto)
		assert.Equal(t, "user-123", dto.ID)
		assert.Equal(t, "test@example.com", dto.Email)
		assert.True(t, dto.EmailVerified)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		mockRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, errors.New("user not found"))

		service := NewService(mockRepo, nil)

		dto, err := service.GetUserByID(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.Nil(t, dto)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestService_UpdateProfile(t *testing.T) {
	t.Run("should update profile successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockBus := new(MockEventBus)

		existingUser, _ := domainUser.NewUser("old@example.com", "ValidPassword123!")
		existingUser.ID = "user-123"

		mockRepo.On("FindByID", mock.Anything, "user-123").Return(existingUser, nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockBus.On("Publish", mock.Anything, mock.AnythingOfType("*user.UserProfileUpdated")).Return(nil)

		service := NewService(mockRepo, mockBus)

		dto, err := service.UpdateProfile(context.Background(), UpdateProfileCommand{
			UserID: "user-123",
			Email:  "new@example.com",
		})

		assert.NoError(t, err)
		assert.NotNil(t, dto)
		assert.Equal(t, "new@example.com", dto.Email)

		mockRepo.AssertExpectations(t)
		mockBus.AssertExpectations(t)
	})

	t.Run("should fail when user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockBus := new(MockEventBus)

		mockRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, errors.New("user not found"))

		service := NewService(mockRepo, mockBus)

		dto, err := service.UpdateProfile(context.Background(), UpdateProfileCommand{
			UserID: "non-existent",
			Email:  "new@example.com",
		})

		assert.Error(t, err)
		assert.Nil(t, dto)
	})

	t.Run("should fail when save fails", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockBus := new(MockEventBus)

		existingUser, _ := domainUser.NewUser("old@example.com", "ValidPassword123!")
		existingUser.ID = "user-123"

		mockRepo.On("FindByID", mock.Anything, "user-123").Return(existingUser, nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("database error"))

		service := NewService(mockRepo, mockBus)

		dto, err := service.UpdateProfile(context.Background(), UpdateProfileCommand{
			UserID: "user-123",
			Email:  "new@example.com",
		})

		assert.Error(t, err)
		assert.Nil(t, dto)
		assert.Contains(t, err.Error(), "database error")
	})
}
