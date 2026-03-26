package auth

import (
	"github.com/gin-gonic/gin"

	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handler"
)

// RefreshTokenHandler 刷新令牌处理器
type RefreshTokenHandler struct {
	authService authApp.AuthService
	respHandler *handler.Handler
}

// NewRefreshTokenHandler 创建刷新令牌处理器
func NewRefreshTokenHandler(
	authService authApp.AuthService,
	respHandler *handler.Handler,
) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		authService: authService,
		respHandler: respHandler,
	}
}

// ServeHTTP 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body authApp.RefreshTokenCommand true "刷新令牌"
// @Success 200 {object} authApp.RefreshTokenResult "刷新成功返回新令牌"
// @Failure 400 {object} handler.APIResponse "请求参数错误"
// @Failure 401 {object} handler.APIResponse "令牌无效"
// @Router /auth/refresh [post]
func (h *RefreshTokenHandler) ServeHTTP(c *gin.Context) {
	var cmd authApp.RefreshTokenCommand
	if !h.respHandler.BindJSON(c, &cmd) {
		return
	}

	// 补充 IP 和 UserAgent 信息
	cmd.IPAddress = c.ClientIP()
	cmd.UserAgent = c.Request.UserAgent()

	result, err := h.authService.RefreshToken(c.Request.Context(), &cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, result)
}
