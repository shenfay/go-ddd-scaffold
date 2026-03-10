package response

import (
	"context"
	"time"

	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/i18n"
)

// BadRequest 返回 400 请求参数错误响应
func BadRequest(ctx context.Context, details any) *Response {
	return &Response{
		Code:    "BadRequest",
		Message: i18n.GetMessageByContext(ctx, "bad_request"),
		Error: &ErrorInfo{
			Code:    "BadRequest",
			Message: i18n.GetMessageByContext(ctx, "bad_request"),
			Details: details,
		},
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Conflict 返回 409 冲突错误响应
func Conflict(ctx context.Context, err *errors.AppError) *Response {
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

// NoContent 返回 204 无内容响应
func NoContent(ctx context.Context) *Response {
	return &Response{
		Code:      "NoContent",
		Message:  i18n.GetMessageByContext(ctx, "no_content"),
		Data:      nil,
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Created 返回 201 创建成功响应
func Created(ctx context.Context, data any) *Response {
	return &Response{
		Code:      "Created",
		Message:  i18n.GetMessageByContext(ctx, "created"),
		Data:      data,
		RequestID: errors.GetRequestID(ctx),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
