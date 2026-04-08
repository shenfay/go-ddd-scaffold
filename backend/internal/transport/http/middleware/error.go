package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// BaseResponse 响应基础结构（所有响应共用）
type BaseResponse struct {
	TraceID   string `json:"trace_id"`  // 链路追踪 ID（用于日志关联和分布式追踪）
	Timestamp string `json:"timestamp"` // 响应时间（RFC3339 格式）
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	BaseResponse             // 嵌套，JSON 自动扁平化
	Code         string      `json:"code"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	BaseResponse             // 嵌套
	Code         string      `json:"code"`
	Message      string      `json:"message"`
	Details      interface{} `json:"details,omitempty"`
}

// ErrorHandling 统一错误处理中间件
// 自动处理通过 c.Error() 设置的错误，并注入 trace_id 和 timestamp
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

// handleAppError 处理应用错误（自动注入 trace_id 和 timestamp）
func handleAppError(c *gin.Context, err error) {
	baseResponse := BaseResponse{
		TraceID:   GetTraceID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// 尝试转换为 AppError
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.HTTPStatus, ErrorResponse{
			BaseResponse: baseResponse,
			Code:         appErr.Code,
			Message:      appErr.Message,
			Details:      appErr.Metadata,
		})
		return
	}

	// 未知错误，返回 500
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		BaseResponse: baseResponse,
		Code:         "SYSTEM.INTERNAL_ERROR",
		Message:      "Internal server error",
	})
}
