package activitylog

import (
	"github.com/gin-gonic/gin"
)

// Handler 活动日志 HTTP 处理器
type Handler struct {
	service *Service
}

// NewHandler 创建活动日志处理器
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.Use(h.authMiddleware())
	{
		users.GET("/:id/activity-logs", h.GetUserActivityLogs)
		users.GET("/me/activity-logs", h.GetCurrentUserActivityLogs)
	}
}

// GetUserActivityLogs 获取指定用户的活动日志
// @Summary Get user activity logs
// @Tags Activity
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} ActivityLog
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users/{id}/activity-logs [get]
func (h *Handler) GetUserActivityLogs(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(400, gin.H{"code": "INVALID_REQUEST", "message": "User ID is required"})
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		// 简单解析，实际应该验证
		// 这里省略验证逻辑
	}
	if o := c.Query("offset"); o != "" {
		// 简单解析，实际应该验证
	}

	logs, err := h.service.GetUserLogs(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}

	c.JSON(200, logs)
}

// GetCurrentUserActivityLogs 获取当前用户的活动日志
// @Summary Get current user activity logs
// @Tags Activity
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} ActivityLog
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users/me/activity-logs [get]
func (h *Handler) GetCurrentUserActivityLogs(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(401, gin.H{"code": "UNAUTHORIZED", "message": "User not authenticated"})
		return
	}

	logs, err := h.service.GetRecentLogs(c.Request.Context(), userID.(string), 10)
	if err != nil {
		c.JSON(500, gin.H{"code": "INTERNAL_ERROR", "message": err.Error()})
		return
	}

	c.JSON(200, logs)
}

// authMiddleware JWT 认证中间件（简化版本）
func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"code": "UNAUTHORIZED", "message": "Missing authorization header"})
			c.Abort()
			return
		}

		// 这里简化处理，实际需要验证 token
		// 完整实现应该调用 tokenService.ValidateAccessToken
		c.Next()
	}
}
