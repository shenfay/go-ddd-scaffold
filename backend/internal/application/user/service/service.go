// Package service 用户应用服务实现
//
// Deprecated: 该文件已废弃，请使用拆分后的服务：
// - AuthenticationService: 认证相关（Register, Login, Logout）
// - UserQueryService: 用户查询相关
// - UserCommandService: 用户命令相关
// - TenantService: 租户管理相关
package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/application/user/dto"
	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	user_event "go-ddd-scaffold/internal/domain/user/event"
	"go-ddd-scaffold/internal/domain/user/repository"
	errPkg "go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/validator"
	"go-ddd-scaffold/pkg/converter"
)

// EventBus 事件总线接口
type EventBus interface {
	Publish(event interface{}) error
}

// Service 用户应用服务实现（向后兼容的包装器）
type Service struct {
	userRepo         repository.UserRepository
	tenantRepo       repository.TenantRepository
	tenantMemberRepo repository.TenantMemberRepository
	jwtService       user_entity.JWTService
	eventBus         EventBus
	converter        converter.Converter
	userValidator    *validator.UserValidator
	authService      AuthenticationService // 认证服务委托
}

// NewService 创建用户服务实例
func NewService(
	userRepo repository.UserRepository,
	tenantRepo repository.TenantRepository,
	tenantMemberRepo repository.TenantMemberRepository,
	jwtService user_entity.JWTService,
	eventBus EventBus,
) *Service {
	svc := &Service{
		userRepo:         userRepo,
		tenantRepo:       tenantRepo,
		tenantMemberRepo: tenantMemberRepo,
		jwtService:       jwtService,
		eventBus:         eventBus,
		converter:        converter.NewConverter(),
	}
	// 初始化 User 校验器
	svc.userValidator = validator.NewUserValidator(nil)
	
	// 创建认证服务（用于 Logout 委托）
	svc.authService = NewAuthenticationService(userRepo, tenantRepo, tenantMemberRepo, jwtService, eventBus, nil)
	
	return svc
}

// Logout 用户登出（委托给 AuthenticationService）
func (s *Service) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	if s.authService != nil {
		return s.authService.Logout(ctx, userID, token)
	}
	return nil
}

