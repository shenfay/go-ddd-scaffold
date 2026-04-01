package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/pkg/response"
)

// Recovery panic 恢复中间件
// 捕获 panic 并返回友好的错误响应
// 依赖 TraceIDMiddleware 提供的 trace_id 进行错误追踪
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 从上下文中获取 TraceID
				traceID := GetTraceID(c)
				if traceID == "" {
					traceID = uuid.New().String()
				}

				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("trace_id", traceID),
					zap.String("path", c.Request.URL.Path),
				)

				resp := response.NewError(common.CodeInternalError, "服务器内部错误")
				resp.WithTraceID(traceID)
				c.JSON(500, resp)
				c.Abort()
			}
		}()
		c.Next()
	}
}
