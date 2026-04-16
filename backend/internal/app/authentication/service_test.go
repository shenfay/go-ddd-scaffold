package authentication

import (
	"context"
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	authErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/auth"
	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestService_Login(t *testing.T) {
	t.Run("should login successfully with valid credentials", func(t *testing.T) {
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
		mockTokenService.On("StoreDeviceInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LoginCommand{
			Email:      email,
			Password:   password,
			IP:         "192.168.1.1",
			UserAgent:  "Mozilla/5.0",
			DeviceType: "web",
		}
		resp, err := service.Login(context.Background(), cmd)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, email, resp.User.Email)
		assert.Equal(t, "access-token", resp.AccessToken)
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

func TestService_GetUserByID(t *testing.T) {
	t.Run("should get user by id successfully", func(t *testing.T) {
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

	t.Run("should return error when user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "nonexistent"
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		result, err := service.GetUserByID(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, userErr.ErrNotFound, err)
	})
}

func TestService_VerifyEmail(t *testing.T) {
	t.Run("should verify email successfully", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		email := "test@example.com"
		token := "valid-token"

		verifyToken, _ := authentication.NewEmailVerificationToken(userID)
		verifyToken.Token = token

		mockUser, _ := user.NewUser(email, "Password123!")
		mockUser.ID = userID

		mockEmailTokenRepo.On("FindByToken", mock.Anything, token).Return(verifyToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)
		mockEmailTokenRepo.On("MarkAsUsed", mock.Anything, verifyToken.ID).Return(nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.VerifyEmail(context.Background(), token)

		assert.NoError(t, err)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockEmailTokenRepo.On("FindByToken", mock.Anything, "invalid-token").Return(nil, authErr.ErrInvalidToken)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.VerifyEmail(context.Background(), "invalid-token")

		assert.Error(t, err)
		assert.Equal(t, authErr.ErrInvalidToken, err)
	})

	t.Run("should fail with expired token", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		token := "expired-token"
		verifyToken, _ := authentication.NewEmailVerificationToken("user-123")
		verifyToken.Token = token
		verifyToken.ExpiresAt = time.Now().Add(-1 * time.Hour)

		mockEmailTokenRepo.On("FindByToken", mock.Anything, token).Return(verifyToken, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.VerifyEmail(context.Background(), token)

		assert.Error(t, err)
		assert.Equal(t, authErr.ErrTokenExpired, err)
	})
}

func TestService_RequestPasswordReset(t *testing.T) {
	t.Run("should request password reset for verified email", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "verified@example.com"
		mockUser, _ := user.NewUser(email, "OldPassword123!")
		mockUser.VerifyEmail()

		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(mockUser, nil)
		mockResetTokenRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RequestPasswordResetCommand{
			Email: email,
		}
		err := service.RequestPasswordReset(context.Background(), cmd)

		assert.NoError(t, err)
	})

	t.Run("should fail with unverified email", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "unverified@example.com"
		mockUser, _ := user.NewUser(email, "Password123!")

		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(mockUser, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RequestPasswordResetCommand{
			Email: email,
		}
		err := service.RequestPasswordReset(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, userErr.ErrEmailNotVerified, err)
	})
}

func TestService_Logout(t *testing.T) {
	t.Run("should logout successfully", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		email := "test@example.com"
		mockUser, _ := user.NewUser(email, "Password123!")
		mockUser.ID = userID

		mockTokenService.On("RevokeToken", mock.Anything, userID).Return(nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LogoutCommand{
			UserID: userID,
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

		// Logout 应该成功,即使用户查找失败(事件发布失败不影响)
		assert.NoError(t, err)
	})
}

func TestService_RefreshToken(t *testing.T) {
	t.Run("should refresh token successfully", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		oldRefreshToken := "old-refresh-token"
		newRefreshToken := "new-refresh-token"
		userID := "user-123"
		email := "test@example.com"

		mockUser, _ := user.NewUser(email, "Password123!")
		mockUser.ID = userID

		mockTokenService.On("ValidateRefreshTokenWithDevice", mock.Anything, oldRefreshToken).Return(&DeviceInfo{
			UserID: userID,
		}, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)
		mockTokenService.On("RevokeDeviceByToken", mock.Anything, oldRefreshToken).Return(nil)
		mockTokenService.On("GenerateTokens", mock.Anything, userID, email).Return(&TokenPair{
			AccessToken:  "new-access-token",
			RefreshToken: newRefreshToken,
			ExpiresIn:    3600,
		}, nil)
		mockTokenService.On("StoreDeviceInfo", mock.Anything, newRefreshToken, mock.Anything).Return(nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		resp, err := service.RefreshToken(context.Background(), RefreshTokenCommand{
			RefreshToken: oldRefreshToken,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "new-access-token", resp.AccessToken)
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

func TestService_ResetPassword(t *testing.T) {
	t.Run("should reset password successfully", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		email := "test@example.com"
		token := "valid-reset-token"

		resetToken, _ := authentication.NewPasswordResetToken(userID)
		resetToken.Token = token

		mockUser, _ := user.NewUser(email, "OldPassword123!")
		mockUser.ID = userID

		mockResetTokenRepo.On("FindByToken", mock.Anything, token).Return(resetToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockResetTokenRepo.On("MarkAsUsed", mock.Anything, resetToken.ID).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       token,
			NewPassword: "NewPassword456!",
		}
		err := service.ResetPassword(context.Background(), cmd)

		assert.NoError(t, err)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		mockResetTokenRepo.On("FindByToken", mock.Anything, "invalid-token").Return(nil, authErr.ErrInvalidResetToken)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       "invalid-token",
			NewPassword: "NewPassword456!",
		}
		err := service.ResetPassword(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, authErr.ErrInvalidResetToken, err)
	})

	t.Run("should fail with expired token", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		token := "expired-token"

		resetToken, _ := authentication.NewPasswordResetToken(userID)
		resetToken.Token = token
		resetToken.ExpiresAt = resetToken.CreatedAt.Add(-1 * time.Hour)

		mockResetTokenRepo.On("FindByToken", mock.Anything, token).Return(resetToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(&user.User{ID: userID}, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       token,
			NewPassword: "NewPassword456!",
		}
		err := service.ResetPassword(context.Background(), cmd)

		assert.Error(t, err)
	})
}

func TestService_SendVerificationEmail(t *testing.T) {
	t.Run("should send verification email", func(t *testing.T) {
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
		mockEmailTokenRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.SendVerificationEmail(context.Background(), userID)

		assert.NoError(t, err)
	})

	t.Run("should skip if email already verified", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-456"
		email := "verified@example.com"
		mockUser, _ := user.NewUser(email, "Password123!")
		mockUser.ID = userID
		mockUser.VerifyEmail()

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.SendVerificationEmail(context.Background(), userID)

		assert.NoError(t, err)
		mockEmailTokenRepo.AssertNotCalled(t, "Create")
	})
}

func TestService_ResetPassword_EdgeCases(t *testing.T) {
	t.Run("should fail when user not found after token validation", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "non-existent-user"
		token := "valid-token"

		resetToken, _ := authentication.NewPasswordResetToken(userID)
		resetToken.Token = token

		mockResetTokenRepo.On("FindByToken", mock.Anything, token).Return(resetToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       token,
			NewPassword: "NewPassword456!",
		}
		err := service.ResetPassword(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, userErr.ErrNotFound, err)
	})

	t.Run("should fail when database update fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		token := "valid-token"

		resetToken, _ := authentication.NewPasswordResetToken(userID)
		resetToken.Token = token

		mockUser, _ := user.NewUser("test@example.com", "OldPassword123!")
		mockUser.ID = userID

		mockResetTokenRepo.On("FindByToken", mock.Anything, token).Return(resetToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(assert.AnError)
		mockResetTokenRepo.On("MarkAsUsed", mock.Anything, resetToken.ID).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       token,
			NewPassword: "NewPassword456!",
		}
		err := service.ResetPassword(context.Background(), cmd)

		assert.Error(t, err)
	})
}

func TestService_VerifyEmail_EdgeCases(t *testing.T) {
	t.Run("should fail with expired verification token", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		token := "expired-token"
		mockEmailTokenRepo.On("FindByToken", mock.Anything, token).Return(nil, authErr.ErrInvalidToken)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.VerifyEmail(context.Background(), token)

		assert.Error(t, err)
		// VerifyEmail returns ErrInvalidToken for any invalid verification token
		assert.True(t, err != nil)
	})

	t.Run("should fail when user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "non-existent"
		token := "valid-token"

		verifyToken := &authentication.EmailVerificationToken{
			ID:     "token-id",
			UserID: userID,
			Token:  token,
		}

		mockEmailTokenRepo.On("FindByToken", mock.Anything, token).Return(verifyToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.VerifyEmail(context.Background(), token)

		assert.Error(t, err)
	})
}
