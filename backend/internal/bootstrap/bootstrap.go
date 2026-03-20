package bootstrap

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	authApp "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/bootstrap/helpers"
	"github.com/shenfay/go-ddd-scaffold/internal/container"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/task_queue"
	http "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
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
	eventBus  kernel.EventBus // 事件总线

	// asynq worker
	asynqServer *asynq.Server

	// === 用户领域组件（按领域分组）===
	user struct {
		service           *userApp.UserServiceImpl
		sideEffectHandler *userEvent.SideEffectHandler
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
		httpDeps:  http.NewDependencies(nil),  // 先创建空的 Dependencies，Handler 在后面赋值
		eventBus:  kernel.NewSimpleEventBus(), // 创建事件总线
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

	// 使用 helpers 注册事件处理器
	_, err := helpers.RegisterEventHandlers(
		b.eventBus,
		b.container,
		b.user.sideEffectHandler,
		b.logger,
	)
	if err != nil {
		return err
	}

	// 初始化 asynq worker
	b.initAsynqWorker()

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

	// 使用 helpers 构建 HTTP 接口
	httpInterfaces, err := helpers.BuildHTTPInterfaces(
		b.config,
		b.logger,
		b.user.service,
		b.auth.authService,
		b.auth.jwtService,
	)
	if err != nil {
		return err
	}

	// 保存 dependencies 供后续使用
	b.httpDeps = httpInterfaces.Deps

	return nil
}

// Start 启动应用
func (b *Bootstrap) Start(ctx context.Context) error {
	b.logger.Info("Starting application...")
	return b.container.Start(ctx)
}

// initAsynqWorker 初始化 asynq worker
func (b *Bootstrap) initAsynqWorker() {
	b.logger.Info("Initializing asynq worker...")

	// 创建 asynq 服务器
	asynqServer := task_queue.NewServer(task_queue.Config{
		RedisAddr:     b.config.Redis.Addr,
		RedisPassword: b.config.Redis.Password,
		RedisDB:       b.config.Redis.DB,
	})

	// 创建处理器
	processor := task_queue.NewProcessor(
		b.logger.Named("asynq"),
		// TODO: 注册具体的事件处理器
	)

	// 创建 mux 来路由不同类型的任务
	mux := asynq.NewServeMux()
	mux.HandleFunc(task_queue.TaskTypeDomainEvent, processor.ProcessTask)

	// 保存引用
	b.asynqServer = asynqServer

	// 后台启动 worker
	go func() {
		b.logger.Info("Starting asynq worker...")
		if err := asynqServer.Run(mux); err != nil {
			b.logger.Error("Asynq worker failed", zap.Error(err))
		}
	}()

	b.logger.Info("Asynq worker initialized successfully")
}

// Stop 停止应用
func (b *Bootstrap) Stop(ctx context.Context) error {
	b.logger.Info("Stopping application...")
	return b.container.Stop(ctx)
}
