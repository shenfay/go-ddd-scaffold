// Package middleware 统一错误处理中间件
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"
)

// ErrorMiddleware 统一错误处理中间件
// 负责将 AppError 转换为标准的 HTTP 响应
// 注意：panic 恢复由 RecoveryWithLogger 中间件处理
func ErrorMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			handleErrors(c, logger)
			return
		}
	}
}

// handleErrors 处理累积的错误
func handleErrors(c *gin.Context, logger *zap.Logger) {
	for _, err := range c.Errors {
		if appErr, ok := err.Err.(*errors.AppError); ok {
			// 获取 HTTP状态码
			httpStatus, _ := errors.GetHTTPStatus(appErr)

			// 记录错误日志
			if logger != nil {
				logger.Error("业务错误",
					zap.String("code", appErr.GetCode()),
					zap.String("message", appErr.GetMessage()),
					zap.String("category", appErr.GetCategory()),
					zap.String("client", c.ClientIP()),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
			}

			// 返回错误响应
			c.JSON(httpStatus, response.Fail(c.Request.Context(), appErr))
			return
		}

		// 非 AppError，作为内部错误处理
		if logger != nil {
			logger.Error("未知错误",
				zap.Any("error", err.Err),
				zap.String("client", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
		}

		c.JSON(http.StatusInternalServerError, response.Fail(c.Request.Context(), errors.InternalError))
	}
}
