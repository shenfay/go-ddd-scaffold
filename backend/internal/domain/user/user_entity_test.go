package user_test

import (
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestUser_NewUser(t *testing.T) {
	email := "test@example.com"
	password := "ValidPassword123!"

	u, err := user.NewUser(email, password)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, email, u.Email)
	assert.NotEmpty(t, u.ID)
	assert.False(t, u.EmailVerified)
	assert.False(t, u.Locked)
	assert.Equal(t, 0, u.FailedAttempts)
	assert.NotEmpty(t, u.Password)
	assert.NotEqual(t, password, u.Password)
}

func TestUser_VerifyPassword(t *testing.T) {
	t.Run("should return true for correct password", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "ValidPassword123!")
		assert.True(t, u.VerifyPassword("ValidPassword123!"))
	})

	t.Run("should return false for incorrect password", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "ValidPassword123!")
		assert.False(t, u.VerifyPassword("WrongPassword"))
	})

	t.Run("should return false for empty password", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "ValidPassword123!")
		assert.False(t, u.VerifyPassword(""))
	})
}

func TestUser_IncrementFailedAttempts(t *testing.T) {
	t.Run("should increment failed attempts", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "ValidPassword123!")
		maxAttempts := 5

		u.IncrementFailedAttempts(maxAttempts)
		assert.Equal(t, 1, u.FailedAttempts)
		assert.False(t, u.Locked)
	})

	t.Run("should lock account after max attempts", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "ValidPassword123!")
		maxAttempts := 5

		for i := 0; i < maxAttempts; i++ {
			u.IncrementFailedAttempts(maxAttempts)
		}

		assert.Equal(t, maxAttempts, u.FailedAttempts)
		assert.True(t, u.Locked)
	})

	t.Run("should not lock before max attempts", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "ValidPassword123!")
		maxAttempts := 5

		for i := 0; i < maxAttempts-1; i++ {
			u.IncrementFailedAttempts(maxAttempts)
		}

		assert.Equal(t, 4, u.FailedAttempts)
		assert.False(t, u.Locked)
	})
}

func TestUser_ResetFailedAttempts(t *testing.T) {
	u, _ := user.NewUser("test@example.com", "ValidPassword123!")
	maxAttempts := 5

	for i := 0; i < 3; i++ {
		u.IncrementFailedAttempts(maxAttempts)
	}
	assert.Equal(t, 3, u.FailedAttempts)

	u.ResetFailedAttempts()
	assert.Equal(t, 0, u.FailedAttempts)
	assert.False(t, u.Locked)
}

func TestUser_UpdateLastLogin(t *testing.T) {
	u, _ := user.NewUser("test@example.com", "ValidPassword123!")

	assert.Nil(t, u.LastLoginAt)

	u.UpdateLastLogin()

	assert.NotNil(t, u.LastLoginAt)
	assert.True(t, u.LastLoginAt.Before(time.Now().Add(1*time.Second)))
}

func TestUser_VerifyEmail(t *testing.T) {
	u, _ := user.NewUser("test@example.com", "ValidPassword123!")

	assert.False(t, u.EmailVerified)

	u.VerifyEmail()

	assert.True(t, u.EmailVerified)
}

func TestUser_ChangePassword(t *testing.T) {
	t.Run("should change password successfully", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "OldPassword123!")

		assert.True(t, u.VerifyPassword("OldPassword123!"))
		assert.False(t, u.VerifyPassword("NewPassword456!"))

		err := u.ChangePassword("NewPassword456!")
		assert.NoError(t, err)

		assert.False(t, u.VerifyPassword("OldPassword123!"))
		assert.True(t, u.VerifyPassword("NewPassword456!"))
	})
}

func TestUser_UpdateEmail(t *testing.T) {
	u, _ := user.NewUser("old@example.com", "ValidPassword123!")

	assert.Equal(t, "old@example.com", u.Email)

	err := u.UpdateEmail("new@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "new@example.com", u.Email)
}

func TestUser_IsLocked(t *testing.T) {
	u, _ := user.NewUser("test@example.com", "ValidPassword123!")
	maxAttempts := 5

	assert.False(t, u.IsLocked())

	for i := 0; i < maxAttempts; i++ {
		u.IncrementFailedAttempts(maxAttempts)
	}

	assert.True(t, u.IsLocked())
}

func TestUser_UpdatedAt(t *testing.T) {
	t.Run("should update UpdatedAt on password change", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "OldPassword123!")
		oldUpdatedAt := u.UpdatedAt

		time.Sleep(10 * time.Millisecond)

		err := u.ChangePassword("NewPassword456!")
		assert.NoError(t, err)

		assert.True(t, u.UpdatedAt.After(oldUpdatedAt))
	})

	t.Run("should update UpdatedAt on email verification", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "ValidPassword123!")
		oldUpdatedAt := u.UpdatedAt

		time.Sleep(10 * time.Millisecond)

		u.VerifyEmail()

		assert.True(t, u.UpdatedAt.After(oldUpdatedAt))
	})
}

func TestUser_PasswordHashing(t *testing.T) {
	u1, _ := user.NewUser("user1@example.com", "SamePassword123!")
	u2, _ := user.NewUser("user2@example.com", "SamePassword123!")

	assert.NotEqual(t, u1.Password, u2.Password)

	assert.True(t, u1.VerifyPassword("SamePassword123!"))
	assert.True(t, u2.VerifyPassword("SamePassword123!"))
}
