package http

import (
	"net/http"
	"strings"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/application/user/service"
	"go-ddd-scaffold/internal/pkg/errors"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService   service.AuthenticationService
	logger         *zap.Logger
	tokenBlacklist service.TokenBlacklistService // Token 黑名单服务
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(
	authService service.AuthenticationService,
	logger*zap.Logger,
	tokenBlacklist service.TokenBlacklistService,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		logger:         logger,
		tokenBlacklist: tokenBlacklist,
	}
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} dto.User"注册成功的用户信息"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 409 {object} response.Response "用户已存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router/api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ValidationFailed.WithDetails(err.Error()))
		return
	}

	ctx := c.Request.Context()
	user, err := h.authService.Register(ctx, &req)
	if err != nil {
		c.Error(err) // 统一交给中间件处理
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
// @Failure 401 {object} response.Response"用户名或密码错误"
// @Failure 403 {object} response.Response "账户被禁用"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router/api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ValidationFailed.WithDetails(err.Error()))
		return
	}

	ctx := c.Request.Context()
	resp, err := h.authService.Login(ctx, &req)
	if err != nil {
		c.Error(err) // 统一交给中间件处理
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, resp))
}

// Logout godoc
// @Summary 用户登出
// @Description 用户登出接口（将 token 加入黑名单）
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response"登出成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router/api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	// 从 Context 获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(errors.ErrUnauthorized.WithDetails("user not authenticated"))
		return
	}

	// 提取 Token（用于加入黑名单）
	token := extractToken(c)

	// 调用登出服务
	err := h.authService.Logout(ctx, userID.(uuid.UUID), token)
	if err != nil {
		c.Error(err) // 统一交给中间件处理
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, nil))
}

// extractToken 从 Authorization Header 中提取 Token
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
