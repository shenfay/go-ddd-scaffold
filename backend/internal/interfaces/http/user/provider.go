package user

import (
	"github.com/gin-gonic/gin"
	http "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
)

// Provider 用户领域提供者
type Provider struct {
	handler *Handler
}

// NewProvider 创建用户领域提供者
func NewProvider(handler *Handler) *Provider {
	return &Provider{
		handler: handler,
	}
}

// RegisterTo 注册用户领域路由到依赖容器
func (p *Provider) RegisterTo(deps *http.Dependencies) {
	// 注册 Handler 到依赖容器
	deps.RegisterProvider("user.handler", p.handler)
}

// RegisterRoutes 注册用户领域路由
func (p *Provider) RegisterRoutes(router *gin.RouterGroup, deps *http.Dependencies) {
	// 获取 Handler（从依赖容器或直接使用）
	handler := p.handler
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
		// GET /users - 列出用户
		userGroup.GET("", handler.ListUsers)

		// GET /users/:user_id - 获取用户详情
		userGroup.GET("/:user_id", handler.GetUser)

		// PUT /users/:user_id - 更新用户
		userGroup.PUT("/:user_id", handler.UpdateUser)

		// PATCH /users/:user_id/activate - 激活用户
		userGroup.PATCH("/:user_id/activate", handler.ActivateUser)

		// PATCH /users/:user_id/deactivate - 禁用用户
		userGroup.PATCH("/:user_id/deactivate", handler.DeactivateUser)

		// POST /users/:user_id/password - 修改密码
		userGroup.POST("/:user_id/password", handler.ChangePassword)
	}
}
