package service

// PasswordPolicy 密码策略接口
// 定义密码强度验证的契约，由基础设施层实现具体策略
type PasswordPolicy interface {
	// Validate 验证密码强度
	// 返回验证错误，如果密码符合策略则返回 nil
	Validate(password string) error

	// GetPolicyDescription 获取策略描述（用于错误提示）
	GetPolicyDescription() string
}

// PasswordPolicyConfig 密码策略配置
type PasswordPolicyConfig struct {
	MinLength           int    // 最小长度
	MaxLength           int    // 最大长度
	RequireUppercase    bool   // 要求大写字母
	RequireLowercase    bool   // 要求小写字母
	RequireDigits       bool   // 要求数字
	RequireSpecialChars bool   // 要求特殊字符
	SpecialChars        string // 允许的特殂字符
	DisallowCommon      bool   // 禁止常见密码
}

// DefaultPasswordPolicyConfig 返回默认密码策略配置
func DefaultPasswordPolicyConfig() PasswordPolicyConfig {
	return PasswordPolicyConfig{
		MinLength:           8,
		MaxLength:           128,
		RequireUppercase:    true,
		RequireLowercase:    true,
		RequireDigits:       true,
		RequireSpecialChars: false,
		SpecialChars:        "!@#$%^&*()_+-=[]{}|;:,.<>?",
		DisallowCommon:      true,
	}
}
