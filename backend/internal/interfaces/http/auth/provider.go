package http

import (
	"github.com/gin-gonic/gin"
)

// AuthRouterProvider 认证路由提供者
type AuthRouterProvider struct {
	handler *AuthHandler
}

// NewAuthRouterProvider 创建认证路由提供者
func NewAuthRouterProvider(handler *AuthHandler) *AuthRouterProvider {
	return &AuthRouterProvider{handler: handler}
}

// ProvidePublicRoutes 注册公开路由（无需认证）
func (p *AuthRouterProvider) ProvidePublicRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", p.handler.Register)
		auth.POST("/login", p.handler.Login)
	}
}

// ProvideProtectedRoutes 注册受保护路由
func (p *AuthRouterProvider) ProvideProtectedRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/logout", p.handler.Logout)
	}
}
