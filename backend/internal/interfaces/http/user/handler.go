package http

import (
	"net/http"

	tenant_dto"go-ddd-scaffold/internal/application/tenant/dto"
	tenant_service "go-ddd-scaffold/internal/application/tenant/service"
	dto"go-ddd-scaffold/internal/application/user/dto"
	user_service "go-ddd-scaffold/internal/application/user/service"
	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserHandler 用户 HTTP 处理器
type UserHandler struct {
	authService    user_service.AuthenticationService
	userQueryService user_service.UserQueryService
	userCommandService user_service.UserCommandService
	tenantService   tenant_service.TenantService
	logger           *zap.Logger
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(
	authService user_service.AuthenticationService,
	userQueryService user_service.UserQueryService,
	userCommandService user_service.UserCommandService,
	tenantService tenant_service.TenantService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		authService:        authService,
		userQueryService:  userQueryService,
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
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/users/{id} [put]
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

// CreateTenant godoc
// @Summary 创建租户
// @Description 创建新的租户（家庭/组织）
// @Tags tenants
// @Accept json
// @Produce json
// @Param request body tenant_dto.CreateTenantRequest true "租户信息"
// @Success 200 {object} response.Response"创建成功的租户信息"
// @Failure 400 {object} response.Response"请求参数错误"
// @Failure 401 {object} response.Response"未授权"
// @Failure 500 {object} response.Response"服务器内部错误"
// @Router/api/tenants [post]
func (h *UserHandler) CreateTenant(c *gin.Context) {
	var req tenant_dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ValidationFailed.WithDetails(err.Error()))
		return
	}

	// 从 Context 获取当前用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(errors.ErrUnauthorized.WithDetails("user not authenticated"))
		return
	}

	ctx := c.Request.Context()
	tenant, err := h.tenantService.CreateTenant(ctx, &req, userID.(uuid.UUID))
	if err != nil {
		c.Error(err) // 统一交给中间件处理
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
// @Param request body dto.UpdateProfileRequest true "个人资料信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
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
