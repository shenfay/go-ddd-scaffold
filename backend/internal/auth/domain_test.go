package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shenfay/go-ddd-scaffold/internal/auth"
)

func TestUser_VerifyPassword(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name     string
		password string
		input    string
		want     bool
	}{
		{
			name:     "correct password should return true",
			password: "Password123!",
			input:    "Password123!",
			want:     true,
		},
		{
			name:     "wrong password should return false",
			password: "Password123!",
			input:    "WrongPassword",
			want:     false,
		},
		{
			name:     "empty password should return false",
			password: "Password123!",
			input:    "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			user, err := auth.NewUser("test@example.com", tt.password)
			require.NoError(t, err)
			require.NotNil(t, user)

			got := user.VerifyPassword(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUser_IsLocked(t *testing.T) {
	t.Parallel()

	t.Run("new user should not be locked", func(t *testing.T) {
		t.Parallel()
		
		user, err := auth.NewUser("test@example.com", "Password123!")
		require.NoError(t, err)
		require.NotNil(t, user)
		
		assert.False(t, user.IsLocked())
	})

	t.Run("locked user should return true", func(t *testing.T) {
		t.Parallel()
		
		user, err := auth.NewUser("test@example.com", "Password123!")
		require.NoError(t, err)
		require.NotNil(t, user)
		
		// 模拟锁定用户
		for i := 0; i < 5; i++ {
			user.IncrementFailedAttempts(5)
		}
		assert.True(t, user.IsLocked())
	})
}

func TestUser_IncrementFailedAttempts(t *testing.T) {
	t.Parallel()

	t.Run("should increment failed attempts", func(t *testing.T) {
		t.Parallel()
		
		user, err := auth.NewUser("test@example.com", "Password123!")
		require.NoError(t, err)

		assert.Equal(t, 0, user.FailedAttempts)

		user.IncrementFailedAttempts(5)
		assert.Equal(t, 1, user.FailedAttempts)

		user.IncrementFailedAttempts(5)
		assert.Equal(t, 2, user.FailedAttempts)
	})

	t.Run("should lock account after max attempts", func(t *testing.T) {
		t.Parallel()
		
		user, err := auth.NewUser("test@example.com", "Password123!")
		require.NoError(t, err)

		maxAttempts := 5
		for i := 0; i < maxAttempts; i++ {
			user.IncrementFailedAttempts(maxAttempts)
		}

		assert.True(t, user.IsLocked())
		assert.Equal(t, maxAttempts, user.FailedAttempts)
	})

	t.Run("should update UpdatedAt timestamp", func(t *testing.T) {
		t.Parallel()
		
		user, err := auth.NewUser("test@example.com", "Password123!")
		require.NoError(t, err)

		beforeUpdate := user.UpdatedAt

		time.Sleep(10 * time.Millisecond)

		user.IncrementFailedAttempts(5)

		assert.True(t, user.UpdatedAt.After(beforeUpdate))
	})
}

func TestUser_ResetFailedAttempts(t *testing.T) {
	t.Parallel()

	user, err := auth.NewUser("test@example.com", "Password123!")
	require.NoError(t, err)

	// 增加失败次数
	user.IncrementFailedAttempts(5)
	user.IncrementFailedAttempts(5)
	assert.Equal(t, 2, user.FailedAttempts)

	// 重置失败次数
	user.ResetFailedAttempts()
	assert.Equal(t, 0, user.FailedAttempts)
	assert.False(t, user.IsLocked())
}

func TestUser_UpdateLastLogin(t *testing.T) {
	t.Parallel()

	user, err := auth.NewUser("test@example.com", "Password123!")
	require.NoError(t, err)

	assert.Nil(t, user.LastLoginAt)

	user.UpdateLastLogin()

	require.NotNil(t, user.LastLoginAt)
	assert.True(t, user.LastLoginAt.Before(time.Now()))
	assert.True(t, user.LastLoginAt.After(user.CreatedAt))
}

func TestUser_VerifyEmail(t *testing.T) {
	t.Parallel()

	user, err := auth.NewUser("test@example.com", "Password123!")
	require.NoError(t, err)

	assert.False(t, user.EmailVerified)

	user.VerifyEmail()

	assert.True(t, user.EmailVerified)
}

func TestUser_ChangePassword(t *testing.T) {
	t.Parallel()

	user, err := auth.NewUser("test@example.com", "Password123!")
	require.NoError(t, err)

	oldHash := user.Password

	err = user.ChangePassword("NewPassword456!")
	require.NoError(t, err)

	assert.NotEqual(t, oldHash, user.Password)
	assert.True(t, user.UpdatedAt.After(user.CreatedAt))

	// 验证新密码有效
	assert.True(t, user.VerifyPassword("NewPassword456!"))
	assert.False(t, user.VerifyPassword("Password123!"))
}

func TestNewUser(t *testing.T) {
	t.Parallel()

	t.Run("should create user with valid email and password", func(t *testing.T) {
		t.Parallel()
		
		email := "test@example.com"
		password := "Password123!"

		user, err := auth.NewUser(email, password)

		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.NotEmpty(t, user.ID)
		assert.NotEqual(t, password, user.Password) // 应该是哈希值
		assert.False(t, user.EmailVerified)
		assert.False(t, user.Locked)
		assert.Equal(t, 0, user.FailedAttempts)
	})

	t.Run("should hash password", func(t *testing.T) {
		t.Parallel()
		
		password := "Password123!"
		user, err := auth.NewUser("test@example.com", password)
		require.NoError(t, err)

		// 验证密码被正确哈希
		assert.True(t, user.VerifyPassword(password))
		assert.False(t, user.VerifyPassword("WrongPassword"))
	})

	t.Run("should set timestamps", func(t *testing.T) {
		t.Parallel()
		
		beforeCreate := time.Now()
		user, err := auth.NewUser("test@example.com", "Password123!")
		require.NoError(t, err)

		assert.True(t, user.CreatedAt.After(beforeCreate) || user.CreatedAt.Equal(beforeCreate))
		assert.True(t, user.CreatedAt.Before(time.Now()))
		assert.Equal(t, user.CreatedAt, user.UpdatedAt)
	})
}

func TestUser_ID_Format(t *testing.T) {
	t.Parallel()

	user, err := auth.NewUser("test@example.com", "Password123!")
	require.NoError(t, err)

	// 验证 ID 格式：user_{ulid}
	assert.Contains(t, user.ID, "user_")
	assert.Greater(t, len(user.ID), len("user_"))
}
