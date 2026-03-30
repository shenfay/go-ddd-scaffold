package usecase

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/model"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
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
	uow             application.UnitOfWorkWithEvents
	activityLogRepo model.ActivityLogRepository
}

// NewUpdateProfileUseCase 创建更新用户资料用例
func NewUpdateProfileUseCase(
	uow application.UnitOfWorkWithEvents,
	activityLogRepo model.ActivityLogRepository,
) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		uow:             uow,
		activityLogRepo: activityLogRepo,
	}
}

// Execute 执行更新用户资料用例
// 优化点：
// 1. 使用 UnitOfWorkWithEvents 自动发布事件
// 2. ActivityLog 在事务内直接写入，保证审计可靠性
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, cmd UpdateProfileCommand) error {
	userRepo := uc.uow.UserRepository()

	// 1. 查找用户
	u, err := userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return common.ErrAggregateNotFound
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

	// 3. 在事务中保存用户并自动发布事件
	err = uc.uow.TransactionWithEvents(ctx, func(ctx context.Context) error {
		// 4. ⚠️ 直接在事务内写入 ActivityLog（同步、可靠）
		auditLog := &model.ActivityLog{
			ID:     idgen.Generate(),
			UserID: u.ID().(vo.UserID).Int64(),
			Action: event.ActivityUserProfileUpdated,
			Status: event.ActivityStatusSuccess,
			Metadata: map[string]interface{}{
				"display_name": u.DisplayName(),
				"first_name":   u.FirstName(),
				"last_name":    u.LastName(),
			},
			OccurredAt: time.Now(),
			CreatedAt:  time.Now(),
		}
		if err := uc.activityLogRepo.Save(ctx, auditLog); err != nil {
			return err
		}

		// 5. 注册聚合根以自动发布事件
		uc.uow.TrackAggregate(u)

		// 6. 保存用户
		return userRepo.Save(ctx, u)
	})

	return err
}
