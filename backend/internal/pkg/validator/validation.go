package validator

import (
	"context"
	"regexp"

	errPkg "go-ddd-scaffold/internal/pkg/errors"
)

// ==================== 业务校验类型定义 ====================

// BusinessValidationError 业务校验错误
type BusinessValidationError struct {
	Field   string
	Message string
}

// BusinessValidationErrors 业务校验错误集合
type BusinessValidationErrors []BusinessValidationError

// Error 实现 error 接口
func (e BusinessValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	result := ""
	for _, err := range e {
		if result != "" {
			result += "; "
		}
		result += err.Field + ": " + err.Message
	}
	return result
}

// HasError 检查是否有错误
func (e BusinessValidationErrors) HasError() bool {
	return len(e) > 0
}

// ==================== 通用校验函数 ====================

// ValidateEmailFormat 校验邮箱格式
func ValidateEmailFormat(email string) error {
	if email == "" {
		return errPkg.ErrInvalidEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errPkg.ErrInvalidEmail
	}

	return nil
}

// ValidatePasswordStrength 校验密码强度
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return errPkg.New("password", "密码长度不能少于6位")
	}

	hasDigit := false
	hasLetter := false

	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
	}

	if !hasDigit || !hasLetter {
		return errPkg.New("password", "密码必须包含字母和数字")
	}

	return nil
}

// ValidateNicknameFormat 校验昵称格式
func ValidateNicknameFormat(nickname string) error {
	if nickname == "" {
		return nil
	}

	if len(nickname) > 50 {
		return errPkg.New("nickname", "昵称长度不能超过50个字符")
	}

	nicknameRegex := regexp.MustCompile(`^[\u4e00-\u9fa5a-zA-Z0-9_]+$`)
	if !nicknameRegex.MatchString(nickname) {
		return errPkg.New("nickname", "昵称只能包含中文、字母、数字和下划线")
	}

	return nil
}

// ValidateTenantLimit 校验租户用户数限制
func ValidateTenantLimit(currentCount, maxCount int) error {
	if currentCount >= maxCount {
		return errPkg.ErrTenantLimitExceed
	}
	return nil
}

// ==================== 通用仓储接口 ====================

// UserRepository 用户仓储接口（业务校验专用）
type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByNickname(ctx context.Context, nickname string) (bool, error)
}
