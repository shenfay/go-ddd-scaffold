package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

const (
	// TraceIDKey Context 键
	TraceIDKey = "trace_id"
	// TraceIDHeader HTTP Header
	TraceIDHeader = "X-Trace-ID"
)

// TraceID 链路追踪 ID 中间件
// 为每个请求生成或传递 trace_id，用于日志关联和分布式追踪
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先使用客户端传入的 trace_id（用于跨服务链路追踪）
		traceID := c.GetHeader(TraceIDHeader)
		if traceID == "" {
			// 生成新的 trace_id（ULID 格式）
			traceID = utils.GenerateID()
		}

		// 存入 Context
		c.Set(TraceIDKey, traceID)

		// 响应头中也包含 trace_id
		c.Header(TraceIDHeader, traceID)

		c.Next()
	}
}

// GetTraceID 从 Context 获取 trace_id
func GetTraceID(c *gin.Context) string {
	if id, exists := c.Get(TraceIDKey); exists {
		return id.(string)
	}
	return ""
}
