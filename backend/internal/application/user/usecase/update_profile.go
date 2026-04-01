package usecase

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/model"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// UpdateProfileUseCase 更新用户资料用例
// 职责：编排用户资料更新流程，保持单一职责和高可测试性
type UpdateProfileUseCase struct {
	uow            application.UnitOfWorkWithEvents
	activityLogger *activityLogger
}

// NewUpdateProfileUseCase 创建更新用户资料用例
func NewUpdateProfileUseCase(
	uow application.UnitOfWorkWithEvents,
	activityLogRepo model.ActivityLogRepository,
) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		uow:            uow,
		activityLogger: newActivityLogger(activityLogRepo),
	}
}

// Execute 执行更新用户资料用例
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, cmd UpdateProfileCommand) (*UpdateProfileResult, error) {
	userRepo := uc.uow.UserRepository()

	// 1. 查找用户
	u, err := userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, common.ErrAggregateNotFound
	}

	// 2. 更新用户信息（使用领域方法）
	if cmd.DisplayName != nil {
		if err := u.SetDisplayName(*cmd.DisplayName); err != nil {
			return nil, err
		}
	}
	if cmd.FirstName != nil {
		if err := u.SetFirstName(*cmd.FirstName); err != nil {
			return nil, err
		}
	}
	if cmd.LastName != nil {
		if err := u.SetLastName(*cmd.LastName); err != nil {
			return nil, err
		}
	}
	if cmd.Gender != nil {
		if err := u.SetGender(*cmd.Gender); err != nil {
			return nil, err
		}
	}
	if cmd.PhoneNumber != nil {
		if err := u.SetPhoneNumber(*cmd.PhoneNumber); err != nil {
			return nil, err
		}
	}

	// 3. 在事务中保存用户并自动发布事件
	err = uc.uow.TransactionWithEvents(ctx, func(ctx context.Context) error {
		// 4. ⚠️ 直接在事务内写入 ActivityLog（同步、可靠）
		if err := uc.activityLogger.LogUserAction(ctx, u.ID().(vo.UserID), event.ActivityUserProfileUpdated, map[string]interface{}{
			"display_name": u.DisplayName(),
			"first_name":   u.FirstName(),
			"last_name":    u.LastName(),
		}); err != nil {
			return err
		}

		// 5. 注册聚合根以自动发布事件
		uc.uow.TrackAggregate(u)

		// 6. 保存用户
		return userRepo.Save(ctx, u)
	})

	if err != nil {
		return nil, err
	}

	return &UpdateProfileResult{Success: true}, nil
}
