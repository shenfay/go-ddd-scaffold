package errors

import "fmt"

// Coder 错误码接口
type Coder interface {
	GetCode() string
	GetMessage() string
}

// AppError 应用错误
type AppError struct {
	CodeVal    string
	MsgVal     string
	CatVal     string
	DetailsVal any
	CauseVal   error
}

// GetCode 返回错误码
func (e *AppError) GetCode() string {
	return e.CodeVal
}

// GetMessage 返回错误消息
func (e *AppError) GetMessage() string {
	return e.MsgVal
}

// GetCategory 返回分类
func (e *AppError) GetCategory() string {
	return e.CatVal
}

// GetDetails 返回详情
func (e *AppError) GetDetails() any {
	return e.DetailsVal
}

// GetCause 返回原因
func (e *AppError) GetCause() error {
	return e.CauseVal
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.CauseVal != nil {
		return fmt.Sprintf("%s: %v", e.MsgVal, e.CauseVal)
	}
	return e.MsgVal
}

// WithDetails 添加详情
func (e *AppError) WithDetails(d any) *AppError {
	e.DetailsVal = d
	return e
}

// WithCause 添加原因
func (e *AppError) WithCause(err error) *AppError {
	e.CauseVal = err
	return e
}

// SetCategory 设置分类
func (e *AppError) SetCategory(c string) *AppError {
	e.CatVal = c
	return e
}