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

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	_ "github.com/shenfay/go-ddd-scaffold/docs/swagger"
	"github.com/shenfay/go-ddd-scaffold/internal/container"
	"github.com/shenfay/go-ddd-scaffold/internal/factory"
	authInfra "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	logger "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/logging"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/task_queue"
	"github.com/shenfay/go-ddd-scaffold/internal/provider"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	env := os.Getenv("ENV_MODE")
	if env == "" {
		env = "development"
	}

	configLoader := config.NewConfigLoader(nil)
	appConfig, err := configLoader.Load(env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 创建正式 logger（双输出模式：控制台 + 文件）
	logConfig := &config.LoggingConfig{
		Level:      appConfig.Logging.Level,
		Format:     appConfig.Logging.Format,
		File:       appConfig.Logging.File,
		MaxSize:    appConfig.Logging.MaxSize,
		MaxBackups: appConfig.Logging.MaxBackups,
		MaxAge:     appConfig.Logging.MaxAge,
	}
	appLogger, err := logger.New(logConfig)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer appLogger.Sync()

	logger := appLogger.Logger

	logger.Info("Starting API server...")

	logger.Info("Configuration loaded",
		zap.String("env", env),
		zap.String("server_port", appConfig.Server.Port),
		zap.String("server_mode", appConfig.Server.Mode))

	ctx := context.Background()

	// 3. 创建容器（系统级基础设施）
	// 这是 Composition Root 的第一步：创建系统级基础设施容器
	cont, err := container.NewContainer(appConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create container", zap.Error(err))
	}

	// 4. 创建领域基础设施提供者（领域特定基础设施）
	// 领域特定的基础设施组件在领域内初始化，而非全局
	domainInfraProvider := provider.NewDomainInfrastructureProvider(appConfig)
	// 注入 Redis 客户端到领域基础设施提供者（用于 JWT 服务）
	domainInfraProvider.SetRedisClient(cont.GetRedis())

	// 5. 创建领域服务（领域层 + 应用层）
	// 使用工厂函数按依赖顺序创建服务
	serviceDeps := &factory.ServiceDependencies{
		Container:            cont,
		DomainInfrastructure: domainInfraProvider,
		Logger:               logger,
	}

	// 5.1 创建用户领域服务（从容器获取系统级基础设施，从 provider 获取领域特定基础设施）
	userDomainServices, err := factory.CreateUserDomainServices(serviceDeps)
	if err != nil {
		logger.Fatal("Failed to create user domain services", zap.Error(err))
	}

	// 5.2 创建认证领域服务（从容器获取系统级基础设施，从 provider 获取领域特定基础设施）
	authDomainServices, err := factory.CreateAuthDomainServices(serviceDeps)
	if err != nil {
		logger.Fatal("Failed to create auth domain services", zap.Error(err))
	}

	// 5. 构建 HTTP 接口（接口层）
	httpInterfaces, err := factory.BuildHTTPInterfaces(
		appConfig,
		logger,
		userDomainServices.UserService,
		authDomainServices.AuthService,
		authDomainServices.AuthJWTService.(*authInfra.JWTService), // 类型断言
	)
	if err != nil {
		logger.Fatal("Failed to build HTTP interfaces", zap.Error(err))
	}

	// 6. 注册事件处理器（基础设施层）
	// TODO: 如果需要，可以在这里注册领域事件处理器
	// if err := factory.RegisterEventHandlers(eventBus, userDomainServices.UserSideEffectHandler, logger); err != nil {
	//     logger.Fatal("Failed to register event handlers", zap.Error(err))
	// }

	// 7. 启动应用（启动所有生命周期钩子）
	if err := cont.Start(ctx); err != nil {
		logger.Fatal("Failed to start application", zap.Error(err))
	}

	// 8. 启动 asynq worker（后台任务处理）
	go func() {
		if err := startAsynqWorker(appConfig, logger); err != nil {
			logger.Error("Failed to start asynq worker", zap.Error(err))
		}
	}()

	// 9. 启动 HTTP 服务器
	go func() {
		addr := ":" + appConfig.Server.Port
		logger.Info("Server listening", zap.String("address", addr))
		if err := http.ListenAndServe(addr, httpInterfaces.Router); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 10. 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 11. 优雅关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := cont.Stop(shutdownCtx); err != nil {
		logger.Error("Failed to stop application", zap.Error(err))
	}
}

// startAsynqWorker 启动 asynq worker（后台任务处理）
func startAsynqWorker(cfg *config.AppConfig, logger *zap.Logger) error {
	logger.Info("Starting asynq worker...")

	// 创建 asynq 服务器
	asynqServer := task_queue.NewServer(task_queue.Config{
		RedisAddr:     cfg.Redis.Addr,
		RedisPassword: cfg.Redis.Password,
		RedisDB:       cfg.Redis.DB,
	})

	// 创建处理器
	processor := task_queue.NewProcessor(
		logger.Named("asynq"),
		// TODO: 注册具体的事件处理器
	)

	// 创建 mux 来路由不同类型的任务
	mux := asynq.NewServeMux()
	mux.HandleFunc(task_queue.TaskTypeDomainEvent, processor.ProcessTask)

	// 启动 worker
	logger.Info("Asynq worker started")
	if err := asynqServer.Run(mux); err != nil {
		return err
	}

	return nil
}
