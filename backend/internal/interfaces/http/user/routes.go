package user

import (
	"github.com/gin-gonic/gin"
	httpiface "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
)

// RegisterRoutes 注册用户领域路由到路由总线
func RegisterRoutes(router *gin.RouterGroup, handler *httpiface.Handler) {
	// 创建处理器（支持真实服务或 Mock）
	h := NewUserHandler(nil, handler)

	// 用户资源路由组
	userRouter := router.Group("/users")
	{
		userRouter.POST("", h.CreateUser)
		userRouter.GET("/:id", h.GetUserByID)
		userRouter.PUT("/:id", h.UpdateUser)
		userRouter.DELETE("/:id", h.DeleteUser)
		userRouter.POST("/:id/activate", h.ActivateUser)
		userRouter.POST("/:id/deactivate", h.DeactivateUser)
	}

	// 认证相关路由
	router.POST("/users/authenticate", h.AuthenticateUser)
}

// init 自动注册到全局路由总线
func init() {
	httpiface.GetRouter().Register(RegisterRoutes)
}
