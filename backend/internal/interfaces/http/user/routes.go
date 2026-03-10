package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户管理路由（仅处理用户管理相关）
func RegisterUserRoutes(router*gin.RouterGroup, handler *UserHandler) {
	users := router.Group("/users")
	{
		users.GET("/:id", handler.GetUser)    // 获取指定用户信息
		users.PUT("/:id", handler.UpdateUser) // 更新指定用户信息
	}
}

// RegisterProfileRoutes 注册个人资料路由
func RegisterProfileRoutes(router *gin.RouterGroup, handler*ProfileHandler) {
	users := router.Group("/users")
	{
		users.GET("/info", handler.GetUserInfo)      // 获取当前用户信息
		users.PUT("/profile", handler.UpdateProfile) // 更新个人资料
	}
}
