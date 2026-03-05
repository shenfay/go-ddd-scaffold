package entity

import (
	"time"

	"go-ddd-scaffold/internal/domain/user/valueobject"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRole 用户角色枚举（用于租户成员等）
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
// 纯领域对象，不包含任何基础设施标签
type User struct {
	ID        uuid.UUID
	Email     valueobject.Email
	Password  HashedPassword
	Nickname  valueobject.Nickname
	Avatar    *string
	Phone     *string
	Bio       *string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// 注意：IsValidRole、IsMember、IsGuest 方法已移除
// 角色检查应通过 Casbin 在租户上下文中进行：
// roles := casbinEnforcer.GetRolesForUserInDomain(userID, tenantID)

// TokenClaims JWT 令牌声明（扩展版）
type TokenClaims struct {
	UserID   uuid.UUID  `json:"userId"`
	TenantID *uuid.UUID `json:"tenantId,omitempty"` // 当前租户 ID（可选，用户可能属于多个租户）
	// 注意：不再包含 Role
	// - 角色在租户上下文中动态查询
	// - 通过 Casbin RBAC 中间件进行权限控制
}

// JWTService JWT 服务接口
type JWTService interface {
	// GenerateToken 生成 JWT 令牌（仅包含用户 ID）
	GenerateToken(userID uuid.UUID) (string, error)
	// GenerateTokenWithTenant 生成带租户上下文的 JWT 令牌
	GenerateTokenWithTenant(userID uuid.UUID, tenantID uuid.UUID) (string, error)
	// ValidateToken 验证 JWT 令牌
	ValidateToken(tokenString string) (*TokenClaims, error)
}

// HashedPassword 已哈希的密码值对象（用于基础设施层）
type HashedPassword string

// NewHashedPassword 对明文密码进行哈希
func NewHashedPassword(plainPassword string) (HashedPassword, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return HashedPassword(string(bytes)), err
}

// Verify 验证密码是否匹配
func (h HashedPassword) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(plainPassword))
	return err == nil
}

// String 返回哈希字符串（注意：不应该在日志中打印）
func (h HashedPassword) String() string {
	return string(h)
}
