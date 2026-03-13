package commands

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// LogoutCommand 登出命令
type LogoutCommand struct {
	UserID      int64
	AccessToken string
	IPAddress   string
	UserAgent   string
}

// LogoutResult 登出结果
type LogoutResult struct {
	Success bool
}

// LogoutHandler 登出命令处理器
type LogoutHandler struct {
	userRepo     user.UserRepository
	tokenService user.TokenService
}

// NewLogoutHandler 创建登出处理器
func NewLogoutHandler(
	userRepo user.UserRepository,
	tokenService user.TokenService,
) *LogoutHandler {
	return &LogoutHandler{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Handle 处理登出命令
func (h *LogoutHandler) Handle(ctx context.Context, cmd *LogoutCommand) (*LogoutResult, error) {
	// 1. 查找用户（验证用户存在）
	_, err := h.userRepo.FindByID(ctx, user.NewUserID(cmd.UserID))
	if err != nil {
		// 用户不存在也返回成功（幂等性）
		return &LogoutResult{Success: true}, nil
	}

	// 2. 可选：将令牌加入黑名单（如果使用 Redis）
	// 这里暂不实现，留给后续扩展
	// if err := h.tokenService.BlacklistToken(cmd.AccessToken); err != nil {
	//     return nil, err
	// }

	// 3. 返回成功
	return &LogoutResult{Success: true}, nil
}
