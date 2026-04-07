package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/handlers"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware"
)

// Router 路由配置
type Router struct {
	engine      *gin.Engine
	authHandler *handlers.AuthHandler
}

// NewRouter 创建路由器
func NewRouter(
	engine *gin.Engine,
	authHandler *handlers.AuthHandler,
) *Router {
	return &Router{
		engine:      engine,
		authHandler: authHandler,
	}
}

// Setup 配置所有路由
func (r *Router) Setup() {
	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API v1 路由组
	v1 := r.engine.Group("/api/v1")
	{
		// 公开路由（无需认证）
		public := v1.Group("")
		{
			r.authHandler.RegisterRoutes(public)
		}

		// 需要认证的路由（后续扩展）
		// protected := v1.Group("")
		// protected.Use(middleware.JWTAuthMiddleware(tokenService))
		// {
		//     protected.GET("/users/me", userHandler.GetProfile)
		// }
	}

	// 注册 Swagger UI 路由（开发环境）
	middleware.RegisterSwagger(r.engine, middleware.DefaultSwaggerConfig())
}
