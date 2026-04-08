package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 成功响应结构
type SuccessResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应 (200 OK)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Code:    "SUCCESS",
		Message: "Request successful",
		Data:    data,
	})
}

// Created 创建成功响应 (201 Created)
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    "CREATED",
		Message: "Resource created successfully",
		Data:    data,
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
