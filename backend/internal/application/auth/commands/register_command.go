package commands

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// RegisterCommand 注册命令
type RegisterCommand struct {
	Username string
	Email    string
	Password string
}

// RegisterResult 注册结果
type RegisterResult struct {
	UserID   string
	Username string
	Email    string
}

// RegisterHandler 注册命令处理器
type RegisterHandler struct {
	userRepo       user.UserRepository
	passwordHasher user.PasswordHasher
	eventPublisher EventPublisher
}

// NewRegisterHandler 创建注册处理器
func NewRegisterHandler(
	userRepo user.UserRepository,
	passwordHasher user.PasswordHasher,
	eventPublisher EventPublisher,
) *RegisterHandler {
	return &RegisterHandler{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		eventPublisher: eventPublisher,
	}
}

// Handle 处理注册命令
func (h *RegisterHandler) Handle(ctx context.Context, cmd *RegisterCommand) (*RegisterResult, error) {
	// 1. 检查用户名是否已存在
	existingUser, _ := h.userRepo.FindByUsername(ctx, cmd.Username)
	if existingUser != nil {
		return nil, ddd.NewBusinessError("USERNAME_EXISTS", "用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	existingUser, _ = h.userRepo.FindByEmail(ctx, cmd.Email)
	if existingUser != nil {
		return nil, ddd.NewBusinessError("EMAIL_EXISTS", "邮箱已被注册")
	}

	// 3. 创建用户实体
	newUser, err := user.NewUser(cmd.Username, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	// 4. 保存用户
	if err := h.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	// 5. 发布领域事件
	events := newUser.GetUncommittedEvents()
	for _, event := range events {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	newUser.ClearUncommittedEvents()

	// 6. 返回结果
	return &RegisterResult{
		UserID:   newUser.ID().(user.UserID).String(),
		Username: newUser.Username().Value(),
		Email:    newUser.Email().Value(),
	}, nil
}
