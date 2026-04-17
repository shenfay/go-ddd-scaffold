package repository

import (
	"context"
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db, nil)

	t.Run("should create user successfully", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "password123")

		err := repo.Create(context.Background(), u)

		assert.NoError(t, err)
		assert.NotEmpty(t, u.ID)
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		u1, _ := user.NewUser("duplicate@example.com", "password123")
		err := repo.Create(context.Background(), u1)
		assert.NoError(t, err)

		u2, _ := user.NewUser("duplicate@example.com", "password456")
		err = repo.Create(context.Background(), u2)

		assert.Error(t, err)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db, nil)

	t.Run("should find user by email", func(t *testing.T) {
		u, _ := user.NewUser("find@example.com", "password123")
		err := repo.Create(context.Background(), u)
		assert.NoError(t, err)

		found, err := repo.FindByEmail(context.Background(), "find@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, u.ID, found.ID)
		assert.Equal(t, "find@example.com", found.Email)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		_, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")

		assert.Error(t, err)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db, nil)

	t.Run("should find user by ID", func(t *testing.T) {
		u, _ := user.NewUser("findbyid@example.com", "password123")
		err := repo.Create(context.Background(), u)
		assert.NoError(t, err)

		found, err := repo.FindByID(context.Background(), u.ID)

		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, u.ID, found.ID)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		_, err := repo.FindByID(context.Background(), "nonexistent-id")

		assert.Error(t, err)
	})
}

func TestUserRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db, nil)

	t.Run("should update user", func(t *testing.T) {
		u, _ := user.NewUser("save@example.com", "password123")
		err := repo.Create(context.Background(), u)
		assert.NoError(t, err)

		// Modify and save
		u.VerifyEmail()
		err = repo.Save(context.Background(), u)

		assert.NoError(t, err)

		// Verify changes persisted
		found, _ := repo.FindByID(context.Background(), u.ID)
		assert.True(t, found.EmailVerified)
	})
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db, nil)

	t.Run("should update user fields", func(t *testing.T) {
		u, _ := user.NewUser("update@example.com", "password123")
		err := repo.Create(context.Background(), u)
		assert.NoError(t, err)

		// Update failed attempts
		u.IncrementFailedAttempts(5)
		err = repo.Update(context.Background(), u)

		assert.NoError(t, err)

		// Verify changes
		found, _ := repo.FindByID(context.Background(), u.ID)
		assert.Equal(t, 1, found.FailedAttempts)
	})
}

func TestUserRepository_Save_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db, nil)

	t.Run("should save user with changes", func(t *testing.T) {
		u, _ := user.NewUser("saveupdate@example.com", "password123")
		err := repo.Create(context.Background(), u)
		assert.NoError(t, err)

		// Modify and save
		u.VerifyEmail()
		err = repo.Save(context.Background(), u)

		assert.NoError(t, err)

		// Verify changes
		found, _ := repo.FindByID(context.Background(), u.ID)
		assert.True(t, found.EmailVerified)
	})
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db, nil)

	t.Run("should return true for existing email", func(t *testing.T) {
		u, _ := user.NewUser("exists@example.com", "password123")
		err := repo.Create(context.Background(), u)
		assert.NoError(t, err)

		exists := repo.ExistsByEmail(context.Background(), "exists@example.com")

		assert.True(t, exists)
	})

	t.Run("should return false for non-existing email", func(t *testing.T) {
		exists := repo.ExistsByEmail(context.Background(), "nonexistent@example.com")

		assert.False(t, exists)
	})
}
