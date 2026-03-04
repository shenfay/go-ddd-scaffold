package service

import (
	"context"
	"time"

	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	user_repo "go-ddd-scaffold/internal/domain/user/repository"
	auth2 "go-ddd-scaffold/internal/infrastructure/auth"

	"github.com/google/uuid"
)

// TenantService 租户服务接口
type TenantService interface {
	// GetUserTenants 获取用户的所有租户
	GetUserTenants(ctx context.Context, userID uuid.UUID) ([]*TenantWithRole, error)
	// CreateTenant 创建租户
	CreateTenant(ctx context.Context, name, description string, ownerID uuid.UUID) (*user_entity.Tenant, error)
}

// TenantWithRole 租户及用户在该租户的角色
type TenantWithRole struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Role        string    `json:"role"` // 用户在该租户的角色
	JoinedAt    time.Time `json:"joinedAt"`
}

type tenantService struct {
	tenantRepo    user_repo.TenantRepository
	memberRepo    user_repo.TenantMemberRepository
	casbinService auth2.CasbinService
}

// NewTenantService 创建租户服务实例
func NewTenantService(
	tenantRepo user_repo.TenantRepository,
	memberRepo user_repo.TenantMemberRepository,
	casbinService auth2.CasbinService,
) TenantService {
	return &tenantService{
		tenantRepo:    tenantRepo,
		memberRepo:    memberRepo,
		casbinService: casbinService,
	}
}

// GetUserTenants 获取用户的所有租户
func (s *tenantService) GetUserTenants(ctx context.Context, userID uuid.UUID) ([]*TenantWithRole, error) {
	// 1. 获取用户的所有租户成员关系
	memberships, err := s.memberRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 获取每个租户的详细信息
	result := make([]*TenantWithRole, 0, len(memberships))
	for _, membership := range memberships {
		tenant, err := s.tenantRepo.GetByID(ctx, membership.TenantID)
		if err != nil {
			continue // 跳过不存在的租户
		}

		result = append(result, &TenantWithRole{
			ID:          tenant.ID.String(),
			Name:        tenant.Name,
			Description: tenant.Description,
			Role:        string(membership.Role),
			JoinedAt:    membership.JoinedAt,
		})
	}

	return result, nil
}

// CreateTenant 创建租户
func (s *tenantService) CreateTenant(ctx context.Context, name, description string, ownerID uuid.UUID) (*user_entity.Tenant, error) {
	tenant := &user_entity.Tenant{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		MaxMembers:  10,                          // 默认最大成员数
		ExpiredAt:   time.Now().AddDate(1, 0, 0), // 默认一年有效期
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 1. 创建租户
	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, err
	}

	// 2. 创建者自动成为 owner
	member := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenant.ID,
		UserID:   ownerID,
		Role:     user_entity.RoleMember, // 默认使用 member 角色
		Status:   user_entity.MemberStatusActive,
		JoinedAt: time.Now(),
	}

	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	// 3. 在 Casbin 中添加角色（使用 AddRoleForUser 方法）
	if err := s.casbinService.AddRoleForUser(ownerID, tenant.ID, string(user_entity.RoleMember)); err != nil {
		return nil, err
	}

	return tenant, nil
}
