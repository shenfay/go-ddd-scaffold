package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRole 用户角色枚举
// 通用角色定义，可根据具体业务场景扩展
type UserRole string

const (
	// 基础角色
	RoleOwner  UserRole = "owner"  // 所有者
	RoleAdmin  UserRole = "admin"  // 管理员
	RoleMember UserRole = "member" // 普通成员
	RoleGuest  UserRole = "guest"  // 访客

	// 管理角色
	RoleSuperAdmin   UserRole = "super_admin"   // 超级管理员
	RoleContentAdmin UserRole = "content_admin" // 内容管理员
	RoleOpsAdmin     UserRole = "ops_admin"     // 运营管理员
)

// UserStatus 用户状态枚举
type UserStatus string

const (
	StatusActive   UserStatus = "active"   // 激活
	StatusInactive UserStatus = "inactive" // 未激活
	StatusLocked   UserStatus = "locked"   // 锁定
)

// User 用户实体（聚合根）
type User struct {
	ID       uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Email    string     `json:"email" gorm:"uniqueIndex;size:255"`
	Password string     `json:"-" gorm:"size:255"` // 不序列化到 JSON
	Nickname string     `json:"nickname" gorm:"size:100"`
	Avatar   *string    `json:"avatar,omitempty" gorm:"size:500"`
	Phone    *string    `json:"phone,omitempty" gorm:"size:20"`
	Bio      *string    `json:"bio,omitempty" gorm:"size:500"`
	Status   UserStatus `json:"status" gorm:"size:20;default:'active'"`
	// 注意：Role 和 TenantID 已移除，改用多租户 +Casbin RBAC 设计
	// - 用户可属于多个租户（通过 tenant_members 表关联）
	// - 用户在每个租户内的角色存储在 tenant_members.role
	// - 权限控制通过 Casbin RBAC 中间件处理

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// 注意：IsValidRole、IsMember、IsGuest 方法已移除
// 角色检查应通过 Casbin 在租户上下文中进行：
// roles := casbinEnforcer.GetRolesForUserInDomain(userID, tenantID)

// TokenClaims JWT 令牌声明（简化版，只包含用户 ID）
type TokenClaims struct {
	UserID uuid.UUID `json:"userId"`
	// 注意：不再包含 Role 和 TenantID
	// - 用户可属于多个租户，角色在租户上下文中动态查询
	// - 通过 Casbin RBAC 中间件进行权限控制
}

// JWTService JWT 服务接口
type JWTService interface {
	// GenerateToken 生成 JWT 令牌（仅包含用户 ID）
	GenerateToken(userID uuid.UUID) (string, error)
	// ValidateToken 验证 JWT 令牌
	ValidateToken(tokenString string) (*TokenClaims, error)
}

// HashPassword 对密码进行bcrypt哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码是否匹配
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
