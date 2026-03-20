package user

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
)

// UserService 用户应用服务接口（核心流程）
type UserService interface {
	// === 核心流程 ===
	RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*RegisterUserResult, error)
	AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticateUserResult, error)
	GetUserByID(ctx context.Context, userID vo.UserID) (*GetUserResult, error)

	// === 辅助功能（可选）===
	UpdateUserProfile(ctx context.Context, cmd *UpdateUserProfileCommand) error
	ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error
}

// UserServiceImpl 用户应用服务实现
type UserServiceImpl struct {
	userRepo       repository.UserRepository
	loginStatsRepo repository.LoginStatsRepository
	eventPublisher kernel.EventPublisher
	passwordHasher service.PasswordHasher
	passwordPolicy service.PasswordPolicy
	tokenService   auth.TokenService
	idGenerator    *snowflake.Node
	// 领域服务
	registrationSvc *service.RegistrationService
	authSvc         *service.AuthenticationService
}

// NewUserService 创建用户应用服务
func NewUserService(
	userRepo repository.UserRepository,
	loginStatsRepo repository.LoginStatsRepository,
	eventPublisher kernel.EventPublisher,
	passwordHasher service.PasswordHasher,
	passwordPolicy service.PasswordPolicy,
	tokenService auth.TokenService,
	idGenerator *snowflake.Node,
) *UserServiceImpl {
	// 创建领域服务
	registrationSvc := service.NewRegistrationService(
		userRepo,
		passwordHasher,
		passwordPolicy,
		func() int64 {
			id, _ := idGenerator.Generate()
			return id
		},
	)
	authSvc := service.NewAuthenticationService(
		userRepo,
		loginStatsRepo,
		passwordHasher,
	)

	return &UserServiceImpl{
		userRepo:        userRepo,
		loginStatsRepo:  loginStatsRepo,
		eventPublisher:  eventPublisher,
		passwordHasher:  passwordHasher,
		passwordPolicy:  passwordPolicy,
		tokenService:    tokenService,
		idGenerator:     idGenerator,
		registrationSvc: registrationSvc,
		authSvc:         authSvc,
	}
}

// ============================================================================
// Service Methods - 核心流程实现
// ============================================================================

// GetUserByID 根据 ID 获取用户
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID vo.UserID) (*GetUserResult, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取登录统计信息
	loginStats, _ := s.loginStatsRepo.FindByUserID(ctx, userID)

	result := &GetUserResult{
		ID:          u.ID().(vo.UserID).Int64(),
		Username:    u.Username().Value(),
		Email:       u.Email().Value(),
		DisplayName: u.DisplayName(),
		FirstName:   u.FirstName(),
		LastName:    u.LastName(),
		Gender:      u.Gender().String(),
		PhoneNumber: u.PhoneNumber(),
		AvatarURL:   u.AvatarURL(),
		Status:      int32(u.Status()),
		CreatedAt:   u.CreatedAt(),
		UpdatedAt:   u.UpdatedAt(),
	}

	// 填充登录统计信息
	if loginStats != nil {
		result.LastLoginAt = loginStats.LastLoginAt()
		result.LoginCount = int32(loginStats.LoginCount())
		result.FailedAttempts = int32(loginStats.FailedAttempts())
		result.LockedUntil = loginStats.LockedUntil()
	}

	return result, nil
}

// AuthenticateUser 认证用户
// 使用 AuthenticationService 领域服务处理认证逻辑
func (s *UserServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticateUserResult, error) {
	// 1. 调用领域服务执行认证
	authResult, err := s.authSvc.Authenticate(ctx, service.AuthenticateRequest{
		Username:  cmd.Username,
		Password:  cmd.Password,
		IPAddress: cmd.IPAddress,
		UserAgent: cmd.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	u := authResult.User

	// 2. 生成 JWT Token
	tokenPair, err := s.tokenService.GenerateTokenPair(
		u.ID().(vo.UserID).Int64(),
		u.Username().Value(),
		u.Email().Value(),
	)
	if err != nil {
		return nil, err
	}

	// 3. 发布领域事件（登录事件已在领域服务中生成）
	events := u.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	u.ClearUncommittedEvents()

	return &AuthenticateUserResult{
		UserID:       u.ID().(vo.UserID).Int64(),
		Username:     u.Username().Value(),
		Email:        u.Email().Value(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// UpdateUserProfile 更新用户资料
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, cmd *UpdateUserProfileCommand) error {
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 更新用户信息
	if cmd.DisplayName != nil {
		u.SetDisplayName(*cmd.DisplayName)
	}
	if cmd.FirstName != nil {
		u.SetFirstName(*cmd.FirstName)
	}
	if cmd.LastName != nil {
		u.SetLastName(*cmd.LastName)
	}
	if cmd.Gender != nil {
		u.SetGender(*cmd.Gender)
	}
	if cmd.PhoneNumber != nil {
		u.SetPhoneNumber(*cmd.PhoneNumber)
	}

	// 保存用户（会自动发布事件）
	return s.userRepo.Save(ctx, u)
}

// ChangePassword 修改密码
func (s *UserServiceImpl) ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error {
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 验证旧密码（使用 PasswordHasher）
	if !s.passwordHasher.Verify(cmd.OldPassword, u.Password().Value()) {
		return kernel.NewBusinessError(kernel.CodeInvalidOldPassword, "原密码错误")
	}

	// 验证新密码强度
	if err := s.passwordPolicy.Validate(cmd.NewPassword); err != nil {
		return err
	}

	// 修改密码
	if err := u.ChangePassword(cmd.NewPassword, cmd.IPAddress); err != nil {
		return err
	}

	// 保存用户（会自动发布事件）
	return s.userRepo.Save(ctx, u)
}

// RegisterUser 注册用户
// 使用 RegistrationService 领域服务处理注册逻辑
func (s *UserServiceImpl) RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*RegisterUserResult, error) {
	// 1. 调用领域服务执行注册
	newUser, err := s.registrationSvc.Register(ctx, service.RegisterRequest{
		Username: cmd.Username,
		Email:    cmd.Email,
		Password: cmd.Password,
	})
	if err != nil {
		return nil, err
	}

	// 2. 保存用户
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	// 3. 创建对应的登录统计记录
	loginStats := aggregate.NewLoginStats(newUser.ID().(vo.UserID))
	if err := s.loginStatsRepo.Save(ctx, loginStats); err != nil {
		return nil, err
	}

	// 4. 发布领域事件
	events := newUser.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	newUser.ClearUncommittedEvents()

	// 5. 返回结果 DTO
	return &RegisterUserResult{
		UserID:   newUser.ID().(vo.UserID).Int64(),
		Username: newUser.Username().Value(),
		Email:    newUser.Email().Value(),
	}, nil
}
