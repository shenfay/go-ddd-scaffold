package valueobject

import (
	"regexp"
	"strings"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
)

// TenantCode 租户编码值对象
type TenantCode struct {
	value string
}

// NewTenantCode 创建租户编码
func NewTenantCode(value string) (*TenantCode, error) {
	code := strings.ToUpper(strings.TrimSpace(value))
	tc := &TenantCode{value: code}
	if err := tc.Validate(); err != nil {
		return nil, err
	}
	return tc, nil
}

// Value 返回租户编码值
func (tc *TenantCode) Value() string {
	return tc.value
}

// Validate 验证租户编码
func (tc *TenantCode) Validate() error {
	if tc.value == "" {
		return &common.ValidationError{
			Field:   "code",
			Message: "tenant code cannot be empty",
		}
	}

	if len(tc.value) < 3 {
		return &common.ValidationError{
			Field:   "code",
			Message: "tenant code must be at least 3 characters long",
		}
	}

	if len(tc.value) > 20 {
		return &common.ValidationError{
			Field:   "code",
			Message: "tenant code cannot exceed 20 characters",
		}
	}

	// 只允许大写字母、数字、下划线和连字符
	validPattern := regexp.MustCompile(`^[A-Z0-9_-]+$`)
	if !validPattern.MatchString(tc.value) {
		return &common.ValidationError{
			Field:   "code",
			Message: "tenant code can only contain uppercase letters, numbers, underscores and hyphens",
		}
	}

	return nil
}

// Equals 比较租户编码是否相等
func (tc *TenantCode) Equals(other *TenantCode) bool {
	if other == nil {
		return false
	}
	return tc.value == other.value
}
