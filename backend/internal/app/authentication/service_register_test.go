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

// TestService_Register 测试用户注册功能
func TestService_Register(t *testing.T) {
	t.Run("should register user successfully", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "test@example.com"
		password := "ValidPassword123!"

		mockUserRepo.On("ExistsByEmail", mock.Anything, email).Return(false)
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockUserRepo.On("FindByID", mock.Anything, mock.Anything).Return(&user.User{
			ID:    "test-user-id",
			Email: email,
		}, nil)
		mockEmailTokenRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		mockTokenService.On("GenerateTokens", mock.Anything, mock.Anything, email).Return(&TokenPair{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		}, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RegisterCommand{
			Email:    email,
			Password: password,
		}
		resp, err := service.Register(context.Background(), cmd)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, email, resp.User.Email)
		assert.Equal(t, "access-token", resp.AccessToken)
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "duplicate@example.com"
		mockUserRepo.On("ExistsByEmail", mock.Anything, email).Return(true)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RegisterCommand{
			Email:    email,
			Password: "ValidPassword123!",
		}
		resp, err := service.Register(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, userErr.ErrEmailAlreadyExists, err)
	})
}

// TestService_Register_EdgeCases 测试注册边缘场景
func TestService_Register_EdgeCases(t *testing.T) {
	t.Run("should fail when user creation fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "test@example.com"
		password := "ValidPassword123!"

		mockUserRepo.On("ExistsByEmail", mock.Anything, email).Return(false)
		mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RegisterCommand{
			Email:    email,
			Password: password,
		}
		resp, err := service.Register(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should succeed even when email token creation fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "test@example.com"
		password := "ValidPassword123!"

		mockUserRepo.On("ExistsByEmail", mock.Anything, email).Return(false)
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockUserRepo.On("FindByID", mock.Anything, mock.Anything).Return(&user.User{
			ID:    "test-user-id",
			Email: email,
		}, nil)
		mockEmailTokenRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)
		mockTokenService.On("GenerateTokens", mock.Anything, mock.Anything, email).Return(&TokenPair{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		}, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RegisterCommand{
			Email:    email,
			Password: password,
		}
		resp, err := service.Register(context.Background(), cmd)

		// Email token creation failure should not affect registration
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, email, resp.User.Email)
	})
}
