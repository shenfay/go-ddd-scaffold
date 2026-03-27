package v1

import (
	"github.com/gin-gonic/gin"
)

// Router HTTP 路由注册器
type Router struct {
	// 空结构体，直接使用 group 注册路由
}

// NewRouter 创建路由注册器
func NewRouter() *Router {
	return &Router{}
}

// AuthHandlers 认证 Handler 集合
type AuthHandlers struct {
	Login          gin.HandlerFunc
	Register       gin.HandlerFunc
	Refresh        gin.HandlerFunc
	Logout         gin.HandlerFunc
	GetCurrentUser gin.HandlerFunc
}

// UserHandlers 用户 Handler 集合
type UserHandlers struct {
	GetUser        gin.HandlerFunc
	UpdateProfile  gin.HandlerFunc
	ChangePassword gin.HandlerFunc
}

// RegisterAuthRoutes 注册认证路由
// handlers: Handler 集合，避免循环依赖
// authMiddleware: JWT 认证中间件，用于保护需要登录的路由（logout、me）
func (r *Router) RegisterAuthRoutes(
	group *gin.RouterGroup,
	handlers AuthHandlers,
	authMiddleware gin.HandlerFunc,
) {
	// 创建 v1 版本分组
	v1Group := group.Group("/v1")
	authGroup := v1Group.Group("/auth")
	{
		// 公开路由（无需认证）
		authGroup.POST("/login", handlers.Login)
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/refresh", handlers.Refresh)

		// 受保护路由（需要认证）
		authGroup.POST("/logout", authMiddleware, handlers.Logout)
		authGroup.GET("/me", authMiddleware, handlers.GetCurrentUser)
	}
}

// RegisterUserRoutes 注册用户路由
// handlers: Handler 集合，避免循环依赖
// 注意：所有用户路由都需要认证保护（建议在 router 注册时统一添加中间件）
func (r *Router) RegisterUserRoutes(
	group *gin.RouterGroup,
	handlers UserHandlers,
) {
	// 创建 v1 版本分组
	v1Group := group.Group("/v1")
	userGroup := v1Group.Group("/users")
	{
		// GET /v1/users/:id - 获取用户详情
		userGroup.GET("/:id", handlers.GetUser)

		// PUT /v1/users/:id - 更新用户资料
		userGroup.PUT("/:id", handlers.UpdateProfile)

		// PUT /v1/users/:id/password - 修改密码
		userGroup.PUT("/:id/password", handlers.ChangePassword)
	}
}
