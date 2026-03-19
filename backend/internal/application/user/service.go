package user

import (
	"context"
	"errors"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
)

// UserService 用户应用服务接口（核心流程）
type UserService interface {
	// === 核心流程 ===
	RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*RegisterUserResult, error)
	AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticateUserResult, error)
	GetUserByID(ctx context.Context, userID user.UserID) (*GetUserResult, error)

	// === 辅助功能（可选）===
	UpdateUserProfile(ctx context.Context, cmd *UpdateUserProfileCommand) error
	ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error
}

// UserServiceImpl 用户应用服务实现
type UserServiceImpl struct {
	userRepo       user.UserRepository
	eventPublisher kernel.EventPublisher
	passwordHasher user.PasswordHasher
	passwordPolicy user.PasswordPolicy
	tokenService   auth.TokenService
	idGenerator    *snowflake.Node
}

// NewUserService 创建用户应用服务
func NewUserService(
	userRepo user.UserRepository,
	eventPublisher kernel.EventPublisher,
	passwordHasher user.PasswordHasher,
	passwordPolicy user.PasswordPolicy,
	tokenService auth.TokenService,
	idGenerator *snowflake.Node,
) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
		passwordHasher: passwordHasher,
		passwordPolicy: passwordPolicy,
		tokenService:   tokenService,
		idGenerator:    idGenerator,
	}
}

// ============================================================================
// Service Methods - 核心流程实现
// ============================================================================

// GetUserByID 根据 ID 获取用户
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID user.UserID) (*GetUserResult, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetUserResult{
		ID:          u.ID().(user.UserID).Int64(),
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
	}, nil
}

// AuthenticateUser 认证用户
func (s *UserServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthenticateUserResult, error) {
	// 1. 查找用户
	u, err := s.userRepo.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, kernel.ErrAggregateNotFound
	}

	// 2. 验证密码
	if !s.passwordHasher.Verify(cmd.Password, u.Password().Value()) {
		// 记录失败登录（可选）
		u.RecordFailedLogin(cmd.IPAddress, cmd.UserAgent, "invalid_password")
		_ = s.userRepo.Save(ctx, u)
		return nil, kernel.NewBusinessError(kernel.CodeInvalidPassword, "密码错误")
	}

	// 3. 检查用户是否可以登录
	if !u.CanLogin() {
		return nil, kernel.NewBusinessError(kernel.CodeUserCannotLogin, "用户无法登录")
	}

	// 4. 生成 JWT Token
	tokenPair, err := s.tokenService.GenerateTokenPair(
		u.ID().(user.UserID).Int64(),
		u.Username().Value(),
		u.Email().Value(),
	)
	if err != nil {
		return nil, err
	}

	// 5. 记录成功登录
	u.RecordLogin(cmd.IPAddress, cmd.UserAgent)
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	// 6. 发布领域事件
	events := u.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	u.ClearUncommittedEvents()

	return &AuthenticateUserResult{
		UserID:       u.ID().(user.UserID).Int64(),
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
func (s *UserServiceImpl) RegisterUser(ctx context.Context, cmd *RegisterUserCommand) (*RegisterUserResult, error) {
	// 1. 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, cmd.Username)
	if err != nil && !errors.Is(err, kernel.ErrAggregateNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, kernel.NewBusinessError(kernel.CodeUsernameExists, "用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	existingUser, err = s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, kernel.ErrAggregateNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, kernel.NewBusinessError(kernel.CodeEmailExists, "邮箱已被注册")
	}

	// 3. 验证密码强度
	if err := s.passwordPolicy.Validate(cmd.Password); err != nil {
		return nil, err
	}

	// 4. 哈希密码
	hashedPassword, err := s.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	// 4. 使用 Snowflake 生成唯一 ID
	userID, err := s.idGenerator.Generate()
	if err != nil {
		return nil, err
	}

	// 5. 创建用户实体
	newUser, err := user.NewUser(cmd.Username, cmd.Email, hashedPassword, func() int64 {
		return userID
	})
	if err != nil {
		return nil, err
	}

	// 5. 保存用户
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	// 6. 发布领域事件
	events := newUser.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	newUser.ClearUncommittedEvents()

	// 7. 返回结果 DTO
	return &RegisterUserResult{
		UserID:   newUser.ID().(user.UserID).Int64(),
		Username: newUser.Username().Value(),
		Email:    newUser.Email().Value(),
	}, nil
}
