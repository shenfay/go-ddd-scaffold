package authentication

import (
	"context"
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/test/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_RefreshToken_Mock(t *testing.T) {
	t.Run("should refresh token successfully", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser()
		oldRefreshToken := "old-refresh-token"
		newTokenPair := f.CreateTokenPair(mockUser.ID, mockUser.Email)

		// Mock ValidateRefreshTokenWithDevice
		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, oldRefreshToken).Return(&DeviceInfo{
			UserID: mockUser.ID,
		}, nil)

		// Mock FindByID
		mockUserRepo.On("FindByID", mock.Anything, mockUser.ID).Return(mockUser, nil)

		// Mock Update for last login
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

		// Mock RevokeDeviceByToken
		mockTokenService.On("RevokeDeviceByToken", mock.Anything, oldRefreshToken).Return(nil)

		// Mock GenerateTokens
		mockTokenService.On("GenerateTokens", mock.Anything, mockUser.ID, mockUser.Email).Return(&TokenPair{
			AccessToken:  newTokenPair.AccessToken,
			RefreshToken: newTokenPair.RefreshToken,
			ExpiresIn:    time.Duration(newTokenPair.ExpiresIn) * time.Second,
		}, nil)

		// Mock StoreDeviceInfo
		mockTokenService.On("StoreDeviceInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RefreshTokenCommand{
			RefreshToken: oldRefreshToken,
		}
		resp, err := service.RefreshToken(context.Background(), cmd)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, newTokenPair.AccessToken, resp.AccessToken)
		assert.Equal(t, newTokenPair.RefreshToken, resp.RefreshToken)

		mockTokenService.AssertExpectations(t)
	})

	t.Run("should fail with expired refresh token", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		expiredToken := "expired-token"

		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, expiredToken).Return((*DeviceInfo)(nil), assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RefreshTokenCommand{
			RefreshToken: expiredToken,
		}
		resp, err := service.RefreshToken(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, resp)

		mockTokenService.AssertExpectations(t)
	})

	t.Run("should fail with user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		refreshToken := "valid-refresh-token"
		userID := "non-existent-user"

		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, refreshToken).Return(&DeviceInfo{
			UserID: userID,
		}, nil)

		mockUserRepo.On("FindByID", mock.Anything, userID).Return((*user.User)(nil), assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RefreshTokenCommand{
			RefreshToken: refreshToken,
		}
		resp, err := service.RefreshToken(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, resp)

		mockTokenService.AssertExpectations(t)
	})
}

func TestService_Logout_Mock(t *testing.T) {
	t.Run("should logout successfully", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser()

		// Mock RevokeToken
		mockTokenService.On("RevokeToken", mock.Anything, mockUser.ID).Return(nil)

		// Mock FindByID for event publishing
		mockUserRepo.On("FindByID", mock.Anything, mockUser.ID).Return(mockUser, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LogoutCommand{
			UserID: mockUser.ID,
		}
		err := service.Logout(context.Background(), cmd)

		assert.NoError(t, err)

		mockTokenService.AssertExpectations(t)
	})

	t.Run("should fail when revoke token fails", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser()

		mockTokenService.On("RevokeToken", mock.Anything, mockUser.ID).Return(assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LogoutCommand{
			UserID: mockUser.ID,
		}
		err := service.Logout(context.Background(), cmd)

		assert.Error(t, err)

		mockTokenService.AssertExpectations(t)
	})
}

func TestService_GetUserDevices_Mock(t *testing.T) {
	t.Run("should get user devices successfully", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser()
		devices := []DeviceInfo{
			{
				UserID:     mockUser.ID,
				IP:         "192.168.1.1",
				DeviceType: "web",
			},
			{
				UserID:     mockUser.ID,
				IP:         "192.168.1.2",
				DeviceType: "mobile",
			},
		}

		mockTokenService.On("GetUserDevices", mock.Anything, mockUser.ID).Return(devices, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		devicesResp, err := service.tokenService.GetUserDevices(context.Background(), mockUser.ID)

		assert.NoError(t, err)
		assert.Len(t, devicesResp, 2)
		assert.Equal(t, "web", devicesResp[0].DeviceType)
		assert.Equal(t, "mobile", devicesResp[1].DeviceType)

		mockTokenService.AssertExpectations(t)
	})

	t.Run("should return empty list when no devices", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser()

		mockTokenService.On("GetUserDevices", mock.Anything, mockUser.ID).Return([]DeviceInfo{}, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		devicesResp, err := service.tokenService.GetUserDevices(context.Background(), mockUser.ID)

		assert.NoError(t, err)
		assert.Empty(t, devicesResp)

		mockTokenService.AssertExpectations(t)
	})
}
