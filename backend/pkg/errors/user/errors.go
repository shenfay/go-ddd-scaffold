package user

import (
	"net/http"

	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// 用户域预定义错误
var (
	// ErrNotFound 用户不存在
	ErrNotFound = &errors.AppError{
		Code:       errors.ErrCodeUserNotFound,
		Message:    "User not found",
		HTTPStatus: http.StatusNotFound,
	}

	// ErrEmailAlreadyExists 邮箱已注册
	ErrEmailAlreadyExists = &errors.AppError{
		Code:       errors.ErrCodeUserEmailAlreadyExists,
		Message:    "Email address is already registered",
		HTTPStatus: http.StatusConflict,
	}

	// ErrEmailNotVerified 邮箱未验证
	ErrEmailNotVerified = &errors.AppError{
		Code:       errors.ErrCodeUserEmailNotVerified,
		Message:    "Email address has not been verified",
		HTTPStatus: http.StatusForbidden,
	}
)

// NewUserError 创建用户域错误（工厂方法）
func NewUserError(code string, message string) *errors.AppError {
	return &errors.AppError{
		Code:       "USER." + code,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}
