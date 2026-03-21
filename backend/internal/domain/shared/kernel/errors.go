package kernel

import (
	"errors"
	"fmt"
)

// ============================================
// 基础错误变量
// ============================================

var (
	// 聚合根相关错误
	ErrAggregateNotFound     = errors.New("aggregate not found")
	ErrAggregateConflict     = errors.New("aggregate version conflict")
	ErrInvalidAggregateState = errors.New("invalid aggregate state")

	// 领域事件相关错误
	ErrEventPublishFailed = errors.New("failed to publish event")
	ErrEventStoreFailed   = errors.New("failed to store event")

	// 仓储相关错误
	ErrRepositorySaveFailed   = errors.New("failed to save to repository")
	ErrRepositoryDeleteFailed = errors.New("failed to delete from repository")

	// 业务规则错误
	ErrBusinessRuleViolation = errors.New("business rule violation")
	ErrConcurrencyConflict   = errors.New("concurrency conflict detected")
)

// ============================================
// 业务错误定义
// ============================================

// BusinessError 业务错误结构
type BusinessError struct {
	Code    int         `json:"code"`              // 错误码
	Message string      `json:"message"`           // 错误消息
	Details interface{} `json:"details,omitempty"` // 详细错误信息
	Field   string      `json:"field,omitempty"`   // 字段信息
	cause   error       // 原始错误
}

// Error 实现error接口
func (be *BusinessError) Error() string {
	if be.Field != "" {
		return fmt.Sprintf("[%05d] %s: %s", be.Code, be.Field, be.Message)
	}
	return fmt.Sprintf("[%05d] %s", be.Code, be.Message)
}

// Unwrap 获取原始错误
func (be *BusinessError) Unwrap() error {
	return be.cause
}

// WithDetails 添加详细错误信息
func (be *BusinessError) WithDetails(details interface{}) *BusinessError {
	be.Details = details
	return be
}

// WithField 添加字段信息
func (be *BusinessError) WithField(field string) *BusinessError {
	be.Field = field
	return be
}

// WithMessage 自定义错误消息
func (be *BusinessError) WithMessage(msg string) *BusinessError {
	be.Message = msg
	return be
}

// WithCause 设置原始错误原因
func (be *BusinessError) WithCause(err error) *BusinessError {
	be.cause = err
	return be
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// ============================================
// 验证错误定义
// ============================================

// ValidationError 验证错误
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// Error 实现error接口
func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", ve.Field, ve.Message)
}

// ValidationErrors 验证错误集合
type ValidationErrors struct {
	Errors []*ValidationError `json:"errors"`
}

// Error 实现error接口
func (ves *ValidationErrors) Error() string {
	if len(ves.Errors) == 0 {
		return "no validation errors"
	}
	if len(ves.Errors) == 1 {
		return ves.Errors[0].Error()
	}
	return fmt.Sprintf("%d validation errors occurred", len(ves.Errors))
}

// Add 添加验证错误
func (ves *ValidationErrors) Add(field, message string, value interface{}) {
	ves.Errors = append(ves.Errors, &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors 检查是否有错误
func (ves *ValidationErrors) HasErrors() bool {
	return len(ves.Errors) > 0
}

// ToBusinessError 转换为业务错误
func (ves *ValidationErrors) ToBusinessError() *BusinessError {
	return NewBusinessError(
		CodeInvalidParam,
		"参数验证失败",
	).WithDetails(ves.Errors)
}

// ============================================
// 并发错误定义
// ============================================

// ConcurrencyError 并发错误
type ConcurrencyError struct {
	AggregateID     interface{} `json:"aggregate_id"`
	ExpectedVersion int         `json:"expected_version"`
	ActualVersion   int         `json:"actual_version"`
	Message         string      `json:"message"`
}

// Error 实现error接口
func (ce *ConcurrencyError) Error() string {
	return fmt.Sprintf("concurrency conflict for aggregate %v: expected version %d, actual version %d - %s",
		ce.AggregateID, ce.ExpectedVersion, ce.ActualVersion, ce.Message)
}

// NewConcurrencyError 创建并发错误
func NewConcurrencyError(aggregateID interface{}, expected, actual int, message string) *ConcurrencyError {
	return &ConcurrencyError{
		AggregateID:     aggregateID,
		ExpectedVersion: expected,
		ActualVersion:   actual,
		Message:         message,
	}
}

// ============================================
// 错误判断辅助函数
// ============================================

// IsBusinessError 检查是否为业务错误
func IsBusinessError(err error) bool {
	var be *BusinessError
	return errors.As(err, &be)
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	var ve *ValidationErrors
	return errors.As(err, &ve)
}

// AsBusinessError 转换为业务错误
func AsBusinessError(err error) *BusinessError {
	var be *BusinessError
	if errors.As(err, &be) {
		return be
	}
	return nil
}

// ============================================
// 错误码常量定义
// ============================================

const (
	// 成功
	CodeSuccess = 0

	// 通用错误 (10000-10099)
	CodeUnknownError  = 10000 // 未知错误
	CodeInvalidParam  = 10001 // 参数无效
	CodeBadRequest    = 10002 // 请求无效
	CodeNotFound      = 10003 // 资源不存在
	CodeConflict      = 10004 // 资源冲突
	CodeUnauthorized  = 10005 // 未授权
	CodeForbidden     = 10006 // 禁止访问
	CodeInternalError = 10010 // 内部错误
	CodeDatabaseError = 10011 // 数据库错误
	CodeCacheError    = 10012 // 缓存错误
	CodeConcurrency   = 10014 // 并发冲突

	// 认证授权 (40000-49999)
	CodeTokenExpired          = 40001 // Token已过期
	CodeTokenInvalid          = 40002 // Token无效
	CodeInvalidCredentials    = 40003 // 凭据无效
	CodeAccountDisabled       = 40004 // 账户已禁用
	CodeAccountCannotLogin    = 40005 // 账户无法登录
	CodeTokenGenerationFailed = 40006 // Token生成失败
	CodeInvalidToken          = 40007 // Token无效
	CodeInvalidUserID         = 40008 // 无效的用户ID
	CodePermissionDenied      = 41001 // 权限不足
)
