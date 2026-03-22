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

	"github.com/gin-gonic/gin"
	_ "github.com/shenfay/go-ddd-scaffold/docs/swagger"
	"github.com/shenfay/go-ddd-scaffold/internal/bootstrap"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/logging"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/internal/module"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
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
	appLogger, err := logging.New(logConfig)
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

	// 3. 创建基础设施（替代 Container）
	infra, cleanup, err := bootstrap.NewInfra(appConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create infrastructure", zap.Error(err))
	}
	defer cleanup()

	// 4. 创建模块（替代 Factory）
	// 4.1 创建用户模块
	userMod := module.NewUserModule(infra)
	logger.Info("User module created", zap.String("module", userMod.Name()))

	// 4.2 创建认证模块
	authMod := module.NewAuthModule(infra)
	logger.Info("Auth module created", zap.String("module", authMod.Name()))

	modules := []bootstrap.Module{authMod, userMod}

	// 4.3 注册事件订阅
	for _, m := range modules {
		if em, ok := m.(bootstrap.EventModule); ok {
			em.RegisterSubscriptions(infra.EventBus)
			logger.Info("Event subscriptions registered", zap.String("module", m.Name()))
		}
	}

	// 5. 构建路由和中间件
	router := gin.New()

	// 5.1 创建中间件工厂并应用全局中间件链
	// 使用 Infra 中的共享组件（ErrorMapper 多领域共享）
	mwFactory := middleware.NewMiddlewareFactory(&middleware.MiddlewareConfig{
		Logger:      logger,
		ErrorMapper: infra.ErrorMapper,
	})

	// 应用全局中间件链（按正确顺序）
	// 顺序：TraceID → Gin Logger → Recovery → Error → LoggerWithTrace
	for _, mw := range mwFactory.Chain() {
		switch v := mw.(type) {
		case gin.HandlerFunc:
			router.Use(v)
		case interface{ Handler() gin.HandlerFunc }:
			// 适配 http.Handler 到 gin.HandlerFunc
			router.Use(v.Handler())
		}
	}

	// 5.2 Health check endpoint (自动注入 TraceID)
	router.GET("/health", func(c *gin.Context) {
		traceID := middleware.GetTraceID(c)
		c.JSON(200, gin.H{
			"status":    "healthy",
			"trace_id":  traceID,
			"timestamp": util.Now().Timestamp(),
		})
	})

	// 5.3 Swagger UI (仅开发环境)
	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", middleware.Swagger())
	}

	// 5.4 创建 API 路由组并注册模块路由
	api := router.Group("/api/v1")

	// TODO(Task4): 当前模块 RegisterHTTP 内部仍创建 Dependencies(nil)，
	// 后续应改为 main 传入 respHandler 或由模块自行管理
	for _, m := range modules {
		if h, ok := m.(bootstrap.HTTPModule); ok {
			h.RegisterHTTP(api)
			logger.Info("HTTP routes registered", zap.String("module", m.Name()))
		}
	}

	// 6. 启动 HTTP 服务器
	addr := ":" + appConfig.Server.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  appConfig.Server.ReadTimeout,
		WriteTimeout: appConfig.Server.WriteTimeout,
	}

	go func() {
		logger.Info("Server listening", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error", zap.Error(err))
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

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown server gracefully", zap.Error(err))
	}

	logger.Info("Server stopped")
}
