package entity

import (
	"time"

	sharedEntity "go-ddd-scaffold/internal/domain/shared/entity"

	"github.com/google/uuid"
)

// Tenant 租户聚合根（多租户 SaaS 场景）
type Tenant struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"size:100"`
	Description string    `json:"description" gorm:"size:500"`

	// 租户配置
	MaxMembers int       `json:"maxMembers" gorm:"default:10"` // 最大成员数
	ExpiredAt  time.Time `json:"expiredAt"`                    // 过期时间

	Members []TenantMember `json:"members,omitempty" gorm:"foreignKey:TenantID"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewTenant 创建新的租户实例
func NewTenant(name string, maxMembers int) *Tenant {
	return &Tenant{
		ID:          uuid.New(),
		Name:        name,
		Description: "",
		MaxMembers:  maxMembers,
		ExpiredAt:   time.Now().AddDate(1, 0, 0), // 默认一年后过期
		Members:     make([]TenantMember, 0),
	}
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

// AddMember 添加成员到租户（聚合根方法）
func (t *Tenant) AddMember(userID uuid.UUID, role sharedEntity.UserRole, invitedBy *uuid.UUID) (*TenantMember, error) {
	if !t.IsValid() {
		return nil, ErrTenantInvalid
	}

	if len(t.Members) >= t.MaxMembers {
		return nil, ErrTenantMemberLimitExceeded
	}

	// 检查用户是否已是成员
	for _, member := range t.Members {
		if member.UserID == userID && member.Status == MemberStatusActive {
			return nil, ErrTenantMemberAlreadyExists
		}
	}

	member := TenantMember{
		ID:        uuid.New(),
		TenantID:  t.ID,
		UserID:    userID,
		Role:      role,
		Status:    MemberStatusActive,
		InvitedBy: invitedBy,
		JoinedAt:  time.Now(),
	}

	t.Members = append(t.Members, member)
	return &member, nil
}

// RemoveMember 移除租户成员（聚合根方法）
func (t *Tenant) RemoveMember(memberID uuid.UUID) error {
	for i, member := range t.Members {
		if member.ID == memberID {
			t.Members[i].Status = MemberStatusRemoved
			leftAt := time.Now()
			t.Members[i].LeftAt = &leftAt
			return nil
		}
	}
	return ErrTenantMemberNotFound
}

// GetActiveMemberCount 获取活跃成员数量
func (t *Tenant) GetActiveMemberCount() int {
	count := 0
	for _, member := range t.Members {
		if member.Status == MemberStatusActive {
			count++
		}
	}
	return count
}
