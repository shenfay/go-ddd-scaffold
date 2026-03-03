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
	Password string     `json:"-" gorm:"size:255"` // 不序列化到JSON
	Nickname string     `json:"nickname" gorm:"size:100"`
	Avatar   string     `json:"avatar" gorm:"size:500"`
	Role     UserRole   `json:"role" gorm:"size:20;index"`
	Status   UserStatus `json:"status" gorm:"size:20;default:'active'"`
	TenantID *uuid.UUID `json:"tenantId,omitempty" gorm:"type:uuid;index"` // 租户 ID（多租户场景）
	// 可扩展字段：根据具体业务需求添加

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// IsValidRole 检查角色是否有效
func (u *User) IsValidRole() bool {
	switch u.Role {
	case RoleMember, RoleGuest, RoleSuperAdmin, RoleContentAdmin, RoleOpsAdmin:
		return true
	default:
		return false
	}
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsMember 检查是否为成员用户
func (u *User) IsMember() bool {
	return u.Role == RoleMember
}

// IsGuest 检查是否为访客用户
func (u *User) IsGuest() bool {
	return u.Role == RoleGuest
}

// TokenClaims JWT令牌声明
type TokenClaims struct {
	UserID   uuid.UUID `json:"userId"`
	Role     UserRole  `json:"role"`
	TenantID uuid.UUID `json:"tenantId"` // 租户ID
}

// JWTService JWT服务接口
type JWTService interface {
	// GenerateToken 生成JWT令牌
	GenerateToken(userID uuid.UUID, role UserRole, tenantID uuid.UUID) (string, error)
	// ValidateToken 验证JWT令牌
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
