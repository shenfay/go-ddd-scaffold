package response

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
)

// TraceIDKey TraceID 在 Context 中的键（与 middleware 包保持一致）
const TraceIDKey = "trace_id"

// FieldError 字段级错误信息
type FieldError struct {
	Field   string `json:"field"`                    // 字段名
	Message string `json:"message"`                  // 错误消息
	Code    string `json:"code,omitempty"`           // 错误代码 (可选)
	Value   string `json:"rejected_value,omitempty"` // 被拒绝的值 (可选)
}

// Response 统一响应结构
type Response struct {
	Code      int          `json:"code"`               // 业务错误码，0 表示成功
	Message   string       `json:"message"`            // 错误消息，成功时为"success"
	Data      interface{}  `json:"data,omitempty"`     // 响应数据
	Errors    []FieldError `json:"errors,omitempty"`   // 字段级错误列表 (验证失败时)
	Details   interface{}  `json:"details,omitempty"`  // 额外详细信息
	TraceID   string       `json:"trace_id,omitempty"` // 请求追踪 ID
	Timestamp int64        `json:"timestamp"`          // 时间戳 (秒)
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
		Code:      common.CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewErrorResponse 创建错误响应 (已废弃，请使用 NewError)
// Deprecated: Use NewError instead
func NewErrorResponse(code int, message string, details interface{}) *Response {
	return &Response{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().Unix(),
	}
}

// NewError 创建错误响应 (推荐)
func NewError(code int, message string) *Response {
	return &Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
}

// NewValidationError 创建验证错误响应 (字段级错误)
func NewValidationError(fieldErrors []FieldError) *Response {
	return &Response{
		Code:      common.CodeInvalidParam,
		Message:   "validation failed",
		Errors:    fieldErrors,
		Timestamp: time.Now().Unix(),
	}
}

// WithTraceID 添加追踪 ID
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

// WithErrors 添加字段错误
func (r *Response) WithErrors(errors []FieldError) *Response {
	r.Errors = errors
	return r
}

// NewPageResponse 创建分页响应
func NewPageResponse(items interface{}, total int64, page, pageSize int) *Response {
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	return &Response{
		Code:    common.CodeSuccess,
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
