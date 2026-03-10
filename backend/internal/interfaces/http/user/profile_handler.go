package http

import (
	"net/http"

	"go-ddd-scaffold/internal/application/user/dto"
	user_service"go-ddd-scaffold/internal/application/user/service"
	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProfileHandler 个人资料 HTTP 处理器
type ProfileHandler struct {
	userQueryService   user_service.UserQueryService
	userCommandService user_service.UserCommandService
}

// NewProfileHandler 创建个人资料处理器实例
func NewProfileHandler(
	userQueryService user_service.UserQueryService,
	userCommandService user_service.UserCommandService,
) *ProfileHandler {
	return &ProfileHandler{
		userQueryService:   userQueryService,
		userCommandService: userCommandService,
	}
}

// GetUserInfo godoc
// @Summary 获取当前用户信息
// @Description 获取登录用户的详细信息
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.User "用户信息"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router/api/users/info [get]
func (h *ProfileHandler) GetUserInfo(c *gin.Context) {
	// 从 JWT token 中获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(errors.ErrUnauthorized.WithDetails("user ID not found in token"))
		return
	}

	// 安全的类型断言
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.Error(errors.ErrUnauthorized.WithDetails("invalid user ID type in context"))
		return
	}

	ctx := c.Request.Context()
	user, err := h.userQueryService.GetUserInfo(ctx, userIDUUID)
	if err != nil {
		c.Error(err) // 统一交给中间件处理
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, user))
}

// UpdateProfile godoc
// @Summary 更新个人资料
// @Description 更新登录用户的个人资料（昵称、手机、简介）
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "个人资料信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response"未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router/api/users/profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ValidationFailed.WithDetails(err.Error()))
		return
	}

	// 从 JWT token 中获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(errors.ErrUnauthorized.WithDetails("user ID not found in token"))
		return
	}

	ctx := c.Request.Context()
	err := h.userCommandService.UpdateProfile(ctx, userID.(uuid.UUID), &req)
	if err != nil {
		c.Error(err) // 统一交给中间件处理
		return
	}

	c.JSON(http.StatusOK, response.OKWithMsg(ctx, nil, "个人资料更新成功"))
}
