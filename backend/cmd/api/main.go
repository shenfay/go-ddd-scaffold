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

	_ "github.com/shenfay/go-ddd-scaffold/docs/swagger"
	"github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	logger "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/logging"
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

	logger := appLogger.Logger // 获取底层的 *zap.Logger

	logger.Info("Starting API server...")

	logger.Info("Configuration loaded",
		zap.String("env", env),
		zap.String("server_port", appConfig.Server.Port),
		zap.String("server_mode", appConfig.Server.Mode))

	// 3. 创建 Bootstrap（Composition Root）
	boot, err := bootstrap.NewBootstrap(appConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create bootstrap", zap.Error(err))
	}

	// 4. 初始化所有组件（Composition Root 的核心）
	ctx := context.Background()
	if err := boot.Initialize(ctx); err != nil {
		logger.Fatal("Failed to initialize components", zap.Error(err))
	}

	// 5. 启动应用
	if err := boot.Start(ctx); err != nil {
		logger.Fatal("Failed to start application", zap.Error(err))
	}

	// 6. 启动 HTTP 服务器
	go func() {
		addr := ":" + appConfig.Server.Port
		logger.Info("Server listening", zap.String("address", addr))
		// 使用 bootstrap 中的 router
		if err := http.ListenAndServe(addr, boot.GetRouter()); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 7. 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 8. 优雅关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := boot.Stop(shutdownCtx); err != nil {
		logger.Error("Failed to stop application", zap.Error(err))
	}
}
