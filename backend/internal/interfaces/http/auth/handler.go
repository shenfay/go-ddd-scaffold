package http

import (
	"net/http"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/application/user/service"
	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService *service.Service
	logger      *zap.Logger
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService *service.Service, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		logger:      logger,
	}
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} dto.User "注册成功的用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 409 {object} map[string]string "用户已存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
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
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} dto.LoginResponse "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Failure 403 {object} response.Response "账户被禁用"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
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

// Logout godoc
// @Summary 用户登出
// @Description 用户登出接口
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "登出成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	// TODO: 可以在这里添加 token 黑名单逻辑
	c.JSON(http.StatusOK, response.OK(ctx, nil))
}
