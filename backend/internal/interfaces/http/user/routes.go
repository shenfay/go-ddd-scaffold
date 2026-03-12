package user

import (
	"github.com/gin-gonic/gin"
	httpiface "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
)

// RegisterRoutes 注册用户路由
func RegisterRoutes(router *gin.RouterGroup, handler *httpiface.Handler) {
	// 创建处理器（支持真实服务或 Mock）
	h := NewUserHandler(nil, handler)

	// 使用 Gin 原生 Group 方式定义路由
	v1 := router.Group("v1/users")
	{
		v1.POST("", h.CreateUser)
		v1.GET("/:id", h.GetUserByID)
		v1.PUT("/:id", h.UpdateUser)
		v1.DELETE("/:id", h.DeleteUser)
		v1.POST("/:id/activate", h.ActivateUser)
		v1.POST("/:id/deactivate", h.DeactivateUser)
	}

	// 认证路由
	router.POST("v1/users/authenticate", h.AuthenticateUser)
}

// init 自动注册到全局路由总线
func init() {
	httpiface.MustGetRouter().Register(RegisterRoutes)
}
