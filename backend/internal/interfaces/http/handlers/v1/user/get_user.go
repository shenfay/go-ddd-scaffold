package user

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/usecase"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers"
)

// GetUserHandler 获取用户处理器
type GetUserHandler struct {
	uc          *usecase.GetUserUseCase
	respHandler *handlers.Handler
}

// NewGetUserHandler 创建获取用户处理器
func NewGetUserHandler(
	uc *usecase.GetUserUseCase,
	respHandler *handlers.Handler,
) *GetUserHandler {
	return &GetUserHandler{
		uc:          uc,
		respHandler: respHandler,
	}
}

// ServeHTTP 获取用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息（从 Token 中获取用户身份）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} interface{} "用户详情"
// @Failure 401 {object} handler.APIResponse "未认证"
// @Failure 404 {object} handler.APIResponse "用户不存在"
// @Router /users/me [get]
func (h *GetUserHandler) Handle(c *gin.Context) {
	// 从 JWT Token 中获取用户 ID（由认证中间件注入）
	userID, exists := c.Get("user_id")
	if !exists {
		h.respHandler.Unauthorized(c, "user not authenticated")
		return
	}

	result, err := h.uc.Execute(c.Request.Context(), vo.NewUserID(userID.(int64)))
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, result)
}
