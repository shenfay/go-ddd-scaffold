package valueobject

import (
	"regexp"
	"strings"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// UserName 用户名值对象
type UserName struct {
	value string
}

// NewUserName 创建用户名
func NewUserName(value string) (*UserName, error) {
	un := &UserName{value: strings.TrimSpace(value)}
	if err := un.Validate(); err != nil {
		return nil, err
	}
	return un, nil
}

// Value 返回用户名值
func (un *UserName) Value() string {
	return un.value
}

// Validate 验证用户名
func (un *UserName) Validate() error {
	if un.value == "" {
		return &kernel.ValidationError{
			Field:   "username",
			Message: "username cannot be empty",
		}
	}

	if len(un.value) < 3 {
		return &kernel.ValidationError{
			Field:   "username",
			Message: "username must be at least 3 characters long",
		}
	}

	if len(un.value) > 50 {
		return &kernel.ValidationError{
			Field:   "username",
			Message: "username cannot exceed 50 characters",
		}
	}

	// 只允许字母、数字、下划线和连字符
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validPattern.MatchString(un.value) {
		return &kernel.ValidationError{
			Field:   "username",
			Message: "username can only contain letters, numbers, underscores and hyphens",
		}
	}

	return nil
}

// Equals 比较用户名是否相等
func (un *UserName) Equals(other *UserName) bool {
	if other == nil {
		return false
	}
	return strings.EqualFold(un.value, other.value)
}
