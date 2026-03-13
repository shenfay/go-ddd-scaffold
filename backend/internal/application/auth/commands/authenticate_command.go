package commands

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// AuthenticateCommand 认证命令
type AuthenticateCommand struct {
	Identifier string // 用户名或邮箱
	Password   string
	IPAddress  string
	UserAgent  string
}

// AuthenticateResult 认证结果
type AuthenticateResult struct {
	UserID       string
	Username     string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // 过期时间（秒）
}

// AuthenticateHandler 认证命令处理器
type AuthenticateHandler struct {
	userRepo       user.UserRepository
	passwordHasher user.PasswordHasher
	tokenService   user.TokenService
	eventPublisher EventPublisher
}

// EventPublisher 事件发布接口
type EventPublisher interface {
	Publish(ctx context.Context, event ddd.DomainEvent) error
}

// NewAuthenticateHandler 创建认证处理器
func NewAuthenticateHandler(
	userRepo user.UserRepository,
	passwordHasher user.PasswordHasher,
	tokenService user.TokenService,
	eventPublisher EventPublisher,
) *AuthenticateHandler {
	return &AuthenticateHandler{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
		eventPublisher: eventPublisher,
	}
}

// Handle 处理认证命令
func (h *AuthenticateHandler) Handle(ctx context.Context, cmd *AuthenticateCommand) (*AuthenticateResult, error) {
	// 1. 查找用户
	var foundUser *user.User
	var err error

	// 先尝试作为邮箱查找
	foundUser, err = h.userRepo.FindByEmail(ctx, cmd.Identifier)
	if err != nil {
		// 再尝试作为用户名查找
		foundUser, err = h.userRepo.FindByUsername(ctx, cmd.Identifier)
		if err != nil {
			return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "用户名或密码错误")
		}
	}

	// 2. 检查账户状态
	if !foundUser.CanLogin() {
		switch foundUser.Status() {
		case user.UserStatusInactive:
			return nil, ddd.NewBusinessError("ACCOUNT_DISABLED", "账户已被禁用")
		case user.UserStatusLocked:
			return nil, ddd.NewBusinessError("ACCOUNT_LOCKED", "账户已被锁定")
		default:
			return nil, ddd.NewBusinessError("ACCOUNT_CANNOT_LOGIN", "账户无法登录")
		}
	}

	// 3. 验证密码
	if !h.passwordHasher.Verify(cmd.Password, foundUser.Password().Value()) {
		// 记录失败登录
		foundUser.RecordFailedLogin(cmd.IPAddress, cmd.UserAgent, "invalid_password")
		_ = h.userRepo.Save(ctx, foundUser)
		return nil, ddd.NewBusinessError("INVALID_CREDENTIALS", "用户名或密码错误")
	}

	// 4. 生成令牌对
	tokenPair, err := h.tokenService.GenerateTokenPair(foundUser.ID().(user.UserID))
	if err != nil {
		return nil, ddd.NewBusinessError("TOKEN_GENERATION_FAILED", "令牌生成失败")
	}

	// 5. 记录成功登录
	foundUser.RecordLogin(cmd.IPAddress, cmd.UserAgent)
	if err := h.userRepo.Save(ctx, foundUser); err != nil {
		return nil, err
	}

	// 6. 发布领域事件
	events := foundUser.GetUncommittedEvents()
	for _, event := range events {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			// 记录错误但不中断主流程
		}
	}
	foundUser.ClearUncommittedEvents()

	// 7. 返回结果
	expiresIn := int64(tokenPair.ExpiresAt.Sub(time.Now()).Seconds())

	return &AuthenticateResult{
		UserID:       foundUser.ID().(user.UserID).String(),
		Username:     foundUser.Username().Value(),
		Email:        foundUser.Email().Value(),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
