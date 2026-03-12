package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// TraceIDKey TraceID 在 Context 中的键
	TraceIDKey = "trace_id"

	// TraceIDHeader HTTP Header 中 TraceID 的键
	TraceIDHeader = "X-Trace-ID"
)

// TraceIDMiddleware 请求追踪 ID 中间件
// 为每个请求生成或获取唯一的追踪 ID，并在请求和响应中传递
//
// 功能说明:
// 1. 优先从请求 Header 中获取 X-Trace-ID(支持客户端传递)
// 2. 如果不存在，则生成新的 UUID 作为 TraceID
// 3. 将 TraceID 设置到 Gin Context 中，供后续中间件和 Handler 使用
// 4. 在响应 Header 中添加 X-Trace-ID，便于客户端追踪
//
// 使用场景:
// - 请求链路追踪
// - 日志关联
// - 问题排查
// - 分布式系统调用链
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 尝试从请求 Header 获取 TraceID
		traceID := c.GetHeader(TraceIDHeader)

		// 2. 如果不存在，生成新的 TraceID
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 3. 设置到 Context 中
		c.Set(TraceIDKey, traceID)

		// 4. 在响应 Header 中添加 TraceID
		c.Header(TraceIDHeader, traceID)

		// 继续处理请求
		c.Next()
	}
}

// GetTraceID 从上下文中获取当前请求的 TraceID
// 可在 Handler、Service 等任何地方使用此方法获取追踪 ID
//
// 使用示例:
//
//	func (h *Handler) CreateUser(c *gin.Context) {
//	    traceID := middleware.GetTraceID(c)
//	    h.logger.Info("Creating user", zap.String("trace_id", traceID))
//	    // ...
//	}
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get(TraceIDKey); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}

// SetTraceID 手动设置 TraceID 到上下文中
// 特殊场景下可能需要手动设置 TraceID，例如：
// - 从其他来源获取 TraceID(如消息队列)
// - 跨 goroutine 传递追踪上下文
func SetTraceID(c *gin.Context, traceID string) {
	c.Set(TraceIDKey, traceID)
}
