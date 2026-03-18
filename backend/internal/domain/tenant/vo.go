package tenant

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/kernel"
)

// TenantID 租户标识
type TenantID struct {
	kernel.Int64Identity
}

// NewTenantID 创建租户标识
func NewTenantID(value int64) TenantID {
	return TenantID{Int64Identity: kernel.NewInt64Identity(value)}
}

// String 返回租户标识字符串
func (tid TenantID) String() string {
	return fmt.Sprintf("tenant-%d", tid.Int64())
}

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
		return &kernel.ValidationError{
			Field:   "code",
			Message: "tenant code cannot be empty",
		}
	}

	if len(tc.value) < 3 {
		return &kernel.ValidationError{
			Field:   "code",
			Message: "tenant code must be at least 3 characters long",
		}
	}

	if len(tc.value) > 20 {
		return &kernel.ValidationError{
			Field:   "code",
			Message: "tenant code cannot exceed 20 characters",
		}
	}

	// 只允许大写字母、数字、下划线和连字符
	validPattern := regexp.MustCompile(`^[A-Z0-9_-]+$`)
	if !validPattern.MatchString(tc.value) {
		return &kernel.ValidationError{
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

// TenantStatus 租户状态枚举
type TenantStatus int

const (
	TenantStatusActive TenantStatus = iota   // 活跃
	TenantStatusInactive                     // 停用
	TenantStatusSuspended                    // 暂停
)

// String 返回状态字符串表示
func (ts TenantStatus) String() string {
	switch ts {
	case TenantStatusActive:
		return "active"
	case TenantStatusInactive:
		return "inactive"
	case TenantStatusSuspended:
		return "suspended"
	default:
		return "unknown"
	}
}

// TenantConfig 租户配置值对象
type TenantConfig struct {
	MaxStorageGB      int               `json:"max_storage_gb"`
	MaxProjects       int               `json:"max_projects"`
	AllowedFeatures   []string          `json:"allowed_features"`
	CustomSettings    map[string]string `json:"custom_settings"`
	RequireMFA        bool              `json:"require_mfa"`
	SessionTimeoutMin int               `json:"session_timeout_min"`
}

// NewDefaultTenantConfig 创建默认租户配置
func NewDefaultTenantConfig() *TenantConfig {
	return &TenantConfig{
		MaxStorageGB:      10,
		MaxProjects:       10,
		AllowedFeatures:   []string{"basic", "api_access"},
		CustomSettings:    make(map[string]string),
		RequireMFA:        false,
		SessionTimeoutMin: 30,
	}
}

// TenantMember 租户成员值对象
type TenantMember struct {
	UserID   user.UserID
	TenantID TenantID
	Role     TenantRole
	JoinedAt string
}

// TenantRole 租户角色枚举
type TenantRole int

const (
	TenantRoleOwner TenantRole = iota  // 所有者
	TenantRoleAdmin                    // 管理员
	TenantRoleMember                   // 成员
	TenantRoleGuest                    // 访客
)

// String 返回角色字符串表示
func (tr TenantRole) String() string {
	switch tr {
	case TenantRoleOwner:
		return "owner"
	case TenantRoleAdmin:
		return "admin"
	case TenantRoleMember:
		return "member"
	case TenantRoleGuest:
		return "guest"
	default:
		return "unknown"
	}
}
