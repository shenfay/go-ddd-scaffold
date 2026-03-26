package usecase

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	ports_auth "github.com/shenfay/go-ddd-scaffold/internal/application/ports/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/aggregate"
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

// LoginUserUseCase 登录用户用例（优化版）
// 职责：编排用户登录的完整流程，保持单一职责和高可测试性
// 架构：使用 UnitOfWorkWithEvents 自动发布事件，ActivityLogWriter 写入审计日志
type LoginUserUseCase struct {
	uow          application.UnitOfWorkWithEvents
	authSvc      *service.AuthenticationService
	tokenService ports_auth.TokenService
	logWriter    *application.ActivityLogWriter
}

// NewLoginUserUseCase 创建登录用户用例（优化版）
func NewLoginUserUseCase(
	uow application.UnitOfWorkWithEvents,
	authSvc *service.AuthenticationService,
	tokenService ports_auth.TokenService,
	logWriter *application.ActivityLogWriter,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		uow:          uow,
		authSvc:      authSvc,
		tokenService: tokenService,
		logWriter:    logWriter,
	}
}

// Execute 执行登录用户用例
// 优化点：
// 1. 使用 UnitOfWorkWithEvents 自动发布事件，无需手动处理
// 2. ActivityLog 在事务内直接写入，保证审计可靠性
func (uc *LoginUserUseCase) Execute(ctx context.Context, cmd LoginUserCommand) (*LoginUserResult, error) {
	var authResult *service.AuthenticateResult

	// 在事务中执行认证，并自动发布事件
	err := uc.uow.TransactionWithEvents(ctx, func(ctx context.Context) error {
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

		// 2. ⚠️ 直接在事务内写入 ActivityLog（同步、可靠）
		//    关键点：ActivityLog 是审计日志，必须在事务内完成
		if err := uc.logWriter.WriteSuccess(
			ctx,
			u.ID().(vo.UserID).Int64(),
			aggregate.ActivityUserLoggedIn,
			map[string]interface{}{
				"username":   u.Username().Value(),
				"ip_address": cmd.IPAddress,
				"user_agent": cmd.UserAgent,
			},
		); err != nil {
			return err
		}

		// 3. 注册聚合根以自动发布事件
		uc.uow.TrackAggregate(u)

		// 4. 保存用户
		userRepo := uc.uow.UserRepository()
		if err := userRepo.Save(ctx, u); err != nil {
			return err
		}

		// 5. 保存登录统计
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

	// 6. 生成 JWT Token
	tokenPair, err := uc.tokenService.GenerateTokenPair(
		u.ID().(vo.UserID).Int64(),
		u.Username().Value(),
		u.Email().Value(),
	)
	if err != nil {
		return nil, err
	}

	return &LoginUserResult{
		UserID:       u.ID().(vo.UserID).Int64(),
		Username:     u.Username().Value(),
		Email:        u.Email().Value(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
