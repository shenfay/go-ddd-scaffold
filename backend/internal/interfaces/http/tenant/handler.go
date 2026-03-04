package tenant

import (
	"net/http"

	"go-ddd-scaffold/internal/application/tenant/service"
	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TenantHandler 租户 HTTP 处理器
type TenantHandler struct {
	tenantService service.TenantService
	logger        *zap.Logger
}

// NewTenantHandler 创建租户处理器实例
func NewTenantHandler(tenantService service.TenantService, logger *zap.Logger) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
		logger:        logger,
	}
}

// GetUserTenants godoc
// @Summary 获取用户的租户列表
// @Description 获取当前登录用户所属的所有租户及其角色
// @Tags tenants
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "租户列表"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/tenants/my-tenants [get]
func (h *TenantHandler) GetUserTenants(c *gin.Context) {
	// 从 Context 获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized(c.Request.Context(), "user not authenticated"))
		return
	}

	ctx := c.Request.Context()
	tenants, err := h.tenantService.GetUserTenants(ctx, userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("获取用户租户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
		return
	}

	c.JSON(http.StatusOK, response.OK(ctx, tenants))
}

// CreateTenant godoc
// @Summary 创建租户
// @Description 创建新的租户，创建者自动成为 owner
// @Tags tenants
// @Accept json
// @Produce json
// @Param request body CreateTenantRequest true "租户信息"
// @Success 201 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/tenants [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidateErr(c.Request.Context(), err.Error()))
		return
	}

	// 从 Context 获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Unauthorized(c.Request.Context(), "user not authenticated"))
		return
	}

	ctx := c.Request.Context()
	tenant, err := h.tenantService.CreateTenant(ctx, req.Name, req.Description, userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("创建租户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
		return
	}

	c.JSON(http.StatusCreated, response.OK(ctx, tenant))
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
}
