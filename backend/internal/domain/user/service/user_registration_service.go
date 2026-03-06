package service

import (
	"context"
	"time"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	errPkg "go-ddd-scaffold/internal/pkg/errors"
	errors "go-ddd-scaffold/internal/pkg/errors"

	"github.com/google/uuid"
)

// UserRegistrationService 用户注册领域服务
type UserRegistrationService struct {
	userRepo repository.UserRepository
	hasher   PasswordHasher
}

// NewUserRegistrationService 创建用户注册领域服务
func NewUserRegistrationService(
	userRepo repository.UserRepository,
	hasher PasswordHasher,
) *UserRegistrationService {
	return &UserRegistrationService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

// RegisterCommand 用户注册命令
type RegisterCommand struct {
	Email    string
	Password string
	Nickname string
	TenantID *uuid.UUID // 可选，如果要在租户上下文中注册
	Role     *string    // 可选，租户成员角色
}

// RegisterUser 注册用户（包含完整业务逻辑）
func (s *UserRegistrationService) RegisterUser(ctx context.Context, cmd RegisterCommand) (*entity.User, error) {
	// 1. 验证邮箱唯一性
	if err := s.validateEmailUnique(ctx, cmd.Email); err != nil {
		return nil, err
	}

	// 2. 验证密码强度
	if err := s.validatePasswordStrength(cmd.Password); err != nil {
		return nil, err
	}

	// 3. 验证昵称
	nickname, err := valueobject.NewNickname(cmd.Nickname)
	if err != nil {
		return nil, errors.InvalidParameter.WithDetails(err.Error())
	}

	// 4. 哈希密码
	hashedPwdStr, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return nil, errPkg.Wrap(err, "HASH_PASSWORD_FAILED", "密码加密失败")
	}

	// 5. 创建值对象
	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return nil, errors.InvalidParameter.WithDetails(err.Error())
	}

	// 6. 创建用户实体（使用工厂方法）
	user := &entity.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  entity.HashedPassword(hashedPwdStr),
		Nickname:  nickname,
		Status:    entity.StatusActive,
		CreatedAt: time.Now(),
	}

	// 7. 如果是租户成员注册，检查租户限制
	if cmd.TenantID != nil && cmd.Role != nil {
		if err := s.checkTenantLimit(ctx, *cmd.TenantID); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// validateEmailUnique 验证邮箱唯一性
func (s *UserRegistrationService) validateEmailUnique(ctx context.Context, email string) error {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// 如果是 NotFound 错误，说明邮箱可用
		if err == errPkg.ErrUserNotFound {
			return nil
		}
		return err
	}

	if existing != nil {
		return errPkg.ErrUserExists
	}

	return nil
}

// validatePasswordStrength 验证密码强度
func (s *UserRegistrationService) validatePasswordStrength(password string) error {
	plainPwd, err := valueobject.NewPlainPassword(password)
	if err != nil {
		return errPkg.ErrInvalidPassword.WithDetails(err.Error())
	}

	if !plainPwd.IsValid() {
		return errPkg.ErrInvalidPassword.WithDetails("密码必须包含字母和数字，长度 6-64 位")
	}

	return nil
}

// checkTenantLimit 检查租户成员数量限制
func (s *UserRegistrationService) checkTenantLimit(ctx context.Context, tenantID uuid.UUID) error {
	// TODO: 需要租户仓储来实现
	// tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	// if err != nil {
	//     return err
	// }
	//
	// count, err := s.tenantMemberRepo.CountByTenant(ctx, tenantID)
	// if err != nil {
	//     return err
	// }
	//
	// if tenant.MaxMembers > 0 && len(count) >= tenant.MaxMembers {
	//     return errPkg.ErrTenantLimitExceed
	// }

	return nil
}
