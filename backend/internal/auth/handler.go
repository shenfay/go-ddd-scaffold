package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/activitylog"
	"github.com/shenfay/go-ddd-scaffold/internal/middleware"
	apperrors "github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// Handler HTTP 处理器
type Handler struct {
	service      *Service
	tokenService *TokenService
	activityLog  *activitylog.Service // 活动日志服务
}

// NewHandler 创建认证处理器
func NewHandler(service *Service, activityLogService *activitylog.Service) *Handler {
	return &Handler{
		service:      service,
		tokenService: service.tokenService,
		activityLog:  activityLogService,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", middleware.LoginRateLimit(), h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshToken)
		// 获取当前用户信息（需要认证）
		auth.GET("/me", h.authMiddleware(), h.GetCurrentUser)
		// 设备管理（需要认证）
		auth.GET("/devices", h.authMiddleware(), h.GetUserDevices)
		auth.DELETE("/devices/:token", h.authMiddleware(), h.RevokeDevice)
		auth.POST("/logout-all", h.authMiddleware(), h.LogoutAllDevices)
	}

	// 用户路由（需要认证）
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

// LogoutRequest 退出登录请求
type LogoutRequest struct {
	UserID string `json:"user_id"`
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
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := RegisterCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.service.Register(c.Request.Context(), cmd)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// 记录活动日志
	if h.activityLog != nil {
		_ = h.activityLog.Record(c.Request.Context(), activitylog.LogParams{
			UserID:      resp.User.ID,
			Email:       resp.User.Email,
			Action:      activitylog.ActivityRegister,
			Status:      activitylog.ActivitySuccess,
			IP:          c.ClientIP(),
			UserAgent:   c.GetHeader("User-Agent"),
			Description: "用户注册成功",
			Metadata:    nil,
		})
	}

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
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := LoginCommand{
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

	// 记录活动日志
	if h.activityLog != nil {
		h.activityLog.RecordAsync(activitylog.LogParams{
			UserID:      resp.User.ID,
			Email:       resp.User.Email,
			Action:      activitylog.ActivityLogin,
			Status:      activitylog.ActivitySuccess,
			IP:          c.ClientIP(),
			UserAgent:   c.GetHeader("User-Agent"),
			Description: "用户登录成功",
			Metadata: map[string]interface{}{
				"login_method": "password",
			},
		})
	}

	c.JSON(http.StatusOK, toAuthResponse(resp))
}

// Logout 处理用户退出
// @Summary User logout
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LogoutRequest true "Logout data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userID := c.GetString("user_id") // 从 JWT 中间件获取
	email := c.GetString("user_email")

	cmd := LogoutCommand{
		UserID: userID,
	}

	if err := h.service.Logout(c.Request.Context(), cmd); err != nil {
		h.handleServiceError(c, err)
		return
	}

	// 记录活动日志
	if h.activityLog != nil {
		h.activityLog.RecordAsync(activitylog.LogParams{
			UserID:      userID,
			Email:       email,
			Action:      activitylog.ActivityLogout,
			Status:      activitylog.ActivitySuccess,
			IP:          c.ClientIP(),
			UserAgent:   c.GetHeader("User-Agent"),
			Description: "用户登出成功",
			Metadata:    nil,
		})
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
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := RefreshTokenCommand{
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
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	// TODO: 使用统一的错误处理中间件
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

// authMiddleware JWT 认证中间件（用于路由组）
func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Missing authorization header",
			})
			c.Abort()
			return
		}

		// 提取 Bearer Token
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

		// 验证 Token
		claims, err := h.tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "INVALID_TOKEN",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}

// GetUserByID 根据 ID 获取用户信息
// @Summary Get user information by ID
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "User ID is required",
		})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
	})
}

// GetCurrentUser 获取当前登录用户信息
// @Summary Get current user information
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/me [get]
func (h *Handler) GetCurrentUser(c *gin.Context) {
	// 获取 Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Missing authorization header",
		})
		return
	}

	// 提取 Bearer Token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Invalid authorization format",
		})
		return
	}

	tokenString := parts[1]

	// 验证 Token
	claims, err := h.tokenService.ValidateAccessToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_TOKEN",
			Message: err.Error(),
		})
		return
	}

	// 从 service 获取用户详细信息
	user, err := h.service.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
	})
}

// toAuthResponse 转换为 HTTP 响应格式
func toAuthResponse(resp *ServiceAuthResponse) *AuthResponse {
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
// @Failure 401 {object} ErrorResponse
// @Router /auth/devices [get]
func (h *Handler) GetUserDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	// 获取设备列表
	devices, err := h.tokenService.GetUserDevices(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get devices",
		})
		return
	}

	// 转换为响应格式
	var deviceResponses []DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, DeviceResponse{
			TokenID:    device.UserID, // 这里需要调整，应该从 token 获取
			DeviceType: device.DeviceType,
			IP:         maskIP(device.IP),
			UserAgent:  device.UserAgent,
			CreatedAt:  device.CreatedAt,
			IsCurrent:  false, // 需要比对当前 token
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
// @Failure 401 {object} ErrorResponse
// @Router /auth/devices/{token} [delete]
func (h *Handler) RevokeDevice(c *gin.Context) {
	userID := c.GetString("user_id")
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Device token is required",
		})
		return
	}

	// 验证 token 属于当前用户
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

	// 撤销设备
	if err := h.tokenService.RevokeDeviceByToken(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to revoke device",
		})
		return
	}

	// 记录活动日志
	if h.activityLog != nil {
		h.activityLog.RecordAsync(activitylog.LogParams{
			UserID:      userID,
			Email:       c.GetString("user_email"),
			Action:      "REVOKE_DEVICE",
			Status:      "SUCCESS",
			IP:          c.ClientIP(),
			UserAgent:   c.GetHeader("User-Agent"),
			Description: "用户踢出设备",
			Metadata:    map[string]interface{}{"revoked_token": token},
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device revoked successfully"})
}

// LogoutAllDevices 退出所有设备
// @Summary Logout from all devices
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Router /auth/logout-all [post]
func (h *Handler) LogoutAllDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	// 撤销所有设备
	if err := h.tokenService.RevokeAllDevices(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to logout from all devices",
		})
		return
	}

	// 记录活动日志
	if h.activityLog != nil {
		h.activityLog.RecordAsync(activitylog.LogParams{
			UserID:      userID,
			Email:       c.GetString("user_email"),
			Action:      "LOGOUT_ALL",
			Status:      "SUCCESS",
			IP:          c.ClientIP(),
			UserAgent:   c.GetHeader("User-Agent"),
			Description: "用户退出所有设备",
			Metadata:    nil,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out from all devices successfully"})
}

// maskIP 脱敏 IP 地址
func maskIP(ip string) string {
	if ip == "" {
		return ""
	}
	// 简单脱敏：保留前两段，后面用 *** 替代
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

	// 简单的设备类型检测（可以根据需要扩展）
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
