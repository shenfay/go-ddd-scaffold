package http

import (
	"github.com/gin-gonic/gin"
)

// UserRouterProvider 用户路由提供者
type UserRouterProvider struct {
	handler *UserHandler
}

// NewUserRouterProvider 创建用户路由提供者
func NewUserRouterProvider(handler *UserHandler) *UserRouterProvider {
	return &UserRouterProvider{handler: handler}
}

// ProvideProtectedRoutes 注册需要认证的路由
func (p *UserRouterProvider) ProvideProtectedRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/info", p.handler.GetUserInfo)
		users.PUT("/profile", p.handler.UpdateProfile)
		users.GET("/:id", p.handler.GetUser)
		users.PUT("/:id", p.handler.UpdateUser)
	}
}
