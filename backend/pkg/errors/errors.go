package errors

import (
	"fmt"
	"net/http"
)

// AppError 应用错误
type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	HTTPStatus int         `json:"-"`
	Err        error       `json:"-"` // 内部错误，不返回给客户端
	Metadata   interface{} `json:"details,omitempty"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Err.Error())
	}
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建应用错误
func NewAppError(code string, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// WithError 设置内部错误
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// WithMetadata 设置元数据
func (e *AppError) WithMetadata(metadata interface{}) *AppError {
	e.Metadata = metadata
	return e
}

// 常见错误码定义
const (
	// 通用错误
	ErrorCodeInternal       = "INTERNAL_ERROR"
	ErrorCodeInvalidRequest = "INVALID_REQUEST"
	ErrorCodeNotFound       = "NOT_FOUND"
	ErrorCodeUnauthorized   = "UNAUTHORIZED"
	ErrorCodeForbidden      = "FORBIDDEN"
	
	// 认证相关错误
	ErrorCodeEmailAlreadyExists     = "EMAIL_ALREADY_EXISTS"
	ErrorCodeInvalidCredentials     = "INVALID_CREDENTIALS"
	ErrorCodeInvalidToken           = "INVALID_TOKEN"
	ErrorCodeTokenExpired           = "TOKEN_EXPIRED"
	ErrorCodeTokenRevoked           = "TOKEN_REVOKED"
	ErrorCodeUserNotFound           = "USER_NOT_FOUND"
	ErrorCodeEmailNotVerified       = "EMAIL_NOT_VERIFIED"
	ErrorCodeAccountLocked          = "ACCOUNT_LOCKED"
	ErrorCodeTooManyLoginAttempts   = "TOO_MANY_LOGIN_ATTEMPTS"
)

// 预定义的错误
var (
	ErrInternal = NewAppError(ErrorCodeInternal, "Internal server error", http.StatusInternalServerError)
	
	ErrInvalidRequest = NewAppError(ErrorCodeInvalidRequest, "Invalid request", http.StatusBadRequest)
	
	ErrUnauthorized = NewAppError(ErrorCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	
	ErrNotFound = NewAppError(ErrorCodeNotFound, "Resource not found", http.StatusNotFound)
	
	ErrEmailAlreadyExists = NewAppError(ErrorCodeEmailAlreadyExists, "Email already exists", http.StatusConflict)
	
	ErrInvalidCredentials = NewAppError(ErrorCodeInvalidCredentials, "Invalid email or password", http.StatusUnauthorized)
	
	ErrInvalidToken = NewAppError(ErrorCodeInvalidToken, "Invalid token", http.StatusUnauthorized)
	
	ErrTokenExpired = NewAppError(ErrorCodeTokenExpired, "Token expired", http.StatusUnauthorized)
	
	ErrUserNotFound = NewAppError(ErrorCodeUserNotFound, "User not found", http.StatusNotFound)
	
	ErrTooManyLoginAttempts = NewAppError(ErrorCodeTooManyLoginAttempts, "Too many login attempts, please try again later", http.StatusTooManyRequests)
	
	ErrAccountLocked = NewAppError(ErrorCodeAccountLocked, "Account locked due to too many failed login attempts", http.StatusLocked)
)
