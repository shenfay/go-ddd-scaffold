package persistence

import (
	"context"
	"database/sql"
)

// DB 数据库操作接口（适配 database/sql）
type DB interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

// DBWrapper 包装 sql.DB 实现 DB 接口
type DBWrapper struct {
	*sql.DB
}

func NewDBWrapper(db *sql.DB) *DBWrapper {
	return &DBWrapper{db}
}

func (w *DBWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return w.DB.QueryRowContext(ctx, query, args...)
}

func (w *DBWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return w.DB.QueryContext(ctx, query, args...)
}

func (w *DBWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.DB.ExecContext(ctx, query, args...)
}

func (w *DBWrapper) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return w.DB.BeginTx(ctx, nil)
}
