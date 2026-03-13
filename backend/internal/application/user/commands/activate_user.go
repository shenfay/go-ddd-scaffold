package commands

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// ActivateUserCommand 激活用户命令
type ActivateUserCommand struct {
	UserID user.UserID
}

// ActivateUserHandler 激活用户命令处理器
type ActivateUserHandler struct {
	userRepo       user.UserRepository
	eventPublisher EventPublisher
}

// NewActivateUserHandler 创建命令处理器
func NewActivateUserHandler(
	userRepo user.UserRepository,
	eventPublisher EventPublisher,
) *ActivateUserHandler {
	return &ActivateUserHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle 处理激活用户命令
func (h *ActivateUserHandler) Handle(ctx context.Context, cmd *ActivateUserCommand) error {
	u, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return ddd.ErrAggregateNotFound
	}

	// 激活用户
	if err := u.Activate(); err != nil {
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
