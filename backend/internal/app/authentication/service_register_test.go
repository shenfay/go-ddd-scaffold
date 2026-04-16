package authentication

import (
	"context"
	"testing"

	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/test/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestService_Register 测试用户注册功能
func TestService_Register(t *testing.T) {
	t.Run("should register user successfully", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser(factory.WithUnverified())
		tokenPair := f.CreateTokenPair(mockUser.ID, mockUser.Email)

		mockUserRepo.On("ExistsByEmail", mock.Anything, mockUser.Email).Return(false)
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockUserRepo.On("FindByID", mock.Anything, mock.Anything).Return(mockUser, nil)
		mockEmailTokenRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		mockTokenService.On("GenerateTokens", mock.Anything, mock.Anything, mockUser.Email).Return(&TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    3600,
		}, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RegisterCommand{
			Email:    mockUser.Email,
			Password: "TestPassword123!",
		}
		resp, err := service.Register(context.Background(), cmd)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mockUser.Email, resp.User.Email)
		assert.Equal(t, tokenPair.AccessToken, resp.AccessToken)
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
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser(factory.WithUnverified())
		tokenPair := f.CreateTokenPair(mockUser.ID, mockUser.Email)

		mockUserRepo.On("ExistsByEmail", mock.Anything, mockUser.Email).Return(false)
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockUserRepo.On("FindByID", mock.Anything, mock.Anything).Return(mockUser, nil)
		mockEmailTokenRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)
		mockTokenService.On("GenerateTokens", mock.Anything, mock.Anything, mockUser.Email).Return(&TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    3600,
		}, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RegisterCommand{
			Email:    mockUser.Email,
			Password: "TestPassword123!",
		}
		resp, err := service.Register(context.Background(), cmd)

		// Email token creation failure should not affect registration
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mockUser.Email, resp.User.Email)
	})
}
