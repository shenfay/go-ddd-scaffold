package entity

import (
	"time"

	"go-ddd-scaffold/internal/domain/user/event"
	"go-ddd-scaffold/internal/domain/user/valueobject"

	"github.com/google/uuid"
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

	ErrInvalidUserStatus UserStatusError = "invalid user status"
)

// UserStatusError 用户状态错误
type UserStatusError string

func (e UserStatusError) Error() string {
	return string(e)
}

// ParseUserStatus 解析用户状态
func ParseUserStatus(s string) (UserStatus, error) {
	status := UserStatus(s)
	switch status {
	case StatusActive, StatusInactive, StatusLocked:
		return status, nil
	default:
		return "", ErrInvalidUserStatus
	}
}

// String 返回状态的字符串表示
func (s UserStatus) String() string {
	return string(s)
}

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

	// 领域事件（临时存储，由 Application Service 发布后清空）
	events []DomainEvent
}

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetEventType() string
	GetEventID() string
	GetAggregateID() uuid.UUID
	GetOccurredAt() time.Time
	GetVersion() int
}

// ErrAlreadyLocked 用户已锁定错误
type ErrAlreadyLocked string

func (e ErrAlreadyLocked) Error() string {
	return string(e)
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// Lock 锁定用户账号
// 用于违规封禁、安全保护等场景
func (u *User) Lock() error {
	if u.Status == StatusLocked {
		return ErrAlreadyLocked("user is already locked")
	}
	u.Status = StatusLocked
	u.addEvent(event.NewUserLockedEvent(u.ID))
	return nil
}

// Activate 激活用户账号
// 用于审核通过、解封等场景
func (u *User) Activate() error {
	if u.Status == StatusActive {
		return nil // 已经是激活状态
	}
	u.Status = StatusActive
	u.addEvent(event.NewUserActivatedEvent(u.ID))
	return nil
}

// UpdateProfile 更新用户资料
// 包含昵称、手机号、个人简介的批量更新
func (u *User) UpdateProfile(nickname valueobject.Nickname, phone *string, bio *string) error {
	u.Nickname = nickname
	u.Phone = phone
	u.Bio = bio
	u.addEvent(event.NewUserProfileUpdatedEvent(
		u.ID,
		nickname.String(),
		phone,
		bio,
	))
	return nil
}

// UpdateEmail 更新用户邮箱
// 重要的账户变更操作，需要验证新旧邮箱
func (u *User) UpdateEmail(newEmail valueobject.Email) error {
	if string(u.Email) == string(newEmail) {
		return nil // 邮箱未变化
	}
	oldEmail := u.Email
	u.Email = newEmail
	u.addEvent(event.NewUserEmailChangedEvent(u.ID, oldEmail.String(), newEmail.String()))
	return nil
}

// Events 获取待发布的领域事件
// Application Service 在持久化后调用此方法获取事件并发布
func (u *User) Events() []DomainEvent {
	return u.events
}

// ClearEvents 清空已发布的领域事件
// 在 Application Service 发布事件后调用
func (u *User) ClearEvents() {
	u.events = nil
}

// addEvent 添加领域事件（内部方法）
func (u *User) addEvent(event DomainEvent) {
	u.events = append(u.events, event)
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

// HashedPassword 已哈希的密码值对象
// 注意：密码加密逻辑已移至 Infrastructure 层（使用 PasswordHasher 接口）
type HashedPassword string

// String 返回哈希字符串（注意：不应该在日志中打印）
func (h HashedPassword) String() string {
	return string(h)
}
