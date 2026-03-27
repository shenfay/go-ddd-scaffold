package auth

import (
	"github.com/gin-gonic/gin"

	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers"
)

// LoginHandler 登录处理器
type LoginHandler struct {
	authService authApp.AuthService
	respHandler *handlers.Handler
}

// NewLoginHandler 创建登录处理器
func NewLoginHandler(
	authService authApp.AuthService,
	respHandler *handlers.Handler,
) *LoginHandler {
	return &LoginHandler{
		authService: authService,
		respHandler: respHandler,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名或邮箱和密码进行登录，获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body authApp.AuthenticateCommand true "登录凭证"
// @Success 200 {object} authApp.AuthenticateResult "登录成功返回令牌和用户信息"
// @Failure 400 {object} handler.APIResponse "请求参数错误"
// @Failure 401 {object} handler.APIResponse "认证失败"
// @Router /auth/login [post]
func (h *LoginHandler) Handle(c *gin.Context) {
	var cmd authApp.AuthenticateCommand
	if !h.respHandler.BindJSON(c, &cmd) {
		return
	}

	// 补充 IP 和 UserAgent 信息
	cmd.IPAddress = c.ClientIP()
	cmd.UserAgent = c.Request.UserAgent()

	result, err := h.authService.AuthenticateUser(c.Request.Context(), &cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, result)
}
