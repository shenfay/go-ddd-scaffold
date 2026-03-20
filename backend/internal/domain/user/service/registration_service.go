package service

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
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

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

// Register 执行用户注册
// 包含所有领域规则验证：唯一性检查、密码强度验证、密码哈希
func (s *RegistrationService) Register(ctx context.Context, req RegisterRequest) (*aggregate.User, error) {
	// 1. 验证密码强度（领域规则）
	if err := s.passwordPolicy.Validate(req.Password); err != nil {
		return nil, err
	}

	// 2. 检查用户名唯一性（领域规则）
	if err := s.ensureUsernameUnique(ctx, req.Username); err != nil {
		return nil, err
	}

	// 3. 检查邮箱唯一性（领域规则）
	if err := s.ensureEmailUnique(ctx, req.Email); err != nil {
		return nil, err
	}

	// 4. 哈希密码（领域服务）
	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, kernel.NewBusinessError(kernel.CodeInternalError, "password hash failed")
	}

	// 5. 创建用户聚合根
	user, err := aggregate.NewUser(req.Username, req.Email, hashedPassword, s.idGenerator)
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
		return kernel.NewBusinessError(kernel.CodeUsernameExists, "用户名已存在")
	}
	if err != kernel.ErrAggregateNotFound {
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
		return kernel.NewBusinessError(kernel.CodeEmailExists, "邮箱已被注册")
	}
	if err != kernel.ErrAggregateNotFound {
		// 其他错误
		return err
	}
	// 未找到用户，邮箱可用
	return nil
}
