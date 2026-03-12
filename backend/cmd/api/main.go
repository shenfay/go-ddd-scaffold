package main

import (
	"net/http"
	"os"

	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	httpiface "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	apperrors "github.com/shenfay/go-ddd-scaffold/shared/errors"
	"go.uber.org/zap"
)

func main() {
	// 创建logger
	mainLogger, _ := zap.NewDevelopment()
	defer mainLogger.Sync()

	mainLogger.Info("Starting API server...")

	// 获取环境变量
	env := os.Getenv("ENV_MODE")
	if env == "" {
		env = "development"
	}

	// 加载配置
	configLoader := config.NewConfigLoader(mainLogger)
	appConfig, err := configLoader.Load(env)
	if err != nil {
		mainLogger.Fatal("Failed to load config", zap.Error(err))
	}

	mainLogger.Info("Configuration loaded",
		zap.String("env", env),
		zap.String("server_port", appConfig.Server.Port),
		zap.String("server_mode", appConfig.Server.Mode))

	// 初始化组件
	errorMapper := apperrors.NewErrorMapper()
	handler := httpiface.NewHandler(errorMapper)

	// 创建带配置的路由
	routerConfig := &httpiface.RouterConfig{
		APIPrefix: "/api",
		Port:      ":" + appConfig.Server.Port,
	}
	router := httpiface.GetRouter(routerConfig)

	// 构建依赖注入容器
	deps := httpiface.NewDependencies(handler)

	// 构建完整路由（自动触发所有已注册的领域路由）
	ginEngine := router.Build(deps)

	// 启动配置监听
	configLoader.WatchConfig(func(newConfig *config.AppConfig) {
		mainLogger.Info("Configuration reloaded",
			zap.String("server_port", newConfig.Server.Port),
			zap.String("server_mode", newConfig.Server.Mode))
	})

	// 启动服务器
	mainLogger.Info("Server listening", zap.String("address", ":"+appConfig.Server.Port))
	if err := ginEngine.Run(":" + appConfig.Server.Port); err != nil && err != http.ErrServerClosed {
		mainLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
