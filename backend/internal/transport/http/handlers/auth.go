package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/response"
	authErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/auth"
	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
	validationErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/validation"
)

// AuthHandler 认证 HTTP 处理器
type AuthHandler struct {
	service      *authentication.Service
	tokenService authentication.TokenService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(service *authentication.Service, tokenService authentication.TokenService) *AuthHandler {
	return &AuthHandler{
		service:      service,
		tokenService: tokenService,
	}
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`   // 用户邮箱
	Password string `json:"password" binding:"required,min=8,max=72"` // 用户密码（8-72字符）
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"` // 用户邮箱
	Password string `json:"password" binding:"required"`    // 用户密码
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // 刷新令牌
}

// Register 处理用户注册
//
// 创建新用户账户并返回认证令牌。
// 邮箱必须在系统中唯一。
//
// @Summary 用户注册
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "用户注册数据"
// @Success 201 {object} middleware.SuccessResponse{data=authentication.AuthResponse} "注册成功"
// @Failure 400 {object} middleware.ErrorResponse "请求参数错误"
// @Failure 409 {object} middleware.ErrorResponse "邮箱已存在"
// @Failure 500 {object} middleware.ErrorResponse "服务器内部错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := authentication.RegisterCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.service.Register(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, authentication.ToAuthResponse(resp))
}

// Login 处理用户登录
//
// 验证用户凭据并返回访问/刷新令牌。
// 跟踪登录失败次数，失败过多时锁定账户。
//
// @Summary 用户登录并返回令牌
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录凭据"
// @Success 200 {object} middleware.SuccessResponse{data=authentication.AuthResponse} "登录成功"
// @Failure 400 {object} middleware.ErrorResponse "请求参数错误"
// @Failure 401 {object} middleware.ErrorResponse "账号或密码错误"
// @Failure 423 {object} middleware.ErrorResponse "账户已锁定"
// @Failure 500 {object} middleware.ErrorResponse "服务器内部错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := authentication.LoginCommand{
		Email:      req.Email,
		Password:   req.Password,
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		DeviceType: detectDeviceType(c.Request.UserAgent()),
	}

	resp, err := h.service.Login(c.Request.Context(), cmd)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, authentication.ToAuthResponse(resp))
}

// Logout 处理用户退出
// @Summary 用户登出
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} middleware.ErrorResponse "未授权"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("user_id")

	cmd := authentication.LogoutCommand{
		UserID: userID,
	}

	if err := h.service.Logout(c.Request.Context(), cmd); err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Logged out successfully"})
}

// RefreshToken 刷新 Access Token
// @Summary 刷新访问令牌
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新令牌"
// @Success 200 {object} authentication.AuthResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse "令牌无效或已过期"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := authentication.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, authentication.ToAuthResponse(resp))
}

// handleServiceError 处理服务层错误
func (h *AuthHandler) handleServiceError(c *gin.Context, err error) {
	switch err {
	case userErr.ErrEmailAlreadyExists:
		c.JSON(http.StatusConflict, middleware.ErrorResponse{
			Code:    userErr.ErrEmailAlreadyExists.Code,
			Message: err.Error(),
		})
	case authErr.ErrInvalidCredentials, authErr.ErrInvalidToken, authErr.ErrTokenExpired:
		c.JSON(http.StatusUnauthorized, middleware.ErrorResponse{
			Code:    authErr.ErrInvalidCredentials.Code,
			Message: err.Error(),
		})
	case authErr.ErrAccountLocked:
		c.JSON(http.StatusLocked, middleware.ErrorResponse{
			Code:    authErr.ErrAccountLocked.Code,
			Message: err.Error(),
		})
	case userErr.ErrNotFound:
		c.JSON(http.StatusNotFound, middleware.ErrorResponse{
			Code:    userErr.ErrNotFound.Code,
			Message: err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Code:    "SYSTEM.INTERNAL_ERROR",
			Message: "Internal server error",
		})
	}
}

// GetUserByID 根据 ID 获取用户信息
// @Summary 根据ID获取用户信息
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "用户ID"
// @Success 200 {object} authentication.UserResponse
// @Failure 401 {object} middleware.ErrorResponse "未授权"
// @Failure 404 {object} middleware.ErrorResponse "用户不存在"
// @Router /users/{id} [get]
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "User ID is required",
		})
		return
	}

	u, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, middleware.ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	response.Success(c, authentication.ToUserResponse(u))
}

// GetCurrentUser 获取当前登录用户信息
// @Summary 获取当前用户信息
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} authentication.UserResponse
// @Failure 401 {object} middleware.ErrorResponse "未授权"
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, middleware.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Missing authorization header",
		})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, middleware.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Invalid authorization format",
		})
		return
	}

	tokenString := parts[1]

	claims, err := h.tokenService.ValidateAccessToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, middleware.ErrorResponse{
			Code:    "INVALID_TOKEN",
			Message: err.Error(),
		})
		return
	}

	u, err := h.service.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, middleware.ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	response.Success(c, authentication.ToUserResponse(u))
}

