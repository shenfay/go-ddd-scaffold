// Package service 用户查询应用服务实现
package service

import (
	"context"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/application/user/dto"
	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/pkg/converter"
)

// userQueryService 用户查询服务实现
type userQueryService struct {
	userRepo         repository.UserRepository
	tenantMemberRepo repository.TenantMemberRepository
	converter        converter.Converter
}

// NewUserQueryService 创建用户查询服务实例
func NewUserQueryService(
	userRepo repository.UserRepository,
	tenantMemberRepo repository.TenantMemberRepository,
) UserQueryService {
	return &userQueryService{
		userRepo:         userRepo,
		tenantMemberRepo: tenantMemberRepo,
		converter:        converter.NewConverter(),
	}
}

// GetUser 根据 ID 获取用户信息
func (s *userQueryService) GetUser(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userDTO := dto.ToUserDTO(userEntity)

	// 尝试获取租户成员信息
	members, err := s.tenantMemberRepo.ListByUser(ctx, userID)
	if err == nil && len(members) > 0 {
		member := members[0]
		userDTO.Role = string(member.Role)
		userDTO.TenantID = s.converter.ToStringPtr(member.TenantID.String())
	}

	return userDTO, nil
}

// GetUserInfo 获取当前用户信息（通过用户 ID）
func (s *userQueryService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToUserDTO(userEntity), nil
}

// ListUsersByTenant 列出租户下所有用户
func (s *userQueryService) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.User, error) {
	members, err := s.tenantMemberRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.User, 0, len(members))
	for _, member := range members {
		userEntity, err := s.userRepo.GetByID(ctx, member.UserID)
		if err != nil {
			continue
		}

		userDTO := dto.ToUserDTO(userEntity)
		userDTO.Role = string(member.Role)
		userDTO.TenantID = s.converter.ToStringPtr(member.TenantID.String())

		dtos = append(dtos, userDTO)
	}
	return dtos, nil
}

// ListMembersByTenant 列出租户下的所有成员
func (s *userQueryService) ListMembersByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.User, error) {
	members, err := s.tenantMemberRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var users []*dto.User
	for _, member := range members {
		if member.Status != user_entity.MemberStatusActive {
			continue
		}

		userEntity, err := s.userRepo.GetByID(ctx, member.UserID)
		if err != nil {
			continue
		}

		userDTO := dto.ToUserDTO(userEntity)
		userDTO.Role = string(member.Role)
		userDTO.TenantID = s.converter.ToStringPtr(member.TenantID.String())
		users = append(users, userDTO)
	}

	return users, nil
}
