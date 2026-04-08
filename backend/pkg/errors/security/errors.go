package security

import (
	"net/http"

	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// 安全域预定义错误
var (
	// ErrTooManyAttempts 过多尝试次数
	ErrTooManyAttempts = &errors.AppError{
		Code:       errors.ErrCodeSecurityTooManyAttempts,
		Message:    "Too many attempts, please try again later",
		HTTPStatus: http.StatusTooManyRequests,
	}
)

// NewSecurityError 创建安全域错误（工厂方法）
func NewSecurityError(code string, message string) *errors.AppError {
	return &errors.AppError{
		Code:       "SECURITY." + code,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}
