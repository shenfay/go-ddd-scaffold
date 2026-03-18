package kernel

import (
	"net/http"
)

// Response 标准响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	TraceID  string      `json:"trace_id,omitempty"`
}

// NewResponse 创建标准响应
func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, details interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    details,
	}
}

// NewPageResponse 创建分页响应
func NewPageResponse(data interface{}, total int64, page, pageSize int) *PageResponse {
	return &PageResponse{
		Code:     CodeSuccess,
		Message:  "success",
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

// WithTraceID 设置 TraceID
func (r *Response) WithTraceID(traceID string) *Response {
	r.TraceID = traceID
	return r
}

// WithTraceID 设置 TraceID（分页响应）
func (r *PageResponse) WithTraceID(traceID string) *PageResponse {
	r.TraceID = traceID
	return r
}

// ============================================
// ErrorMapper 错误映射器
// ============================================

// ErrorMapper HTTP错误映射器
type ErrorMapper struct {
	mappings map[int]int
}

// NewErrorMapper 创建错误映射器
func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{
		mappings: map[int]int{
			CodeSuccess:          http.StatusOK,
			CodeInvalidParam:     http.StatusBadRequest,
			CodeNotFound:         http.StatusNotFound,
			CodeConflict:         http.StatusConflict,
			CodeUnauthorized:     http.StatusUnauthorized,
			CodeForbidden:        http.StatusForbidden,
			CodeInternalError:    http.StatusInternalServerError,
			CodeUserNotFound:     http.StatusNotFound,
			CodeUserExists:       http.StatusConflict,
			CodeInvalidPassword:  http.StatusUnauthorized,
			CodeAccountLocked:    http.StatusForbidden,
			CodeTenantNotFound:   http.StatusNotFound,
			CodeTokenExpired:     http.StatusUnauthorized,
			CodeTokenInvalid:     http.StatusUnauthorized,
			CodePermissionDenied: http.StatusForbidden,
		},
	}
}

// Map 将业务错误映射为 HTTP 状态码和响应
func (m *ErrorMapper) Map(err error) (httpStatus int, code int, message string, details interface{}) {
	if err == nil {
		return http.StatusOK, CodeSuccess, "success", nil
	}

	// 检查是否是业务错误
	if be := AsBusinessError(err); be != nil {
		httpStatus, exists := m.mappings[be.Code]
		if !exists {
			httpStatus = http.StatusInternalServerError
		}
		return httpStatus, be.Code, be.Message, be.Details
	}

	// 默认返回 500
	return http.StatusInternalServerError, CodeInternalError, err.Error(), nil
}
