package user_test

import (
	"strings"
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "ValidPassword123!"

		hashed, err := user.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)
	})

	t.Run("should generate different hashes for same password", func(t *testing.T) {
		password := "SamePassword123!"

		hashed1, err := user.HashPassword(password)
		assert.NoError(t, err)

		hashed2, err := user.HashPassword(password)
		assert.NoError(t, err)

		assert.NotEqual(t, hashed1, hashed2)
	})

	t.Run("should handle empty password", func(t *testing.T) {
		hashed, err := user.HashPassword("")

		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})

	t.Run("should handle long password", func(t *testing.T) {
		password := strings.Repeat("A", 100)

		hashed, err := user.HashPassword(password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds 72 bytes")
		assert.Empty(t, hashed)
	})

	t.Run("should handle max length password (72 bytes)", func(t *testing.T) {
		password := strings.Repeat("A", 72)

		hashed, err := user.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})

	t.Run("should handle special characters", func(t *testing.T) {
		password := "!@#$%^&*()_+-=[]{}|;':\",./<>?"

		hashed, err := user.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})

	t.Run("should produce bcrypt hash format", func(t *testing.T) {
		password := "TestPassword123!"

		hashed, err := user.HashPassword(password)
		assert.NoError(t, err)

		assert.True(t, strings.HasPrefix(hashed, "$2a$"),
			"Bcrypt hash should start with $2a$, got: %s", hashed[:4])
	})
}
