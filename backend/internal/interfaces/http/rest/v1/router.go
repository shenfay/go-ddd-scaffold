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

// RegisterAuthRoutes 注册认证路由
// handlerProvider 提供 Handler 的函数，避免循环依赖
func (r *Router) RegisterAuthRoutes(
	group *gin.RouterGroup,
	handlerProvider func() (login, register, refresh, logout, getCurrentUser gin.HandlerFunc),
	authMiddleware gin.HandlerFunc, // 新增认证中间件参数
) {
	authGroup := group.Group("/auth")
	{
		loginHandler, registerHandler, refreshHandler, logoutHandler, getCurrentUserHandler := handlerProvider()

		// POST /auth/login - 用户登录
		authGroup.POST("/login", loginHandler)

		// POST /auth/register - 用户注册
		authGroup.POST("/register", registerHandler)

		// POST /auth/refresh - 刷新令牌（不需要认证中间件，通过 refresh_token 验证）
		authGroup.POST("/refresh", refreshHandler)

		// POST /auth/logout - 用户登出（需要认证）
		authGroup.POST("/logout", authMiddleware, logoutHandler)

		// GET /auth/me - 获取当前用户（需要认证）
		authGroup.GET("/me", authMiddleware, getCurrentUserHandler)
	}
}

// RegisterUserRoutes 注册用户路由
// handlerProvider 提供 Handler 的函数，避免循环依赖
func (r *Router) RegisterUserRoutes(
	group *gin.RouterGroup,
	handlerProvider func() (getUser, updateProfile, changePassword gin.HandlerFunc),
) {
	userGroup := group.Group("/users")
	{
		getUserHandler, updateProfileHandler, changePasswordHandler := handlerProvider()

		// GET /users/:id - 获取用户详情
		userGroup.GET("/:id", getUserHandler)

		// PUT /users/:id - 更新用户资料
		userGroup.PUT("/:id", updateProfileHandler)

		// PUT /users/:id/password - 修改密码
		userGroup.PUT("/:id/password", changePasswordHandler)
	}
}
