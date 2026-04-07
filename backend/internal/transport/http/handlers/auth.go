package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"

	"github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware"
	apperrors "github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// AuthHandler 认证HTTP处理器
type AuthHandler struct {
	service      *authentication.Service
	tokenService authentication.TokenService
	asynqClient  *asynq.Client
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(service *authentication.Service, tokenService authentication.TokenService, asynqClient *asynq.Client) *AuthHandler {
	return &AuthHandler{
		service:      service,
		tokenService: tokenService,
		asynqClient:  asynqClient,
	}
}

// RegisterRoutes 注册路由
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", middleware.LoginRateLimit(), h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshToken)
		auth.GET("/me", h.authMiddleware(), h.GetCurrentUser)
		auth.GET("/devices", h.authMiddleware(), h.GetUserDevices)
		auth.DELETE("/devices/:token", h.authMiddleware(), h.RevokeDevice)
		auth.POST("/logout-all", h.authMiddleware(), h.LogoutAllDevices)
	}

	users := router.Group("/users")
	users.Use(h.authMiddleware())
	{
		users.GET("/:id", h.GetUserByID)
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Register 处理用户注册
// @Summary Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := authentication.RegisterCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.service.Register(c.Request.Context(), cmd)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// 发布活动日志
	h.publishActivityLog(c, resp.User.ID, resp.User.Email, "USER.REGISTERED", "用户注册")

	c.JSON(http.StatusCreated, toAuthResponse(resp))
}

// Login 处理用户登录
// @Summary User login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 423 {object} ErrorResponse "Account locked"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
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

	// 发布活动日志
	h.publishActivityLog(c, resp.User.ID, resp.User.Email, "USER.LOGIN", "用户登录")

	c.JSON(http.StatusOK, toAuthResponse(resp))
}

// Logout 处理用户退出
// @Summary User logout
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse "Unauthorized"
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

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken 刷新 Access Token
// @Summary Refresh access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Invalid or expired token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := authentication.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), cmd)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, toAuthResponse(resp))
}

// handleServiceError 处理服务层错误
func (h *AuthHandler) handleServiceError(c *gin.Context, err error) {
	switch err {
	case apperrors.ErrEmailAlreadyExists:
		c.JSON(http.StatusConflict, ErrorResponse{
			Code:    apperrors.ErrorCodeEmailAlreadyExists,
			Message: err.Error(),
		})
	case apperrors.ErrInvalidCredentials, apperrors.ErrInvalidToken, apperrors.ErrTokenExpired:
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    apperrors.ErrorCodeInvalidCredentials,
			Message: err.Error(),
		})
	case apperrors.ErrAccountLocked:
		c.JSON(http.StatusLocked, ErrorResponse{
			Code:    apperrors.ErrorCodeAccountLocked,
			Message: err.Error(),
		})
	case apperrors.ErrUserNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    apperrors.ErrorCodeUserNotFound,
			Message: err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    apperrors.ErrorCodeInternal,
			Message: "Internal server error",
		})
	}
}

// authMiddleware JWT 认证中间件
func (h *AuthHandler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Missing authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := h.tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "INVALID_TOKEN",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}

// GetUserByID 根据 ID 获取用户信息
// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "User ID is required",
		})
		return
	}

	u, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
	})
}

// GetCurrentUser 获取当前登录用户信息
// @Summary Get current user information
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Missing authorization header",
		})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Invalid authorization format",
		})
		return
	}

	tokenString := parts[1]

	claims, err := h.tokenService.ValidateAccessToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_TOKEN",
			Message: err.Error(),
		})
		return
	}

	u, err := h.service.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
	})
}

// toAuthResponse 转换为 HTTP 响应格式
func toAuthResponse(resp *authentication.ServiceAuthResponse) *AuthResponse {
	return &AuthResponse{
		User: &UserResponse{
			ID:            resp.User.ID,
			Email:         resp.User.Email,
			EmailVerified: resp.User.EmailVerified,
			LastLoginAt:   resp.User.LastLoginAt,
			CreatedAt:     resp.User.CreatedAt,
		},
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    int64(resp.ExpiresIn / time.Second),
	}
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
// @Summary Get user's logged-in devices
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DevicesResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/devices [get]
func (h *AuthHandler) GetUserDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	devices, err := h.tokenService.GetUserDevices(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
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
// @Summary Revoke a specific device
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Param token path string true "Device token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/devices/{token} [delete]
func (h *AuthHandler) RevokeDevice(c *gin.Context) {
	userID := c.GetString("user_id")
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Device token is required",
		})
		return
	}

	deviceInfo, err := h.tokenService.ValidateRefreshTokenWithDevice(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_TOKEN",
			Message: "Invalid or expired device token",
		})
		return
	}

	if deviceInfo.UserID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    "FORBIDDEN",
			Message: "You can only revoke your own devices",
		})
		return
	}

	if err := h.tokenService.RevokeDeviceByToken(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to revoke device",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device revoked successfully"})
}

// LogoutAllDevices 退出所有设备
// @Summary Logout from all devices
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.tokenService.RevokeAllDevices(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to logout from all devices",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out from all devices successfully"})
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

// publishActivityLog 发布活动日志任务
func (h *AuthHandler) publishActivityLog(c *gin.Context, userID, email, action, description string) {
	if h.asynqClient == nil {
		return
	}

	userAgent := c.GetHeader("User-Agent")
	device := detectDeviceType(userAgent)

	payload := map[string]interface{}{
		"user_id":     userID,
		"email":       email,
		"action":      action,
		"status":      "SUCCESS",
		"ip":          c.ClientIP(),
		"user_agent":  userAgent,
		"device":      device,
		"description": description,
		"metadata":    nil,
	}

	payloadBytes, _ := json.Marshal(payload)
	h.asynqClient.EnqueueContext(c.Request.Context(),
		asynq.NewTask("activity:record", payloadBytes),
		asynq.Queue("default"),
	)
}
