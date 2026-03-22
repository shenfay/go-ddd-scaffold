package user

import (
	"github.com/gin-gonic/gin"
)

// Routes 用户领域路由
type Routes struct {
	handler *Handler
}

// NewRoutes 创建用户领域路由
func NewRoutes(handler *Handler) *Routes {
	return &Routes{
		handler: handler,
	}
}

// Register 注册用户领域路由
func (r *Routes) Register(routerGroup *gin.RouterGroup) {
	// 用户资源路由
	userGroup := routerGroup.Group("/users")
	{
		// GET /users/:user_id - 获取用户详情
		userGroup.GET("/:user_id", r.handler.GetUser)

		// PUT /users/:user_id - 更新用户
		userGroup.PUT("/:user_id", r.handler.UpdateUser)

		// POST /users/:user_id/password - 修改密码
		userGroup.POST("/:user_id/password", r.handler.ChangePassword)
	}
}
