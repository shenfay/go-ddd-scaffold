package commands

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID      user.UserID
	OldPassword string
	NewPassword string
	IPAddress   string
}

// ChangePasswordHandler 修改密码命令处理器
type ChangePasswordHandler struct {
	userRepo       user.UserRepository
	passwordHasher user.PasswordHasher
	eventPublisher EventPublisher
}

// NewChangePasswordHandler 创建命令处理器
func NewChangePasswordHandler(
	userRepo user.UserRepository,
	passwordHasher user.PasswordHasher,
	eventPublisher EventPublisher,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		eventPublisher: eventPublisher,
	}
}

// Handle 处理修改密码命令
func (h *ChangePasswordHandler) Handle(ctx context.Context, cmd *ChangePasswordCommand) error {
	u, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return ddd.ErrAggregateNotFound
	}

	// 修改密码（内部会验证旧密码）
	if err := u.ChangePassword(cmd.OldPassword, cmd.NewPassword, cmd.IPAddress); err != nil {
		return err
	}

	// 保存用户
	if err := h.userRepo.Save(ctx, u); err != nil {
		return err
	}

	// 发布领域事件
	events := u.GetUncommittedEvents()
	for _, event := range events {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断
		}
	}
	u.ClearUncommittedEvents()

	return nil
}