// Register 用户注册
func (s *Service) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
	// 1. 业务校验：密码强度
	if err := validator.ValidatePasswordStrength(req.Password); err != nil {
		return nil, errPkg.ErrInvalidPassword
	}

	// 2. 检查邮箱是否已存在
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errPkg.ErrUserExists
	}

	// 3. 将 string 的 TenantID 转换为 uuid.UUID
	var tenantID *uuid.UUID
	if req.TenantID != nil {
		var err error
		tenantID, err = s.converter.ToUUIDPtr(*req.TenantID)
		if err != nil {
			return nil, err
		}
	}

	// 4. 验证租户限制（如果是成员注册）
	if req.Role != nil && *req.Role == "member" {
		if tenantID == nil {
			return nil, errPkg.ErrUnauthorized // 成员必须指定租户
		}

		// 检查租户总用户数限制
		count, err := s.tenantMemberRepo.ListByTenant(ctx, *tenantID)
		if err != nil {
			return nil, err
		}

		tenant, err := s.tenantRepo.GetByID(ctx, *tenantID)
		if err != nil {
			return nil, err
		}

		// 检查租户最大用户数限制
		if len(count) >= tenant.MaxMembers {
			return nil, errPkg.ErrTenantLimitExceed
		}
	}

	// 4. 密码加密
	hashedPassword, err := user_entity.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 5. 创建用户实体（仅包含基础信息，不包含角色和租户）
	newUser := &user_entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   user_entity.StatusActive,
	}

	// 6. 保存用户到仓储
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// 6. 如果指定了租户，则创建租户成员关系（默认角色为 member）
	if tenantID != nil {
		tenantMember := &user_entity.TenantMember{
			ID:       uuid.New(),
			TenantID: *tenantID,
			UserID:   newUser.ID,
			Role:     "member", // 默认角色
			Status:   user_entity.MemberStatusActive,
			JoinedAt: time.Now(),
		}

		if err := s.tenantMemberRepo.Create(ctx, tenantMember); err != nil {
			return nil, err
		}
	}

	// 7. 发布用户注册事件
	event := &user_event.UserRegisteredEvent{
		UserID:    newUser.ID,
		Email:     newUser.Email,
		TenantID:  tenantID,
		Timestamp: time.Now(),
	}
	s.eventBus.Publish(event)

	// 9. 转换为 DTO 返回
	userDTO := dto.ToUserDTO(newUser)
	if tenantID != nil {
		member, err := s.tenantMemberRepo.GetByUserAndTenant(ctx, newUser.ID, *tenantID)
		if err == nil {
			userDTO.Role = string(member.Role)
			userDTO.TenantID = s.converter.ToStringPtr(member.TenantID.String())
		}
	}
	return userDTO, nil
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. 查找用户
	userEntity, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errPkg.ErrUserNotFound
	}

	// 2. 验证密码
	if !user_entity.CheckPassword(req.Password, userEntity.Password) {
		return nil, errPkg.ErrInvalidPassword
	}

	// 3. 检查用户状态
	if userEntity.Status != user_entity.StatusActive {
		return nil, errPkg.ErrUnauthorized
	}

	// 4. 生成 JWT 令牌（仅包含用户 ID）
	token, err := s.jwtService.GenerateToken(userEntity.ID)
	if err != nil {
		return nil, err
	}

	// 6. 发布登录事件
	event := &user_event.UserLoggedInEvent{
		UserID:    userEntity.ID,
		Timestamp: time.Now(),
	}
	s.eventBus.Publish(event)

	// 7. 返回登录响应
	response := &dto.LoginResponse{
		User:        dto.ToUserDTO(userEntity),
		AccessToken: token,
	}

	return response, nil
}

// GetUser 获取用户信息
func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
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

// UpdateUser 更新用户信息
func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) error {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if req.Email != nil {
		userEntity.Email = *req.Email
	}
	if req.Password != nil {
		hashedPassword, err := user_entity.HashPassword(*req.Password)
		if err != nil {
			return err
		}
		userEntity.Password = hashedPassword
	}
	if req.Status != nil {
		userEntity.Status = user_entity.UserStatus(*req.Status)
	}

	return s.userRepo.Update(ctx, userEntity)
}

// GetUserInfo 获取当前用户信息（通过用户 ID）
func (s *Service) GetUserInfo(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToUserDTO(userEntity), nil
}

// UpdateProfile 更新用户个人资料
func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 更新昵称
	if req.Nickname != nil {
		userEntity.Nickname = *req.Nickname
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
func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

// CreateTenant 创建租户
func (s *Service) CreateTenant(ctx context.Context, req *dto.CreateTenantRequest, ownerID uuid.UUID) (*dto.Tenant, error) {
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
func (s *Service) GetTenant(ctx context.Context, tenantID uuid.UUID) (*dto.Tenant, error) {
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

// ListUsersByTenant 列出租户下所有用户
func (s *Service) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.User, error) {
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
func (s *Service) ListMembersByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.User, error) {
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

// CheckPermission 检查用户权限（基于角色的简单权限检查）
// 注意：这是简化的权限检查示例，实际项目中建议使用 Casbin 等权限框架
func (s *Service) CheckPermission(ctx context.Context, userID uuid.UUID, requiredRole user_entity.UserRole) bool {
	members, err := s.tenantMemberRepo.ListByUser(ctx, userID)
	if err != nil || len(members) == 0 {
		return false
	}

	for _, member := range members {
		// 超级管理员拥有所有权限
		if member.Role == user_entity.RoleSuperAdmin {
			return true
		}
		// 角色匹配检查
		if member.Role == requiredRole {
			return true
		}
		// 成员角色可以访问访客资源
		if member.Role == user_entity.RoleMember && requiredRole == user_entity.RoleGuest {
			return true
		}
	}

	return false
}
