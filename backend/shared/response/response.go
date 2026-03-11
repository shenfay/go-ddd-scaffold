package response

import (
	"time"

	apperrors "github.com/shenfay/go-ddd-scaffold/shared/errors"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`               // 业务错误码，0表示成功
	Message   string      `json:"message"`            // 错误消息，成功时为"success"
	Data      interface{} `json:"data,omitempty"`     // 响应数据
	Details   interface{} `json:"details,omitempty"`  // 详细错误信息
	TraceID   string      `json:"trace_id,omitempty"` // 请求追踪ID
	Timestamp int64       `json:"timestamp"`          // 时间戳(秒)
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code      int         `json:"code"`               // 错误码
	Message   string      `json:"message"`            // 错误消息
	Details   interface{} `json:"details,omitempty"`  // 详细错误信息
	TraceID   string      `json:"trace_id,omitempty"` // 请求追踪ID
	Timestamp int64       `json:"timestamp"`          // 时间戳(秒)
}

// PageData 分页数据
type PageData struct {
	Items     interface{} `json:"items"`      // 数据列表
	Total     int64       `json:"total"`      // 总数
	Page      int         `json:"page"`       // 当前页码
	PageSize  int         `json:"page_size"`  // 每页数量
	TotalPage int         `json:"total_page"` // 总页数
}

// NewResponse 创建成功响应
func NewResponse(data interface{}) *Response {
	return &Response{
		Code:      apperrors.CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, details interface{}) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().Unix(),
	}
}

// WithTraceID 添加追踪ID
func (r *Response) WithTraceID(traceID string) *Response {
	r.TraceID = traceID
	return r
}

// WithMessage 自定义消息
func (r *Response) WithMessage(msg string) *Response {
	r.Message = msg
	return r
}

// WithDetails 添加详情
func (r *Response) WithDetails(details interface{}) *Response {
	r.Details = details
	return r
}

// WithTraceID 添加追踪ID
func (e *ErrorResponse) WithTraceID(traceID string) *ErrorResponse {
	e.TraceID = traceID
	return e
}

// NewPageResponse 创建分页响应
func NewPageResponse(items interface{}, total int64, page, pageSize int) *Response {
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	return &Response{
		Code:    apperrors.CodeSuccess,
		Message: "success",
		Data: PageData{
			Items:     items,
			Total:     total,
			Page:      page,
			PageSize:  pageSize,
			TotalPage: totalPage,
		},
		Timestamp: time.Now().Unix(),
	}
}
