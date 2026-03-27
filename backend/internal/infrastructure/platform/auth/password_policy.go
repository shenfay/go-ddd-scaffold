package auth

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
)

// DefaultPasswordPolicy 默认密码策略实现
type DefaultPasswordPolicy struct {
	config service.PasswordPolicyConfig
	// 常见弱密码列表（实际项目中可以从配置文件或数据库加载）
	commonPasswords map[string]bool
}

// NewDefaultPasswordPolicy 创建默认密码策略
func NewDefaultPasswordPolicy(config service.PasswordPolicyConfig) *DefaultPasswordPolicy {
	return &DefaultPasswordPolicy{
		config: config,
		commonPasswords: map[string]bool{
			"password":    true,
			"123456":      true,
			"12345678":    true,
			"qwerty":      true,
			"abc123":      true,
			"password123": true,
			"admin":       true,
			"letmein":     true,
			"welcome":     true,
			"monkey":      true,
			"dragon":      true,
			"master":      true,
			"sunshine":    true,
			"princess":    true,
			"football":    true,
			"baseball":    true,
			"iloveyou":    true,
			"trustno1":    true,
			"1234567":     true,
			"123456789":   true,
			"admin123":    true,
			"welcome123":  true,
			"password1":   true,
			"123123":      true,
			"111111":      true,
		},
	}
}

// Validate 验证密码强度
func (p *DefaultPasswordPolicy) Validate(password string) error {
	// 1. 长度验证
	if err := p.validateLength(password); err != nil {
		return err
	}

	// 2. 字符类型验证
	if err := p.validateCharacterTypes(password); err != nil {
		return err
	}

	// 3. 常见密码检查
	if err := p.validateNotCommon(password); err != nil {
		return err
	}

	return nil
}

// validateLength 验证密码长度
func (p *DefaultPasswordPolicy) validateLength(password string) error {
	if len(password) < p.config.MinLength {
		return common.NewBusinessError(
			common.CodeInvalidParam,
			fmt.Sprintf("密码长度不能少于 %d 位", p.config.MinLength),
		)
	}
	if p.config.MaxLength > 0 && len(password) > p.config.MaxLength {
		return common.NewBusinessError(
			common.CodeInvalidParam,
			fmt.Sprintf("密码长度不能超过 %d 位", p.config.MaxLength),
		)
	}
	return nil
}

// validateCharacterTypes 验证字符类型
func (p *DefaultPasswordPolicy) validateCharacterTypes(password string) error {
	hasUpper, hasLower, hasDigit, hasSpecial := p.analyzeCharacterTypes(password)

	if p.config.RequireUppercase && !hasUpper {
		return common.NewBusinessError(
			common.CodeInvalidParam,
			"密码必须包含至少一个大写字母",
		)
	}

	if p.config.RequireLowercase && !hasLower {
		return common.NewBusinessError(
			common.CodeInvalidParam,
			"密码必须包含至少一个小写字母",
		)
	}

	if p.config.RequireDigits && !hasDigit {
		return common.NewBusinessError(
			common.CodeInvalidParam,
			"密码必须包含至少一个数字",
		)
	}

	if p.config.RequireSpecialChars && !hasSpecial {
		return common.NewBusinessError(
			common.CodeInvalidParam,
			fmt.Sprintf("密码必须包含至少一个特殊字符 (%s)", p.config.SpecialChars),
		)
	}

	return nil
}

// analyzeCharacterTypes 分析密码中的字符类型
func (p *DefaultPasswordPolicy) analyzeCharacterTypes(password string) (hasUpper, hasLower, hasDigit, hasSpecial bool) {
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(p.config.SpecialChars, char):
			hasSpecial = true
		}
	}
	return
}

// validateNotCommon 检查是否为常见弱密码
func (p *DefaultPasswordPolicy) validateNotCommon(password string) error {
	if !p.config.DisallowCommon {
		return nil
	}

	lowerPassword := strings.ToLower(password)
	if p.commonPasswords[lowerPassword] {
		return common.NewBusinessError(
			common.CodeInvalidParam,
			"该密码过于常见，请使用更复杂的密码",
		)
	}

	return nil
}

// GetPolicyDescription 获取策略描述
func (p *DefaultPasswordPolicy) GetPolicyDescription() string {
	var requirements []string

	requirements = append(requirements,
		fmt.Sprintf("长度 %d-%d 位", p.config.MinLength, p.config.MaxLength),
	)

	if p.config.RequireUppercase {
		requirements = append(requirements, "包含大写字母")
	}
	if p.config.RequireLowercase {
		requirements = append(requirements, "包含小写字母")
	}
	if p.config.RequireDigits {
		requirements = append(requirements, "包含数字")
	}
	if p.config.RequireSpecialChars {
		requirements = append(requirements, "包含特殊字符")
	}

	return "密码要求：" + strings.Join(requirements, "，")
}
