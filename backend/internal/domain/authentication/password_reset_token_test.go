package authentication_test

import (
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"github.com/stretchr/testify/assert"
)

func TestPasswordResetToken_Creation(t *testing.T) {
	// 测试正常创建令牌
	userID := "test-user-123"
	token, err := authentication.NewPasswordResetToken(userID)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, userID, token.UserID)
	assert.False(t, token.Used)
	assert.Len(t, token.Token, 64) // 32 字节的十六进制字符串
	assert.True(t, token.ExpiresAt.After(time.Now()))
	assert.True(t, token.ExpiresAt.Before(time.Now().Add(2*time.Hour)))
}

func TestPasswordResetToken_IsExpired(t *testing.T) {
	token, _ := authentication.NewPasswordResetToken("user-123")

	// 新创建的令牌不应过期
	assert.False(t, token.IsExpired())
}

func TestPasswordResetToken_MarkAsUsed(t *testing.T) {
	token, _ := authentication.NewPasswordResetToken("user-123")

	// 初始状态未使用
	assert.False(t, token.Used)

	// 标记为已使用
	token.MarkAsUsed()
	assert.True(t, token.Used)
}

func TestPasswordResetToken_SecureToken(t *testing.T) {
	// 测试生成的令牌具有足够的随机性
	token1, _ := authentication.NewPasswordResetToken("user-123")
	token2, _ := authentication.NewPasswordResetToken("user-123")

	// 两次生成的令牌应该不同
	assert.NotEqual(t, token1.Token, token2.Token)
}
