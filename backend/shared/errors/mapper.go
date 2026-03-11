package apperrors

import (
	"errors"
	"net/http"

	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// ============================================
// 错误映射器
// ============================================

// ErrorMapper 错误映射器
type ErrorMapper struct {
	customMappings map[error]CodeInfo
}

// NewErrorMapper 创建错误映射器
func NewErrorMapper() *ErrorMapper {
	m := &ErrorMapper{
		customMappings: make(map[error]CodeInfo),
	}
	m.registerDefaults()
	return m
}

// 注册默认映射
func (m *ErrorMapper) registerDefaults() {
	// DDD领域错误映射
	m.Register(ddd.ErrAggregateNotFound, CodeNotFound, "资源不存在")
	m.Register(ddd.ErrBusinessRuleViolation, CodeInvalidParam, "业务规则违反")
	m.Register(ddd.ErrValidationFailed, CodeInvalidParam, "参数验证失败")
	m.Register(ddd.ErrConcurrencyConflict, CodeConcurrency, "并发冲突")
	m.Register(ddd.ErrInvalidOperation, CodeInvalidParam, "无效操作")
}

// Register 注册错误映射
func (m *ErrorMapper) Register(err error, code int, message string) {
	m.customMappings[err] = CodeInfo{
		Code:       code,
		HTTPStatus: GetHTTPStatus(code),
		Message:    message,
	}
}

// Map 映射错误到响应信息
func (m *ErrorMapper) Map(err error) (int, int, string, interface{}) {
	if err == nil {
		return http.StatusOK, CodeSuccess, "success", nil
	}

	// 检查具体错误类型
	switch e := err.(type) {
	case *BusinessError:
		return m.mapBusinessError(e)
	case *ValidationErrors:
		return http.StatusBadRequest, CodeInvalidParam, "参数验证失败", e.Errors
	case *ConcurrencyError:
		return http.StatusConflict, CodeConcurrency, "并发冲突", map[string]interface{}{
			"aggregate_id":     e.AggregateID,
			"expected_version": e.ExpectedVersion,
			"actual_version":   e.ActualVersion,
		}
	}

	// 检查DDD错误类型
	var dddErr *ddd.BusinessError
	if errors.As(err, &dddErr) {
		return m.mapDDDBusinessError(dddErr)
	}

	var validationErrs *ddd.ValidationErrors
	if errors.As(err, &validationErrs) {
		return http.StatusBadRequest, CodeInvalidParam, "参数验证失败", validationErrs.Errors
	}

	var concurrencyErr *ddd.ConcurrencyError
	if errors.As(err, &concurrencyErr) {
		return http.StatusConflict, CodeConcurrency, "并发冲突", map[string]interface{}{
			"aggregate_id":     concurrencyErr.AggregateID,
			"expected_version": concurrencyErr.ExpectedVersion,
			"actual_version":   concurrencyErr.ActualVersion,
		}
	}

	// 查找预定义映射
	if info, ok := m.customMappings[err]; ok {
		return info.HTTPStatus, info.Code, info.Message, nil
	}

	// 默认内部错误
	return http.StatusInternalServerError, CodeUnknownError, "服务器内部错误", nil
}

// mapBusinessError 映射业务错误
func (m *ErrorMapper) mapBusinessError(err *BusinessError) (int, int, string, interface{}) {
	// 根据业务错误码映射HTTP状态码
	httpStatus := GetHTTPStatus(err.Code)
	if httpStatus == 0 {
		httpStatus = http.StatusBadRequest
	}

	details := map[string]interface{}{
		"error_code": err.Code,
	}
	if err.Field != "" {
		details["field"] = err.Field
	}
	if err.Details != nil {
		details["details"] = err.Details
	}

	return httpStatus, err.Code, err.Message, details
}

// mapDDDBusinessError 映射DDD业务错误
func (m *ErrorMapper) mapDDDBusinessError(err *ddd.BusinessError) (int, int, string, interface{}) {
	// 根据DDD错误码映射到5位错误码
	code := m.mapDDDCodeToAppCode(err.Code)

	httpStatus := GetHTTPStatus(code)
	if httpStatus == 0 {
		httpStatus = http.StatusBadRequest
	}

	details := map[string]interface{}{
		"error_code": err.Code,
	}
	if err.Field != "" {
		details["field"] = err.Field
	}
	if err.Details != "" {
		details["details"] = err.Details
	}

	return httpStatus, code, err.Message, details
}

// mapDDDCodeToAppCode 将DDD错误码映射到5位应用错误码
func (m *ErrorMapper) mapDDDCodeToAppCode(dddCode string) int {
	// DDD错误码到应用错误码的映射
	mapping := map[string]int{
		"NOT_FOUND":             CodeNotFound,
		"USER_NOT_FOUND":        CodeUserNotFound,
		"TENANT_NOT_FOUND":      CodeTenantNotFound,
		"USERNAME_EXISTS":       CodeUserExists,
		"EMAIL_EXISTS":          CodeEmailExists,
		"INVALID_PASSWORD":      CodeInvalidPassword,
		"INVALID_EMAIL":         CodeInvalidEmail,
		"INVALID_USERNAME":      CodeInvalidUsername,
		"ACCOUNT_LOCKED":        CodeAccountLocked,
		"ACCOUNT_DISABLED":      CodeAccountDisabled,
		"ACCOUNT_NOT_ACTIVATED": CodeAccountNotActivated,
		"NOT_TENANT_MEMBER":     CodeNotTenantMember,
		"NOT_OWNER":             CodeNotTenantOwner,
		"TOKEN_EXPIRED":         CodeTokenExpired,
		"TOKEN_INVALID":         CodeTokenInvalid,
		"PERMISSION_DENIED":     CodePermissionDenied,
		"CONFLICT":              CodeConflict,
	}

	if appCode, ok := mapping[dddCode]; ok {
		return appCode
	}

	return CodeInvalidParam
}

// MustMap 映射错误（panic on error）
func (m *ErrorMapper) MustMap(err error) (int, int, string, interface{}) {
	httpStatus, code, message, details := m.Map(err)
	if httpStatus == 0 {
		panic("ErrorMapper.Map returned zero httpStatus")
	}
	return httpStatus, code, message, details
}
