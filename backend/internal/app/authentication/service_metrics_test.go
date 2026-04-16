package authentication

import (
	"context"
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestService_WithoutMetrics 测试没有指标服务时的正常运行
func TestService_WithoutMetrics(t *testing.T) {
	t.Run("should work normally without metrics service", func(t *testing.T) {
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

		// 传入 nil metrics
		service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

		cmd := LoginCommand{
			Email:    email,
			Password: password,
		}
		resp, err := service.Login(context.Background(), cmd)

		// 应该正常工作,不因为 metrics 为 nil 而 panic
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
