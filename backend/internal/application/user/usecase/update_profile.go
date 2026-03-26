package usecase

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UpdateProfileCommand 更新用户资料命令
type UpdateProfileCommand struct {
	UserID      vo.UserID
	DisplayName *string
	FirstName   *string
	LastName    *string
	Gender      *vo.UserGender // 指针类型支持可选参数
	PhoneNumber *string
}

// UpdateProfileUseCase 更新用户资料用例
// 职责：编排用户资料更新流程，保持单一职责和高可测试性
type UpdateProfileUseCase struct {
	uow application.UnitOfWork
}

// NewUpdateProfileUseCase 创建更新用户资料用例
func NewUpdateProfileUseCase(uow application.UnitOfWork) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		uow: uow,
	}
}

// Execute 执行更新用户资料用例
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, cmd UpdateProfileCommand) error {
	userRepo := uc.uow.UserRepository()

	// 1. 查找用户
	u, err := userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 更新用户信息（使用领域方法）
	if cmd.DisplayName != nil {
		if err := u.SetDisplayName(*cmd.DisplayName); err != nil {
			return err
		}
	}
	if cmd.FirstName != nil {
		if err := u.SetFirstName(*cmd.FirstName); err != nil {
			return err
		}
	}
	if cmd.LastName != nil {
		if err := u.SetLastName(*cmd.LastName); err != nil {
			return err
		}
	}
	if cmd.Gender != nil {
		if err := u.SetGender(*cmd.Gender); err != nil {
			return err
		}
	}
	if cmd.PhoneNumber != nil {
		if err := u.SetPhoneNumber(*cmd.PhoneNumber); err != nil {
			return err
		}
	}

	// 3. 保存用户（会自动发布事件）
	return userRepo.Save(ctx, u)
}
