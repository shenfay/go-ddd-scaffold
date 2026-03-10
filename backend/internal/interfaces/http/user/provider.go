package http

import (
	"github.com/gin-gonic/gin"
)

// UserRouterProvider 用户路由提供者（仅处理用户管理）
type UserRouterProvider struct {
	handler *UserHandler
}

// NewUserRouterProvider 创建用户路由提供者
func NewUserRouterProvider(handler *UserHandler) *UserRouterProvider {
	return &UserRouterProvider{handler: handler}
}

// ProvideProtectedRoutes 注册需要认证的路由
func (p *UserRouterProvider) ProvideProtectedRoutes(router*gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/:id", p.handler.GetUser)      // 获取指定用户信息
		users.PUT("/:id", p.handler.UpdateUser)   // 更新指定用户信息
	}
}

// ProfileRouterProvider 个人资料路由提供者
type ProfileRouterProvider struct {
	handler *ProfileHandler
}

// NewProfileRouterProvider 创建个人资料路由提供者
func NewProfileRouterProvider(handler*ProfileHandler) *ProfileRouterProvider {
	return &ProfileRouterProvider{handler: handler}
}

// ProvideProtectedRoutes 注册需要认证的个人资料路由
func (p *ProfileRouterProvider) ProvideProtectedRoutes(router*gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/info", p.handler.GetUserInfo)      // 获取当前用户信息
		users.PUT("/profile", p.handler.UpdateProfile) // 更新个人资料
	}
}
