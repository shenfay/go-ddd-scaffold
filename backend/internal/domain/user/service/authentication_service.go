package service

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
)

// 类型别名和错误码
const (
	CodeInvalidPassword = user.CodeInvalidPassword
	CodeUserCannotLogin = user.CodeUserCannotLogin
	CodeAccountLocked   = user.CodeAccountLocked
)

type UserID = user.UserID

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

// AuthenticateParams 用户认证参数（值对象）
// 封装认证所需的基本数据，不包含任何业务逻辑
type AuthenticateParams struct {
	Username  string
	Password  string
	IPAddress string
	UserAgent string
}

// AuthenticateResult 用户认证结果（值对象）
// 封装认证成功的结果数据
type AuthenticateResult struct {
	User       *user.User
	LoginStats *user.LoginStats
	Success    bool
}

// Authenticate 执行用户认证
// 返回认证结果和可能的错误
func (s *AuthenticationService) Authenticate(ctx context.Context, params AuthenticateParams) (*AuthenticateResult, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByUsername(ctx, params.Username)
	if err != nil {
		if err == common.ErrAggregateNotFound {
			return nil, common.NewBusinessError(CodeInvalidPassword, "用户名或密码错误")
		}
		return nil, err
	}

	// 2. 获取或创建登录统计
	loginStats, err := s.getOrCreateLoginStats(ctx, user.ID().(UserID))
	if err != nil {
		return nil, err
	}

	// 3. 检查是否可以登录（用户状态 + 登录统计）
	if !user.CanLogin() {
		return nil, common.NewBusinessError(CodeUserCannotLogin, "用户无法登录")
	}
	if !loginStats.CanLogin() {
		return nil, common.NewBusinessError(CodeUserCannotLogin, "账户已被锁定")
	}

	// 4. 验证密码
	if !s.passwordHasher.Verify(params.Password, user.Password().Value()) {
		// 记录失败登录
		loginStats.RecordFailedLogin()
		// 检查是否需要锁定（连续 5 次失败）
		if loginStats.FailedAttempts() >= 5 {
			loginStats.Lock(30 * time.Minute)
		}
		// 保存登录统计
		if err := s.loginStatsRepo.Save(ctx, loginStats); err != nil {
			return nil, err
		}
		return nil, common.NewBusinessError(CodeInvalidPassword, "用户名或密码错误")
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
func (s *AuthenticationService) getOrCreateLoginStats(ctx context.Context, userID UserID) (*user.LoginStats, error) {
	loginStats, err := s.loginStatsRepo.FindByUserID(ctx, userID)
	if err != nil {
		if err == common.ErrAggregateNotFound {
			// 创建新的登录统计
			loginStats = user.NewLoginStats(userID)
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
func (s *AuthenticationService) LockAccount(ctx context.Context, userID UserID, duration time.Duration) error {
	loginStats, err := s.loginStatsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	loginStats.Lock(duration)
	return s.loginStatsRepo.Save(ctx, loginStats)
}

// UnlockAccount 解锁账户
func (s *AuthenticationService) UnlockAccount(ctx context.Context, userID UserID) error {
	loginStats, err := s.loginStatsRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	loginStats.Unlock()
	return s.loginStatsRepo.Save(ctx, loginStats)
}
