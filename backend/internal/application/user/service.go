package user

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
)

// UserService 用户应用服务接口（核心流程）
type UserService interface {
	// === 核心流程 ===
	Register(ctx context.Context, req *RegisterRequest) (*UserDTO, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthenticateUserResult, error)
	GetUserByID(ctx context.Context, userID vo.UserID) (*UserDTO, error)

	// === 辅助功能（可选）===
	UpdateProfile(ctx context.Context, userID vo.UserID, req *UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID vo.UserID, req *ChangePasswordRequest) error
}

// UserServiceImpl 用户应用服务实现
type UserServiceImpl struct {
	uow            application.UnitOfWork
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
	uow application.UnitOfWork,
	eventPublisher kernel.EventPublisher,
	passwordHasher service.PasswordHasher,
	passwordPolicy service.PasswordPolicy,
	tokenService auth.TokenService,
	idGenerator *snowflake.Node,
) *UserServiceImpl {
	// 从 UoW 获取仓储
	userRepo := uow.UserRepository()
	loginStatsRepo := uow.LoginStatsRepository()

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
		uow:             uow,
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
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID vo.UserID) (*UserDTO, error) {
	userRepo := s.uow.UserRepository()

	u, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO 返回
	return ConvertUserToDTO(u), nil
}

// Login 用户登录
// 使用 AuthenticationService 领域服务处理认证逻辑，并使用 UnitOfWork 管理事务
func (s *UserServiceImpl) Login(ctx context.Context, req *LoginRequest) (*AuthenticateUserResult, error) {
	var authResult *service.AuthenticateResult

	// 在事务中执行认证
	err := s.uow.Transaction(ctx, func(ctx context.Context) error {
		var err error

		// 1. 调用领域服务执行认证
		authResult, err = s.authSvc.Authenticate(ctx, service.AuthenticateRequest{
			Username:  req.Username,
			Password:  req.Password,
			IPAddress: req.IPAddress,
			UserAgent: req.UserAgent,
		})
		if err != nil {
			return err
		}

		u := authResult.User

		// 2. 保存用户（Repository 内部会保存事件）
		userRepo := s.uow.UserRepository()
		if err := userRepo.Save(ctx, u); err != nil {
			return err
		}

		// 3. 保存登录统计
		loginStatsRepo := s.uow.LoginStatsRepository()
		if err := loginStatsRepo.Save(ctx, authResult.LoginStats); err != nil {
			return err // 回滚整个事务
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	u := authResult.User

	// 4. 生成 JWT Token
	tokenPair, err := s.tokenService.GenerateTokenPair(
		u.ID().(vo.UserID).Int64(),
		u.Username().Value(),
		u.Email().Value(),
	)
	if err != nil {
		return nil, err
	}

	// 5. 异步发布领域事件（事务成功后）
	events := u.GetUncommittedEvents()
	go s.publishEventsAsync(events)

	return &AuthenticateUserResult{
		UserID:       u.ID().(vo.UserID).Int64(),
		Username:     u.Username().Value(),
		Email:        u.Email().Value(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID vo.UserID, req *UpdateProfileRequest) error {
	userRepo := s.uow.UserRepository()

	u, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 更新用户信息
	if req.DisplayName != nil {
		u.SetDisplayName(*req.DisplayName)
	}
	if req.FirstName != nil {
		u.SetFirstName(*req.FirstName)
	}
	if req.LastName != nil {
		u.SetLastName(*req.LastName)
	}
	if req.Gender != nil {
		u.SetGender(*req.Gender)
	}
	if req.PhoneNumber != nil {
		u.SetPhoneNumber(*req.PhoneNumber)
	}

	// 保存用户（会自动发布事件）
	return userRepo.Save(ctx, u)
}

// ChangePassword 修改密码
func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID vo.UserID, req *ChangePasswordRequest) error {
	userRepo := s.uow.UserRepository()

	u, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 验证旧密码（使用 PasswordHasher）
	if !s.passwordHasher.Verify(req.OldPassword, u.Password().Value()) {
		return kernel.NewBusinessError(kernel.CodeInvalidOldPassword, "原密码错误")
	}

	// 验证新密码强度
	if err := s.passwordPolicy.Validate(req.NewPassword); err != nil {
		return err
	}

	// 修改密码
	if err := u.ChangePassword(req.NewPassword, req.IPAddress); err != nil {
		return err
	}

	// 保存用户（会自动发布事件）
	return userRepo.Save(ctx, u)
}

// Register 用户注册
// 使用 RegistrationService 领域服务处理注册逻辑，并使用 UnitOfWork 管理事务
func (s *UserServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*UserDTO, error) {
	var newUser *aggregate.User

	// 在事务中执行注册
	err := s.uow.Transaction(ctx, func(ctx context.Context) error {
		var err error

		// 1. 调用领域服务执行注册（领域逻辑）
		newUser, err = s.registrationSvc.Register(ctx, service.RegisterRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			return err
		}

		// 2. 保存用户（Repository 内部会保存事件）
		userRepo := s.uow.UserRepository()
		if err := userRepo.Save(ctx, newUser); err != nil {
			return err
		}

		// 3. 保存登录统计
		loginStatsRepo := s.uow.LoginStatsRepository()
		loginStats := aggregate.NewLoginStats(newUser.ID().(vo.UserID))
		if err := loginStatsRepo.Save(ctx, loginStats); err != nil {
			return err // 回滚整个事务
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. 异步发布领域事件（事务成功后）
	events := newUser.GetUncommittedEvents()
	go s.publishEventsAsync(events)

	// 5. 返回 DTO
	return ConvertUserToDTO(newUser), nil
}

// publishEventsAsync 异步发布领域事件
func (s *UserServiceImpl) publishEventsAsync(events []kernel.DomainEvent) {
	for _, event := range events {
		if err := s.eventPublisher.Publish(context.Background(), event); err != nil {
			// 记录日志，可以重试或发送到死信队列
			// TODO: 实现事件重试机制
		}
	}
}
