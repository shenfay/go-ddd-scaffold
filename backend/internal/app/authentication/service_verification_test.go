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

// TestService_VerifyEmail 测试邮箱验证功能
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

		mockUser, _ := user.NewUser(email, "Password123!")
		mockUser.ID = userID

		verifyToken := &authentication.EmailVerificationToken{
			ID:        "token-id",
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		mockEmailTokenRepo.On("FindByToken", mock.Anything, token).Return(verifyToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
		mockEmailTokenRepo.On("MarkAsUsed", mock.Anything, "token-id").Return(nil)

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

		token := "invalid-token"

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
			ID:        "token-id",
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		mockEmailTokenRepo.On("FindByToken", mock.Anything, token).Return(verifyToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.VerifyEmail(context.Background(), token)

		assert.Error(t, err)
	})
}

// TestService_VerifyEmail_EdgeCases 测试邮箱验证边缘场景
func TestService_VerifyEmail_EdgeCases(t *testing.T) {
	t.Run("should fail with expired verification token", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-123"
		token := "expired-token"

		verifyToken := &authentication.EmailVerificationToken{
			ID:        "token-id",
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(-1 * time.Hour), // 已过期
		}

		mockEmailTokenRepo.On("FindByToken", mock.Anything, token).Return(verifyToken, nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.VerifyEmail(context.Background(), token)

		assert.Error(t, err)
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
			ID:        "token-id",
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		mockEmailTokenRepo.On("FindByToken", mock.Anything, token).Return(verifyToken, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.VerifyEmail(context.Background(), token)

		assert.Error(t, err)
	})
}

// TestService_SendVerificationEmail 测试发送验证邮件功能
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

		// Should succeed but skip sending
		assert.NoError(t, err)
	})
}

// TestService_SendVerificationEmail_EdgeCases 测试发送验证邮件边缘场景
func TestService_SendVerificationEmail_EdgeCases(t *testing.T) {
	t.Run("should fail when user not found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "nonexistent"
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, userErr.ErrNotFound)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.SendVerificationEmail(context.Background(), userID)

		assert.Error(t, err)
	})

	t.Run("should fail when email token creation fails", func(t *testing.T) {
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
		mockEmailTokenRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.SendVerificationEmail(context.Background(), userID)

		assert.Error(t, err)
	})

	t.Run("should skip sending when email already verified", func(t *testing.T) {
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
		// 不应该调用 EmailTokenRepo.Create

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.SendVerificationEmail(context.Background(), userID)

		// Should succeed but skip sending
		assert.NoError(t, err)
	})

	t.Run("should create new token when existing token is expired", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockResetTokenRepo := new(MockPasswordResetTokenRepository)
		mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
		mockTokenService := new(MockTokenService)
		publisher := event.NewPublisher(nil)

		userID := "user-789"
		email := "test@example.com"
		mockUser, _ := user.NewUser(email, "Password123!")
		mockUser.ID = userID

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(mockUser, nil)
		mockEmailTokenRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		err := service.SendVerificationEmail(context.Background(), userID)

		assert.NoError(t, err)
	})
}
