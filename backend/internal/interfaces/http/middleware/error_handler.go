package middleware

import (
	"net/http"

	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误处理中间件
//
// 职责:
// 1. 捕获 Handler 中通过 c.Error() 记录的错误
// 2. 将 AppError 转换为统一的 HTTP 响应
// 3. 记录错误日志
// 4. 避免在每个 Handler 中重复错误处理逻辑
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果没有错误，直接返回
		if len(c.Errors) == 0 {
			return
		}

		// 处理所有错误（通常只有一个）
		for _, err := range c.Errors {
			handleError(c, err)
		}
	}
}

// handleError 处理单个错误
func handleError(c *gin.Context, ginErr *gin.Error) {
	ctx := c.Request.Context()

	// 尝试将错误转换为 AppError
	if appErr, ok := ginErr.Err.(*errors.AppError); ok {
		// 根据错误类型返回相应的 HTTP 状态码
		statusCode, _ := errors.GetHTTPStatus(appErr)

		// 记录错误日志（如果是服务器内部错误）
		if statusCode >= http.StatusInternalServerError {
			// 这里可以集成日志中间件
			// logger.Error("Internal error", zap.Error(appErr))
		}

		// 返回统一的错误响应
		c.JSON(statusCode, response.Fail(ctx, appErr))
		return
	}

	// 非 AppError，视为内部错误
	c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
}

// HandlePanic 恢复 panic 并转换为错误
func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rv := recover(); rv != nil {
				// 创建 AppError
				appErr := errors.InternalError.WithDetails(rv)
				c.JSON(http.StatusInternalServerError, response.Fail(c.Request.Context(), appErr))
				c.Abort()
			}
		}()
		c.Next()
	}
}
