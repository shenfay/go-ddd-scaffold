package authentication

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	args := m.Called(ctx, email)
	return args.Bool(0)
}

// MockPasswordResetTokenRepository 模拟密码重置令牌仓储
type MockPasswordResetTokenRepository struct {
	mock.Mock
}

func (m *MockPasswordResetTokenRepository) Create(ctx context.Context, token *authentication.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockPasswordResetTokenRepository) FindByToken(ctx context.Context, token string) (*authentication.PasswordResetToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authentication.PasswordResetToken), args.Error(1)
}

func (m *MockPasswordResetTokenRepository) MarkAsUsed(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}

func (m *MockPasswordResetTokenRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockEmailVerificationTokenRepository 模拟邮箱验证令牌仓储
type MockEmailVerificationTokenRepository struct {
	mock.Mock
}

func (m *MockEmailVerificationTokenRepository) Create(ctx context.Context, token *authentication.EmailVerificationToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockEmailVerificationTokenRepository) FindByToken(ctx context.Context, token string) (*authentication.EmailVerificationToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authentication.EmailVerificationToken), args.Error(1)
}

func (m *MockEmailVerificationTokenRepository) FindByUserID(ctx context.Context, userID string) (*authentication.EmailVerificationToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authentication.EmailVerificationToken), args.Error(1)
}

func (m *MockEmailVerificationTokenRepository) MarkAsUsed(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}

func (m *MockEmailVerificationTokenRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockTokenService 模拟Token服务
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateTokens(ctx context.Context, userID, email string) (*TokenPair, error) {
	args := m.Called(ctx, userID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TokenPair), args.Error(1)
}

func (m *MockTokenService) RevokeToken(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}

func (m *MockTokenService) ValidateRefreshTokenWithDevice(ctx context.Context, token string) (*DeviceInfo, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DeviceInfo), args.Error(1)
}

func (m *MockTokenService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*JWTClaims), args.Error(1)
}

func (m *MockTokenService) StoreDeviceInfo(ctx context.Context, token string, deviceInfo DeviceInfo) error {
	args := m.Called(ctx, token, deviceInfo)
	return args.Error(0)
}

func (m *MockTokenService) RevokeDeviceByToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenService) RevokeAllDevices(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTokenService) GetUserDevices(ctx context.Context, userID string) ([]DeviceInfo, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DeviceInfo), args.Error(1)
}
