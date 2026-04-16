package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware"
)

// Success 返回成功响应（200 OK）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, middleware.SuccessResponse{
		BaseResponse: newBaseResponse(c),
		Code:         "SUCCESS",
		Message:      "Request successful",
		Data:         data,
	})
}

// Created 返回创建成功响应（201 Created）
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, middleware.SuccessResponse{
		BaseResponse: newBaseResponse(c),
		Code:         "CREATED",
		Message:      "Resource created successfully",
		Data:         data,
	})
}

// NoContent 返回无内容响应（204 No Content）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 返回错误响应（交由错误处理中间件处理）
func Error(c *gin.Context, err error) {
	_ = c.Error(err) // 错误处理中间件会处理这个错误
}

// newBaseResponse 创建基础响应结构
func newBaseResponse(c *gin.Context) middleware.BaseResponse {
	return middleware.BaseResponse{
		TraceID:   middleware.GetTraceID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
