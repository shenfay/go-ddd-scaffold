// Package valueobject 值对象相关错误定义
package valueobject

import "errors"

var (
	ErrInvalidEmail = errors.New("invalid email format")
)
