package ddd

import (
	"errors"
	"fmt"
)

// 企业级DDD错误定义
var (
	// 聚合根相关错误
	ErrAggregateNotFound     = errors.New("aggregate not found")
	ErrAggregateConflict     = errors.New("aggregate version conflict")
	ErrInvalidAggregateState = errors.New("invalid aggregate state")

	// 领域事件相关错误
	ErrEventPublishFailed    = errors.New("failed to publish event")
	ErrEventStoreFailed      = errors.New("failed to store event")
	ErrEventReplayFailed     = errors.New("failed to replay events")

	// 仓储相关错误
	ErrRepositorySaveFailed  = errors.New("failed to save to repository")
	ErrRepositoryDeleteFailed = errors.New("failed to delete from repository")
	ErrRepositoryQueryFailed = errors.New("failed to query repository")

	// 业务规则错误
	ErrBusinessRuleViolation = errors.New("business rule violation")
	ErrInvalidOperation      = errors.New("invalid operation")
	ErrConcurrencyConflict   = errors.New("concurrency conflict detected")

	// 验证错误
	ErrValidationFailed      = errors.New("validation failed")
	ErrInvalidParameter      = errors.New("invalid parameter")
)

// BusinessError 业务错误结构
type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Field   string `json:"field,omitempty"`
}

// Error 实现error接口
func (be *BusinessError) Error() string {
	if be.Field != "" {
		return fmt.Sprintf("[%s] %s: %s", be.Code, be.Field, be.Message)
	}
	return fmt.Sprintf("[%s] %s", be.Code, be.Message)
}

// NewBusinessError 创建业务错误
func NewBusinessError(code, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// NewBusinessErrorWithDetails 创建带详情的业务错误
func NewBusinessErrorWithDetails(code, message, details string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewBusinessErrorWithField 创建带字段的业务错误
func NewBusinessErrorWithField(code, message, field string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Field:   field,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
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

// ConcurrencyError 并发错误
type ConcurrencyError struct {
	AggregateID interface{} `json:"aggregate_id"`
	ExpectedVersion int     `json:"expected_version"`
	ActualVersion   int     `json:"actual_version"`
	Message         string  `json:"message"`
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

// IsConcurrencyError 检查是否为并发错误
func IsConcurrencyError(err error) bool {
	_, ok := err.(*ConcurrencyError)
	return ok
}