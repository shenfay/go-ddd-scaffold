package auth

import (
	"github.com/gin-gonic/gin"

	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handler"
)

// GetCurrentUserHandler 获取当前用户处理器
type GetCurrentUserHandler struct {
	authService authApp.AuthService
	respHandler *handler.Handler
}

// NewGetCurrentUserHandler 创建获取当前用户处理器
func NewGetCurrentUserHandler(
	authService authApp.AuthService,
	respHandler *handler.Handler,
) *GetCurrentUserHandler {
	return &GetCurrentUserHandler{
		authService: authService,
		respHandler: respHandler,
	}
}

// ServeHTTP 获取当前登录用户的信息
// @Summary 获取当前用户
// @Description 获取当前登录用户的信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} authApp.UserInfoResult "当前用户信息"
// @Failure 401 {object} handler.APIResponse "未授权"
// @Router /auth/me [get]
func (h *GetCurrentUserHandler) ServeHTTP(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		h.respHandler.Error(c, kernel.NewBusinessError(kernel.CodeNotFound, "用户未找到"))
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		h.respHandler.Error(c, kernel.NewBusinessError(kernel.CodeInvalidUserID, "无效的用户 ID"))
		return
	}

	// 从仓储获取完整用户信息
	ctx := c.Request.Context()
	foundUser, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, foundUser)
}
