package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
)

// Routes 认证模块路由
type Routes struct {
	handler      *Handler
	tokenService auth.TokenService
}

// NewRoutes 创建路由
func NewRoutes(handler *Handler, tokenService auth.TokenService) *Routes {
	return &Routes{
		handler:      handler,
		tokenService: tokenService,
	}
}

// Register 注册路由到 RouterGroup（会自动继承全局中间件）
func (r *Routes) Register(routerGroup *gin.RouterGroup) {
	// 公开端点（无需认证）
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/login", r.handler.Login)
		auth.POST("/register", r.handler.Register)
		auth.POST("/refresh", r.handler.RefreshToken)
	}

	// 需要认证的端点
	protected := routerGroup.Group("/auth")
	protected.Use(middleware.AuthMiddleware(r.tokenService))
	{
		protected.POST("/logout", r.handler.Logout)
		protected.GET("/me", r.handler.GetCurrentUser)
	}
}
