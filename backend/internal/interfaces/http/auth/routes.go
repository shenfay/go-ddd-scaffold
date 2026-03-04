package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(router *gin.RouterGroup, handler *AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/logout", handler.Logout)
	}
}

// API 路由定义常量
const (
	AuthRegisterPath = "/auth/register"
	AuthLoginPath    = "/auth/login"
	AuthLogoutPath   = "/auth/logout"
)

// RouteInfo 路由信息结构
type RouteInfo struct {
	Method      string
	Path        string
	Description string
}

// GetAuthRoutes 获取认证路由信息列表
func GetAuthRoutes() []RouteInfo {
	return []RouteInfo{
		{"POST", AuthRegisterPath, "用户注册"},
		{"POST", AuthLoginPath, "用户登录"},
		{"POST", AuthLogoutPath, "用户登出"},
	}
}
