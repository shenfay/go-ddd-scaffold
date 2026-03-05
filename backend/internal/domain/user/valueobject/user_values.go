package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Email 邮箱值对象
type Email string

func NewEmail(email string) (Email, error) {
	e := Email(email)
	if !e.IsValid() {
		return "", errors.New("邮箱格式不正确")
	}
	return e, nil
}

func (e Email) IsValid() bool {
	if len(string(e)) == 0 || len(string(e)) > 255 {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(string(e))
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Domain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// NewEmailFromString 从字符串创建 Email（不验证，用于从数据库读取已持久化的数据）
func NewEmailFromString(s string) Email {
	return Email(s)
}

// NewNicknameFromString 从字符串创建 Nickname（不验证，用于从数据库读取已持久化的数据）
func NewNicknameFromString(s string) Nickname {
	return Nickname(s)
}

// PlainPassword 明文密码值对象（仅用于输入验证）
type PlainPassword string

func NewPlainPassword(password string) (PlainPassword, error) {
	p := PlainPassword(password)
	if !p.IsValid() {
		return "", errors.New("密码不符合要求")
	}
	return p, nil
}

func (p PlainPassword) IsValid() bool {
	password := string(p)

	// 长度检查：至少 6 位（根据项目规范）
	if len(password) < 6 {
		return false
	}

	// 最长 64 位
	if len(password) > 64 {
		return false
	}

	// 必须包含数字
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	// 必须包含字母
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)

	return hasNumber && hasLetter
}

func (p PlainPassword) String() string {
	return string(p)
}

func (p PlainPassword) Strength() string {
	password := string(p)

	// 简单强度评估
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	strength := 0
	if hasLower {
		strength++
	}
	if hasUpper {
		strength++
	}
	if hasNumber {
		strength++
	}
	if hasSpecial {
		strength++
	}

	if len(password) >= 12 {
		strength++
	}

	if strength >= 4 {
		return "强"
	} else if strength >= 3 {
		return "中"
	} else {
		return "弱"
	}
}

// HashedPassword 已哈希的密码值对象（用于存储）
type HashedPassword string

func NewHashedPassword(hashed string) HashedPassword {
	return HashedPassword(hashed)
}

func (h HashedPassword) String() string {
	return string(h)
}

func (h HashedPassword) Verify(plainPassword PlainPassword) bool {
	// 实际验证由 entity.HashedPassword 处理
	return true
}

// Nickname 昵称值对象
type Nickname string

func NewNickname(nickname string) (Nickname, error) {
	n := Nickname(nickname)
	if !n.IsValid() {
		return "", errors.New("昵称格式不正确")
	}
	return n, nil
}

func (n Nickname) IsValid() bool {
	nickname := string(n)

	// 长度检查：至少2个字符
	if utf8.RuneCountInString(nickname) < 2 {
		return false
	}

	// 最长20个字符
	if utf8.RuneCountInString(nickname) > 20 {
		return false
	}

	// 不能包含特殊字符（允许中文、英文、数字、下划线）
	validPattern := regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_]+$`)
	return validPattern.MatchString(nickname)
}

func (n Nickname) String() string {
	return string(n)
}

func (n Nickname) Length() int {
	return utf8.RuneCountInString(string(n))
}

// Permission 权限值对象
type Permission struct {
	Resource    string // 资源标识
	Action      string // 操作类型：read, create, update, delete
	Effect      bool   // 是否允许
	TenantLevel bool   // 是否为租户级权限
}

func NewPermission(resource, action string, effect, tenantLevel bool) (Permission, error) {
	if resource == "" || action == "" {
		return Permission{}, errors.New("资源和操作不能为空")
	}

	// 验证操作类型
	validActions := map[string]bool{
		"read":   true,
		"create": true,
		"update": true,
		"delete": true,
		"manage": true, // 完全管理权限
	}

	if !validActions[action] {
		return Permission{}, fmt.Errorf("无效的操作类型: %s", action)
	}

	return Permission{
		Resource:    resource,
		Action:      action,
		Effect:      effect,
		TenantLevel: tenantLevel,
	}, nil
}

func (p Permission) String() string {
	effect := "deny"
	if p.Effect {
		effect = "allow"
	}
	return fmt.Sprintf("%s:%s:%s", p.Resource, p.Action, effect)
}

func (p Permission) IsAllowed() bool {
	return p.Effect
}

func (p Permission) Matches(resource, action string) bool {
	return p.Resource == resource && p.Action == action
}

// TenantName 租户名称值对象
type TenantName string

func NewTenantName(name string) (TenantName, error) {
	tn := TenantName(name)
	if !tn.IsValid() {
		return "", errors.New("租户名称格式不正确")
	}
	return tn, nil
}

func (tn TenantName) IsValid() bool {
	name := string(tn)

	// 长度检查：至少2个字符
	if utf8.RuneCountInString(name) < 2 {
		return false
	}

	// 最长50个字符
	if utf8.RuneCountInString(name) > 50 {
		return false
	}

	return true
}

func (tn TenantName) String() string {
	return string(tn)
}
