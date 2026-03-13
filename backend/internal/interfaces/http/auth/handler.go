package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application/auth/commands"
	httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// Handler 认证 HTTP处理器
type Handler struct {
	authenticateHandler *commands.AuthenticateHandler
	registerHandler     *commands.RegisterHandler
	refreshTokenHandler *commands.RefreshTokenHandler
	logoutHandler       *commands.LogoutHandler
	respHandler         *httpShared.Handler
}

// NewHandler 创建处理器
func NewHandler(
	authenticateHandler *commands.AuthenticateHandler,
	registerHandler *commands.RegisterHandler,
	refreshTokenHandler *commands.RefreshTokenHandler,
	logoutHandler *commands.LogoutHandler,
	respHandler *httpShared.Handler,
) *Handler {
	return &Handler{
		authenticateHandler: authenticateHandler,
		registerHandler:     registerHandler,
		refreshTokenHandler: refreshTokenHandler,
		logoutHandler:       logoutHandler,
		respHandler:         respHandler,
	}
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !h.respHandler.BindJSON(c, &req) {
		return
	}

	cmd := &commands.AuthenticateCommand{
		Identifier: req.UsernameOrEmail,
		Password:   req.Password,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
	}

	result, err := h.authenticateHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
		"token_type":    "Bearer",
		"user": map[string]interface{}{
			"id":       result.UserID,
			"username": result.Username,
			"email":    result.Email,
		},
	})
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if !h.respHandler.BindJSON(c, &req) {
		return
	}

	cmd := &commands.RegisterCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.registerHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, map[string]interface{}{
		"user_id":  result.UserID,
		"username": result.Username,
		"email":    result.Email,
	})
}

// RefreshToken 刷新令牌
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if !h.respHandler.BindJSON(c, &req) {
		return
	}

	cmd := &commands.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
	}

	result, err := h.refreshTokenHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
		"token_type":    "Bearer",
	})
}

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

	cmd := &commands.LogoutCommand{
		UserID:      userID,
		AccessToken: accessToken,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	}

	_, err := h.logoutHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}

// GetCurrentUser 获取当前用户
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.respHandler.Error(c, ddd.NewBusinessError("USER_NOT_FOUND", "用户未找到"))
		return
	}

	// TODO: 从仓储获取用户信息并返回
	h.respHandler.Success(c, map[string]interface{}{
		"id": userID,
	})
}

// LoginRequest 登录请求
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// extractBearerToken 从 Authorization Header 提取 Bearer Token
func extractBearerToken(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
