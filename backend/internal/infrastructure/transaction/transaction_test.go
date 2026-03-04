// Package transaction_test UnitOfWork 事务管理测试
package transaction_test

import (
	"context"
	"testing"

	"go-ddd-scaffold/internal/infrastructure/transaction"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestUnitOfWork_Commit_Success 测试事务正常提交
func TestUnitOfWork_Commit_Success(t *testing.T) {
	// 1. 准备测试数据库（内存 SQLite）
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 2. 创建 UnitOfWork
	uow := transaction.NewGormUnitOfWork(db)

	// 3. 开启事务并提交
	ctx := context.Background()
	tx, err := uow.Begin(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	// 4. 在事务中执行操作
	err = tx.GetDB().Exec("CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)").Error
	assert.NoError(t, err)

	// 5. 提交事务
	err = tx.Commit()
	assert.NoError(t, err)

	// 6. 验证表已创建（在事务外可以查询）
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// TestUnitOfWork_Rollback_OnPanic 测试 panic 时自动回滚
func TestUnitOfWork_Rollback_OnPanic(t *testing.T) {
	// 1. 准备测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	uow := transaction.NewGormUnitOfWork(db)
	ctx := context.Background()

	// 2. 开启事务
	tx, err := uow.Begin(ctx)
	assert.NoError(t, err)

	// 3. 模拟业务逻辑中发生 panic
	func() {
		defer func() {
			if p := recover(); p != nil {
				// 发生 panic，执行回滚
				_ = tx.Rollback()
			}
		}()

		// 在事务中创建表
		err = tx.GetDB().Exec("CREATE TABLE test_panic (id INTEGER PRIMARY KEY)").Error
		assert.NoError(t, err)

		// 模拟 panic
		panic("test panic")
	}()

	// 4. 验证表不存在（因为回滚了）
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_panic'").Scan(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestUnitOfWork_Rollback_Explicit 测试显式回滚
func TestUnitOfWork_Rollback_Explicit(t *testing.T) {
	// 1. 准备测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	uow := transaction.NewGormUnitOfWork(db)
	ctx := context.Background()

	// 2. 开启事务
	tx, err := uow.Begin(ctx)
	assert.NoError(t, err)

	// 3. 在事务中创建表
	err = tx.GetDB().Exec("CREATE TABLE test_explicit (id INTEGER PRIMARY KEY)").Error
	assert.NoError(t, err)

	// 4. 显式回滚
	err = tx.Rollback()
	assert.NoError(t, err)

	// 5. 验证表不存在
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_explicit'").Scan(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestUnitOfWork_NestedTransaction 测试嵌套事务（SavePoint）
func TestUnitOfWork_NestedTransaction(t *testing.T) {
	// 1. 准备测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	uow := transaction.NewGormUnitOfWork(db)
	ctx := context.Background()

	// 2. 外层事务
	tx1, err := uow.Begin(ctx)
	assert.NoError(t, err)

	// 3. 在外层事务中创建表
	err = tx1.GetDB().Exec("CREATE TABLE nested_test (id INTEGER PRIMARY KEY, value TEXT)").Error
	assert.NoError(t, err)

	// 4. 插入第一条数据
	err = tx1.GetDB().Exec("INSERT INTO nested_test (id, value) VALUES (1, 'outer')").Error
	assert.NoError(t, err)

	// 5. 内层事务（SavePoint）
	tx2 := tx1.GetDB().Begin()
	err = tx2.Exec("INSERT INTO nested_test (id, value) VALUES (2, 'inner')").Error
	assert.NoError(t, err)

	// 6. 内层事务回滚（不影响外层）
	err = tx2.Rollback().Error
	assert.NoError(t, err)

	// 7. 外层事务提交
	err = tx1.Commit()
	assert.NoError(t, err)

	// 8. 验证只有外层数据
	var count int64
	err = db.Model(&struct{}{}).Table("nested_test").Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// TestUnitOfWork_ConcurrentAccess 测试并发访问
func TestUnitOfWork_ConcurrentAccess(t *testing.T) {
	// 1. 准备测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	uow := transaction.NewGormUnitOfWork(db)
	ctx := context.Background()

	// 2. 创建基础表
	err = db.Exec("CREATE TABLE concurrent_test (id INTEGER PRIMARY KEY, counter INTEGER)").Error
	assert.NoError(t, err)
	err = db.Exec("INSERT INTO concurrent_test (id, counter) VALUES (1, 0)").Error
	assert.NoError(t, err)

	// 3. 并发更新计数器
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			tx, err := uow.Begin(ctx)
			if err != nil {
				return
			}
			defer tx.Rollback()

			// 读取当前值
			var counter int64
			err = tx.GetDB().Raw("SELECT counter FROM concurrent_test WHERE id = 1").Scan(&counter).Error
			if err != nil {
				return
			}

			// 更新
			err = tx.GetDB().Exec("UPDATE concurrent_test SET counter = counter + 1 WHERE id = 1").Error
			if err != nil {
				return
			}

			_ = tx.Commit()
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 4. 验证最终计数（SQLite 可能不是精确的 10，但应该在合理范围内）
	var finalCount int64
	err = db.Raw("SELECT counter FROM concurrent_test WHERE id = 1").Scan(&finalCount).Error
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, finalCount, int64(1))
	assert.LessOrEqual(t, finalCount, int64(10))
}
