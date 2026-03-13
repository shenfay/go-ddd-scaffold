package http

import (
	"sync"

	"github.com/gin-gonic/gin"
	docs "github.com/shenfay/go-ddd-scaffold/docs/swagger"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
	apperrors "github.com/shenfay/go-ddd-scaffold/shared/errors"
	"go.uber.org/zap"
)

// @title Go DDD Scaffold API
// @version 1.0
// @description Go DDD Scaffold API 文档 - 基于 DDD 和 CQRS 的企业级脚手架
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在 Header 中输入：Bearer {token}

// RouteRegistrar 路由注册接口 - 各领域必须实现此接口
type RouteRegistrar func(router *gin.RouterGroup, handler *Handler, deps *Dependencies)

// RouterConfig 路由配置
type RouterConfig struct {
	APIPrefix string // API 路径前缀，默认 "/api/v1"
	Port      string // 服务端口，默认 ":8080"
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
	return &Router{
		config:     config,
		ginEngine:  gin.New(),
		registrars: make([]RouteRegistrar, 0),
	}
}

// Register 注册领域路由到总线
// 各领域的 init() 函数会调用此方法自动注册
func (r *Router) Register(registrar RouteRegistrar) {
	if r == nil {
		return
	}
	// 如果全局路由已初始化，直接注册到全局路由
	if globalRouter != nil && r == globalRouter {
		r.registrars = append(r.registrars, registrar)
		return
	}
	// 否则暂存到 pendingRegs，等待初始化后再注册
	pendingRegs = append(pendingRegs, registrar)
}

// Build 构建完整路由，注册所有领域的路由并返回 Gin 引擎
// 该方法只执行一次，确保路由不会重复注册
func (r *Router) Build(deps *Dependencies, logger *zap.Logger) *gin.Engine {
	if !r.initialized {
		// 设置 Handler 并触发所有已注册的领域路由
		r.handler = deps.Handler

		// 应用全局中间件链（按正确顺序）
		// 顺序：TraceID → Gin Logger → Recovery → Error → Custom Logger with TraceID
		r.ginEngine.Use(
			middleware.TraceIDMiddleware(), // ① TraceID 追踪中间件
			gin.Logger(),                   // ② Gin 默认彩色日志中间件
			middleware.Recovery(logger),    // ③ Panic 恢复中间件
			middleware.Error(
				apperrors.NewErrorMapper(),
				logger,
			), // ④ 错误处理中间件
			middleware.LoggerWithTrace(logger), // ⑤ 带 TraceID 的自定义日志
		)

		// 创建 API 路由组（中间件已应用，所以路由组会继承中间件）
		apiGroup := r.ginEngine.Group(r.config.APIPrefix)
		for _, registrar := range r.registrars {
			registrar(apiGroup, r.handler, deps)
		}

		// Health check endpoint (自动注入 TraceID)
		r.ginEngine.GET("/health", func(c *gin.Context) {
			traceID := middleware.GetTraceID(c)
			c.JSON(200, gin.H{
				"status":    "healthy",
				"trace_id":  traceID,
				"timestamp": util.Now().Timestamp(),
			})
		})

		// Swagger UI (仅开发环境)
		if gin.Mode() == gin.DebugMode {
			r.setupSwagger()
		}

		r.initialized = true
	}

	return r.ginEngine
}

// GetEngine 获取底层的 Gin 引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.ginEngine
}

// setupSwagger 设置 Swagger UI（仅开发环境）
func (r *Router) setupSwagger() {
	// 设置 Swagger 基础路径
	docs.SwaggerInfo.BasePath = r.config.APIPrefix

	// 使用中间件方式注册 Swagger
	r.ginEngine.GET("/swagger/*any", middleware.Swagger())
}

// Dependencies 路由依赖注入
type Dependencies struct {
	Handler *Handler
	// 使用 map 存储各领域依赖，避免结构体字段膨胀
	providers map[string]interface{}
}

// NewDependencies 创建依赖注入容器
func NewDependencies(handler *Handler) *Dependencies {
	return &Dependencies{
		Handler:   handler,
		providers: make(map[string]interface{}),
	}
}

// RegisterProvider 注册领域依赖提供者
func (d *Dependencies) RegisterProvider(name string, provider interface{}) {
	if d.providers == nil {
		d.providers = make(map[string]interface{})
	}
	d.providers[name] = provider
}

// GetProvider 获取领域依赖提供者
func (d *Dependencies) GetProvider(name string) interface{} {
	if d.providers == nil {
		return nil
	}
	return d.providers[name]
}

// globalRouter 全局路由总线实例（单例）
var (
	globalRouter *Router
	routerOnce   sync.Once
	pendingRegs  []RouteRegistrar // 存储 init 时注册的函数，延迟初始化时使用
)

// GetRouter 获取全局路由总线实例（单例）
// config 参数仅在首次调用时生效，必须由 main.go 提供配置
func GetRouter(config *RouterConfig) *Router {
	routerOnce.Do(func() {
		globalRouter = NewRouter(config)

		// 将 pendingRegs 中的注册函数转移到 globalRouter
		for _, reg := range pendingRegs {
			globalRouter.registrars = append(globalRouter.registrars, reg)
		}
		pendingRegs = nil // 清理内存
	})
	return globalRouter
}

// MustGetRouter 获取全局路由总线实例（用于模块注册）
// 如果尚未初始化，会将注册函数暂存到 pendingRegs
func MustGetRouter() *Router {
	// 如果已经初始化，直接返回
	if globalRouter != nil {
		return globalRouter
	}

	// 返回一个临时的 Router 用于收集注册函数
	return &Router{
		registrars: make([]RouteRegistrar, 0),
	}
}
