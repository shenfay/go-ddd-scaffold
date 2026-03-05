// Package service 用户命令应用服务实现
package service

import (
	"context"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	errPkg "go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/pkg/converter"
)

// userCommandService 用户命令服务实现
type userCommandService struct {
	userRepo         repository.UserRepository
	tenantMemberRepo repository.TenantMemberRepository
	converter        converter.Converter
}

// NewUserCommandService 创建用户命令服务实例
func NewUserCommandService(
	userRepo repository.UserRepository,
	tenantMemberRepo repository.TenantMemberRepository,
) UserCommandService {
	return &userCommandService{
		userRepo:         userRepo,
		tenantMemberRepo: tenantMemberRepo,
		converter:        converter.NewConverter(),
	}
}

// UpdateUser 更新用户信息
func (s *userCommandService) UpdateUser(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) error {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if req.Email != nil {
		email, err := valueobject.NewEmail(*req.Email)
		if err != nil {
			return errPkg.ErrInvalidEmail
		}
		userEntity.Email = email
	}
	if req.Password != nil {
		plainPassword, err := valueobject.NewPlainPassword(*req.Password)
		if err != nil {
			return errPkg.ErrInvalidPassword
		}
		hashedPassword, err := entity.NewHashedPassword(plainPassword.String())
		if err != nil {
			return err
		}
		userEntity.Password = hashedPassword
	}
	if req.Status != nil {
		userEntity.Status = entity.UserStatus(*req.Status)
	}

	return s.userRepo.Update(ctx, userEntity)
}

// UpdateProfile 更新用户个人资料
func (s *userCommandService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 更新昵称
	if req.Nickname != nil {
		nickname, err := valueobject.NewNickname(*req.Nickname)
		if err != nil {
			return err
		}
		userEntity.Nickname = nickname
	}

	// 更新手机号
	if req.Phone != nil {
		userEntity.Phone = req.Phone
	}

	// 更新个人简介
	if req.Bio != nil {
		userEntity.Bio = req.Bio
	}

	return s.userRepo.Update(ctx, userEntity)
}

// DeleteUser 删除用户
func (s *userCommandService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}
