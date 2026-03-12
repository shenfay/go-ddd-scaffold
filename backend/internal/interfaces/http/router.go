package http

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// RouteRegistrar 路由注册接口 - 各领域必须实现此接口
type RouteRegistrar func(router *gin.RouterGroup, handler *Handler)

// RouterConfig 路由配置
type RouterConfig struct {
	APIPrefix string // API 路径前缀，默认 "/api/v1"
	Port      string // 服务端口，默认 ":8080"
}

// DefaultRouterConfig 返回默认配置
func DefaultRouterConfig() *RouterConfig {
	return &RouterConfig{
		APIPrefix: "/api/v1",
		Port:      ":8080",
	}
}

// Router 路由总线，负责收集和注册所有领域的路由
type Router struct {
	config      *RouterConfig
	ginEngine   *gin.Engine
	registrars  []RouteRegistrar
	handler     *Handler
	initialized bool // 防止重复初始化
}

// NewRouter 创建路由总线
func NewRouter(config *RouterConfig) *Router {
	if config == nil {
		config = DefaultRouterConfig()
	}
	return &Router{
		config:     config,
		ginEngine:  gin.New(),
		registrars: make([]RouteRegistrar, 0),
	}
}

// Register 注册领域路由到总线
// 各领域的 init() 函数会调用此方法自动注册
func (r *Router) Register(registrar RouteRegistrar) {
	r.registrars = append(r.registrars, registrar)
}

// Build 构建完整路由，注册所有领域的路由并返回 Gin 引擎
// 该方法只执行一次，确保路由不会重复注册
func (r *Router) Build(deps *Dependencies) *gin.Engine {
	if !r.initialized {
		// 设置 Handler 并触发所有已注册的领域路由
		r.handler = deps.Handler

		apiGroup := r.ginEngine.Group(r.config.APIPrefix)
		for _, registrar := range r.registrars {
			registrar(apiGroup, r.handler)
		}

		// 应用全局中间件
		r.ginEngine.Use(
			gin.Recovery(),
			gin.Logger(),
		)

		// Health check endpoint
		r.ginEngine.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "healthy",
			})
		})

		r.initialized = true
	}

	return r.ginEngine
}

// GetEngine 获取底层的 Gin 引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.ginEngine
}

// Dependencies 路由依赖注入
type Dependencies struct {
	Handler *Handler
	// 后续可以添加更多依赖，如：
	// UserService    *app.UserApplicationService
	// TenantService  *app.TenantApplicationService
	// AuthMiddleware gin.HandlerFunc
}

// NewDependencies 创建依赖注入容器
func NewDependencies(handler *Handler) *Dependencies {
	return &Dependencies{
		Handler: handler,
	}
}

// globalRouter 全局路由总线实例（单例）
var (
	globalRouter *Router
	routerOnce   sync.Once
)

// GetRouter 获取全局路由总线实例（单例）
// 各领域的 init() 函数通过此函数获取路由器并注册路由
func GetRouter() *Router {
	routerOnce.Do(func() {
		globalRouter = NewRouter(nil)
	})
	return globalRouter
}
