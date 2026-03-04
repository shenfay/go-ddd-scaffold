// Package service 租户应用服务实现
package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/application/user/dto"
	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/pkg/converter"
)

// tenantService 租户应用服务实现
type tenantService struct {
	tenantRepo       repository.TenantRepository
	tenantMemberRepo repository.TenantMemberRepository
	converter        converter.Converter
}

// NewTenantService 创建租户服务实例
func NewTenantService(
	tenantRepo repository.TenantRepository,
	tenantMemberRepo repository.TenantMemberRepository,
) TenantService {
	return &tenantService{
		tenantRepo:       tenantRepo,
		tenantMemberRepo: tenantMemberRepo,
		converter:        converter.NewConverter(),
	}
}

// CreateTenant 创建租户
func (s *tenantService) CreateTenant(ctx context.Context, req *dto.CreateTenantRequest, ownerID uuid.UUID) (*dto.Tenant, error) {
	tenant := &user_entity.Tenant{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: "",
		MaxMembers:  req.MaxMembers,
		ExpiredAt:   time.Now().AddDate(1, 0, 0),
	}

	if req.Description != nil {
		tenant.Description = *req.Description
	}

	err := s.tenantRepo.Create(ctx, tenant)
	if err != nil {
		return nil, err
	}

	// 自动添加创建者为租户成员（owner 角色）
	member := &user_entity.TenantMember{
		ID:       uuid.New(),
		TenantID: tenant.ID,
		UserID:   ownerID,
		Role:     user_entity.RoleOwner,
		Status:   user_entity.MemberStatusActive,
		JoinedAt: time.Now(),
	}

	if err := s.tenantMemberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	return dto.ToTenantDTO(tenant, 1), nil
}

// GetTenant 获取租户信息
func (s *tenantService) GetTenant(ctx context.Context, tenantID uuid.UUID) (*dto.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	members, err := s.tenantMemberRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	activeCount := int64(0)
	for _, member := range members {
		if member.Status == user_entity.MemberStatusActive {
			activeCount++
		}
	}

	return dto.ToTenantDTO(tenant, activeCount), nil
}
