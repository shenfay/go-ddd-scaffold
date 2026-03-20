package user

import (
	"github.com/gin-gonic/gin"
	http "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
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

// RegisterTo 注册用户领域路由到依赖容器
func (r *Routes) RegisterTo(deps *http.Dependencies) {
	// 注册 Handler 到依赖容器
	deps.RegisterProvider("user.handler", r.handler)
}

// Register 注册用户领域路由
func (r *Routes) Register(router *gin.RouterGroup, deps *http.Dependencies) {
	// 获取 Handler（从依赖容器或直接使用）
	handler := r.handler
	if handler == nil {
		if h := deps.GetProvider("user.handler"); h != nil {
			handler = h.(*Handler)
		}
	}

	if handler == nil || handler.respHandler == nil {
		return
	}

	// 用户资源路由
	userGroup := router.Group("/users")
	{
		// GET /users/:user_id - 获取用户详情
		userGroup.GET("/:user_id", handler.GetUser)

		// PUT /users/:user_id - 更新用户
		userGroup.PUT("/:user_id", handler.UpdateUser)

		// POST /users/:user_id/password - 修改密码
		userGroup.POST("/:user_id/password", handler.ChangePassword)
	}
}
