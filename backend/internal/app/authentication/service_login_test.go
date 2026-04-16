package authentication

import (
	"context"
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/test/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestService_Login 测试用户登录功能
func TestService_Login(t *testing.T) {
	t.Run("should login successfully with valid credentials", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser()
		tokenPair := f.CreateTokenPair(mockUser.ID, mockUser.Email)
		deviceInfo := f.CreateDeviceInfo()

		mockUserRepo.On("FindByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockTokenService.On("GenerateTokens", mock.Anything, mockUser.ID, mockUser.Email).Return(&TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    3600,
		}, nil)
		mockTokenService.On("StoreDeviceInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LoginCommand{
			Email:      mockUser.Email,
			Password:   "TestPassword123!",
			IP:         deviceInfo.IP,
			UserAgent:  deviceInfo.UserAgent,
			DeviceType: deviceInfo.DeviceType,
		}
		resp, err := service.Login(context.Background(), cmd)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mockUser.Email, resp.User.Email)
		assert.Equal(t, tokenPair.AccessToken, resp.AccessToken)
	})

	t.Run("should fail with invalid credentials", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "test@example.com"
		mockUser, _ := user.NewUser(email, "CorrectPassword123!")
		mockUser.ID = "user-123"

		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LoginCommand{
			Email:    email,
			Password: "WrongPassword456!",
		}
		resp, err := service.Login(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should fail with locked account", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "locked@example.com"
		mockUser, _ := user.NewUser(email, "Password123!")
		mockUser.ID = "user-locked"
		for i := 0; i < 5; i++ {
			mockUser.IncrementFailedAttempts(5)
		}

		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(mockUser, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LoginCommand{
			Email:    email,
			Password: "Password123!",
		}
		resp, err := service.Login(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should fail when user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "nonexistent@example.com"
		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LoginCommand{
			Email:    email,
			Password: "Password123!",
		}
		resp, err := service.Login(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_Login_EdgeCases 测试登录边缘场景
func TestService_Login_EdgeCases(t *testing.T) {
	t.Run("should fail when token generation fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "test@example.com"
		password := "ValidPassword123!"
		mockUser, _ := user.NewUser(email, password)
		mockUser.ID = "user-123"

		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockTokenService.On("GenerateTokens", mock.Anything, mockUser.ID, email).Return(nil, assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LoginCommand{
			Email:    email,
			Password: password,
		}
		resp, err := service.Login(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should succeed when device info storage fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "test@example.com"
		password := "ValidPassword123!"
		mockUser, _ := user.NewUser(email, password)
		mockUser.ID = "user-123"

		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockTokenService.On("GenerateTokens", mock.Anything, mockUser.ID, email).Return(&TokenPair{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		}, nil)
		mockTokenService.On("StoreDeviceInfo", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LoginCommand{
			Email:      email,
			Password:   password,
			IP:         "192.168.1.1",
			UserAgent:  "Mozilla/5.0",
			DeviceType: "web",
		}
		resp, err := service.Login(context.Background(), cmd)

		// StoreDeviceInfo failure should not affect login
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, email, resp.User.Email)
	})
}
