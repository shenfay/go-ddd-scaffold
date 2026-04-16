package authentication

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
)

// BenchmarkService_Register 测试注册流程性能
func BenchmarkService_Register(b *testing.B) {
	mockUserRepo := new(MockUserRepository)
	mockResetTokenRepo := new(MockPasswordResetTokenRepository)
	mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
	mockTokenService := new(MockTokenService)
	publisher := event.NewPublisher(nil)

	service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := "bench@example.com"
		password := "BenchmarkPass123!"
		userID := "user-123"

		mockUserRepo.On("ExistsByEmail", mock.Anything, email).Return(false).Once()
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(&user.User{
			ID:            userID,
			Email:         email,
			EmailVerified: false,
		}, nil).Once()
		mockEmailTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*EmailVerificationToken")).Return(nil).Once()
		mockTokenService.On("GenerateTokens", mock.Anything, mock.Anything, mock.Anything).Return(&TokenPair{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		}, nil).Once()

		cmd := RegisterCommand{
			Email:    email,
			Password: password,
		}
		_, _ = service.Register(context.Background(), cmd)

		mockUserRepo.ExpectedCalls = nil
		mockEmailTokenRepo.ExpectedCalls = nil
		mockTokenService.ExpectedCalls = nil
	}
}

// BenchmarkService_Login 测试登录流程性能
func BenchmarkService_Login(b *testing.B) {
	mockUserRepo := new(MockUserRepository)
	mockResetTokenRepo := new(MockPasswordResetTokenRepository)
	mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
	mockTokenService := new(MockTokenService)
	publisher := event.NewPublisher(nil)

	email := "bench@example.com"
	password := "BenchmarkPass123!"
	mockUser, _ := user.NewUser(email, password)
	mockUser.ID = "user-123"
	mockUser.VerifyEmail()

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
		Email:    email,
		Password: password,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Login(context.Background(), cmd)
		mockUserRepo.ExpectedCalls = nil
		mockTokenService.ExpectedCalls = nil
	}
}

// BenchmarkService_RefreshToken 测试Token刷新性能
func BenchmarkService_RefreshToken(b *testing.B) {
	mockUserRepo := new(MockUserRepository)
	mockResetTokenRepo := new(MockPasswordResetTokenRepository)
	mockEmailTokenRepo := new(MockEmailVerificationTokenRepository)
	mockTokenService := new(MockTokenService)
	publisher := event.NewPublisher(nil)

	userID := "user-123"
	email := "bench@example.com"

	mockTokenService.On("ValidateRefreshToken", mock.Anything, "refresh-token").Return(userID, email, nil)
	mockTokenService.On("GenerateTokens", mock.Anything, userID, email).Return(&TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresIn:    3600,
	}, nil)

	service := NewService(mockUserRepo, mockResetTokenRepo, mockEmailTokenRepo, mockTokenService, publisher, nil)

	cmd := RefreshTokenCommand{
		RefreshToken: "refresh-token",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.RefreshToken(context.Background(), cmd)
		mockTokenService.ExpectedCalls = nil
	}
}

// BenchmarkPassword_Hash 测试密码哈希性能(bcrypt)
func BenchmarkPassword_Hash(b *testing.B) {
	password := "BenchmarkPass123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = user.NewUser("test@example.com", password)
	}
}

// BenchmarkPassword_Verify 测试密码验证性能
func BenchmarkPassword_Verify(b *testing.B) {
	password := "BenchmarkPass123!"
	testUser, _ := user.NewUser("test@example.com", password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = testUser.VerifyPassword(password)
	}
}

// BenchmarkToken_Generation 测试Token生成性能
func BenchmarkToken_Generation(b *testing.B) {
	mockTokenService := new(MockTokenService)
	userID := "user-123"
	email := "bench@example.com"

	mockTokenService.On("GenerateTokens", mock.Anything, userID, email).Return(&TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockTokenService.GenerateTokens(context.Background(), userID, email)
		mockTokenService.ExpectedCalls = nil
	}
}
