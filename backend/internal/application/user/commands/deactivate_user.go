package commands

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// DeactivateUserCommand 禁用用户命令
type DeactivateUserCommand struct {
	UserID user.UserID
	Reason string
}

// DeactivateUserHandler 禁用用户命令处理器
type DeactivateUserHandler struct {
	userRepo       user.UserRepository
	eventPublisher EventPublisher
}

// NewDeactivateUserHandler 创建命令处理器
func NewDeactivateUserHandler(
	userRepo user.UserRepository,
	eventPublisher EventPublisher,
) *DeactivateUserHandler {
	return &DeactivateUserHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle 处理禁用用户命令
func (h *DeactivateUserHandler) Handle(ctx context.Context, cmd *DeactivateUserCommand) error {
	u, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return ddd.ErrAggregateNotFound
	}

	// 禁用用户
	if err := u.Deactivate(cmd.Reason); err != nil {
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
