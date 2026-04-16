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

// TestService_RequestPasswordReset 测试密码重置请求功能
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

// TestService_RequestPasswordReset_EdgeCases 测试密码重置请求边缘场景
func TestService_RequestPasswordReset_EdgeCases(t *testing.T) {
	t.Run("should succeed even when user not found (security)", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "nonexistent@example.com"
		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RequestPasswordResetCommand{
			Email: email,
		}
		err := service.RequestPasswordReset(context.Background(), cmd)

		// Security: should not reveal whether email exists
		assert.NoError(t, err)
	})

	t.Run("should fail when reset token creation fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		email := "verified@example.com"
		mockUser, _ := user.NewUser(email, "OldPassword123!")
		mockUser.VerifyEmail()

		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(mockUser, nil)
		mockResetTokenRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := RequestPasswordResetCommand{
			Email: email,
		}
		err := service.RequestPasswordReset(context.Background(), cmd)

		assert.Error(t, err)
	})
}

// TestService_ResetPassword 测试密码重置功能
func TestService_ResetPassword(t *testing.T) {
	t.Run("should reset password successfully", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		email := "test@example.com"
		token := "valid-token"
		newPassword := "NewPassword123!"

		mockUser, _ := user.NewUser(email, "OldPassword123!")
		mockUser.ID = userID

		resetToken := &authentication.PasswordResetToken{
			ID:        "token-id",
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		mockResetTokenRepo.On("FindByToken", mock.Anything, token).Return(resetToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockResetTokenRepo.On("MarkAsUsed", mock.Anything, "token-id").Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       token,
			NewPassword: newPassword,
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

		token := "invalid-token"

		mockResetTokenRepo.On("FindByToken", mock.Anything, token).Return(nil, authErr.ErrInvalidToken)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       token,
			NewPassword: "NewPassword123!",
		}
		err := service.ResetPassword(context.Background(), cmd)

		assert.Error(t, err)
	})
}

// TestService_ResetPassword_EdgeCases 测试密码重置边缘场景
func TestService_ResetPassword_EdgeCases(t *testing.T) {
	t.Run("should fail when user not found after token validation", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "nonexistent-user"
		token := "valid-token"

		resetToken := &authentication.PasswordResetToken{
			ID:        "token-id",
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		mockResetTokenRepo.On("FindByToken", mock.Anything, token).Return(resetToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       token,
			NewPassword: "NewPassword123!",
		}
		err := service.ResetPassword(context.Background(), cmd)

		assert.Error(t, err)
	})

	t.Run("should fail when database update fails", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		email := "test@example.com"
		token := "valid-token"
		newPassword := "NewPassword123!"

		mockUser, _ := user.NewUser(email, "OldPassword123!")
		mockUser.ID = userID

		resetToken := &authentication.PasswordResetToken{
			ID:        "token-id",
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		mockResetTokenRepo.On("FindByToken", mock.Anything, token).Return(resetToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(assert.AnError)
		mockResetTokenRepo.On("MarkAsUsed", mock.Anything, "token-id").Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := ResetPasswordCommand{
			Token:       token,
			NewPassword: newPassword,
		}
		err := service.ResetPassword(context.Background(), cmd)

		assert.Error(t, err)
	})
}
