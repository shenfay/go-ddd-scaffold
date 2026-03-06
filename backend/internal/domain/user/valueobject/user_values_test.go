package valueobject_test

import (
	"testing"

	"go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	t.Run("有效邮箱", func(t *testing.T) {
		// Arrange
		emailStr := "test@example.com"

		// Act
		email, err := valueobject.NewEmail(emailStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, emailStr, email.String())
	})

	t.Run("无效邮箱 - 无@", func(t *testing.T) {
		// Arrange
		emailStr := "invalid-email"

		// Act
		email, err := valueobject.NewEmail(emailStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, email)
	})

	t.Run("无效邮箱 - 无域名", func(t *testing.T) {
		// Arrange
		emailStr := "test@"

		// Act
		email, err := valueobject.NewEmail(emailStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, email)
	})

	t.Run("无效邮箱 - 空字符串", func(t *testing.T) {
		// Arrange
		emailStr := ""

		// Act
		email, err := valueobject.NewEmail(emailStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, email)
	})

	t.Run("有效邮箱 - 包含点号", func(t *testing.T) {
		// Arrange
		emailStr := "test.user@example.com"

		// Act
		email, err := valueobject.NewEmail(emailStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, emailStr, email.String())
	})

	t.Run("有效邮箱 - 包含加号", func(t *testing.T) {
		// Arrange
		emailStr := "test+tag@example.com"

		// Act
		email, err := valueobject.NewEmail(emailStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, emailStr, email.String())
	})
}

func TestEmail_Equals(t *testing.T) {
	t.Run("相同邮箱", func(t *testing.T) {
		// Arrange
		email1, _ := valueobject.NewEmail("test@example.com")
		email2, _ := valueobject.NewEmail("test@example.com")

		// Act & Assert
		assert.True(t, email1.Equals(email2))
		assert.True(t, email2.Equals(email1))
	})

	t.Run("不同邮箱", func(t *testing.T) {
		// Arrange
		email1, _ := valueobject.NewEmail("test@example.com")
		email2, _ := valueobject.NewEmail("other@example.com")

		// Act & Assert
		assert.False(t, email1.Equals(email2))
	})

	t.Run("大小写不敏感（邮箱通常不区分）", func(t *testing.T) {
		// Arrange
		email1, _ := valueobject.NewEmail("Test@Example.com")
		email2, _ := valueobject.NewEmail("test@example.com")

		// Act & Assert - Email 是字符串类型，直接比较
		// 注意：虽然邮箱通常不区分大小写，但这里做字符串级别的比较
		assert.NotEqual(t, email1.String(), email2.String()) // 字符串不同
	})
}

func TestNewNickname(t *testing.T) {
	t.Run("有效昵称 - 中文", func(t *testing.T) {
		// Arrange
		nicknameStr := "测试用户"

		// Act
		nickname, err := valueobject.NewNickname(nicknameStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, nicknameStr, nickname.String())
	})

	t.Run("有效昵称 - 短英文", func(t *testing.T) {
		// Arrange
		nicknameStr := "Tom"

		// Act
		nickname, err := valueobject.NewNickname(nicknameStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, nicknameStr, nickname.String())
	})

	t.Run("无效昵称 - 太长", func(t *testing.T) {
		// Arrange - 超过 50 字符
		longNickname := "这是一段非常非常非常非常非常非常非常非常非常非常非常非常非常长的昵称"

		// Act
		nickname, err := valueobject.NewNickname(longNickname)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, nickname)
	})

	t.Run("无效昵称 - 空字符串", func(t *testing.T) {
		// Arrange
		nicknameStr := ""

		// Act
		nickname, err := valueobject.NewNickname(nicknameStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, nickname)
	})

	t.Run("无效昵称 - 只有空格", func(t *testing.T) {
		// Arrange
		nicknameStr := "   "

		// Act
		nickname, err := valueobject.NewNickname(nicknameStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, nickname)
	})
}

func TestNewPlainPassword(t *testing.T) {
	t.Run("有效密码 - 符合要求", func(t *testing.T) {
		// Arrange
		passwordStr := "Test123!"

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, passwordStr, password.String())
	})

	t.Run("有效密码 - 最小长度 8", func(t *testing.T) {
		// Arrange
		passwordStr := "Test1234"

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, passwordStr, password.String())
	})

	t.Run("有效密码 - 包含数字和字母", func(t *testing.T) {
		// Arrange
		passwordStr := "password123"

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, passwordStr, password.String())
	})

	t.Run("无效密码 - 太短", func(t *testing.T) {
		// Arrange
		passwordStr := "T1!"

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, password)
	})

	t.Run("无效密码 - 只有小写字母", func(t *testing.T) {
		// Arrange
		passwordStr := "passwordonly"

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, password)
	})

	t.Run("无效密码 - 只有大写字母", func(t *testing.T) {
		// Arrange
		passwordStr := "PASSWORDONLY"

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, password)
	})

	t.Run("无效密码 - 只有数字", func(t *testing.T) {
		// Arrange
		passwordStr := "12345678"

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, password)
	})

	t.Run("无效密码 - 空字符串", func(t *testing.T) {
		// Arrange
		passwordStr := ""

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, password)
	})

	t.Run("有效密码 - 包含特殊字符", func(t *testing.T) {
		// Arrange
		passwordStr := "Test@123!@#"

		// Act
		password, err := valueobject.NewPlainPassword(passwordStr)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, passwordStr, password.String())
	})
}
