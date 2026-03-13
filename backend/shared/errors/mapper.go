package apperrors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

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

// convertValidatorError 将 validator 错误安全地转换为自定义 ValidationErrors
func (m *ErrorMapper) convertValidatorError(err error) *ValidationErrors {
	// 尝试转换为 validator.ValidationErrors
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		ve := &ValidationErrors{}
		for _, fe := range validationErrs {
			ve.Add(fe.Field(), getValidationMessage(fe), nil)
		}
		return ve
	}

	// 尝试转换为 gin.ErrorTypeBind 错误
	var ginErr *gin.Error
	if errors.As(err, &ginErr) && ginErr.Type == gin.ErrorTypeBind {
		// 使用 errors.As 安全转换
		if errors.As(ginErr.Err, &validationErrs) {
			ve := &ValidationErrors{}
			for _, fe := range validationErrs {
				ve.Add(fe.Field(), getValidationMessage(fe), nil)
			}
			return ve
		}
	}

	return nil
}

// getValidationMessage 根据验证器标签生成友好的错误消息
func getValidationMessage(fe validator.FieldError) string {
	tag := fe.Tag()
	field := fe.Field()

	messages := map[string]string{
		"required": "不能为空",
		"email":    "必须是有效的邮箱地址",
		"min":      "长度不能小于 " + fe.Param(),
		"max":      "长度不能超过 " + fe.Param(),
		"len":      "长度必须为 " + fe.Param(),
		"oneof":    "必须是以下值之一：" + fe.Param(),
		"url":      "必须是有效的 URL",
		"numeric":  "必须是数字",
		"number":   "必须是数字",
		"boolean":  "必须是布尔值",
	}

	if msg, ok := messages[tag]; ok {
		return field + msg
	}

	return field + "验证失败 (" + tag + ")"
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

	// 检查 DDD 错误类型
	var dddErr *ddd.BusinessError
	if errors.As(err, &dddErr) {
		return m.mapDDDBusinessError(dddErr)
	}

	var validationErrs *ddd.ValidationErrors
	if errors.As(err, &validationErrs) {
		return http.StatusBadRequest, CodeInvalidParam, "参数验证失败", validationErrs.Errors
	}

	// 安全地将 validator 错误转换为自定义 ValidationErrors
	if convertedErr := m.convertValidatorError(err); convertedErr != nil {
		return http.StatusBadRequest, CodeInvalidParam, "参数验证失败", convertedErr.Errors
	}

	// 查找预定义映射（使用安全的比较方式）
	for mappedErr, info := range m.customMappings {
		if errors.Is(err, mappedErr) || err.Error() == mappedErr.Error() {
			return info.HTTPStatus, info.Code, info.Message, nil
		}
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
