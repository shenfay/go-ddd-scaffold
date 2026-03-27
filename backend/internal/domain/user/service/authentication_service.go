package service

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// AuthenticationService 用户认证领域服务
// 封装登录认证的业务逻辑，协调User和LoginStats两个聚合根
type AuthenticationService struct {
	userRepo       repository.UserRepository
	loginStatsRepo repository.LoginStatsRepository
	passwordHasher PasswordHasher
}

// NewAuthenticationService 创建认证领域服务
func NewAuthenticationService(
	userRepo repository.UserRepository,
	loginStatsRepo repository.LoginStatsRepository,
	passwordHasher PasswordHasher,
) *AuthenticationService {
	return &AuthenticationService{
		userRepo:       userRepo,
		loginStatsRepo: loginStatsRepo,
		passwordHasher: passwordHasher,
	}
}

// AuthenticateRequest 认证请求参数
type AuthenticateRequest struct {
	Username  string
	Password  string
	IPAddress string
	UserAgent string
}

// AuthenticateResult 认证结果
type AuthenticateResult struct {
	User       *aggregate.User
	LoginStats *aggregate.LoginStats
	Success    bool
}

// Authenticate 执行用户认证
// 返回认证结果和可能的错误
func (s *AuthenticationService) Authenticate(ctx context.Context, req AuthenticateRequest) (*AuthenticateResult, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if err == common.ErrAggregateNotFound {
			return nil, common.NewBusinessError(aggregate.CodeInvalidPassword, "用户名或密码错误")
		}
		return nil, err
	}

	// 2. 获取或创建登录统计
	loginStats, err := s.getOrCreateLoginStats(ctx, user.ID().(vo.UserID))
	if err != nil {
		return nil, err
	}

	// 3. 检查是否可以登录（用户状态 + 登录统计）
	if !user.CanLogin() {
		return nil, common.NewBusinessError(aggregate.CodeUserCannotLogin, "用户无法登录")
	}
	if !loginStats.CanLogin() {
		return nil, common.NewBusinessError(aggregate.CodeUserCannotLogin, "账户已被锁定")
	}

	// 4. 验证密码
	if !s.passwordHasher.Verify(req.Password, user.Password().Value()) {
		// 记录失败登录
		loginStats.RecordFailedLogin()
		// 检查是否需要锁定（连续5次失败）
		if loginStats.FailedAttempts() >= 5 {
			loginStats.Lock(30 * time.Minute)
		}
		// 保存登录统计
		if err := s.loginStatsRepo.Save(ctx, loginStats); err != nil {
			return nil, err
		}
		return nil, common.NewBusinessError(aggregate.CodeInvalidPassword, "用户名或密码错误")
	}

	// 5. 记录成功登录
	loginStats.RecordLogin()

	// 6. 保存登录统计
	if err := s.loginStatsRepo.Save(ctx, loginStats); err != nil {
		return nil, err
	}

	return &AuthenticateResult{
		User:       user,
		LoginStats: loginStats,
		Success:    true,
	}, nil
}

// getOrCreateLoginStats 获取或创建登录统计
func (s *AuthenticationService) getOrCreateLoginStats(ctx context.Context, userID vo.UserID) (*aggregate.LoginStats, error) {
	loginStats, err := s.loginStatsRepo.FindByUserID(ctx, userID)
	if err != nil {
		if err == common.ErrAggregateNotFound {
			// 创建新的登录统计
			loginStats = aggregate.NewLoginStats(userID)
			if err := s.loginStatsRepo.Save(ctx, loginStats); err != nil {
				return nil, err
			}
			return loginStats, nil
		}
		return nil, err
	}
	return loginStats, nil
}

// LockAccount 锁定账户
func (s *AuthenticationService) LockAccount(ctx context.Context, userID vo.UserID, duration time.Duration) error {
	loginStats, err := s.loginStatsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	loginStats.Lock(duration)
	return s.loginStatsRepo.Save(ctx, loginStats)
}

// UnlockAccount 解锁账户
func (s *AuthenticationService) UnlockAccount(ctx context.Context, userID vo.UserID) error {
	loginStats, err := s.loginStatsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	loginStats.Unlock()
	return s.loginStatsRepo.Save(ctx, loginStats)
}
