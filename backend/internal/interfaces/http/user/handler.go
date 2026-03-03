package http

import (
	"net/http"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/application/user/service"
	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	userService *service.Service
	logger      *zap.Logger
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *service.Service, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} dto.User "注册成功的用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 409 {object} map[string]string "用户已存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(c.Request.Context(), errors.ErrInvalidEmail.WithDetails(err.Error())))
		return
	}

	ctx := c.Request.Context()
	user, err := h.userService.Register(ctx, &req)
	if err != nil {
		switch err {
		case errors.ErrUserExists:
			c.JSON(http.StatusConflict, response.Fail(ctx, errors.ErrUserExists))
		case errors.ErrTenantLimitExceed:
			c.JSON(http.StatusConflict, response.Fail(ctx, errors.ErrTenantLimitExceed))
		case errors.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, response.Fail(ctx, errors.ErrUnauthorized))
		default:
			h.logger.Error("注册失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
		}
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, user))
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录接口
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} dto.LoginResponse "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Failure 403 {object} response.Response "账户被禁用"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidateErr(c.Request.Context(), err.Error()))
		return
	}

	ctx := c.Request.Context()
	resp, err := h.userService.Login(ctx, &req)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			c.JSON(http.StatusUnauthorized, response.Fail(ctx, errors.ErrUserNotFound))
		case errors.ErrInvalidPassword:
			c.JSON(http.StatusUnauthorized, response.Fail(ctx, errors.ErrInvalidPassword))
		case errors.ErrUnauthorized:
			c.JSON(http.StatusForbidden, response.Fail(ctx, errors.ErrUnauthorized))
		default:
			h.logger.Error("登录失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
		}
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, resp))
}

// GetUser godoc
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} dto.User "用户信息"
// @Failure 400 {object} response.Response "无效的用户ID"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(c.Request.Context(), errors.InvalidParameter.WithDetails("无效的用户ID格式")))
		return
	}

	ctx := c.Request.Context()
	user, err := h.userService.GetUser(ctx, userID)
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
// @Param id path string true "用户ID"
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
		c.JSON(http.StatusBadRequest, response.Fail(c.Request.Context(), errors.InvalidParameter.WithDetails("无效的用户ID格式")))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidateErr(c.Request.Context(), err.Error()))
		return
	}

	ctx := c.Request.Context()
	err = h.userService.UpdateUser(ctx, userID, &req)
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

	ctx := c.Request.Context()
	tenant, err := h.userService.CreateTenant(ctx, &req)
	if err != nil {
		h.logger.Error("创建租户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, tenant))
}
