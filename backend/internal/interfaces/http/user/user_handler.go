package http

import (
	"net/http"

	dto"go-ddd-scaffold/internal/application/user/dto"
	user_service"go-ddd-scaffold/internal/application/user/service"
	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler 用户 HTTP 处理器（仅处理用户管理相关）
type UserHandler struct {
	userQueryService user_service.UserQueryService
	userCommandService user_service.UserCommandService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(
	userQueryService user_service.UserQueryService,
	userCommandService user_service.UserCommandService,
) *UserHandler {
	return &UserHandler{
		userQueryService:  userQueryService,
		userCommandService: userCommandService,
	}
}

// GetUser godoc
// @Summary 获取用户信息
// @Description 根据用户 ID 获取用户详细信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 200 {object} dto.User "用户信息"
// @Failure 400 {object} response.Response "无效的用户 ID"
// @Failure 404 {object} response.Response"用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router/api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(errors.InvalidParameter.WithDetails("无效的用户 ID 格式"))
		return
	}

	ctx := c.Request.Context()
	user, err := h.userQueryService.GetUser(ctx, userID)
	if err != nil {
		c.Error(err) // 统一交给中间件处理
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, user))
}

// UpdateUser godoc
// @Summary 更新用户信息
// @Description 更新指定用户的部分信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body dto.UpdateUserRequest true "更新信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response"请求参数错误"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router/api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(errors.InvalidParameter.WithDetails("无效的用户 ID 格式"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ValidationFailed.WithDetails(err.Error()))
		return
	}

	ctx := c.Request.Context()
	err = h.userCommandService.UpdateUser(ctx, userID, &req)
	if err != nil {
		c.Error(err) // 统一交给中间件处理
		return
	}

	c.JSON(http.StatusOK, response.OKWithMsg(ctx, nil, "用户信息更新成功"))
}
