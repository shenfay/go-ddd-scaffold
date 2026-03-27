package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/pkg/response"
)

// Handler HTTP 响应处理器（端口层）
// 提供统一的 HTTP 响应处理方法，符合 DDD Ports & Adapters 模式
type Handler struct {
	errorMapper *common.ErrorMapper
}

// NewHandler 创建响应处理器
func NewHandler(mapper *common.ErrorMapper) *Handler {
	return &Handler{errorMapper: mapper}
}

// NewResponseWithTraceID 创建带 TraceID 的成功响应
// 自动从 Gin Context 中获取 TraceID
func NewResponseWithTraceID(c *gin.Context, data interface{}) *response.Response {
	resp := response.NewResponse(data)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	return resp
}

// NewErrorResponseWithTraceID 创建带 TraceID 的错误响应
// 自动从 Gin Context 中获取 TraceID
func NewErrorResponseWithTraceID(c *gin.Context, code int, message string, details interface{}) *response.ErrorResponse {
	resp := response.NewErrorResponse(code, message, details)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	return resp
}

// Success 成功响应（自动注入 TraceID）
func (h *Handler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewResponseWithTraceID(c, data))
}

// Created 创建成功响应（201，自动注入 TraceID）
func (h *Handler) Created(c *gin.Context, data interface{}) {
	resp := NewResponseWithTraceID(c, data)
	resp.Code = 0
	resp.Message = "created"
	c.JSON(http.StatusCreated, resp)
}

// Accepted 接受响应（202，自动注入 TraceID）
func (h *Handler) Accepted(c *gin.Context, data interface{}) {
	resp := NewResponseWithTraceID(c, data)
	resp.Code = 0
	resp.Message = "accepted"
	c.JSON(http.StatusAccepted, resp)
}

// NoContent 无内容响应（204）
func (h *Handler) NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 错误响应（自动注入 TraceID）
func (h *Handler) Error(c *gin.Context, err error) {
	httpStatus, code, message, details := h.errorMapper.Map(err)
	c.JSON(httpStatus, NewErrorResponseWithTraceID(c, code, message, details))
}

// Page 分页响应（自动注入 TraceID）
func (h *Handler) Page(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	resp := response.NewPageResponse(items, total, page, pageSize)
	if traceID := middleware.GetTraceID(c); traceID != "" {
		resp.WithTraceID(traceID)
	}
	c.JSON(http.StatusOK, resp)
}

// BadRequest 400 错误
func (h *Handler) BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, response.NewErrorResponse(
		common.CodeBadRequest,
		message,
		nil,
	))
}

// Unauthorized 401 错误
func (h *Handler) Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.NewErrorResponse(
		common.CodeUnauthorized,
		message,
		nil,
	))
}

// Forbidden 403 错误
func (h *Handler) Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, response.NewErrorResponse(
		common.CodeForbidden,
		message,
		nil,
	))
}

// NotFound 404 错误
func (h *Handler) NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, response.NewErrorResponse(
		common.CodeNotFound,
		message,
		nil,
	))
}

// Conflict 409 错误
func (h *Handler) Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, response.NewErrorResponse(
		common.CodeConflict,
		message,
		nil,
	))
}

// InternalServerError 500 错误
func (h *Handler) InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
		common.CodeInternalError,
		message,
		nil,
	))
}

// BindJSON 绑定 JSON 并处理错误
func (h *Handler) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.Error(c, err)
		return false
	}
	return true
}

// BindQuery 绑定 Query 参数并处理错误
func (h *Handler) BindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		h.Error(c, err)
		return false
	}
	return true
}

// BindUri 绑定 URI 参数并处理错误
func (h *Handler) BindUri(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		h.Error(c, err)
		return false
	}
	return true
}
