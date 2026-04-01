package user

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/usecase"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers"
)

// ChangePasswordHandler 修改密码处理器
type ChangePasswordHandler struct {
	uc          *usecase.ChangePasswordUseCase
	respHandler *handlers.Handler
}

// NewChangePasswordHandler 创建修改密码处理器
func NewChangePasswordHandler(
	uc *usecase.ChangePasswordUseCase,
	respHandler *handlers.Handler,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		uc:          uc,
		respHandler: respHandler,
	}
}

// ServeHTTP 修改用户密码
// @Summary 修改用户密码
// @Description 修改当前登录用户的密码（从 Token 中获取用户身份）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body usecase.ChangePasswordCommand true "原密码和新密码"
// @Success 204 {object} handler.APIResponse "修改成功"
// @Failure 400 {object} handler.APIResponse "请求参数错误"
// @Failure 401 {object} handler.APIResponse "未认证"
// @Failure 404 {object} handler.APIResponse "用户不存在"
// @Router /users/password [put]
func (h *ChangePasswordHandler) Handle(c *gin.Context) {
	// 从 JWT Token 中获取用户 ID（由认证中间件注入）
	userID, exists := c.Get("user_id")
	if !exists {
		h.respHandler.Unauthorized(c, "user not authenticated")
		return
	}

	var cmd usecase.ChangePasswordCommand
	cmd.UserID = vo.NewUserID(userID.(int64))

	result, err := h.uc.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	if !result.Success {
		h.respHandler.Error(c, fmt.Errorf("change password failed"))
		return
	}

	h.respHandler.NoContent(c)
}
