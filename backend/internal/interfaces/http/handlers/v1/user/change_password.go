package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/usecase"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
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
// @Description 修改指定用户的登录密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body usecase.ChangePasswordCommand true "原密码和新密码"
// @Success 204 {object} handler.APIResponse "修改成功"
// @Failure 400 {object} handler.APIResponse "请求参数错误"
// @Failure 404 {object} handler.APIResponse "用户不存在"
// @Router /users/{id}/password [put]
func (h *ChangePasswordHandler) Handle(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	var cmd usecase.ChangePasswordCommand
	cmd.UserID = vo.NewUserID(userIDInt)

	err = h.uc.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}
