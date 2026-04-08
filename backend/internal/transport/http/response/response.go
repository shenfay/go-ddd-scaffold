package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware"
)

// Success 成功响应 (200 OK)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, middleware.SuccessResponse{
		BaseResponse: newBaseResponse(c),
		Code:         "SUCCESS",
		Message:      "Request successful",
		Data:         data,
	})
}

// Created 创建成功响应 (201 Created)
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, middleware.SuccessResponse{
		BaseResponse: newBaseResponse(c),
		Code:         "CREATED",
		Message:      "Resource created successfully",
		Data:         data,
	})
}

// NoContent 无内容响应 (204 No Content)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 错误响应（交给错误处理中间件）
func Error(c *gin.Context, err error) {
	c.Error(err)
}

// newBaseResponse 创建基础响应
func newBaseResponse(c *gin.Context) middleware.BaseResponse {
	return middleware.BaseResponse{
		TraceID:   middleware.GetTraceID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
