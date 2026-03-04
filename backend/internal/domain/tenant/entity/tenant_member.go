package entity

import (
	"time"

	"github.com/google/uuid"
	sharedEntity "go-ddd-scaffold/internal/domain/shared/entity"
)

// MemberStatus 成员状态枚举
type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"   // 激活
	MemberStatusInactive MemberStatus = "inactive" // 未激活
	MemberStatusRemoved  MemberStatus = "removed"  // 已移除
)

// TenantMember 租户成员实体（作为 Tenant 的子实体）
type TenantMember struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID    `json:"tenantId" gorm:"type:uuid;index"`
	UserID    uuid.UUID    `json:"userId" gorm:"type:uuid;index"`
	Role      sharedEntity.UserRole     `json:"role" gorm:"size:20;index"`
	Status    MemberStatus `json:"status" gorm:"size:20;default:'active'"`
	InvitedBy *uuid.UUID   `json:"invitedBy,omitempty" gorm:"type:uuid;index"` // 邀请人 ID
	JoinedAt  time.Time    `json:"joinedAt"`
	LeftAt    *time.Time   `json:"leftAt,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewTenantMember 创建新的租户成员实例
func NewTenantMember(tenantID, userID uuid.UUID, role sharedEntity.UserRole) TenantMember {
	return TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     role,
		Status:   MemberStatusActive,
		JoinedAt: time.Now(),
	}
}

// IsActive 检查成员是否活跃
func (m *TenantMember) IsActive() bool {
	return m.Status == MemberStatusActive
}

// IsRemoved 检查成员是否已被移除
func (m *TenantMember) IsRemoved() bool {
	return m.Status == MemberStatusRemoved
}

// Remove 移除成员
func (m *TenantMember) Remove() {
	m.Status = MemberStatusRemoved
	leftAt := time.Now()
	m.LeftAt = &leftAt
}
