package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	testcontainerspkg "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDatabase 测试数据库容器
type TestDatabase struct {
	Container *postgres.PostgresContainer
	DB        *gorm.DB
}

// TestRedis 测试 Redis 容器
type TestRedis struct {
	Container *redis.RedisContainer
}

// SetupTestDatabase 启动测试数据库容器
func SetupTestDatabase(t *testing.T) *TestDatabase {
	t.Helper()

	ctx := context.Background()

	dbName := "testdb"
	dbUser := "testuser"
	dbPassword := "testpass"

	dbContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainerspkg.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30000)),
	)
	require.NoError(t, err, "Failed to start postgres container")

	// 清理容器
	t.Cleanup(func() {
		if err := dbContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate postgres container: %v", err)
		}
	})

	connStr, err := dbContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get connection string")

	// 连接数据库
	db, err := gorm.Open(pgdriver.Open(connStr), &gorm.Config{
		Logger: logger.Discard, // 禁用测试日志
	})
	require.NoError(t, err, "Failed to connect to database")

	// 自动迁移
	err = db.AutoMigrate()
	require.NoError(t, err, "Failed to auto migrate")

	return &TestDatabase{
		Container: dbContainer,
		DB:        db,
	}
}

// SetupTestRedis 启动测试 Redis 容器
func SetupTestRedis(t *testing.T) *TestRedis {
	t.Helper()

	ctx := context.Background()

	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelWarning),
	)
	require.NoError(t, err, "Failed to start redis container")

	// 清理容器
	t.Cleanup(func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate redis container: %v", err)
		}
	})

	return &TestRedis{
		Container: redisContainer,
	}
}

// GetRedisAddr 获取 Redis 连接地址
func (r *TestRedis) GetRedisAddr(t *testing.T) string {
	ctx := context.Background()
	addr, err := r.Container.Endpoint(ctx, "")
	require.NoError(t, err, "Failed to get redis endpoint")
	return addr
}
