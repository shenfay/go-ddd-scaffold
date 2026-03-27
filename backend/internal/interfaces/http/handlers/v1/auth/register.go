package auth

import (
	"github.com/gin-gonic/gin"

	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers"
)

// RegisterHandler 注册处理器
type RegisterHandler struct {
	authService authApp.AuthService
	respHandler *handlers.Handler
}

// NewRegisterHandler 创建注册处理器
func NewRegisterHandler(
	authService authApp.AuthService,
	respHandler *handlers.Handler,
) *RegisterHandler {
	return &RegisterHandler{
		authService: authService,
		respHandler: respHandler,
	}
}

// ServeHTTP 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body authApp.RegisterCommand true "用户注册信息"
// @Success 201 {object} authApp.RegisterResult "注册成功返回用户信息"
// @Failure 400 {object} handler.APIResponse "请求参数错误"
// @Failure 409 {object} handler.APIResponse "用户已存在"
// @Router /auth/register [post]
func (h *RegisterHandler) Handle(c *gin.Context) {
	var cmd authApp.RegisterCommand
	if !h.respHandler.BindJSON(c, &cmd) {
		return
	}

	result, err := h.authService.RegisterUser(c.Request.Context(), &cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Created(c, result)
}
