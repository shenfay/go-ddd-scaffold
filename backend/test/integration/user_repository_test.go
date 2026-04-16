package integration

import (
	"context"
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/test/integration/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Create and FindByEmail", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()

		userRepo := repository.NewUserRepository(testDB.DB, nil)

		email := "test@example.com"
		password := "ValidPassword123!"
		testUser, err := user.NewUser(email, password)
		require.NoError(t, err)

		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)
		assert.NotEmpty(t, testUser.ID)

		foundUser, err := userRepo.FindByEmail(ctx, email)
		require.NoError(t, err)
		assert.Equal(t, email, foundUser.Email)
		assert.Equal(t, testUser.ID, foundUser.ID)
	})

	t.Run("ExistsByEmail", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		userRepo := repository.NewUserRepository(testDB.DB, nil)

		email := "exists@example.com"
		testUser, _ := user.NewUser(email, "Password123!")
		_ = userRepo.Create(ctx, testUser)

		exists := userRepo.ExistsByEmail(ctx, email)
		assert.True(t, exists)

		notExists := userRepo.ExistsByEmail(ctx, "nonexistent@example.com")
		assert.False(t, notExists)
	})

	t.Run("Update user", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		userRepo := repository.NewUserRepository(testDB.DB, nil)

		email := "update@example.com"
		testUser, _ := user.NewUser(email, "Password123!")
		_ = userRepo.Create(ctx, testUser)

		testUser.VerifyEmail()
		err := userRepo.Update(ctx, testUser)
		require.NoError(t, err)

		foundUser, _ := userRepo.FindByEmail(ctx, email)
		assert.True(t, foundUser.EmailVerified)
	})

	t.Run("FindByID", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		userRepo := repository.NewUserRepository(testDB.DB, nil)

		email := "findbyid@example.com"
		testUser, _ := user.NewUser(email, "Password123!")
		_ = userRepo.Create(ctx, testUser)

		foundUser, err := userRepo.FindByID(ctx, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, foundUser.ID)
		assert.Equal(t, email, foundUser.Email)
	})
}
