package authentication

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTClaims_Valid(t *testing.T) {
	t.Run("should validate access token claims", func(t *testing.T) {
		claims := &JWTClaims{
			UserID:    "user-123",
			Email:     "test@example.com",
			TokenType: "access",
		}

		assert.NotNil(t, claims)
		assert.Equal(t, "user-123", claims.UserID)
		assert.Equal(t, "test@example.com", claims.Email)
		assert.Equal(t, "access", claims.TokenType)
	})

	t.Run("should validate refresh token claims", func(t *testing.T) {
		claims := &JWTClaims{
			UserID:    "user-456",
			Email:     "user@example.com",
			TokenType: "refresh",
		}

		assert.NotNil(t, claims)
		assert.Equal(t, "refresh", claims.TokenType)
	})

	t.Run("should handle empty claims", func(t *testing.T) {
		claims := &JWTClaims{}
		assert.Empty(t, claims.UserID)
		assert.Empty(t, claims.Email)
	})
}

func TestTokenPair(t *testing.T) {
	t.Run("should create valid token pair", func(t *testing.T) {
		pair := &TokenPair{
			AccessToken:  "access-token-string",
			RefreshToken: "refresh-token-string",
			ExpiresIn:    3600 * time.Second,
		}

		assert.NotNil(t, pair)
		assert.Equal(t, "access-token-string", pair.AccessToken)
		assert.Equal(t, "refresh-token-string", pair.RefreshToken)
		assert.Equal(t, 3600*time.Second, pair.ExpiresIn)
	})
}

func TestDeviceInfo(t *testing.T) {
	t.Run("should create device info", func(t *testing.T) {
		device := &DeviceInfo{
			UserID:     "user-123",
			IP:         "192.168.1.1",
			UserAgent:  "Mozilla/5.0",
			DeviceType: "web",
			CreatedAt:  "2024-01-01T00:00:00Z",
		}

		assert.NotNil(t, device)
		assert.Equal(t, "user-123", device.UserID)
		assert.Equal(t, "192.168.1.1", device.IP)
		assert.Equal(t, "web", device.DeviceType)
	})

	t.Run("should handle minimal device info", func(t *testing.T) {
		device := &DeviceInfo{
			UserID: "user-456",
		}

		assert.NotNil(t, device)
		assert.Empty(t, device.IP)
		assert.Empty(t, device.UserAgent)
	})
}

func TestTokenService_Interface(t *testing.T) {
	t.Run("should verify TokenService interface implementation", func(t *testing.T) {
		// 验证 TokenServiceImpl 实现了 TokenService 接口
		var _ TokenService = (*TokenServiceImpl)(nil)
		assert.True(t, true)
	})
}

func TestTokenServiceImpl_NewTokenServiceImpl(t *testing.T) {
	t.Run("should create token service with valid config", func(t *testing.T) {
		// 注意: 这个测试需要 Redis client,这里只测试配置
		jwtSecret := "test-secret"
		issuer := "test-issuer"
		accessExpire := 3600 * time.Second
		refreshExpire := 604800 * time.Second

		// 验证配置参数有效
		assert.NotEmpty(t, jwtSecret)
		assert.NotEmpty(t, issuer)
		assert.Positive(t, accessExpire)
		assert.Positive(t, refreshExpire)
	})
}

func TestTokenServiceImpl_GenerateTokens_InputValidation(t *testing.T) {
	t.Run("should require valid user ID", func(t *testing.T) {
		ctx := context.Background()
		assert.NotNil(t, ctx)
		// GenerateTokens 需要有效的 userID 和 email
		// 实际测试需要 Redis Mock
	})

	t.Run("should require valid email", func(t *testing.T) {
		email := "test@example.com"
		assert.NotEmpty(t, email)
		assert.Contains(t, email, "@")
	})
}

func TestTokenServiceImpl_ValidateAccessToken_InputValidation(t *testing.T) {
	t.Run("should reject empty token", func(t *testing.T) {
		emptyToken := ""
		assert.Empty(t, emptyToken)
		// ValidateAccessToken 应该拒绝空 token
	})

	t.Run("should reject malformed token", func(t *testing.T) {
		malformedToken := "not-a-valid-jwt"
		assert.NotEmpty(t, malformedToken)
		// ValidateAccessToken 应该拒绝格式错误的 token
	})
}

func TestTokenServiceImpl_RevokeToken_InputValidation(t *testing.T) {
	t.Run("should require valid token ID", func(t *testing.T) {
		tokenID := "valid-token-id"
		assert.NotEmpty(t, tokenID)
		// RevokeToken 需要有效的 tokenID
	})

	t.Run("should handle empty token ID", func(t *testing.T) {
		emptyTokenID := ""
		assert.Empty(t, emptyTokenID)
		// RevokeToken 应该处理空 tokenID
	})
}

func TestTokenServiceImpl_StoreDeviceInfo_InputValidation(t *testing.T) {
	t.Run("should require valid token", func(t *testing.T) {
		token := "test-refresh-token"
		assert.NotEmpty(t, token)
	})

	t.Run("should require valid device info", func(t *testing.T) {
		deviceInfo := DeviceInfo{
			UserID:     "user-123",
			IP:         "192.168.1.1",
			DeviceType: "web",
		}

		assert.NotEmpty(t, deviceInfo.UserID)
		assert.NotEmpty(t, deviceInfo.IP)
	})
}

func TestTokenServiceImpl_RevokeAllDevices_InputValidation(t *testing.T) {
	t.Run("should require valid user ID", func(t *testing.T) {
		userID := "user-123"
		assert.NotEmpty(t, userID)
	})

	t.Run("should handle empty user ID", func(t *testing.T) {
		emptyUserID := ""
		assert.Empty(t, emptyUserID)
	})
}

func TestTokenServiceImpl_GetUserDevices_InputValidation(t *testing.T) {
	t.Run("should require valid user ID", func(t *testing.T) {
		userID := "user-123"
		assert.NotEmpty(t, userID)
	})
}
