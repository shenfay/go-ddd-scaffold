package validation

import (
	"net/http"

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
