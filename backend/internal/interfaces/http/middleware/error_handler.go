package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	apperrors "github.com/shenfay/go-ddd-scaffold/shared/errors"
	"github.com/shenfay/go-ddd-scaffold/shared/response"
)

// ErrorHandler 错误处理中间件
func ErrorHandler(mapper *apperrors.ErrorMapper, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) == 0 {
			return
		}

		// 获取最后一个错误
		err := c.Errors.Last().Err
		traceID := c.GetString("trace_id")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 映射错误
		httpStatus, code, message, details := mapper.Map(err)

		// 记录错误日志
		if httpStatus >= 500 {
			logger.Error("server error",
				zap.Error(err),
				zap.String("trace_id", traceID),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", httpStatus),
			)
		} else {
			logger.Warn("client error",
				zap.Error(err),
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
			Timestamp: time.Now().Unix(),
		})
	}
}

// TraceIDMiddleware 请求追踪ID中间件
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				traceID, _ := c.Get("trace_id")
				traceIDStr, _ := traceID.(string)
				if traceIDStr == "" {
					traceIDStr = uuid.New().String()
				}

				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("trace_id", traceIDStr),
					zap.String("path", c.Request.URL.Path),
				)

				c.JSON(500, response.ErrorResponse{
					Code:      apperrors.CodeInternalError,
					Message:   "服务器内部错误",
					TraceID:   traceIDStr,
					Timestamp: time.Now().Unix(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
