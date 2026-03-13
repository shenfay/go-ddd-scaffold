package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
)

// Provider 认证模块提供者
type Provider struct {
	handler      *Handler
	tokenService user.TokenService
}

// NewProvider 创建提供者
func NewProvider(handler *Handler, tokenService user.TokenService) *Provider {
	return &Provider{
		handler:      handler,
		tokenService: tokenService,
	}
}

// ProvideRoutes 注册路由
func (p *Provider) ProvideRoutes(router *gin.Engine) {
	// 公开端点（无需认证）
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", p.handler.Login)
		auth.POST("/register", p.handler.Register)
		auth.POST("/refresh", p.handler.RefreshToken)
	}

	// 需要认证的端点
	protected := router.Group("/api/v1/auth")
	protected.Use(middleware.AuthMiddleware(p.tokenService))
	{
		protected.POST("/logout", p.handler.Logout)
		protected.GET("/me", p.handler.GetCurrentUser)
	}
}
