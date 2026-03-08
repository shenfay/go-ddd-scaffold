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

// TestUnitOfWork_ErrorPropagation 测试错误传播和回滚
func TestUnitOfWork_ErrorPropagation(t *testing.T) {
	db := setupTestDB(t)
	uow := transaction.NewGormUnitOfWork(db)

	executionOrder := make([]string, 0)

	err := uow.WithTransaction(context.Background(), func(ctx context.Context) error {
		executionOrder = append(executionOrder, "step1")

		// 第一步成功
		if err := uow.WithTransaction(ctx, func(innerCtx context.Context) error {
			executionOrder = append(executionOrder, "step2")
			return nil
		}); err != nil {
			return err
		}

		executionOrder = append(executionOrder, "step3")

		// 第二步失败，应该回滚所有
		return errors.New("step3 failed")
	})

	assert.Error(t, err)
	assert.Equal(t, []string{"step1", "step2", "step3"}, executionOrder)
}

// TestUnitOfWork_MultipleOperations 测试多操作原子性
func TestUnitOfWork_MultipleOperations(t *testing.T) {
	db := setupTestDB(t)
	uow := transaction.NewGormUnitOfWork(db)

	operations := make([]int, 0)

	err := uow.WithTransaction(context.Background(), func(ctx context.Context) error {
		operations = append(operations, 1)
		operations = append(operations, 2)

		// 模拟第三个操作失败
		return errors.New("third operation failed")
	})

	assert.Error(t, err)
	assert.Equal(t, []int{1, 2}, operations)
	// 验证没有后续操作执行
	assert.Len(t, operations, 2)
}
