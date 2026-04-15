package auth

import (
	"net/http"

	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// 认证域预定义错误
var (
	// ErrInvalidCredentials 凭证无效（邮箱或密码错误）
	ErrInvalidCredentials = &errors.AppError{
		Code:       errors.ErrCodeAuthInvalidCredentials,
		Message:    "Invalid email or password",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrInvalidToken 无效 Token
	ErrInvalidToken = &errors.AppError{
		Code:       errors.ErrCodeAuthInvalidToken,
		Message:    "Invalid token",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrTokenExpired Token 已过期
	ErrTokenExpired = &errors.AppError{
		Code:       errors.ErrCodeAuthTokenExpired,
		Message:    "Token has expired",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrTokenRevoked Token 已被撤销
	ErrTokenRevoked = &errors.AppError{
		Code:       errors.ErrCodeAuthTokenRevoked,
		Message:    "Token has been revoked",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrAccountLocked 账户已锁定
	ErrAccountLocked = &errors.AppError{
		Code:       errors.ErrCodeAuthAccountLocked,
		Message:    "Account locked due to too many failed login attempts",
		HTTPStatus: http.StatusLocked,
	}

	// ErrInvalidResetToken 无效或已过期的密码重置令牌
	ErrInvalidResetToken = &errors.AppError{
		Code:       errors.ErrCodeAuthInvalidResetToken,
		Message:    "Invalid or expired password reset token",
		HTTPStatus: http.StatusBadRequest,
	}
)

// NewAuthError 创建认证域错误（工厂方法）
func NewAuthError(code string, message string) *errors.AppError {
	return &errors.AppError{
		Code:       "AUTH." + code,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}
