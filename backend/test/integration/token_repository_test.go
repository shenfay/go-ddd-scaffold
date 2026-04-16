package integration

import (
	"context"
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/test/integration/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordResetTokenRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Create and FindByToken", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()

		tokenRepo := repository.NewPasswordResetTokenRepository(testDB.DB)

		userID := "user-123"
		token, err := authentication.NewPasswordResetToken(userID)
		require.NoError(t, err)

		err = tokenRepo.Create(ctx, token)
		require.NoError(t, err)
		assert.NotEmpty(t, token.ID)

		foundToken, err := tokenRepo.FindByToken(ctx, token.Token)
		require.NoError(t, err)
		assert.Equal(t, userID, foundToken.UserID)
		assert.Equal(t, token.Token, foundToken.Token)
		assert.False(t, foundToken.Used)
	})

	t.Run("MarkAsUsed", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		tokenRepo := repository.NewPasswordResetTokenRepository(testDB.DB)

		userID := "user-456"
		token, _ := authentication.NewPasswordResetToken(userID)
		_ = tokenRepo.Create(ctx, token)

		err := tokenRepo.MarkAsUsed(ctx, token.ID)
		require.NoError(t, err)

		foundToken, _ := tokenRepo.FindByToken(ctx, token.Token)
		assert.True(t, foundToken.Used)
	})

	t.Run("DeleteExpired", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		tokenRepo := repository.NewPasswordResetTokenRepository(testDB.DB)

		userID := "user-789"
		token, _ := authentication.NewPasswordResetToken(userID)
		_ = tokenRepo.Create(ctx, token)

		err := tokenRepo.DeleteExpired(ctx)
		require.NoError(t, err)

		foundToken, err := tokenRepo.FindByToken(ctx, token.Token)
		assert.NoError(t, err)
		assert.NotNil(t, foundToken)
	})

	t.Run("Token expiry check", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		tokenRepo := repository.NewPasswordResetTokenRepository(testDB.DB)

		userID := "user-expiry"
		token, _ := authentication.NewPasswordResetToken(userID)
		token.ExpiresAt = time.Now().Add(-1 * time.Hour)
		_ = tokenRepo.Create(ctx, token)

		foundToken, _ := tokenRepo.FindByToken(ctx, token.Token)
		assert.True(t, foundToken.IsExpired())
	})
}

func TestEmailVerificationTokenRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Create and FindByToken", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()

		tokenRepo := repository.NewEmailVerificationTokenRepository(testDB.DB)

		userID := "user-123"
		token, err := authentication.NewEmailVerificationToken(userID)
		require.NoError(t, err)

		err = tokenRepo.Create(ctx, token)
		require.NoError(t, err)
		assert.NotEmpty(t, token.ID)

		foundToken, err := tokenRepo.FindByToken(ctx, token.Token)
		require.NoError(t, err)
		assert.Equal(t, userID, foundToken.UserID)
		assert.Equal(t, token.Token, foundToken.Token)
		assert.False(t, foundToken.Used)
	})

	t.Run("FindByUserID", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		tokenRepo := repository.NewEmailVerificationTokenRepository(testDB.DB)

		userID := "user-find"
		token, _ := authentication.NewEmailVerificationToken(userID)
		_ = tokenRepo.Create(ctx, token)

		foundToken, err := tokenRepo.FindByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, userID, foundToken.UserID)
	})

	t.Run("MarkAsUsed", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		tokenRepo := repository.NewEmailVerificationTokenRepository(testDB.DB)

		userID := "user-mark"
		token, _ := authentication.NewEmailVerificationToken(userID)
		_ = tokenRepo.Create(ctx, token)

		err := tokenRepo.MarkAsUsed(ctx, token.ID)
		require.NoError(t, err)

		foundToken, _ := tokenRepo.FindByToken(ctx, token.Token)
		assert.True(t, foundToken.Used)
	})

	t.Run("DeleteExpired", func(t *testing.T) {
		testDB := testutil.SetupTestDatabase(t)
		ctx := context.Background()
		tokenRepo := repository.NewEmailVerificationTokenRepository(testDB.DB)

		userID := "user-delete"
		token, _ := authentication.NewEmailVerificationToken(userID)
		_ = tokenRepo.Create(ctx, token)

		err := tokenRepo.DeleteExpired(ctx)
		require.NoError(t, err)

		foundToken, err := tokenRepo.FindByToken(ctx, token.Token)
		assert.NoError(t, err)
		assert.NotNil(t, foundToken)
	})
}
