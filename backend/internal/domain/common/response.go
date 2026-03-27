package common

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
			CodeSuccess:       http.StatusOK,
			CodeInvalidParam:  http.StatusBadRequest,
			CodeNotFound:      http.StatusNotFound,
			CodeConflict:      http.StatusConflict,
			CodeUnauthorized:  http.StatusUnauthorized,
			CodeForbidden:     http.StatusForbidden,
			CodeInternalError: http.StatusInternalServerError,
			// 用户模块错误码 (20000-29999)
			20001: http.StatusNotFound,     // CodeUserNotFound
			20002: http.StatusConflict,     // CodeUserExists
			21001: http.StatusUnauthorized, // CodeInvalidPassword
			21004: http.StatusForbidden,    // CodeAccountLocked
			// 租户模块错误码 (30000-39999)
			30001: http.StatusNotFound, // CodeTenantNotFound
			// 认证授权错误码 (40000-49999)
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
