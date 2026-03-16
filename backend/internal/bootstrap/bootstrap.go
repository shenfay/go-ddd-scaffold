package bootstrap

import (
	"context"

	"github.com/gin-gonic/gin"
	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/container"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	http "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	authHttp "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/auth"
	userHttp "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/user"
	apperrors "github.com/shenfay/go-ddd-scaffold/shared/errors"
	"go.uber.org/zap"
)

// Bootstrap 应用启动器
// 负责创建和注册所有基础设施组件、领域服务、应用服务
// 是应用的 Composition Root（组合根）
type Bootstrap struct {
	container container.Container
	config    *config.AppConfig
	logger    *zap.Logger
	httpDeps  *http.Dependencies

	// === 用户领域组件（按领域分组）===
	user struct {
		service      *userApp.UserServiceImpl
		eventHandler *userApp.UserEventHandler
	}

	// === 认证领域组件（按领域分组）===
	auth struct {
		jwtService  *auth.JWTService
		authService authApp.AuthService
	}
}

// NewBootstrap 创建应用启动器
func NewBootstrap(cfg *config.AppConfig, logger *zap.Logger) (*Bootstrap, error) {
	// 创建容器
	c, err := container.NewContainer(cfg, logger)
	if err != nil {
		return nil, err
	}

	return &Bootstrap{
		container: c,
		config:    cfg,
		logger:    logger,
		httpDeps:  http.NewDependencies(nil), // 先创建空的 Dependencies，Handler 在后面赋值
	}, nil
}

// GetContainer 获取容器实例
func (b *Bootstrap) GetContainer() container.Container {
	return b.container
}

// GetRouter 获取路由引擎
func (b *Bootstrap) GetRouter() *gin.Engine {
	router := http.GetRouter(nil)
	if router == nil {
		return gin.New()
	}
	return router.GetEngine()
}

// Initialize 初始化所有组件
// 这是 Composition Root 的核心方法，按顺序创建并注册所有依赖
func (b *Bootstrap) Initialize(ctx context.Context) error {
	b.logger.Info("Initializing application components...")

	// === 1. 初始化基础设施层 ===
	if err := b.initializeInfrastructure(ctx); err != nil {
		return err
	}

	// === 2. 初始化领域层 ===
	if err := b.initializeDomains(ctx); err != nil {
		return err
	}

	// === 3. 初始化应用层 ===
	if err := b.initializeApplication(ctx); err != nil {
		return err
	}

	// === 4. 初始化接口层 ===
	if err := b.initializeInterfaces(ctx); err != nil {
		return err
	}

	b.logger.Info("All components initialized successfully")
	return nil
}

// initializeInfrastructure 初始化基础设施层
func (b *Bootstrap) initializeInfrastructure(ctx context.Context) error {
	b.logger.Info("Initializing infrastructure layer...")

	// 基础设施组件已经在 container.NewContainer() 中创建
	// 包括：Database, Redis, Cache, Logger, Router

	// TODO: 如果需要额外的基础设施组件，在这里添加
	// 例如：消息队列、文件存储、外部 API 客户端等

	return nil
}

// initializeDomains 初始化领域层
func (b *Bootstrap) initializeDomains(ctx context.Context) error {
	b.logger.Info("Initializing domain layer...")

	// === 用户领域 ===
	if err := b.initUserDomain(ctx); err != nil {
		return err
	}

	// === 认证领域 ===
	if err := b.initAuthDomain(ctx); err != nil {
		return err
	}

	// TODO: 其他领域
	// if err := b.initTenantDomain(ctx); err != nil {
	//     return err
	// }

	return nil
}

// initializeApplication 初始化应用层
func (b *Bootstrap) initializeApplication(ctx context.Context) error {
	b.logger.Info("Initializing application layer...")

	// 应用层服务（协调器）在领域初始化时已经创建并注册
	// 这里可以进行额外的应用级配置

	return nil
}

// initializeInterfaces 初始化接口层
func (b *Bootstrap) initializeInterfaces(ctx context.Context) error {
	b.logger.Info("Initializing interface layer...")

	// 创建 HTTP Handler（响应处理）
	respHandler := http.NewHandler(apperrors.NewErrorMapper())

	// 创建用户领域 HTTP Handler（业务处理）
	// 直接使用 Bootstrap 中持有的领域组件，类型安全
	userHandler := userHttp.NewHandler(b.user.service)

	// 获取 router 并构建路由
	router := http.GetRouter(&http.RouterConfig{
		APIPrefix: "/api/v1",
		Port:      b.config.Server.Port,
	})

	// 创建依赖容器
	deps := http.NewDependencies(respHandler)

	// 注册用户领域提供者并注册路由
	userProvider := userHttp.NewProvider(userHandler)
	userProvider.RegisterTo(deps)

	// 手动注册用户路由（替代 init 自动注册）
	router.Register(func(routerGroup *gin.RouterGroup, handler *http.Handler, deps *http.Dependencies) {
		userProvider.RegisterRoutes(routerGroup, deps)
	})

	// === 注册认证领域路由 ===
	authHandler := authHttp.NewHandler(
		b.auth.authService,
		respHandler,
	)
	authProvider := authHttp.NewProvider(authHandler, b.auth.jwtService)

	// 手动注册认证路由
	router.Register(func(routerGroup *gin.RouterGroup, handler *http.Handler, deps *http.Dependencies) {
		authProvider.ProvideRoutes(routerGroup)
	})

	// 构建路由（触发所有领域的注册）
	_ = router.Build(deps, b.logger)

	return nil
}

// Start 启动应用
func (b *Bootstrap) Start(ctx context.Context) error {
	b.logger.Info("Starting application...")
	return b.container.Start(ctx)
}

// Stop 停止应用
func (b *Bootstrap) Stop(ctx context.Context) error {
	b.logger.Info("Stopping application...")
	return b.container.Stop(ctx)
}
