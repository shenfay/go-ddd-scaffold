package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
)

// Provider 认证模块提供者
type Provider struct {
	handler      *Handler
	tokenService auth.TokenService
}

// NewProvider 创建提供者
func NewProvider(handler *Handler, tokenService auth.TokenService) *Provider {
	return &Provider{
		handler:      handler,
		tokenService: tokenService,
	}
}

// ProvideRoutes 注册路由到 RouterGroup（会自动继承全局中间件）
func (p *Provider) ProvideRoutes(routerGroup *gin.RouterGroup) {
	// 公开端点（无需认证）
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/login", p.handler.Login)
		auth.POST("/register", p.handler.Register)
		auth.POST("/refresh", p.handler.RefreshToken)
	}

	// 需要认证的端点
	protected := routerGroup.Group("/auth")
	protected.Use(middleware.AuthMiddleware(p.tokenService))
	{
		protected.POST("/logout", p.handler.Logout)
		protected.GET("/me", p.handler.GetCurrentUser)
	}
}
