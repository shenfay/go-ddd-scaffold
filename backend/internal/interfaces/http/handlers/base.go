// Package handlers provides HTTP response handlers for the API layer.
// It implements the Port pattern from DDD Ports & Adapters architecture,
// offering unified HTTP response handling methods with automatic TraceID injection.
//
// Usage:
//
//	handler := handlers.NewHandler(errorMapper)
//	handler.Success(c, data)          // 200 OK
//	handler.Created(c, data)          // 201 Created
//	handler.Error(c, err)             // Error response with trace ID
//	handler.Page(c, items, total, ...) // Paginated response
//
// The package separates concerns:
//   - pkg/response: Data structures (Response, FieldError, PageData)
//   - handlers/base.go: HTTP-specific handler logic (this package)
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	httpinfra "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/http"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/pkg/response"
)

// Handler HTTP 响应处理器（端口层）
// 提供统一的 HTTP 响应处理方法，符合 DDD Ports & Adapters 模式
type Handler struct {
	errorMapper *httpinfra.ErrorMapper
}

// NewHandler 创建响应处理器
func NewHandler(mapper *httpinfra.ErrorMapper) *Handler {
	return &Handler{errorMapper: mapper}
}

// Success 成功响应 (自动注入 TraceID)
func (h *Handler) Success(c *gin.Context, data interface{}) {
	resp := response.NewResponse(data)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusOK, resp)
}

// Created 创建成功响应 (201，自动注入 TraceID)
func (h *Handler) Created(c *gin.Context, data interface{}) {
	resp := response.NewResponse(data)
	resp.Message = "created"
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusCreated, resp)
}

// Accepted 接受响应 (202，自动注入 TraceID)
func (h *Handler) Accepted(c *gin.Context, data interface{}) {
	resp := response.NewResponse(data)
	resp.Message = "accepted"
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusAccepted, resp)
}

// NoContent 无内容响应（204）
func (h *Handler) NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 错误响应 (自动注入 TraceID)
// 自动映射业务错误为对应的 HTTP 状态码和响应格式
func (h *Handler) Error(c *gin.Context, err error) {
	httpStatus, code, message, details := h.errorMapper.Map(err)
	resp := response.NewError(code, message)
	if details != nil {
		resp.WithDetails(details)
	}
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(httpStatus, resp)
}

// Page 分页响应（自动注入 TraceID）
func (h *Handler) Page(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	resp := response.NewPageResponse(items, total, page, pageSize)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusOK, resp)
}

// ValidationError 验证错误响应 (字段级错误，400 Bad Request)
func (h *Handler) ValidationError(c *gin.Context, fieldErrors []response.FieldError) {
	resp := response.NewValidationError(fieldErrors)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusBadRequest, resp)
}

// Unauthorized 401 错误
func (h *Handler) Unauthorized(c *gin.Context, message string) {
	resp := response.NewError(common.CodeUnauthorized, message)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusUnauthorized, resp)
}

// Forbidden 403 错误
func (h *Handler) Forbidden(c *gin.Context, message string) {
	resp := response.NewError(common.CodeForbidden, message)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusForbidden, resp)
}

// NotFound 404 错误
func (h *Handler) NotFound(c *gin.Context, message string) {
	resp := response.NewError(common.CodeNotFound, message)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusNotFound, resp)
}

// Conflict 409 错误
func (h *Handler) Conflict(c *gin.Context, message string) {
	resp := response.NewError(common.CodeConflict, message)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusConflict, resp)
}

// InternalError 500 错误
func (h *Handler) InternalError(c *gin.Context, message string) {
	resp := response.NewError(common.CodeInternalError, message)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusInternalServerError, resp)
}

// BindJSON 绑定 JSON 并处理验证错误
func (h *Handler) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		// 提取字段级错误
		fieldErrors := extractFieldErrors(err)
		if len(fieldErrors) > 0 {
			h.ValidationError(c, fieldErrors)
		} else {
			h.Error(c, err)
		}
		return false
	}
	return true
}

// BindQuery 绑定 Query 参数并处理验证错误
func (h *Handler) BindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		fieldErrors := extractFieldErrors(err)
		if len(fieldErrors) > 0 {
			h.ValidationError(c, fieldErrors)
		} else {
			h.Error(c, err)
		}
		return false
	}
	return true
}

// BindUri 绑定 URI 参数并处理验证错误
func (h *Handler) BindUri(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		fieldErrors := extractFieldErrors(err)
		if len(fieldErrors) > 0 {
			h.ValidationError(c, fieldErrors)
		} else {
			h.Error(c, err)
		}
		return false
	}
	return true
}

// extractFieldErrors 从 gin 的错误中提取字段级错误
// 支持 validator.ValidationErrors 和 gin.Error
func extractFieldErrors(err error) []response.FieldError {
	var fieldErrors []response.FieldError

	// 处理 validator.ValidationErrors (主要验证库)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			fieldErrors = append(fieldErrors, response.FieldError{
				Field:   formatFieldName(e.Field()),
				Message: getValidationMessage(e),
				Code:    e.Tag(),
				Value:   "",
			})
		}
		return fieldErrors
	}

	// 处理 gin.Error
	if ginErr, ok := err.(*gin.Error); ok {
		fieldErrors = append(fieldErrors, response.FieldError{
			Field:   "",
			Message: ginErr.Error(),
		})
		return fieldErrors
	}

	// 其他错误不提取字段信息
	return fieldErrors
}

// formatFieldName 格式化字段名 (将 PascalCase 转为 snake_case)
func formatFieldName(fieldName string) string {
	var result strings.Builder
	for i, r := range fieldName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// getValidationMessage 根据验证标签生成友好的错误消息
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " 是必填项"
	case "email":
		return e.Field() + " 必须是有效的邮箱地址"
	case "url":
		return e.Field() + " 必须是有效的 URL"
	case "min":
		return e.Field() + " 的最小值是 " + e.Param()
	case "max":
		return e.Field() + " 的最大值是 " + e.Param()
	case "len":
		return e.Field() + " 的长度必须是 " + e.Param()
	case "oneof":
		return e.Field() + " 必须是以下值之一：" + e.Param()
	case "mobile":
		return e.Field() + " 必须是有效的手机号"
	case "contains":
		return e.Field() + " 必须包含 " + e.Param()
	default:
		return e.Field() + " 验证失败 (" + e.Tag() + ")"
	}
}
