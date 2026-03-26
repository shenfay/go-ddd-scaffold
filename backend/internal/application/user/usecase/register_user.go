package usecase

import (
	"context"
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// RegisterUserCommand 注册用户命令
type RegisterUserCommand struct {
	Username string
	Email    string
	Password string
}

// RegisterUserResult 注册用户结果
type RegisterUserResult struct {
	UserID   int64
	Username string
	Email    string
}

// RegisterUserUseCase 注册用户用例
// 职责：编排用户注册的完整流程，保持单一职责和高可测试性
type RegisterUserUseCase struct {
	uow             application.UnitOfWork
	registrationSvc *service.RegistrationService
	eventPublisher  kernel.EventPublisher
}

// NewRegisterUserUseCase 创建注册用户用例
func NewRegisterUserUseCase(
	uow application.UnitOfWork,
	registrationSvc *service.RegistrationService,
	eventPublisher kernel.EventPublisher,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		uow:             uow,
		registrationSvc: registrationSvc,
		eventPublisher:  eventPublisher,
	}
}

// Execute 执行注册用户用例
func (uc *RegisterUserUseCase) Execute(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error) {
	var newUser *aggregate.User

	// 在事务中执行注册
	err := uc.uow.Transaction(ctx, func(ctx context.Context) error {
		// 1. 调用领域服务执行注册（领域逻辑）
		var err error
		newUser, err = uc.registrationSvc.Register(ctx, service.RegisterRequest{
			Username: cmd.Username,
			Email:    cmd.Email,
			Password: cmd.Password,
		})
		if err != nil {
			return err
		}

		// 2. 保存用户
		userRepo := uc.uow.UserRepository()
		if err := userRepo.Save(ctx, newUser); err != nil {
			return err
		}

		// 3. 保存登录统计
		loginStatsRepo := uc.uow.LoginStatsRepository()
		loginStats := aggregate.NewLoginStats(newUser.ID().(vo.UserID))
		if err := loginStatsRepo.Save(ctx, loginStats); err != nil {
			return err
		}

		// 4. 发布领域事件（在同一事务中）
		events := newUser.GetUncommittedEvents()
		for _, event := range events {
			if err := uc.eventPublisher.Publish(ctx, event); err != nil {
				return fmt.Errorf("failed to publish event: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 5. 返回结果
	return &RegisterUserResult{
		UserID:   newUser.ID().(vo.UserID).Int64(),
		Username: newUser.Username().Value(),
		Email:    newUser.Email().Value(),
	}, nil
}
