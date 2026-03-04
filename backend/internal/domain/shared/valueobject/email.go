// Package valueobject 提供领域驱动设计中的值对象实现
package valueobject

import (
	"encoding/json"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email Email 值对象
type Email struct {
	value string
}

// NewEmail 创建新的 Email 值对象
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !IsValidEmail(email) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: email}, nil
}

// Value 获取底层字符串值
func (e Email) Value() string {
	return e.value
}

// String 转换为字符串
func (e Email) String() string {
	return e.value
}

// Equals 判断是否相等
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// MarshalJSON 实现 json.Marshaler 接口
func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.value)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (e *Email) UnmarshalJSON(data []byte) error {
	var emailStr string
	if err := json.Unmarshal(data, &emailStr); err != nil {
		return err
	}
	email, err := NewEmail(emailStr)
	if err != nil {
		return err
	}
	e.value = email.value
	return nil
}
