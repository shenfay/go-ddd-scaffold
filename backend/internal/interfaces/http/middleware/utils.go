package middleware

import "github.com/gin-gonic/gin"

// GetTraceIDFromContext 从 Gin Context 中获取 TraceID
// 可在 handler 等外部包中使用
func GetTraceIDFromContext(c *gin.Context) string {
	if c == nil {
		return ""
	}

	// 尝试从上下文中获取 trace_id
	if val, exists := c.Get(TraceIDKey); exists {
		if traceID, ok := val.(string); ok {
			return traceID
		}
	}
	return ""
}
