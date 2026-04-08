package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorHandling 统一错误处理中间件
// 自动处理通过 c.Error() 设置的错误
func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理 c.Error() 设置的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleAppError(c, err)
			c.Abort()
		}
	}
}

// handleAppError 处理应用错误
func handleAppError(c *gin.Context, err error) {
	// 尝试转换为 AppError
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.HTTPStatus, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Metadata,
		})
		return
	}

	// 未知错误，返回 500
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "SYSTEM.INTERNAL_ERROR",
		Message: "Internal server error",
	})
}
