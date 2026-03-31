package http

import (
	"net/http"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
)

// ============================================
// ErrorMapper 错误映射器
// ============================================

// ErrorMapper HTTP 错误映射器
// 将领域错误映射为 HTTP 状态码
// 注意：这是基础设施层的实现，依赖于 domain 层的 BusinessError
type ErrorMapper struct {
	mappings map[int]int
}

// NewErrorMapper 创建错误映射器
func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{
		mappings: map[int]int{
			// 通用错误码
			common.CodeSuccess:       http.StatusOK,
			common.CodeInvalidParam:  http.StatusBadRequest,
			common.CodeNotFound:      http.StatusNotFound,
			common.CodeConflict:      http.StatusConflict,
			common.CodeUnauthorized:  http.StatusUnauthorized,
			common.CodeForbidden:     http.StatusForbidden,
			common.CodeInternalError: http.StatusInternalServerError,

			// 用户模块错误码 (20000-29999)
			20001: http.StatusNotFound,     // CodeUserNotFound
			20002: http.StatusConflict,     // CodeUserExists
			21001: http.StatusUnauthorized, // CodeInvalidPassword
			21004: http.StatusForbidden,    // CodeAccountLocked

			// 租户模块错误码 (30000-39999)
			30001: http.StatusNotFound, // CodeTenantNotFound

			// 认证授权错误码 (40000-49999)
			common.CodeTokenExpired:     http.StatusUnauthorized,
			common.CodeTokenInvalid:     http.StatusUnauthorized,
			common.CodePermissionDenied: http.StatusForbidden,
		},
	}
}

// Map 将业务错误映射为 HTTP 状态码和响应
func (m *ErrorMapper) Map(err error) (httpStatus int, code int, message string, details interface{}) {
	if err == nil {
		return http.StatusOK, common.CodeSuccess, "success", nil
	}

	// 检查是否是业务错误
	if be := common.AsBusinessError(err); be != nil {
		httpStatus, exists := m.mappings[be.Code]
		if !exists {
			httpStatus = http.StatusInternalServerError
		}
		return httpStatus, be.Code, be.Message, be.Details
	}

	// 默认返回 500
	return http.StatusInternalServerError, common.CodeInternalError, err.Error(), nil
}
