package usecase

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	ports_auth "github.com/shenfay/go-ddd-scaffold/internal/application/ports/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/service"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// LoginUserCommand 登录用户命令
type LoginUserCommand struct {
	Username  string
	Password  string
	IPAddress string
	UserAgent string
}

// LoginUserResult 登录用户结果
type LoginUserResult struct {
	UserID       int64
	Username     string
	Email        string
	AccessToken  string
	RefreshToken string
}

// LoginUserUseCase 登录用户用例
// 职责：编排用户登录的完整流程，保持单一职责和高可测试性
type LoginUserUseCase struct {
	uow            application.UnitOfWork
	authSvc        *service.AuthenticationService
	tokenService   ports_auth.TokenService
	eventPublisher kernel.EventPublisher
}

// NewLoginUserUseCase 创建登录用户用例
func NewLoginUserUseCase(
	uow application.UnitOfWork,
	authSvc *service.AuthenticationService,
	tokenService ports_auth.TokenService,
	eventPublisher kernel.EventPublisher,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		uow:            uow,
		authSvc:        authSvc,
		tokenService:   tokenService,
		eventPublisher: eventPublisher,
	}
}

// Execute 执行登录用户用例
func (uc *LoginUserUseCase) Execute(ctx context.Context, cmd LoginUserCommand) (*LoginUserResult, error) {
	var authResult *service.AuthenticateResult

	// 在事务中执行认证
	err := uc.uow.Transaction(ctx, func(ctx context.Context) error {
		var err error

		// 1. 调用领域服务执行认证
		authResult, err = uc.authSvc.Authenticate(ctx, service.AuthenticateRequest{
			Username:  cmd.Username,
			Password:  cmd.Password,
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
		})
		if err != nil {
			return err
		}

		u := authResult.User

		// 2. 保存用户（Repository 内部会保存事件）
		userRepo := uc.uow.UserRepository()
		if err := userRepo.Save(ctx, u); err != nil {
			return err
		}

		// 3. 保存登录统计
		loginStatsRepo := uc.uow.LoginStatsRepository()
		if err := loginStatsRepo.Save(ctx, authResult.LoginStats); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	u := authResult.User

	// 4. 生成 JWT Token
	tokenPair, err := uc.tokenService.GenerateTokenPair(
		u.ID().(vo.UserID).Int64(),
		u.Username().Value(),
		u.Email().Value(),
	)
	if err != nil {
		return nil, err
	}

	// 5. 异步发布领域事件（事务成功后）
	events := u.GetUncommittedEvents()
	go uc.publishEventsAsync(events)

	return &LoginUserResult{
		UserID:       u.ID().(vo.UserID).Int64(),
		Username:     u.Username().Value(),
		Email:        u.Email().Value(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// publishEventsAsync 异步发布领域事件
func (uc *LoginUserUseCase) publishEventsAsync(events []kernel.DomainEvent) {
	for _, event := range events {
		if err := uc.eventPublisher.Publish(context.Background(), event); err != nil {
			// TODO: 实现事件重试机制和死信队列
		}
	}
}
