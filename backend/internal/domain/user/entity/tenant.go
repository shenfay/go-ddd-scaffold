package entity

import (
	"time"

	"github.com/google/uuid"
)

// MemberStatus 成员状态枚举
type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"   // 激活
	MemberStatusInactive MemberStatus = "inactive" // 未激活
	MemberStatusRemoved  MemberStatus = "removed"  // 已移除
)

// TenantMember 租户成员实体
type TenantMember struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID    `json:"tenantId" gorm:"type:uuid;index"`
	UserID    uuid.UUID    `json:"userId" gorm:"type:uuid;index"`
	Role      UserRole     `json:"role" gorm:"size:20;index"`
	Status    MemberStatus `json:"status" gorm:"size:20;default:'active'"`
	InvitedBy *uuid.UUID   `json:"invitedBy,omitempty" gorm:"type:uuid;index"` // 邀请人ID
	JoinedAt  time.Time    `json:"joinedAt"`
	LeftAt    *time.Time   `json:"leftAt,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Tenant 租户实体（多租户 SaaS 场景）
type Tenant struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"size:100"`
	Description string    `json:"description" gorm:"size:500"`

	// 租户配置
	MaxMembers int       `json:"maxMembers" gorm:"default:10"` // 最大成员数
	ExpiredAt  time.Time `json:"expiredAt"`                    // 过期时间

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// IsValid 检查租户是否有效
func (t *Tenant) IsValid() bool {
	return t.ExpiredAt.After(time.Now()) && t.MaxMembers > 0
}

// IsExpired 检查租户是否过期
func (t *Tenant) IsExpired() bool {
	return t.ExpiredAt.Before(time.Now())
}

// CanAddMoreMembers 检查是否还能添加更多成员
func (t *Tenant) CanAddMoreMembers(currentCount int) bool {
	return currentCount < t.MaxMembers
}
