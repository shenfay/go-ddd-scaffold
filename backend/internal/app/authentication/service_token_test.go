package authentication

import (
	"context"
	"testing"

	authErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/auth"
	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/test/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestService_RefreshToken 测试 Token 刷新功能
func TestService_RefreshToken(t *testing.T) {
	t.Run("should refresh token successfully", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		oldRefreshToken := "old-refresh-token"
		mockUser := f.CreateUser()
		newTokenPair := f.CreateTokenPair(mockUser.ID, mockUser.Email)

		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, oldRefreshToken).Return(&DeviceInfo{
			UserID: mockUser.ID,
		}, nil)
		mockUserRepo.On("FindByID", mock.Anything, mockUser.ID).Return(mockUser, nil)
		mockTokenService.On("RevokeDeviceByToken", mock.Anything, oldRefreshToken).Return(nil)
		mockTokenService.On("GenerateTokens", mock.Anything, mockUser.ID, mockUser.Email).Return(&TokenPair{
			AccessToken:  newTokenPair.AccessToken,
			RefreshToken: newTokenPair.RefreshToken,
			ExpiresIn:    3600,
		}, nil)
		mockTokenService.On("StoreDeviceInfo", mock.Anything, newTokenPair.RefreshToken, mock.Anything).Return(nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		resp, err := service.RefreshToken(context.Background(), RefreshTokenCommand{
			RefreshToken: oldRefreshToken,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, newTokenPair.AccessToken, resp.AccessToken)
	})

	t.Run("should fail with invalid refresh token", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, "invalid-token").Return(nil, authErr.ErrInvalidToken)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		resp, err := service.RefreshToken(context.Background(), RefreshTokenCommand{
			RefreshToken: "invalid-token",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should fail when device validation fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, "tampered-token").Return(nil, authErr.ErrInvalidToken)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		resp, err := service.RefreshToken(context.Background(), RefreshTokenCommand{
			RefreshToken: "tampered-token",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_RefreshToken_EdgeCases 测试 Token 刷新边缘场景
func TestService_RefreshToken_EdgeCases(t *testing.T) {
	t.Run("should fail when user not found after token validation", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "nonexistent-user"
		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, "valid-token").Return(&DeviceInfo{
			UserID: userID,
		}, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		resp, err := service.RefreshToken(context.Background(), RefreshTokenCommand{
			RefreshToken: "valid-token",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should fail when token generation fails after validation", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser()

		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, "valid-token").Return(&DeviceInfo{
			UserID: mockUser.ID,
		}, nil)
		mockUserRepo.On("FindByID", mock.Anything, mockUser.ID).Return(mockUser, nil)
		mockTokenService.On("RevokeDeviceByToken", mock.Anything, "valid-token").Return(nil)
		mockTokenService.On("GenerateTokens", mock.Anything, mockUser.ID, mockUser.Email).Return(nil, assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		resp, err := service.RefreshToken(context.Background(), RefreshTokenCommand{
			RefreshToken: "valid-token",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_Logout 测试用户登出功能
func TestService_Logout(t *testing.T) {
	t.Run("should logout successfully", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockUser := f.CreateUser()

		mockTokenService.On("RevokeToken", mock.Anything, mockUser.ID).Return(nil)
		mockUserRepo.On("FindByID", mock.Anything, mockUser.ID).Return(mockUser, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LogoutCommand{
			UserID: mockUser.ID,
		}
		err := service.Logout(context.Background(), cmd)

		assert.NoError(t, err)
	})

	t.Run("should logout even when user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"

		mockTokenService.On("RevokeToken", mock.Anything, userID).Return(nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LogoutCommand{
			UserID: userID,
		}
		err := service.Logout(context.Background(), cmd)

		// Logout should succeed even when user lookup fails
		assert.NoError(t, err)
	})
}

// TestService_Logout_EdgeCases 测试登出边缘场景
func TestService_Logout_EdgeCases(t *testing.T) {
	t.Run("should fail when token revoke fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockTokenService.On("RevokeToken", mock.Anything, mock.Anything).Return(assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LogoutCommand{
			UserID: "user-123",
		}
		err := service.Logout(context.Background(), cmd)

		// Logout should fail if token revoke fails
		assert.Error(t, err)
	})
}
