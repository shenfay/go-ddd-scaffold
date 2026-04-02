package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	apperrors "github.com/shenfay/go-ddd-scaffold/pkg/errors"
	"github.com/shenfay/go-ddd-scaffold/internal/middleware"
)

// Handler HTTP 处理器
type Handler struct {
	service *Service
}

// NewHandler 创建认证处理器
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
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
// @Router /api/v1/auth/register [post]
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
// @Router /api/v1/auth/login [post]
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
		Email:     req.Email,
		Password:  req.Password,
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	
	resp, err := h.service.Login(c.Request.Context(), cmd)
	if err != nil {
		h.handleServiceError(c, err)
		return
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
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userID := c.GetString("user_id") // 从 JWT 中间件获取
	
	cmd := LogoutCommand{
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
// @Router /api/v1/auth/refresh [post]
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
