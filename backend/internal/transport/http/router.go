package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/handlers"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware"
)

// Router 路由配置
type Router struct {
	engine       *gin.Engine
	authHandler  *handlers.AuthHandler
	tokenService authentication.TokenService
}

// NewRouter 创建路由器
func NewRouter(
	engine *gin.Engine,
	authHandler *handlers.AuthHandler,
	tokenService authentication.TokenService,
) *Router {
	return &Router{
		engine:       engine,
		authHandler:  authHandler,
		tokenService: tokenService,
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
		r.setupAuthRoutes(v1)
		r.setupUserRoutes(v1)
	}

	// 注册 Swagger UI 路由（开发环境）
	middleware.RegisterSwagger(r.engine, middleware.DefaultSwaggerConfig())
}

// setupAuthRoutes 配置认证相关路由
func (r *Router) setupAuthRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	{
		// 公开路由
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", middleware.LoginRateLimit(), r.authHandler.Login)
		auth.POST("/logout", r.authHandler.Logout)
		auth.POST("/refresh", r.authHandler.RefreshToken)

		// 需要认证的路由
		authMiddleware := middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
			TokenService: r.tokenService,
		})
		auth.GET("/me", authMiddleware, r.authHandler.GetCurrentUser)
		auth.GET("/devices", authMiddleware, r.authHandler.GetUserDevices)
		auth.DELETE("/devices/:token", authMiddleware, r.authHandler.RevokeDevice)
		auth.POST("/logout-all", authMiddleware, r.authHandler.LogoutAllDevices)
	}
}

// setupUserRoutes 配置用户相关路由
func (r *Router) setupUserRoutes(v1 *gin.RouterGroup) {
	users := v1.Group("/users")
	users.Use(middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
		TokenService: r.tokenService,
	}))
	{
		users.GET("/:id", r.authHandler.GetUserByID)
	}
}
