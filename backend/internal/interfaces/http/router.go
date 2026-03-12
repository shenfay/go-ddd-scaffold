package http

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
	apperrors "github.com/shenfay/go-ddd-scaffold/shared/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RouteRegistrar 路由注册接口 - 各领域必须实现此接口
type RouteRegistrar func(router *gin.RouterGroup, handler *Handler)

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
	if r == nil || r.registrars == nil {
		return
	}
	r.registrars = append(r.registrars, registrar)
}

// Build 构建完整路由，注册所有领域的路由并返回 Gin 引擎
// 该方法只执行一次，确保路由不会重复注册
func (r *Router) Build(deps *Dependencies) *gin.Engine {
	if !r.initialized {
		// 设置 Handler 并触发所有已注册的领域路由
		r.handler = deps.Handler

		// 创建 logger
		logger := createLogger()

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
			registrar(apiGroup, r.handler)
		}

		// Health check endpoint (自动注入 TraceID)
		r.ginEngine.GET("/health", func(c *gin.Context) {
			traceID := middleware.GetTraceID(c)
			c.JSON(200, gin.H{
				"status":    "healthy",
				"trace_id":  traceID,
				"timestamp": time.Now().Unix(),
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
	pendingRegs  []func(*Router) // 存储 init 时注册的函数，延迟初始化时使用
)

// GetRouter 获取全局路由总线实例（单例）
// config 参数仅在首次调用时生效，必须由 main.go 提供配置
func GetRouter(config *RouterConfig) *Router {
	routerOnce.Do(func() {
		globalRouter = NewRouter(config)

		// 应用所有在 init 中注册的函数
		for _, regFunc := range pendingRegs {
			regFunc(globalRouter)
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

	// 否则返回一个临时的 Router 用于收集注册函数
	tempRouter := &Router{
		registrars: make([]RouteRegistrar, 0),
	}

	// 包装 Register 方法，使其能延迟执行
	wrappedReg := func(r *Router) {
		for _, reg := range tempRouter.registrars {
			r.Register(reg)
		}
	}

	pendingRegs = append(pendingRegs, wrappedReg)
	return tempRouter
}

// createLogger 创建开发环境 logger（支持彩色文本输出）
func createLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
