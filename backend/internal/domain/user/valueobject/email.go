package valueobject

import (
	"regexp"
	"strings"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// Email 邮箱值对象
type Email struct {
	value string
}

// NewEmail 创建邮箱
func NewEmail(value string) (*Email, error) {
	email := &Email{value: strings.TrimSpace(strings.ToLower(value))}
	if err := email.Validate(); err != nil {
		return nil, err
	}
	return email, nil
}

// Value 返回邮箱值
func (e *Email) Value() string {
	return e.value
}

// Validate 验证邮箱格式
func (e *Email) Validate() error {
	if e.value == "" {
		return &kernel.ValidationError{
			Field:   "email",
			Message: "email cannot be empty",
		}
	}

	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(e.value) {
		return &kernel.ValidationError{
			Field:   "email",
			Message: "invalid email format",
		}
	}

	return nil
}

// Equals 比较邮箱是否相等
func (e *Email) Equals(other *Email) bool {
	if other == nil {
		return false
	}
	return strings.EqualFold(e.value, other.value)
}
