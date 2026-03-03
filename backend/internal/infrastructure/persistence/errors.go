package persistence

import "errors"

var (
	// ErrNotFound 数据未找到错误
	ErrNotFound = errors.New("record not found")
)
