package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/usecase"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handler"
)

// UpdateProfileHandler 更新用户资料处理器
type UpdateProfileHandler struct {
	uc          *usecase.UpdateProfileUseCase
	respHandler *handler.Handler
}

// NewUpdateProfileHandler 创建更新用户资料处理器
func NewUpdateProfileHandler(
	uc *usecase.UpdateProfileUseCase,
	respHandler *handler.Handler,
) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		uc:          uc,
		respHandler: respHandler,
	}
}

// ServeHTTP 更新用户信息
// @Summary 更新用户信息
// @Description 更新指定用户的详细信息（部分更新）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body usecase.UpdateProfileCommand true "用户更新信息"
// @Success 200 {object} handler.APIResponse "更新成功"
// @Failure 400 {object} handler.APIResponse "请求参数错误"
// @Failure 404 {object} handler.APIResponse "用户不存在"
// @Router /users/{id} [put]
func (h *UpdateProfileHandler) ServeHTTP(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	var cmd usecase.UpdateProfileCommand
	cmd.UserID = vo.NewUserID(userIDInt)

	err = h.uc.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, nil)
}
