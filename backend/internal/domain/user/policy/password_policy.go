package policy

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// PasswordPolicyConfig 密码策略配置
type PasswordPolicyConfig struct {
	MinLength           int    // 最小长度
	MaxLength           int    // 最大长度
	RequireUppercase    bool   // 要求大写字母
	RequireLowercase    bool   // 要求小写字母
	RequireDigits       bool   // 要求数字
	RequireSpecialChars bool   // 要求特殊字符
	SpecialChars        string // 允许的特殊字符
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

// PasswordHasher 密码哈希接口
type PasswordHasher interface {
	// Hash 哈希密码
	Hash(password string) (string, error)
	// Verify 验证密码
	Verify(password, hash string) bool
}

// BcryptPasswordHasher 基于 bcrypt 的密码哈希实现
type BcryptPasswordHasher struct {
	cost int // bcrypt 成本因子
}

// NewBcryptPasswordHasher 创建 bcrypt 密码哈希器
func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptPasswordHasher{cost: cost}
}

// Hash 哈希密码
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Verify 验证密码
func (h *BcryptPasswordHasher) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// PasswordPolicyImpl 密码策略实现
type PasswordPolicyImpl struct {
	config PasswordPolicyConfig
}

// NewPasswordPolicy 创建密码策略实例
func NewPasswordPolicy(config PasswordPolicyConfig) *PasswordPolicyImpl {
	return &PasswordPolicyImpl{
		config: config,
	}
}

// Validate 验证密码强度
func (p *PasswordPolicyImpl) Validate(password string) error {
	// 检查长度
	if len(password) < p.config.MinLength {
		return fmt.Errorf("密码长度不能少于 %d 个字符", p.config.MinLength)
	}
	if len(password) > p.config.MaxLength {
		return fmt.Errorf("密码长度不能超过 %d 个字符", p.config.MaxLength)
	}

	// 检查大写字母
	if p.config.RequireUppercase && !containsUppercase(password) {
		return fmt.Errorf("密码必须包含至少一个大写字母")
	}

	// 检查小写字母
	if p.config.RequireLowercase && !containsLowercase(password) {
		return fmt.Errorf("密码必须包含至少一个小写字母")
	}

	// 检查数字
	if p.config.RequireDigits && !containsDigit(password) {
		return fmt.Errorf("密码必须包含至少一个数字")
	}

	// 检查特殊字符
	if p.config.RequireSpecialChars && !containsSpecialChar(password, p.config.SpecialChars) {
		return fmt.Errorf("密码必须包含至少一个特殊字符 (%s)", p.config.SpecialChars)
	}

	// 检查常见密码
	if p.config.DisallowCommon && isCommonPassword(password) {
		return fmt.Errorf("密码过于简单，请使用更复杂的密码")
	}

	return nil
}

// GetPolicyDescription 获取策略描述
func (p *PasswordPolicyImpl) GetPolicyDescription() string {
	var rules []string

	rules = append(rules, fmt.Sprintf("长度在 %d-%d 个字符之间", p.config.MinLength, p.config.MaxLength))

	if p.config.RequireUppercase {
		rules = append(rules, "包含至少一个大写字母")
	}
	if p.config.RequireLowercase {
		rules = append(rules, "包含至少一个小写字母")
	}
	if p.config.RequireDigits {
		rules = append(rules, "包含至少一个数字")
	}
	if p.config.RequireSpecialChars {
		rules = append(rules, fmt.Sprintf("包含至少一个特殊字符 (%s)", p.config.SpecialChars))
	}
	if p.config.DisallowCommon {
		rules = append(rules, "不使用常见密码")
	}

	return strings.Join(rules, "，") + "。"
}

// containsUppercase 检查是否包含大写字母
func containsUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// containsLowercase 检查是否包含小写字母
func containsLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// containsDigit 检查是否包含数字
func containsDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// containsSpecialChar 检查是否包含特殊字符
func containsSpecialChar(s string, specialChars string) bool {
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}

// isCommonPassword 检查是否为常见密码
func isCommonPassword(password string) bool {
	// 简单实现：检查一些常见密码
	commonPasswords := map[string]bool{
		"123456":      true,
		"password":    true,
		"12345678":    true,
		"qwerty":      true,
		"abc123":      true,
		"111111":      true,
		"1234567890":  true,
		"iloveyou":    true,
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
		"superman":    true,
		"trustno1":    true,
		"password1":   true,
		"1234567":     true,
		"123456789":   true,
		"12345678910": true,
		"qazwsx":      true,
		"1q2w3e4r":    true,
		"1q2w3e":      true,
		"123qwe":      true,
		"asdfghjkl":   true,
		"zaq12wsx":    true,
		"passw0rd":    true,
		"shadow":      true,
		"hidden":      true,
		"internet":    true,
		"batman":      true,
		"michael":     true,
		"pepper":      true,
		"hockey":      true,
		"harley":      true,
		"ranger":      true,
		"winter":      true,
		"summer":      true,
		"spring":      true,
		"autumn":      true,
		"fall":        true,
		"love":        true,
		"baby":        true,
		"angel":       true,
		"junior":      true,
		"thunder":     true,
		"matrix":      true,
		"pokemon":     true,
		"spongebob":   true,
		"starwars":    true,
	}

	lowerPassword := strings.ToLower(password)
	return commonPasswords[lowerPassword] || regexp.MustCompile(`^(.)\1+$`).MatchString(password)
}
