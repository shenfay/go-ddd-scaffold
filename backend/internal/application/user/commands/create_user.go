package commands

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
	Username string
	Email    string
	Password string
}

// CreateUserResult 创建用户结果
type CreateUserResult struct {
	UserID    string
	Username  string
	Email     string
	Status    string
	CreatedAt string
}

// CreateUserHandler 创建用户命令处理器
type CreateUserHandler struct {
	userRepo       user.UserRepository
	passwordHasher user.PasswordHasher
	eventPublisher EventPublisher
}

// NewCreateUserHandler 创建命令处理器
func NewCreateUserHandler(
	userRepo user.UserRepository,
	passwordHasher user.PasswordHasher,
	eventPublisher EventPublisher,
) *CreateUserHandler {
	return &CreateUserHandler{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		eventPublisher: eventPublisher,
	}
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

// EventPublisher 事件发布接口
type EventPublisher interface {
	Publish(ctx context.Context, event ddd.DomainEvent) error
}

// Handle 处理创建用户命令
func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) (*CreateUserResult, error) {
	// 检查仓储是否初始化
	if h.userRepo == nil {
		return nil, ddd.NewBusinessError("SERVICE_UNAVAILABLE", "user repository not initialized")
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
	newUser, err := user.NewUser(cmd.Username, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	// 保存用户
	if err := h.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	// 发布领域事件
	events := newUser.GetUncommittedEvents()
	for _, event := range events {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断
		}
	}
	newUser.ClearUncommittedEvents()

	return &CreateUserResult{
		UserID:    newUser.ID().(user.UserID).String(),
		Username:  newUser.Username().Value(),
		Email:     newUser.Email().Value(),
		Status:    newUser.Status().String(),
		CreatedAt: newUser.CreatedAt().Format(time.RFC3339),
	}, nil
}
