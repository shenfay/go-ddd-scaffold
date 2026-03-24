package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL 驱动
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/support/config"
	"go.uber.org/zap"
)

// Database 数据库连接管理器
type Database struct {
	*sql.DB
	config *config.DatabaseConfig
	logger *zap.Logger
}

// NewDatabase 创建新的数据库连接
func NewDatabase(cfg *config.DatabaseConfig, logger *zap.Logger) (*Database, error) {
	dsn := cfg.GetDSN()

	// 打开数据库连接
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 配置连接池
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
		zap.Int("max_open_conns", cfg.MaxOpenConns))

	return &Database{
		DB:     db,
		config: cfg,
		logger: logger,
	}, nil
}

// Close 关闭数据库连接
func (db *Database) Close() error {
	if err := db.DB.Close(); err != nil {
		db.logger.Error("failed to close database connection", zap.Error(err))
		return err
	}

	db.logger.Info("database connection closed")
	return nil
}

// Stats 返回数据库连接池统计信息
func (db *Database) Stats() sql.DBStats {
	return db.DB.Stats()
}
