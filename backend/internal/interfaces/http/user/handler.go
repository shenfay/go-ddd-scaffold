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
func NewHandler(userService userApp.UserService, respHandler *httpShared.Handler) *Handler {
	return &Handler{
		userService: userService,
		mapper:      NewMapper(),
		respHandler: respHandler,
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

	userID, err := h.mapper.ParseUserID(uriReq.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	req := h.mapper.ToUpdateProfileRequest(&bodyReq)

	err = h.userService.UpdateProfile(c.Request.Context(), userID, req)
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

	userID, err := h.mapper.ParseUserID(uriReq.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	req := h.mapper.ToChangePasswordRequest(&bodyReq)

	err = h.userService.ChangePassword(c.Request.Context(), userID, req)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}

// toUserResponse 将 Application DTO 转换为 Response DTO
func toUserResponse(dto *userApp.UserDTO) *UserResponse {
	return &UserResponse{
		ID:          dto.ID,
		Username:    dto.Username,
		Email:       dto.Email,
		DisplayName: stringPtr(dto.DisplayName),
		FirstName:   nil, // UserDTO 中没有这些字段
		LastName:    nil,
		Gender:      nil,
		PhoneNumber: nil,
		AvatarURL:   nil,
		Status:      0, // 需要转换 status
		CreatedAt:   dto.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   "", // UserDTO 中没有 UpdatedAt
	}
}

// stringPtr 辅助函数：将 string 转换为 *string
func stringPtr(s string) *string {
	return &s
}
