package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_ChangePassword(t *testing.T) {
	t.Run("should change password successfully", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "OldPassword123!")
		
		err := user.ChangePassword("NewPassword456!")
		
		assert.NoError(t, err)
		assert.True(t, user.VerifyPassword("NewPassword456!"))
		assert.False(t, user.VerifyPassword("OldPassword123!"))
	})

	t.Run("should handle short password", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		
		// bcrypt 可以哈希任何密码,包括短密码
		err := user.ChangePassword("123")
		
		assert.NoError(t, err)
		assert.True(t, user.VerifyPassword("123"))
	})
}

func TestUser_UpdateEmail(t *testing.T) {
	t.Run("should update email successfully", func(t *testing.T) {
		user, _ := NewUser("old@example.com", "ValidPassword123!")
		oldUpdatedAt := user.UpdatedAt
		
		err := user.UpdateEmail("new@example.com")
		
		assert.NoError(t, err)
		assert.Equal(t, "new@example.com", user.Email)
		assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
	})

	t.Run("should handle empty email", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		
		err := user.UpdateEmail("")
		
		assert.NoError(t, err)
		assert.Empty(t, user.Email)
	})
}

func TestUser_VerifyPassword(t *testing.T) {
	t.Run("should return true for correct password", func(t *testing.T) {
		password := "CorrectPassword123!"
		user, _ := NewUser("test@example.com", password)
		
		assert.True(t, user.VerifyPassword(password))
	})

	t.Run("should return false for incorrect password", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "CorrectPassword123!")
		
		assert.False(t, user.VerifyPassword("WrongPassword456!"))
	})

	t.Run("should return false for empty password", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		
		assert.False(t, user.VerifyPassword(""))
	})
}

func TestUser_IncrementFailedAttempts(t *testing.T) {
	t.Run("should increment failed attempts", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		
		user.IncrementFailedAttempts(5)
		
		assert.Equal(t, 1, user.FailedAttempts)
		assert.False(t, user.IsLocked())
	})

	t.Run("should lock account after max attempts", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		
		for i := 0; i < 5; i++ {
			user.IncrementFailedAttempts(5)
		}
		
		assert.Equal(t, 5, user.FailedAttempts)
		assert.True(t, user.IsLocked())
	})

	t.Run("should not lock if below max attempts", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		
		user.IncrementFailedAttempts(3)
		user.IncrementFailedAttempts(3)
		
		assert.Equal(t, 2, user.FailedAttempts)
		assert.False(t, user.IsLocked())
	})
}

func TestUser_ResetFailedAttempts(t *testing.T) {
	t.Run("should reset failed attempts to zero", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		user.FailedAttempts = 3
		user.Locked = true
		
		user.ResetFailedAttempts()
		
		assert.Equal(t, 0, user.FailedAttempts)
		assert.True(t, user.Locked) // ResetFailedAttempts 不应解锁
	})
}

func TestUser_VerifyEmail(t *testing.T) {
	t.Run("should mark email as verified", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		assert.False(t, user.EmailVerified)
		
		user.VerifyEmail()
		
		assert.True(t, user.EmailVerified)
	})
}

func TestUser_IsLocked(t *testing.T) {
	t.Run("should return true when locked", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		user.Locked = true
		
		assert.True(t, user.IsLocked())
	})

	t.Run("should return false when not locked", func(t *testing.T) {
		user, _ := NewUser("test@example.com", "ValidPassword123!")
		
		assert.False(t, user.IsLocked())
	})
}
