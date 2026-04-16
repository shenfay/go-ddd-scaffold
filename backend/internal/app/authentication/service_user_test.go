package authentication

import (
	"context"
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestService_GetUserByID 测试根据 ID 获取用户功能
func TestService_GetUserByID(t *testing.T) {
	t.Run("should return user when found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		email := "test@example.com"
		mockUser, _ := user.NewUser(email, "Password123!")
		mockUser.ID = userID

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		result, err := service.GetUserByID(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, email, result.Email)
	})

	t.Run("should fail when user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "nonexistent-user"
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		result, err := service.GetUserByID(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, userErr.ErrNotFound, err)
	})
}

// TestService_GetUserByID_EdgeCases 测试获取用户边缘场景
func TestService_GetUserByID_EdgeCases(t *testing.T) {
	t.Run("should fail when database error occurs", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		result, err := service.GetUserByID(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, assert.AnError, err)
	})
}
