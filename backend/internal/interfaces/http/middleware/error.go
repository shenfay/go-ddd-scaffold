package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/pkg/response"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
)

// Error 错误处理中间件
// 负责捕获业务错误，统一映射并返回标准错误响应
// 依赖 TraceIDMiddleware 提供的 trace_id 进行错误追踪
func Error(mapper *common.ErrorMapper, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) == 0 {
			return
		}

		// 获取最后一个错误
		err := c.Errors.Last().Err

		// 从上下文中获取 TraceID
		traceID := GetTraceID(c)
		if traceID == "" {
			// 理论上不应该发生，因为 TraceIDMiddleware 应该已经设置了
			traceID = uuid.New().String()
		}

		// 映射错误
		httpStatus, code, message, details := mapper.Map(err)

		// 记录错误日志（添加更多上下文）
		if httpStatus >= 500 {
			logger.Error("server error",
				zap.String("error_type", fmt.Sprintf("%T", err)),
				zap.String("error_message", err.Error()),
				zap.Stack("stack"),
				zap.String("trace_id", traceID),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", httpStatus),
			)
		} else {
			logger.Warn("client error",
				zap.String("error_type", fmt.Sprintf("%T", err)),
				zap.String("error_message", err.Error()),
				zap.String("trace_id", traceID),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", httpStatus),
			)
		}

		// 返回统一错误响应
		c.JSON(httpStatus, response.ErrorResponse{
			Code:      code,
			Message:   message,
			Details:   details,
			TraceID:   traceID,
			Timestamp: util.Now().Timestamp(),
		})
	}
}
