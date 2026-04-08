package validation

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// 校验域预定义错误
var (
	// ErrFieldRequired 字段必填
	ErrFieldRequired = &errors.AppError{
		Code:       "VALIDATION.FIELD_REQUIRED",
		Message:    "Field is required",
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrFieldInvalid 字段格式无效
	ErrFieldInvalid = &errors.AppError{
		Code:       "VALIDATION.FIELD_INVALID",
		Message:    "Field format is invalid",
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrFieldTooShort 字段长度太短
	ErrFieldTooShort = &errors.AppError{
		Code:       "VALIDATION.FIELD_TOO_SHORT",
		Message:    "Field length is too short",
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrFieldTooLong 字段长度太长
	ErrFieldTooLong = &errors.AppError{
		Code:       "VALIDATION.FIELD_TOO_LONG",
		Message:    "Field length is too long",
		HTTPStatus: http.StatusBadRequest,
	}
)

// NewValidationError 创建校验域错误（工厂方法）
func NewValidationError(code string, message string) *errors.AppError {
	return &errors.AppError{
		Code:       "VALIDATION." + code,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// FromGinError 从 Gin 绑定错误创建验证错误
func FromGinError(err error) *errors.AppError {
	if err == nil {
		return nil
	}

	// 检查是否是验证错误
	if ve, ok := err.(validator.ValidationErrors); ok {
		messages := make([]string, 0, len(ve))
		for _, fe := range ve {
			messages = append(messages, formatFieldError(fe))
		}
		return &errors.AppError{
			Code:       "VALIDATION.INVALID_REQUEST",
			Message:    strings.Join(messages, "; "),
			HTTPStatus: http.StatusBadRequest,
		}
	}

	// 其他错误
	return &errors.AppError{
		Code:       "VALIDATION.INVALID_REQUEST",
		Message:    err.Error(),
		HTTPStatus: http.StatusBadRequest,
	}
}

// formatFieldError 格式化字段错误
func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
