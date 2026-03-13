package commands

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// UpdateUserCommand 更新用户命令
type UpdateUserCommand struct {
	UserID      user.UserID
	DisplayName *string
	FirstName   *string
	LastName    *string
	Gender      *user.UserGender
	PhoneNumber *string
	AvatarURL   *string
}

// UpdateUserResult 更新用户结果
type UpdateUserResult struct {
	UserID      string
	Username    string
	Email       string
	DisplayName string
	FirstName   string
	LastName    string
	Gender      string
	PhoneNumber string
	AvatarURL   string
	Status      string
	UpdatedAt   string
}

// UpdateUserHandler 更新用户命令处理器
type UpdateUserHandler struct {
	userRepo       user.UserRepository
	eventPublisher EventPublisher
}

// NewUpdateUserHandler 创建命令处理器
func NewUpdateUserHandler(
	userRepo user.UserRepository,
	eventPublisher EventPublisher,
) *UpdateUserHandler {
	return &UpdateUserHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle 处理更新用户命令
func (h *UpdateUserHandler) Handle(ctx context.Context, cmd *UpdateUserCommand) (*UpdateUserResult, error) {
	u, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, ddd.ErrAggregateNotFound
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
	if cmd.AvatarURL != nil {
		u.SetAvatarURL(*cmd.AvatarURL)
	}

	// 保存用户
	if err := h.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	// 发布领域事件
	events := u.GetUncommittedEvents()
	for _, event := range events {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断
		}
	}
	u.ClearUncommittedEvents()

	return &UpdateUserResult{
		UserID:      u.ID().(user.UserID).String(),
		Username:    u.Username().Value(),
		Email:       u.Email().Value(),
		DisplayName: u.DisplayName(),
		FirstName:   u.FirstName(),
		LastName:    u.LastName(),
		Gender:      u.Gender().String(),
		PhoneNumber: u.PhoneNumber(),
		AvatarURL:   u.AvatarURL(),
		Status:      u.Status().String(),
		UpdatedAt:   u.UpdatedAt().Format(time.RFC3339),
	}, nil
}
