package apperrors

import (
	"errors"
	"fmt"
)

// ============================================
// 业务错误定义
// ============================================

// BusinessError 业务错误结构
type BusinessError struct {
	Code    int         `json:"code"`              // 5位错误码
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

// NewBusinessErrorWithDetails 创建带详情的业务错误
func NewBusinessErrorWithDetails(code int, message string, details interface{}) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewBusinessErrorWithField 创建带字段的业务错误
func NewBusinessErrorWithField(code int, message, field string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Field:   field,
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

// FirstError 返回第一个错误
func (ves *ValidationErrors) FirstError() error {
	if len(ves.Errors) == 0 {
		return nil
	}
	return ves.Errors[0]
}

// ToBusinessError 转换为业务错误
func (ves *ValidationErrors) ToBusinessError() *BusinessError {
	return NewBusinessErrorWithDetails(
		CodeInvalidParam,
		"参数验证失败",
		ves.Errors,
	)
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
	cause           error
}

// Error 实现error接口
func (ce *ConcurrencyError) Error() string {
	return fmt.Sprintf("concurrency conflict for aggregate %v: expected version %d, actual version %d - %s",
		ce.AggregateID, ce.ExpectedVersion, ce.ActualVersion, ce.Message)
}

// Unwrap 获取原始错误
func (ce *ConcurrencyError) Unwrap() error {
	return ce.cause
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

// IsConcurrencyError 检查是否为并发错误
func IsConcurrencyError(err error) bool {
	var ce *ConcurrencyError
	return errors.As(err, &ce)
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

// AsValidationError 转换为验证错误
func AsValidationError(err error) *ValidationErrors {
	var ve *ValidationErrors
	if errors.As(err, &ve) {
		return ve
	}
	return nil
}
