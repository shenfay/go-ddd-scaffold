package kernel

import (
	"net/http"
)

// ============================================
// ErrorMapper 错误映射器
// ============================================

// ErrorMapper HTTP 错误映射器
// 将领域错误映射为 HTTP 状态码
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
