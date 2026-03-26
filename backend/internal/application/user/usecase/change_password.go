package usecase

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID      vo.UserID
	OldPassword string
	NewPassword string
	IPAddress   string
}

// ChangePasswordUseCase 修改密码用例
// 职责：编排密码修改流程，包含密码验证和强度检查
type ChangePasswordUseCase struct {
	uow            application.UnitOfWork
	passwordHasher service.PasswordHasher
	passwordPolicy service.PasswordPolicy
}

// NewChangePasswordUseCase 创建修改密码用例
func NewChangePasswordUseCase(
	uow application.UnitOfWork,
	passwordHasher service.PasswordHasher,
	passwordPolicy service.PasswordPolicy,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		uow:            uow,
		passwordHasher: passwordHasher,
		passwordPolicy: passwordPolicy,
	}
}

// Execute 执行修改密码用例
func (uc *ChangePasswordUseCase) Execute(ctx context.Context, cmd ChangePasswordCommand) error {
	userRepo := uc.uow.UserRepository()

	// 1. 查找用户
	u, err := userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 验证旧密码（使用 PasswordHasher）
	if !uc.passwordHasher.Verify(cmd.OldPassword, u.Password().Value()) {
		return kernel.NewBusinessError(aggregate.CodeInvalidOldPassword, "原密码错误")
	}

	// 3. 验证新密码强度
	if err := uc.passwordPolicy.Validate(cmd.NewPassword); err != nil {
		return err
	}

	// 4. 修改密码（领域方法）
	if err := u.ChangePassword(cmd.NewPassword, cmd.IPAddress); err != nil {
		return err
	}

	// 5. 保存用户（会自动发布事件）
	return userRepo.Save(ctx, u)
}
