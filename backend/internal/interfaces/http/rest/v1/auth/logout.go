package auth

import (
	"strings"

	"github.com/gin-gonic/gin"

	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handler"
)

// LogoutHandler 登出处理器
type LogoutHandler struct {
	authService authApp.AuthService
	respHandler *handler.Handler
}

// NewLogoutHandler 创建登出处理器
func NewLogoutHandler(
	authService authApp.AuthService,
	respHandler *handler.Handler,
) *LogoutHandler {
	return &LogoutHandler{
		authService: authService,
		respHandler: respHandler,
	}
}

// ServeHTTP 用户登出
// @Summary 用户登出
// @Description 使当前访问令牌失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 204 {object} handler.APIResponse "登出成功"
// @Failure 401 {object} handler.APIResponse "未授权"
// @Router /auth/logout [post]
func (h *LogoutHandler) ServeHTTP(c *gin.Context) {
	// 从上下文获取用户 ID（由认证中间件注入）
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

	// 从 Header 获取 Access Token
	authHeader := c.GetHeader("Authorization")
	accessToken := extractBearerToken(authHeader)

	cmd := &authApp.LogoutCommand{
		UserID:      userID,
		AccessToken: accessToken,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	}

	_, err := h.authService.Logout(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}

// extractBearerToken 从 Authorization Header 提取 Bearer Token
func extractBearerToken(authHeader string) string {
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}
