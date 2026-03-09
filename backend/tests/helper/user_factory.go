// Package helper 测试辅助工具包
package helper

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserFactory 用户测试工厂
type UserFactory struct {
	t *testing.T
}

// NewUserFactory 创建用户工厂
func NewUserFactory(t *testing.T) *UserFactory {
	return &UserFactory{t: t}
}

// CreateUser 创建测试用户
// 支持自定义选项函数
func (f *UserFactory) CreateUser(opts ...func(*user_entity.User)) *user_entity.User {
	email, err := valueobject.NewEmail("test@example.com")
	require.NoError(f.t, err)

	nickname, err := valueobject.NewNickname("TestUser")
	require.NoError(f.t, err)

	user := &user_entity.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  user_entity.HashedPassword("$2a$12$..."),
		Nickname:  nickname,
		Status:    user_entity.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 应用自定义选项
	for _, opt := range opts {
		opt(user)
	}

	return user
}

// CreateTenantMember 创建测试租户成员
func (f *UserFactory) CreateTenantMember(tenantID, userID uuid.UUID, role user_entity.UserRole) *user_entity.TenantMember {
	return &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     role,
		Status:   user_entity.MemberStatusActive,
		JoinedAt: time.Now(),
	}
}

// ==================== 自定义选项函数 ====================

// WithEmail 设置邮箱
func WithEmail(email string) func(*user_entity.User) {
	return func(u *user_entity.User) {
		e, err := valueobject.NewEmail(email)
		if err == nil {
			u.Email = e
		}
	}
}

// WithNickname 设置昵称
func WithNickname(nickname string) func(*user_entity.User) {
	return func(u *user_entity.User) {
		n, err := valueobject.NewNickname(nickname)
		if err == nil {
			u.Nickname = n
		}
	}
}

// WithPassword 设置密码哈希
func WithPassword(password string) func(*user_entity.User) {
	return func(u *user_entity.User) {
		u.Password = user_entity.HashedPassword(password)
	}
}

// WithStatus 设置状态
func WithStatus(status user_entity.UserStatus) func(*user_entity.User) {
	return func(u *user_entity.User) {
		u.Status = status
	}
}

// WithID 设置 ID
func WithID(id uuid.UUID) func(*user_entity.User) {
	return func(u *user_entity.User) {
		u.ID = id
	}
}

// WithPhone 设置手机号
func WithPhone(phone string) func(*user_entity.User) {
	return func(u *user_entity.User) {
		u.Phone = &phone
	}
}

// WithBio 设置个人简介
func WithBio(bio string) func(*user_entity.User) {
	return func(u *user_entity.User) {
		u.Bio = &bio
	}
}

// ==================== Tenant 工厂 ====================

// TenantFactory 租户测试工厂
type TenantFactory struct {
	t *testing.T
}

// NewTenantFactory 创建租户工厂
func NewTenantFactory(t *testing.T) *TenantFactory {
	return &TenantFactory{t: t}
}

// CreateTenant 创建测试租户
func (f *TenantFactory) CreateTenant(opts ...func(*user_entity.Tenant)) *user_entity.Tenant {
	tenant := &user_entity.Tenant{
		ID:          uuid.New(),
		Name:        "Test Tenant",
		Description: "Test Description",
		MaxMembers:  10,
		ExpiredAt:   time.Now().AddDate(1, 0, 0), // 默认 1 年有效期
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 应用自定义选项
	for _, opt := range opts {
		opt(tenant)
	}

	return tenant
}

// WithName 设置租户名称
func WithName(name string) func(*user_entity.Tenant) {
	return func(t *user_entity.Tenant) {
		t.Name = name
	}
}

// WithDescription 设置描述
func WithDescription(desc string) func(*user_entity.Tenant) {
	return func(t *user_entity.Tenant) {
		t.Description = desc
	}
}

// WithMaxMembers 设置最大成员数
func WithMaxMembers(max int) func(*user_entity.Tenant) {
	return func(t *user_entity.Tenant) {
		t.MaxMembers = max
	}
}

// WithExpiredAt 设置过期时间
func WithExpiredAt(expiredAt time.Time) func(*user_entity.Tenant) {
	return func(t *user_entity.Tenant) {
		t.ExpiredAt = expiredAt
	}
}

// WithTenantID 设置 ID
func WithTenantID(id uuid.UUID) func(*user_entity.Tenant) {
	return func(t *user_entity.Tenant) {
		t.ID = id
	}
}
