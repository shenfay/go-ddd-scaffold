package vo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// ============================================================================
// Identity - 身份标识
// ============================================================================

// UserID 用户标识
type UserID struct {
	kernel.Int64Identity
}

// NewUserID 创建用户标识
func NewUserID(value int64) UserID {
	return UserID{Int64Identity: kernel.NewInt64Identity(value)}
}

// String 返回用户标识字符串
func (uid UserID) String() string {
	return fmt.Sprintf("%d", uid.Int64())
}

// ============================================================================
// UserName - 用户名值对象
// ============================================================================

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

// ============================================================================
// Email - 邮箱值对象
// ============================================================================

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

// ============================================================================
// HashedPassword - 加密密码值对象
// ============================================================================

// HashedPassword 加密密码值对象
type HashedPassword struct {
	value string
}

// NewHashedPassword 创建加密密码
func NewHashedPassword(hashedValue string) *HashedPassword {
	return &HashedPassword{value: hashedValue}
}

// Value 返回加密密码值
func (hp *HashedPassword) Value() string {
	return hp.value
}

// ============================================================================
// UserStatus - 用户状态枚举
// ============================================================================

// UserStatus 用户状态枚举
type UserStatus int

const (
	UserStatusPending  UserStatus = iota // 待激活
	UserStatusActive                     // 激活
	UserStatusInactive                   // 禁用
	UserStatusLocked                     // 锁定
)

// String 返回状态字符串表示
func (us UserStatus) String() string {
	switch us {
	case UserStatusPending:
		return "pending"
	case UserStatusActive:
		return "active"
	case UserStatusInactive:
		return "inactive"
	case UserStatusLocked:
		return "locked"
	default:
		return "unknown"
	}
}

// ============================================================================
// UserGender - 用户性别枚举
// ============================================================================

// UserGender 用户性别枚举
type UserGender int

const (
	UserGenderUnknown UserGender = iota
	UserGenderMale
	UserGenderFemale
	UserGenderOther
)

// String 返回性别字符串表示
func (ug UserGender) String() string {
	switch ug {
	case UserGenderMale:
		return "male"
	case UserGenderFemale:
		return "female"
	case UserGenderOther:
		return "other"
	default:
		return "unknown"
	}
}
