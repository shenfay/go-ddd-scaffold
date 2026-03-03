package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(router *gin.RouterGroup, handler *UserHandler) {
	users := router.Group("/users")
	{
		users.POST("/register", handler.Register)
		users.POST("/login", handler.Login)
		users.GET("/:id", handler.GetUser)
		users.PUT("/:id", handler.UpdateUser)
	}

	tenants := router.Group("/tenants")
	{
		tenants.POST("", handler.CreateTenant)
	}
}

// API路由定义常量
const (
	UserRegisterPath = "/users/register"
	UserLoginPath    = "/users/login"
	UserGetPath      = "/users/:id"
	UserUpdatePath   = "/users/:id"
	TenantCreatePath = "/tenants"
)

// RouteInfo 路由信息结构
type RouteInfo struct {
	Method      string
	Path        string
	Description string
}

// GetUserRoutes 获取用户路由信息列表
func GetUserRoutes() []RouteInfo {
	return []RouteInfo{
		{"POST", UserRegisterPath, "用户注册"},
		{"POST", UserLoginPath, "用户登录"},
		{"GET", UserGetPath, "获取用户信息"},
		{"PUT", UserUpdatePath, "更新用户信息"},
		{"POST", TenantCreatePath, "创建租户"},
	}
}
