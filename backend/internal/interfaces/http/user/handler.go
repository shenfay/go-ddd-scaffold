package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/application/user/usecase"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
)

// Handler HTTP 处理器
type Handler struct {
	getUserUC        *usecase.GetUserUseCase
	updateProfileUC  *usecase.UpdateProfileUseCase
	changePasswordUC *usecase.ChangePasswordUseCase
	respHandler      *httpShared.Handler
}

// NewHandler 创建处理器
func NewHandler(
	getUserUC *usecase.GetUserUseCase,
	updateProfileUC *usecase.UpdateProfileUseCase,
	changePasswordUC *usecase.ChangePasswordUseCase,
	respHandler *httpShared.Handler,
) *Handler {
	return &Handler{
		getUserUC:        getUserUC,
		updateProfileUC:  updateProfileUC,
		changePasswordUC: changePasswordUC,
		respHandler:      respHandler,
	}
}

// @Summary 获取用户详情
// @Description 根据用户 ID 获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 200 {object} github_com_shenfay_go-ddd-scaffold_internal_domain_user_aggregate.User "用户详情"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id} [get]
// GetUser 获取用户
func (h *Handler) GetUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}
	userID := vo.NewUserID(userIDInt)

	result, err := h.getUserUC.Execute(c.Request.Context(), userID)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, result)
}

// @Summary 更新用户信息
// @Description 更新指定用户的详细信息（部分更新）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body github_com_shenfay_go-ddd-scaffold_internal_application_user_usecase.UpdateProfileCommand true "用户更新信息"
// @Success 200 {object} httpShared.APIResponse "更新成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id} [put]
// UpdateUser 更新用户
func (h *Handler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	var cmd usecase.UpdateProfileCommand
	cmd.UserID = vo.NewUserID(userIDInt)

	err = h.updateProfileUC.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, nil)
}

// @Summary 修改用户密码
// @Description 修改指定用户的登录密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body github_com_shenfay_go-ddd-scaffold_internal_application_user_usecase.ChangePasswordCommand true "原密码和新密码"
// @Success 204 {object} httpShared.APIResponse "修改成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id}/password [put]
// ChangePassword 修改密码
func (h *Handler) ChangePassword(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	var cmd usecase.ChangePasswordCommand
	cmd.UserID = vo.NewUserID(userIDInt)

	err = h.changePasswordUC.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}
