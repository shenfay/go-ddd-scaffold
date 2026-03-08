package transaction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/infrastructure/transaction"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	return db
}

func TestUnitOfWork_Commit(t *testing.T) {
	db := setupTestDB(t)
	uow := transaction.NewGormUnitOfWork(db)

	executed := false
	err := uow.WithTransaction(context.Background(), func(ctx context.Context) error {
		executed = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestUnitOfWork_Rollback(t *testing.T) {
	db := setupTestDB(t)
	uow := transaction.NewGormUnitOfWork(db)

	executed := false
	err := uow.WithTransaction(context.Background(), func(ctx context.Context) error {
		executed = true
		return errors.New("force rollback")
	})

	assert.Error(t, err)
	assert.Equal(t, "force rollback", err.Error())
	assert.True(t, executed)
}

func TestUnitOfWork_Begin(t *testing.T) {
	db := setupTestDB(t)
	uow := transaction.NewGormUnitOfWork(db)

	tx, err := uow.Begin(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	// 手动提交
	err = tx.Commit()
	assert.NoError(t, err)
}

func TestUnitOfWork_PanicRollback(t *testing.T) {
	db := setupTestDB(t)
	uow := transaction.NewGormUnitOfWork(db)

	assert.Panics(t, func() {
		_ = uow.WithTransaction(context.Background(), func(ctx context.Context) error {
			panic("force panic")
		})
	})
}
