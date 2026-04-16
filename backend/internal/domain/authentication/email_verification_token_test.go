package authentication_test

import (
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"github.com/stretchr/testify/assert"
)

func TestEmailVerificationToken_Creation(t *testing.T) {
	// 测试正常创建邮箱验证令牌
	userID := "test-user-456"
	token, err := authentication.NewEmailVerificationToken(userID)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, userID, token.UserID)
	assert.False(t, token.Used)
	assert.Len(t, token.Token, 64) // 32 字节的十六进制字符串
	assert.True(t, token.ExpiresAt.After(time.Now()))
	// 邮箱验证令牌有效期为 24 小时
	assert.True(t, token.ExpiresAt.Before(time.Now().Add(25*time.Hour)))
	assert.True(t, token.ExpiresAt.After(time.Now().Add(23*time.Hour)))
}

func TestEmailVerificationToken_IsExpired(t *testing.T) {
	token, _ := authentication.NewEmailVerificationToken("user-456")

	// 新创建的令牌不应过期
	assert.False(t, token.IsExpired())
}

func TestEmailVerificationToken_MarkAsUsed(t *testing.T) {
	token, _ := authentication.NewEmailVerificationToken("user-456")

	// 初始状态未使用
	assert.False(t, token.Used)

	// 标记为已使用
	token.MarkAsUsed()
	assert.True(t, token.Used)

	// 确保不可逆(再次标记仍为 true)
	token.MarkAsUsed()
	assert.True(t, token.Used)
}

func TestEmailVerificationToken_SecureToken(t *testing.T) {
	// 测试生成的令牌具有足够的随机性
	token1, _ := authentication.NewEmailVerificationToken("user-456")
	token2, _ := authentication.NewEmailVerificationToken("user-456")

	// 两次生成的令牌应该不同
	assert.NotEqual(t, token1.Token, token2.Token)
	assert.NotEqual(t, token1.ID, token2.ID)
}

func TestEmailVerificationToken_DifferentUsers(t *testing.T) {
	// 测试不同用户的令牌独立性
	userID1 := "user-001"
	userID2 := "user-002"

	token1, err := authentication.NewEmailVerificationToken(userID1)
	assert.NoError(t, err)

	token2, err := authentication.NewEmailVerificationToken(userID2)
	assert.NoError(t, err)

	assert.Equal(t, userID1, token1.UserID)
	assert.Equal(t, userID2, token2.UserID)
	assert.NotEqual(t, token1.Token, token2.Token)
}

func TestEmailVerificationToken_TokenFormat(t *testing.T) {
	// 测试令牌格式(应为十六进制字符串)
	token, _ := authentication.NewEmailVerificationToken("user-789")

	// 验证令牌只包含十六进制字符
	for _, char := range token.Token {
		isHex := (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')
		assert.True(t, isHex, "Token should only contain hex characters, got: %c", char)
	}
}

func TestEmailVerificationToken_ExpiryBoundary(t *testing.T) {
	// 测试边界情况: 令牌在过期前后的状态
	token, _ := authentication.NewEmailVerificationToken("user-boundary")

	// 刚创建时未过期
	assert.False(t, token.IsExpired())

	// 模拟过期: 手动设置过期时间为过去
	token.ExpiresAt = time.Now().Add(-1 * time.Hour)
	assert.True(t, token.IsExpired())

	// 模拟刚好到期
	token.ExpiresAt = time.Now().Add(-1 * time.Millisecond)
	assert.True(t, token.IsExpired())

	// 模拟还未到期
	token.ExpiresAt = time.Now().Add(1 * time.Millisecond)
	assert.False(t, token.IsExpired())
}
