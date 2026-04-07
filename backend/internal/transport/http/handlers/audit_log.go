package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/repository"
	"github.com/shenfay/go-ddd-scaffold/pkg/utils"
)

// AuditLogHandler 审计日志HTTP处理器（查询）
type AuditLogHandler struct {
	auditLogRepo repository.AuditLogRepository
}

// NewAuditLogHandler 创建审计日志处理器
func NewAuditLogHandler(auditLogRepo repository.AuditLogRepository) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogRepo: auditLogRepo,
	}
}

// RegisterRoutes 注册审计日志路由
func (h *AuditLogHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auditLogs := rg.Group("/audit-logs")
	{
		auditLogs.GET("", h.ListAuditLogs)
		auditLogs.GET("/user/:user_id", h.GetUserAuditLogs)
	}
}

// ListAuditLogs 查询审计日志列表
func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	limit := utils.ToInt(c.DefaultQuery("limit", "20"))
	offset := utils.ToInt(c.DefaultQuery("offset", "0"))

	logs, err := h.auditLogRepo.FindByUserID(c.Request.Context(), "", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to fetch audit logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUserAuditLogs 查询用户审计日志
func (h *AuditLogHandler) GetUserAuditLogs(c *gin.Context) {
	userID := c.Param("user_id")
	limit := utils.ToInt(c.DefaultQuery("limit", "20"))
	offset := utils.ToInt(c.DefaultQuery("offset", "0"))

	logs, err := h.auditLogRepo.FindByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to fetch user audit logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"limit":  limit,
		"offset": offset,
	})
}
