package errors

import (
	"context"
)

// ctxKey 上下文键类型
type ctxKey string

const reqIDKey ctxKey = "requestID"

// WithRequestID 设置请求ID到上下文
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, reqIDKey, id)
}

// GetRequestID 从上下文获取请求ID
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// 从标准context获取
	if id, ok := ctx.Value(reqIDKey).(string); ok && id != "" {
		return id
	}

	// 兼容从 gin context key 获取 (middleware 使用的 key)
	if id, ok := ctx.Value("requestID").(string); ok && id != "" {
		return id
	}

	return ""
}
