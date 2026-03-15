package commands

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// RefreshTokenCommand 刷新令牌命令
type RefreshTokenCommand struct {
	RefreshToken string
	IPAddress    string
	UserAgent    string
}

// RefreshTokenResult 刷新令牌结果
type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// RefreshTokenHandler 刷新令牌处理器
type RefreshTokenHandler struct {
	userRepo     user.UserRepository
	tokenService auth.TokenService
}

// NewRefreshTokenHandler 创建刷新令牌处理器
func NewRefreshTokenHandler(
	userRepo user.UserRepository,
	tokenService auth.TokenService,
) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Handle 处理刷新令牌命令
func (h *RefreshTokenHandler) Handle(ctx context.Context, cmd *RefreshTokenCommand) (*RefreshTokenResult, error) {
	// 1. 验证刷新令牌并提取用户 ID
	claims, err := h.tokenService.ValidateToken(cmd.RefreshToken)
	if err != nil {
		return nil, ddd.NewBusinessError("INVALID_TOKEN", "无效的刷新令牌")
	}

	// 2. 查找用户
	foundUser, err := h.userRepo.FindByID(ctx, user.NewUserID(claims.UserID))
	if err != nil {
		return nil, ddd.NewBusinessError("USER_NOT_FOUND", "用户不存在")
	}

	// 3. 检查账户状态
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

	// 4. 生成新的令牌对
	tokenPair, err := h.tokenService.GenerateTokenPair(foundUser.ID().(user.UserID).Int64(), foundUser.Username().Value(), foundUser.Email().Value())
	if err != nil {
		return nil, ddd.NewBusinessError("TOKEN_GENERATION_FAILED", "令牌生成失败")
	}

	// 5. 返回结果
	expiresIn := int64(tokenPair.ExpiresAt.Sub(time.Now()).Seconds())

	return &RefreshTokenResult{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
