package usecase

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
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

// ChangePasswordUseCase 修改密码用例（优化版）
// 职责：编排密码修改流程，包含密码验证和强度检查
// 架构：使用 UnitOfWorkWithEvents 自动发布事件，ActivityLogWriter 写入审计日志
type ChangePasswordUseCase struct {
	uow            application.UnitOfWorkWithEvents
	passwordHasher service.PasswordHasher
	passwordPolicy service.PasswordPolicy
	logWriter      *application.ActivityLogWriter
}

// NewChangePasswordUseCase 创建修改密码用例（优化版）
func NewChangePasswordUseCase(
	uow application.UnitOfWorkWithEvents,
	passwordHasher service.PasswordHasher,
	passwordPolicy service.PasswordPolicy,
	logWriter *application.ActivityLogWriter,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		uow:            uow,
		passwordHasher: passwordHasher,
		passwordPolicy: passwordPolicy,
		logWriter:      logWriter,
	}
}

// Execute 执行修改密码用例
// 优化点：
// 1. 使用 UnitOfWorkWithEvents 自动发布事件
// 2. ActivityLog 在事务内直接写入，保证审计可靠性
func (uc *ChangePasswordUseCase) Execute(ctx context.Context, cmd ChangePasswordCommand) error {
	userRepo := uc.uow.UserRepository()

	// 1. 查找用户
	u, err := userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 验证旧密码（使用 PasswordHasher）
	if !uc.passwordHasher.Verify(cmd.OldPassword, u.Password().Value()) {
		return kernel.NewBusinessError(kernel.CodeInvalidCredentials, "原密码错误")
	}

	// 3. 验证新密码强度
	if err := uc.passwordPolicy.Validate(cmd.NewPassword); err != nil {
		return err
	}

	// 4. 在事务中修改密码并自动发布事件
	err = uc.uow.TransactionWithEvents(ctx, func(ctx context.Context) error {
		// 5. ⚠️ 直接在事务内写入 ActivityLog（同步、可靠）
		if err := uc.logWriter.WriteSuccess(
			ctx,
			u.ID().(vo.UserID).Int64(),
			aggregate.ActivityUserPasswordChanged,
			map[string]interface{}{
				"ip_address": cmd.IPAddress,
			},
		); err != nil {
			return err
		}

		// 6. 修改密码（领域方法）
		if err := u.ChangePassword(cmd.NewPassword, cmd.IPAddress); err != nil {
			return err
		}

		// 7. 注册聚合根以自动发布事件
		uc.uow.TrackAggregate(u)

		// 8. 保存用户
		return userRepo.Save(ctx, u)
	})

	return err
}
