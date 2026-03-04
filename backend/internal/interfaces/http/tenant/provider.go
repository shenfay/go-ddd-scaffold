package tenant

import (
	"github.com/gin-gonic/gin"
)

// TenantRouterProvider 租户路由提供者
type TenantRouterProvider struct {
	handler *TenantHandler
}

// NewTenantRouterProvider 创建租户路由提供者
func NewTenantRouterProvider(handler *TenantHandler) *TenantRouterProvider {
	return &TenantRouterProvider{handler: handler}
}

// ProvideProtectedRoutes 注册需要认证的路由
func (p *TenantRouterProvider) ProvideProtectedRoutes(router *gin.RouterGroup) {
	// 租户创建路由
	tenants := router.Group("/tenants")
	{
		tenants.POST("", p.handler.CreateTenant)
	}
	
	// 获取我的租户列表
	router.GET("/tenants/my-tenants", p.handler.GetUserTenants)
}
