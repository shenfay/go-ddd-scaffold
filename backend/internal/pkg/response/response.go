package response

import (
	"context"
	"time"

	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/i18n"
)

// Response 统一响应结构
type Response struct {
	Code      string     `json:"code"`
	Message   string     `json:"message"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorInfo `json:"error,omitempty"`
	RequestID string     `json:"requestId"`
	Timestamp string     `json:"timestamp"`
}

// ErrorInfo 错误信息结构
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// OK 创建成功响应
func OK(ctx context.Context, data any) *Response {
	return &Response{
		Code:      "Success",
		Message:   i18n.GetMessageByContext(ctx, "success"),
		Data:      data,
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// OKWithMsg 创建成功响应（带自定义消息）
func OKWithMsg(ctx context.Context, data any, msg string) *Response {
	return &Response{
		Code:      "Success",
		Message:   msg,
		Data:      data,
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Fail 创建错误响应
func Fail(ctx context.Context, err *errors.AppError) *Response {
	resp := &Response{
		Code:      err.GetCode(),
		Message:   err.GetMessage(),
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err.GetDetails() != nil {
		resp.Error = &ErrorInfo{
			Code:    err.GetCode(),
			Message: err.GetMessage(),
			Details: err.GetDetails(),
		}
	}

	return resp
}

// ValidateErr 创建参数校验错误响应
func ValidateErr(ctx context.Context, details any) *Response {
	return &Response{
		Code:    "ValidationFailed",
		Message: i18n.GetMessageByContext(ctx, "validation_failed"),
		Error: &ErrorInfo{
			Code:    "ValidationFailed",
			Message: i18n.GetMessageByContext(ctx, "validation_failed"),
			Details: details,
		},
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NotFound 创建资源不存在错误响应
func NotFound(ctx context.Context, resource string) *Response {
	return &Response{
		Code:    "NotFound",
		Message: i18n.GetMessageByContext(ctx, "not_found"),
		Error: &ErrorInfo{
			Code:    "NotFound",
			Message: i18n.GetMessageByContext(ctx, "not_found"),
			Details: resource,
		},
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// ServerErr 创建服务器内部错误响应
func ServerErr(ctx context.Context) *Response {
	return &Response{
		Code:      "System.InternalError",
		Message:   i18n.GetMessageByContext(ctx, "system_error"),
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// BusinessErr 创建业务错误响应
func BusinessErr(ctx context.Context, code, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Page 创建分页数据响应
func Page(ctx context.Context, items any, total int64, page, pageSize int) *Response {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	hasNext := page < totalPages
	hasPrev := page > 1

	return &Response{
		Code:    "Success",
		Message: i18n.GetMessageByContext(ctx, "success"),
		Data: PageData{
			Items:      items,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Custom 创建自定义响应
func Custom(ctx context.Context, code, msg string, data any) *Response {
	return &Response{
		Code:      code,
		Message:   msg,
		Data:      data,
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Unauthorized 返回 401 未授权响应
func Unauthorized(ctx context.Context, message string) *Response {
	return &Response{
		Code:    "Unauthorized",
		Message: message,
		Error: &ErrorInfo{
			Code:    "Unauthorized",
			Message: message,
		},
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Forbidden 返回 403 禁止访问响应
func Forbidden(ctx context.Context, message string) *Response {
	return &Response{
		Code:    "Forbidden",
		Message: message,
		Error: &ErrorInfo{
			Code:    "Forbidden",
			Message: message,
		},
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
