package commands

import (
	"context"

	domainUser "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	appUser "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// RegisterUserCommandHandler 用户注册命令处理器
type RegisterUserCommandHandler struct {
	userRepo       domainUser.UserRepository
	passwordService PasswordService
	eventPublisher EventPublisher
}

// NewRegisterUserCommandHandler 创建用户注册命令处理器
func NewRegisterUserCommandHandler(
	userRepo domainUser.UserRepository,
	passwordService PasswordService,
	eventPublisher EventPublisher,
) *RegisterUserCommandHandler {
	return &RegisterUserCommandHandler{
		userRepo:        userRepo,
		passwordService: passwordService,
		eventPublisher:  eventPublisher,
	}
}

// Handle 处理用户注册命令
func (h *RegisterUserCommandHandler) Handle(ctx context.Context, cmd *appUser.RegisterUserCommand) (*domainUser.User, error) {
	// 验证输入
	if err := h.validateRegisterCommand(cmd); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	if _, err := h.userRepo.FindByUsername(ctx, cmd.Username); err == nil {
		return nil, ddd.NewBusinessError("USERNAME_EXISTS", "username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := h.userRepo.FindByEmail(ctx, cmd.Email); err == nil {
		return nil, ddd.NewBusinessError("EMAIL_EXISTS", "email already exists")
	}

	// 创建用户
	newUser, err := domainUser.NewUser(cmd.Username, cmd.Email, cmd.Password)
	if err != nil {
		return nil, ddd.NewBusinessErrorWithDetails("USER_CREATION_FAILED", "failed to create user", err.Error())
	}

	// 保存用户
	if err := h.userRepo.Save(ctx, newUser); err != nil {
		return nil, ddd.NewBusinessErrorWithDetails("USER_SAVE_FAILED", "failed to save user", err.Error())
	}

	// 发布领域事件
	domainEvent := domainUser.NewUserRegisteredEvent(newUser.ID().(domainUser.UserID), cmd.Username, cmd.Email)
	if err := h.eventPublisher.Publish(ctx, domainEvent); err != nil {
		// 记录错误但不中断主流程
	}

	return newUser, nil
}

// validateRegisterCommand 验证注册命令
func (h *RegisterUserCommandHandler) validateRegisterCommand(cmd *appUser.RegisterUserCommand) error {
	validationErrors := &ddd.ValidationErrors{}

	if len(cmd.Username) < 3 {
		validationErrors.Add("username", "username must be at least 3 characters", cmd.Username)
	}

	if len(cmd.Password) < 8 {
		validationErrors.Add("password", "password must be at least 8 characters", "***")
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

// ActivateUserCommandHandler 用户激活命令处理器
type ActivateUserCommandHandler struct {
	userRepo       domainUser.UserRepository
	eventPublisher EventPublisher
}

// NewActivateUserCommandHandler 创建用户激活命令处理器
func NewActivateUserCommandHandler(
	userRepo domainUser.UserRepository,
	eventPublisher EventPublisher,
) *ActivateUserCommandHandler {
	return &ActivateUserCommandHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle 处理用户激活命令
func (h *ActivateUserCommandHandler) Handle(ctx context.Context, cmd *appUser.ActivateUserCommand) error {
	u, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return ddd.ErrAggregateNotFound
	}

	if err := u.Activate(); err != nil {
		return err
	}

	if err := h.userRepo.Save(ctx, u); err != nil {
		return err
	}

	// 发布领域事件
	domainEvent := domainUser.NewUserActivatedEvent(cmd.UserID)
	return h.eventPublisher.Publish(ctx, domainEvent)
}

// 依赖接口定义（这些应该放在 shared/kernel 中）
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) bool
	ValidatePasswordStrength(password string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event ddd.DomainEvent) error
}
