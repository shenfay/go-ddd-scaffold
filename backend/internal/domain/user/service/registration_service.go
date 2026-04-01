package service

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
)

// RegistrationService 用户注册领域服务
// 封装用户注册的业务逻辑，确保领域规则的完整性
type RegistrationService struct {
	userRepo       repository.UserRepository
	passwordHasher PasswordHasher
	passwordPolicy PasswordPolicy
	idGenerator    func() int64
}

// NewRegistrationService 创建注册领域服务
func NewRegistrationService(
	userRepo repository.UserRepository,
	passwordHasher PasswordHasher,
	passwordPolicy PasswordPolicy,
	idGenerator func() int64,
) *RegistrationService {
	return &RegistrationService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		passwordPolicy: passwordPolicy,
		idGenerator:    idGenerator,
	}
}

// RegisterUserParams 用户注册参数（值对象）
// 封装注册所需的基本数据，不包含任何业务逻辑
type RegisterUserParams struct {
	Username string
	Email    string
	Password string
}

// Register 执行用户注册
// 包含所有领域规则验证：唯一性检查、密码强度验证、密码哈希
func (s *RegistrationService) Register(ctx context.Context, params RegisterUserParams) (*user.User, error) {
	// 1. 验证密码强度（领域规则）
	if err := s.passwordPolicy.Validate(params.Password); err != nil {
		return nil, err
	}

	// 2. 检查用户名唯一性（领域规则）
	if err := s.ensureUsernameUnique(ctx, params.Username); err != nil {
		return nil, err
	}

	// 3. 检查邮箱唯一性（领域规则）
	if err := s.ensureEmailUnique(ctx, params.Email); err != nil {
		return nil, err
	}

	// 4. 哈希密码（领域服务）
	hashedPassword, err := s.passwordHasher.Hash(params.Password)
	if err != nil {
		return nil, common.NewBusinessError(common.CodeInternalError, "password hash failed")
	}

	// 5. 创建用户聚合根
	user, err := user.NewUser(params.Username, params.Email, hashedPassword, s.idGenerator)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ensureUsernameUnique 确保用户名唯一
func (s *RegistrationService) ensureUsernameUnique(ctx context.Context, username string) error {
	_, err := s.userRepo.FindByUsername(ctx, username)
	if err == nil {
		// 找到用户，说明用户名已存在
		return common.NewBusinessError(user.CodeUsernameExists, "用户名已存在")
	}
	if err != common.ErrAggregateNotFound {
		// 其他错误
		return err
	}
	// 未找到用户，用户名可用
	return nil
}

// ensureEmailUnique 确保邮箱唯一
func (s *RegistrationService) ensureEmailUnique(ctx context.Context, email string) error {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// 找到用户，说明邮箱已存在
		return common.NewBusinessError(user.CodeEmailExists, "邮箱已被注册")
	}
	if err != common.ErrAggregateNotFound {
		// 其他错误
		return err
	}
	// 未找到用户，邮箱可用
	return nil
}
