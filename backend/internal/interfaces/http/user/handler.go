package user

import (
	"time"

	"github.com/gin-gonic/gin"
	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
)

// Handler HTTP处理器
type Handler struct {
	userService userApp.UserService
	mapper      *Mapper
	respHandler *httpShared.Handler
}

// NewHandler 创建处理器
func NewHandler(userService userApp.UserService) *Handler {
	return &Handler{
		userService: userService,
		mapper:      NewMapper(),
		respHandler: httpShared.NewHandler(nil),
	}
}

// @Summary 获取用户详情
// @Description 根据用户 ID 获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 200 {object} UserResponse "用户详情"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id} [get]
// GetUser 获取用户
func (h *Handler) GetUser(c *gin.Context) {
	var req GetUserRequest
	if !h.respHandler.BindUri(c, &req) {
		return
	}

	userID, err := h.mapper.ParseUserID(req.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	result, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, toUserResponse(result))
}

// @Summary 更新用户信息
// @Description 更新指定用户的详细信息（部分更新）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body UpdateUserRequest true "用户更新信息"
// @Success 200 {object} UserResponse "更新成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id} [put]
// UpdateUser 更新用户
func (h *Handler) UpdateUser(c *gin.Context) {
	var uriReq GetUserRequest
	if !h.respHandler.BindUri(c, &uriReq) {
		return
	}

	var bodyReq UpdateUserRequest
	if !h.respHandler.BindJSON(c, &bodyReq) {
		return
	}

	cmd, err := h.mapper.ToUpdateUserCommand(&bodyReq, uriReq.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	err = h.userService.UpdateUserProfile(c.Request.Context(), cmd)
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
// @Param request body ChangePasswordRequest true "原密码和新密码"
// @Success 204 {object} httpShared.APIResponse "修改成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id}/password [put]
// ChangePassword 修改密码
func (h *Handler) ChangePassword(c *gin.Context) {
	var uriReq GetUserRequest
	if !h.respHandler.BindUri(c, &uriReq) {
		return
	}

	var bodyReq ChangePasswordRequest
	if !h.respHandler.BindJSON(c, &bodyReq) {
		return
	}

	cmd, err := h.mapper.ToChangePasswordCommand(&bodyReq, uriReq.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	err = h.userService.ChangePassword(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}

// toUserResponse 将 Application Result 转换为 Response DTO
func toUserResponse(result *userApp.GetUserResult) *UserResponse {
	return &UserResponse{
		ID:          result.ID,
		Username:    result.Username,
		Email:       result.Email,
		DisplayName: stringPtr(result.DisplayName),
		FirstName:   stringPtr(result.FirstName),
		LastName:    stringPtr(result.LastName),
		Gender:      stringPtr(result.Gender),
		PhoneNumber: stringPtr(result.PhoneNumber),
		AvatarURL:   stringPtr(result.AvatarURL),
		Status:      result.Status,
		CreatedAt:   result.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   result.UpdatedAt.Format(time.RFC3339),
	}
}

// stringPtr 辅助函数：将 string 转换为 *string
func stringPtr(s string) *string {
	return &s
}
