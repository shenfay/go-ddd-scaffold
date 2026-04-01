package usecase

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/model"
	userAggregate "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
)

// RegisterUserUseCase 注册用户用例
// 职责：编排用户注册的完整流程，保持单一职责和高可测试性
type RegisterUserUseCase struct {
	uow             application.UnitOfWorkWithEvents
	registrationSvc *service.RegistrationService
	activityLogRepo model.ActivityLogRepository
}

// NewRegisterUserUseCase 创建注册用户用例
func NewRegisterUserUseCase(
	uow application.UnitOfWorkWithEvents,
	registrationSvc *service.RegistrationService,
	activityLogRepo model.ActivityLogRepository,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		uow:             uow,
		registrationSvc: registrationSvc,
		activityLogRepo: activityLogRepo,
	}
}

// Execute 执行注册用户用例
// 优化点：
// 1. 无需手动发布事件，只需注册聚合根即可
// 2. ActivityLog 在事务内直接写入，保证审计可靠性
func (uc *RegisterUserUseCase) Execute(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error) {
	var newUser *userAggregate.User

	// 在事务中执行注册，并自动发布事件
	err := uc.uow.TransactionWithEvents(ctx, func(ctx context.Context) error {
		// 1. 调用领域服务执行注册（领域逻辑）
		var err error
		newUser, err = uc.registrationSvc.Register(ctx, service.RegisterUserParams{
			Username: cmd.Username,
			Email:    cmd.Email,
			Password: cmd.Password,
		})
		if err != nil {
			return err
		}

		// 2. ⚠️ 直接在事务内写入 ActivityLog（同步、可靠）
		//    关键点：ActivityLog 是审计日志，必须在事务内完成，不能通过事件异步处理
		auditLog := &model.ActivityLog{
			ID:     idgen.Generate(),
			UserID: newUser.ID().(vo.UserID).Int64(),
			Action: event.ActivityUserRegistered,
			Status: event.ActivityStatusSuccess,
			Metadata: map[string]interface{}{
				"username": newUser.Username().Value(),
				"email":    newUser.Email().Value(),
			},
			OccurredAt: time.Now(),
			CreatedAt:  time.Now(),
		}
		if err := uc.activityLogRepo.Save(ctx, auditLog); err != nil {
			return err
		}

		// 3. 注册聚合根以自动发布事件
		//    关键点：注册后，UnitOfWork 会在事务提交时自动发布所有领域事件
		uc.uow.TrackAggregate(newUser)

		// 4. 保存用户
		userRepo := uc.uow.UserRepository()
		if err := userRepo.Save(ctx, newUser); err != nil {
			return err
		}

		// 5. 保存登录统计
		loginStatsRepo := uc.uow.LoginStatsRepository()
		loginStats := userAggregate.NewLoginStats(newUser.ID().(vo.UserID))
		if err := loginStatsRepo.Save(ctx, loginStats); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 6. 返回结果
	return &RegisterUserResult{
		UserID:   newUser.ID().(vo.UserID).Int64(),
		Username: newUser.Username().Value(),
		Email:    newUser.Email().Value(),
	}, nil
}

// Benefits of V2:
//
// 1. Simplified Code:
//    - Removed manual event publishing loop (lines 80-85 in v1)
//    - Single line: uc.uow.TrackAggregate(newUser)
//
// 2. Better Separation of Concerns:
//    - UseCase focuses on business logic, not infrastructure concerns
//    - Event publishing is an infrastructure responsibility
//
// 3. Reduced Error-Prone:
//    - Cannot forget to publish events
//    - Cannot forget to clear uncommitted events
//    - Event publishing is guaranteed in transaction
//
// 4. Consistent Pattern:
//    - All UseCases can follow the same pattern
//    - Easier to maintain and test
//
// 5. Transaction Safety:
//    - Events are published only if transaction succeeds
//    - No partial event publishing
//
// Usage Example:
//
//     // Create infrastructure components
//     eventPublisher := asynq.NewEventPublisherAdapter(query, taskPublisher, logger)
//     uow := application.NewUnitOfWork(db, query,
//         application.WithEventPublisher(eventPublisher))
//
//     // Create use case
//     registrationSvc := service.NewRegistrationService(...)
//     useCase := usecase.NewRegisterUserUseCase(uow, registrationSvc)
//
//     // Execute (events are automatically published)
//     result, err := useCase.Execute(ctx, cmd)
