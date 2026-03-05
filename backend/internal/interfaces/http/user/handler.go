package http

import (
	"fmt"
	"net/http"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/application/user/service"
	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserHandler 用户 HTTP 处理器
type UserHandler struct {
	authService        service.AuthenticationService
	userQueryService   service.UserQueryService
	userCommandService service.UserCommandService
	tenantService      service.TenantService
	logger             *zap.Logger
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(
	authService service.AuthenticationService,
	userQueryService service.UserQueryService,
	userCommandService service.UserCommandService,
	tenantService service.TenantService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		authService:        authService,
		userQueryService:   userQueryService,
		userCommandService: userCommandService,
		tenantService:      tenantService,
		logger:             logger,
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
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(c.Request.Context(), errors.InvalidParameter.WithDetails("无效的用户 ID 格式")))
		return
	}

	ctx := c.Request.Context()
	user, err := h.userQueryService.GetUser(ctx, userID)
	if err != nil {
		if err == errors.ErrUserNotFound {
			c.JSON(http.StatusNotFound, response.Fail(ctx, errors.ErrUserNotFound))
			return
		}
		h.logger.Error("获取用户信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
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
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(c.Request.Context(), errors.InvalidParameter.WithDetails("无效的用户 ID 格式")))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidateErr(c.Request.Context(), err.Error()))
		return
	}

	ctx := c.Request.Context()
	err = h.userCommandService.UpdateUser(ctx, userID, &req)
	if err != nil {
		if err == errors.ErrUserNotFound {
			c.JSON(http.StatusNotFound, response.Fail(ctx, errors.ErrUserNotFound))
			return
		}
		h.logger.Error("更新用户信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
		return
	}

	c.JSON(http.StatusOK, response.OKWithMsg(ctx, nil, "用户信息更新成功"))
}

// CreateTenant godoc
// @Summary 创建租户
// @Description 创建新的租户（家庭/组织）
// @Tags tenants
// @Accept json
// @Produce json
// @Param request body dto.CreateTenantRequest true "租户信息"
// @Success 200 {object} response.Response "创建成功的租户信息"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/tenants [post]
func (h *UserHandler) CreateTenant(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidateErr(c.Request.Context(), err.Error()))
		return
	}

	// 从 Context 获取当前用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		h.logger.Error("CreateTenant: user ID not found in context")
		c.JSON(http.StatusUnauthorized, response.Unauthorized(c.Request.Context(), "user not authenticated"))
		return
	}

	ctx := c.Request.Context()
	tenant, err := h.tenantService.CreateTenant(ctx, &req, userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("创建租户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, tenant))
}

// GetUserInfo godoc
// @Summary 获取当前用户信息
// @Description 获取登录用户的详细信息
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.User "用户信息"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/users/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从 JWT token 中获取用户 ID
	h.logger.Debug("GetUserInfo: 尝试获取 userID")
	userID, exists := c.Get("userID")
	h.logger.Debug("GetUserInfo: userID exists?", zap.Bool("exists", exists), zap.Any("userID", userID))

	if !exists {
		h.logger.Warn("GetUserInfo: user ID not found in context")
		c.JSON(http.StatusUnauthorized, response.Fail(c.Request.Context(), errors.ErrUnauthorized.WithDetails("user ID not found in token")))
		return
	}

	// 安全的类型断言
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Error("GetUserInfo: invalid user ID type in context", zap.Any("userID", userID), zap.String("type", fmt.Sprintf("%T", userID)))
		c.JSON(http.StatusInternalServerError, response.Fail(c.Request.Context(), errors.ErrUnauthorized.WithDetails("invalid user ID type in context")))
		return
	}

	ctx := c.Request.Context()
	user, err := h.userQueryService.GetUserInfo(ctx, userIDUUID)
	if err != nil {
		h.logger.Error("获取用户信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
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
// @Param request body dto.UpdateProfileRequest true "个人资料信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidateErr(c.Request.Context(), err.Error()))
		return
	}

	// 从 JWT token 中获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Fail(c.Request.Context(), errors.ErrUnauthorized.WithDetails("user ID not found in token")))
		return
	}

	ctx := c.Request.Context()
	err := h.userCommandService.UpdateProfile(ctx, userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("更新个人资料失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
		return
	}

	c.JSON(http.StatusOK, response.OKWithMsg(ctx, nil, "个人资料更新成功"))
}
