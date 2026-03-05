// Package service 认证应用服务实现
package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/domain/user/entity"
	user_event "go-ddd-scaffold/internal/domain/user/event"
	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	"go-ddd-scaffold/internal/infrastructure/auth"
	errPkg "go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/validator"
	"go-ddd-scaffold/pkg/converter"
)

// TokenBlacklistService Token 黑名单服务接口（导出供其他包使用）
type TokenBlacklistService = auth.TokenBlacklistService

// globalTokenBlacklistService 全局 Token 黑名单服务单例
var globalTokenBlacklistService TokenBlacklistService

// SetGlobalTokenBlacklistService 设置全局 Token 黑名单服务（由 main.go 初始化时调用）
func SetGlobalTokenBlacklistService(svc TokenBlacklistService) {
	globalTokenBlacklistService = svc
}

// GetTokenBlacklistService 获取全局 Token 黑名单服务
func GetTokenBlacklistService() TokenBlacklistService {
	return globalTokenBlacklistService
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(event interface{}) error
}

// authenticationService 认证服务实现
type authenticationService struct {
	userRepo         repository.UserRepository
	tenantRepo       repository.TenantRepository
	tenantMemberRepo repository.TenantMemberRepository
	jwtService       entity.JWTService
	eventBus         EventBus
	converter        converter.Converter
	userValidator    *validator.UserValidator
	tokenBlacklist   TokenBlacklistService // Token 黑名单服务
}

// NewAuthenticationService 创建认证服务实例
func NewAuthenticationService(
	userRepo repository.UserRepository,
	tenantRepo repository.TenantRepository,
	tenantMemberRepo repository.TenantMemberRepository,
	jwtService entity.JWTService,
	eventBus EventBus,
	tokenBlacklist TokenBlacklistService,
) AuthenticationService {
	svc := &authenticationService{
		userRepo:         userRepo,
		tenantRepo:       tenantRepo,
		tenantMemberRepo: tenantMemberRepo,
		jwtService:       jwtService,
		eventBus:         eventBus,
		converter:        converter.NewConverter(),
		tokenBlacklist:   tokenBlacklist,
	}
	svc.userValidator = validator.NewUserValidator(nil)
	return svc
}

// Register 用户注册
func (s *authenticationService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
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

	// 5. 验证密码强度并创建 HashedPassword
	plainPassword, err := valueobject.NewPlainPassword(req.Password)
	if err != nil {
		return nil, errPkg.ErrInvalidPassword
	}

	hashedPassword, err := entity.NewHashedPassword(plainPassword.String())
	if err != nil {
		return nil, err
	}

	// 6. 创建用户实体（仅包含基础信息，不包含角色和租户）
	email, _ := valueobject.NewEmail(req.Email)
	nickname, _ := valueobject.NewNickname(req.Nickname)

	newUser := &entity.User{
		ID:       uuid.New(),
		Email:    email,
		Password: hashedPassword,
		Nickname: nickname,
		Status:   entity.StatusActive,
	}

	// 7. 保存用户到仓储
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// 8. 如果指定了租户，则创建租户成员关系（默认角色为 member）
	if tenantID != nil {
		tenantMember := &entity.TenantMember{
			ID:       uuid.New(),
			TenantID: *tenantID,
			UserID:   newUser.ID,
			Role:     "member", // 默认角色
			Status:   entity.MemberStatusActive,
			JoinedAt: time.Now(),
		}

		if err := s.tenantMemberRepo.Create(ctx, tenantMember); err != nil {
			return nil, err
		}
	}

	// 9. 发布用户注册事件
	event := &user_event.UserRegisteredEvent{
		UserID:    newUser.ID,
		Email:     newUser.Email.String(),
		TenantID:  tenantID,
		Timestamp: time.Now(),
	}
	s.eventBus.Publish(event)

	// 10. 转换为 DTO 返回
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
func (s *authenticationService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. 查找用户
	userEntity, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errPkg.ErrUserNotFound
	}

	// 2. 验证密码
	if !userEntity.Password.Verify(req.Password) {
		return nil, errPkg.ErrInvalidPassword
	}

	// 3. 检查用户状态
	if userEntity.Status != entity.StatusActive {
		return nil, errPkg.ErrUnauthorized
	}

	// 4. 生成 JWT 令牌（仅包含用户 ID）
	token, err := s.jwtService.GenerateToken(userEntity.ID)
	if err != nil {
		return nil, err
	}

	// 5. 发布登录事件
	event := &user_event.UserLoggedInEvent{
		UserID:    userEntity.ID,
		Timestamp: time.Now(),
	}
	s.eventBus.Publish(event)

	// 6. 返回登录响应
	response := &dto.LoginResponse{
		User:        dto.ToUserDTO(userEntity),
		AccessToken: token,
	}

	return response, nil
}

// Logout 用户登出
func (s *authenticationService) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	// 如果配置了 Token 黑名单服务，将 token 加入黑名单
	if s.tokenBlacklist != nil && token != "" {
		// 先验证 token 获取过期时间
		_, err := s.jwtService.ValidateToken(token)
		if err == nil {
			// 计算 token 的过期时间（从 JWT claims 中获取）
			// 简化处理：使用配置的过期时间
			expireAt := time.Now().Add(24 * time.Hour) // 默认 24 小时

			// 加入黑名单
			err = s.tokenBlacklist.AddToBlacklist(ctx, token, expireAt)
			if err != nil {
				// 记录错误但不影响登出流程
				// TODO: 使用 logger 记录
			}
		}
	}

	return nil
}
