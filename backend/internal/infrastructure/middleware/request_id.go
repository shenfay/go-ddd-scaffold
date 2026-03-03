package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID 生成请求ID中间件（避免与logger.go冲突）
// 负责生成/读取request-id，存入context，设置响应头
func RequestID(headerKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取或生成request-id
		requestID := c.GetHeader(headerKey)
		if requestID == "" {
			requestID = generateUUID()
		}

		// 2. 存入Gin Context（供handler使用）
		c.Set("requestId", requestID)

		// 3. 存入标准context（供errors/response包使用）
		ctx := context.WithValue(c.Request.Context(), "requestID", requestID)
		c.Request = c.Request.WithContext(ctx)

		// 4. 设置响应头
		c.Header(headerKey, requestID)

		c.Next()
	}
}

// generateUUID 生成UUID v4格式的request-id
func generateUUID() string {
	return uuid.New().String()
}

// GetReqIDFromGin 从Gin Context获取request-id（供其他中间件使用）
func GetReqIDFromGin(c *gin.Context) string {
	if id, exists := c.Get("requestId"); exists {
		return id.(string)
	}
	return ""
}

// GetRequestIDFromStdContext 从标准context获取request-id
func GetRequestIDFromStdContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, ok := ctx.Value("requestID").(string)
	if !ok || id == "" {
		return "unknown"
	}
	return id
}