// maskIP 脱敏 IP 地址
func maskIP(ip string) string {
	if ip == "" {
		return ""
	}
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + ".***"
	}
	return ip[:len(ip)/2] + "***"
}

// DeviceResponse 设备响应
type DeviceResponse struct {
	TokenID    string `json:"token_id"`
	DeviceType string `json:"device_type"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	CreatedAt  string `json:"created_at"`
	IsCurrent  bool   `json:"is_current"`
}

// DevicesResponse 设备列表响应
type DevicesResponse struct {
	Devices []DeviceResponse `json:"devices"`
}

// GetUserDevices 获取当前用户的所有登录设备
// @Summary 获取用户登录设备列表
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DevicesResponse
// @Failure 401 {object} middleware.ErrorResponse "未授权"
// @Router /auth/devices [get]
func (h *AuthHandler) GetUserDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	devices, err := h.tokenService.GetUserDevices(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get devices",
		})
		return
	}

	var deviceResponses []DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, DeviceResponse{
			TokenID:    device.UserID,
			DeviceType: device.DeviceType,
			IP:         maskIP(device.IP),
			UserAgent:  device.UserAgent,
			CreatedAt:  device.CreatedAt,
			IsCurrent:  false,
		})
	}

	c.JSON(http.StatusOK, DevicesResponse{
		Devices: deviceResponses,
	})
}

// RevokeDevice 踢出指定设备
// @Summary 撤销指定设备
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Param token path string true "设备令牌"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse "未授权"
// @Router /auth/devices/{token} [delete]
func (h *AuthHandler) RevokeDevice(c *gin.Context) {
	userID := c.GetString("user_id")
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Device token is required",
		})
		return
	}

	deviceInfo, err := h.tokenService.ValidateRefreshTokenWithDevice(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, middleware.ErrorResponse{
			Code:    "INVALID_TOKEN",
			Message: "Invalid or expired device token",
		})
		return
	}

	if deviceInfo.UserID != userID {
		c.JSON(http.StatusForbidden, middleware.ErrorResponse{
			Code:    "FORBIDDEN",
			Message: "You can only revoke your own devices",
		})
		return
	}

	if err := h.tokenService.RevokeDeviceByToken(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to revoke device",
		})
		return
	}

	response.Success(c, gin.H{"message": "Device revoked successfully"})
}

// LogoutAllDevices 退出所有设备
// @Summary 登出所有设备
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} middleware.ErrorResponse "未授权"
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.tokenService.RevokeAllDevices(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to logout from all devices",
		})
		return
	}

	response.Success(c, gin.H{"message": "Logged out from all devices successfully"})
}

// detectDeviceType 根据 User-Agent 检测设备类型
func detectDeviceType(userAgent string) string {
	ua := userAgent
	if ua == "" {
		return "unknown"
	}

	if containsAny(ua, []string{"Mobile", "Android", "iPhone", "iPad"}) {
		if containsAny(ua, []string{"iPad", "Tablet"}) {
			return "tablet"
		}
		return "mobile"
	}

	return "desktop"
}

// containsAny 检查字符串是否包含任意一个子串
func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// RequestPasswordResetRequest 密码重置请求 DTO
type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 执行密码重置 DTO
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// RequestPasswordReset 处理密码重置请求
// @Summary 请求密码重置
// @Description 发送密码重置邮件到指定邮箱
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RequestPasswordResetRequest true "密码重置请求"
// @Success 200 {object} middleware.SuccessResponse "请求成功"
// @Failure 400 {object} middleware.ErrorResponse "请求参数错误"
// @Failure 500 {object} middleware.ErrorResponse "服务器内部错误"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := authentication.RequestPasswordResetCommand{
		Email: req.Email,
	}

	if err := h.service.RequestPasswordReset(c.Request.Context(), cmd); err != nil {
		response.Error(c, err)
		return
	}

	// 安全考虑:统一返回成功,不暴露邮箱是否存在
	response.Success(c, gin.H{
		"message": "If the email exists, a reset link has been sent",
	})
}

// ResetPassword 执行密码重置
// @Summary 执行密码重置
// @Description 使用重置令牌设置新密码
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "密码重置数据"
// @Success 200 {object} middleware.SuccessResponse "重置成功"
// @Failure 400 {object} middleware.ErrorResponse "令牌无效或密码强度不足"
// @Failure 500 {object} middleware.ErrorResponse "服务器内部错误"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := authentication.ResetPasswordCommand{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}

	if err := h.service.ResetPassword(c.Request.Context(), cmd); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Password reset successfully",
	})
}
