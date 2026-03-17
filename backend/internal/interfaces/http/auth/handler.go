package auth

import (
	"github.com/gin-gonic/gin"

	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// Handler 认证 HTTP处理器
type Handler struct {
	authService authApp.AuthService
	respHandler *httpShared.Handler
}

// NewHandler 创建处理器
func NewHandler(
	authService authApp.AuthService,
	respHandler *httpShared.Handler,
) *Handler {
	return &Handler{
		authService: authService,
		respHandler: respHandler,
	}
}

// @Summary 用户登录
// @Description 使用用户名或邮箱和密码进行登录，获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录凭证"
// @Success 200 {object} LoginResponse "登录成功返回令牌和用户信息"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 401 {object} httpShared.APIResponse "认证失败"
// @Router /auth/login [post]
// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !h.respHandler.BindJSON(c, &req) {
		return
	}

	cmd := &authApp.AuthenticateCommand{
		Identifier: req.UsernameOrEmail,
		Password:   req.Password,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
	}

	result, err := h.authService.AuthenticateUser(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	response := LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
		User: UserInfo{
			ID:       result.UserID,
			Username: result.Username,
			Email:    result.Email,
		},
	}

	h.respHandler.Success(c, response)
}

// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "用户注册信息"
// @Success 201 {object} RegisterResponse "注册成功返回用户信息"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 409 {object} httpShared.APIResponse "用户已存在"
// @Router /auth/register [post]
// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if !h.respHandler.BindJSON(c, &req) {
		return
	}

	cmd := &authApp.RegisterCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.authService.RegisterUser(c.Request.Context(), cmd)
	if err != nil {
		// 添加详细错误日志
		h.respHandler.Error(c, err)
		return
	}

	response := RegisterResponse{
		UserID:   result.UserID,
		Username: result.Username,
		Email:    result.Email,
	}

	h.respHandler.Created(c, response)
}

// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新令牌"
// @Success 200 {object} RefreshTokenResponse "刷新成功返回新令牌"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 401 {object} httpShared.APIResponse "令牌无效"
// @Router /auth/refresh [post]
// RefreshToken 刷新令牌
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if !h.respHandler.BindJSON(c, &req) {
		return
	}

	cmd := &authApp.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	response := RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
	}

	h.respHandler.Success(c, response)
}

// @Summary 用户登出
// @Description 使当前访问令牌失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 204 {object} httpShared.APIResponse "登出成功"
// @Failure 401 {object} httpShared.APIResponse "未授权"
// @Router /auth/logout [post]
// Logout 用户登出
func (h *Handler) Logout(c *gin.Context) {
	// 从上下文获取用户 ID（由认证中间件注入）
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		h.respHandler.Error(c, ddd.NewBusinessError("USER_NOT_FOUND", "用户未找到"))
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		h.respHandler.Error(c, ddd.NewBusinessError("INVALID_USER_ID", "无效的用户 ID"))
		return
	}

	// 从 Header 获取 Access Token
	authHeader := c.GetHeader("Authorization")
	accessToken := extractBearerToken(authHeader)

	cmd := &authApp.LogoutCommand{
		UserID:      userID,
		AccessToken: accessToken,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	}

	_, err := h.authService.Logout(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}

// @Summary 获取当前用户
// @Description 获取当前登录用户的信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CurrentUserResponse "当前用户信息"
// @Failure 401 {object} httpShared.APIResponse "未授权"
// @Router /auth/me [get]
// GetCurrentUser 获取当前用户
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		h.respHandler.Error(c, ddd.NewBusinessError("USER_NOT_FOUND", "用户未找到"))
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		h.respHandler.Error(c, ddd.NewBusinessError("INVALID_USER_ID", "无效的用户 ID"))
		return
	}

	// 从仓储获取完整用户信息
	ctx := c.Request.Context()
	foundUser, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, CurrentUserResponse{
		ID:          foundUser.ID,
		Username:    foundUser.Username,
		Email:       foundUser.Email,
		DisplayName: foundUser.DisplayName,
		Status:      foundUser.Status,
	})
}

// extractBearerToken 从 Authorization Header 提取 Bearer Token
func extractBearerToken(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
