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
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/eventstore"
	logging "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/log"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/internal/module"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
	"go.uber.org/zap"
)

// Application API 应用结构
type Application struct {
	config  *config.AppConfig
	logger  *zap.Logger
	infra   *bootstrap.Infra
	cleanup func()
	modules []bootstrap.Module
}

// NewApplication 创建并初始化应用
func NewApplication() *Application {
	app := &Application{}
	app.init()
	return app
}

// init 初始化应用（加载配置、创建基础设施）
func (a *Application) init() {
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
	a.config = appConfig

	// 2. 创建 logger
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

	a.logger = appLogger.Logger

	a.logger.Info("Starting API server...")
	a.logger.Info("Configuration loaded",
		zap.String("env", env),
		zap.String("server_port", appConfig.Server.Port),
		zap.String("server_mode", appConfig.Server.Mode))

	// 3. 创建基础设施
	infra, cleanup, err := bootstrap.NewInfra(appConfig, a.logger)
	if err != nil {
		a.logger.Fatal("Failed to create infrastructure", zap.Error(err))
	}
	a.infra = infra
	a.cleanup = cleanup
}

// Cleanup 清理资源
func (a *Application) Cleanup() {
	if a.cleanup != nil {
		a.cleanup()
	}
}

// CreateModules 创建并注册模块
func (a *Application) CreateModules() {
	// 创建用户模块
	userMod := module.NewUserModule(a.infra)
	a.logger.Info("User module created", zap.String("module", userMod.Name()))

	// 创建认证模块
	authMod := module.NewAuthModule(a.infra)
	a.logger.Info("Auth module created", zap.String("module", authMod.Name()))

	a.modules = []bootstrap.Module{authMod, userMod}

	// 注册事件订阅
	for _, m := range a.modules {
		if em, ok := m.(bootstrap.EventModule); ok {
			em.RegisterSubscriptions(a.infra.EventBus)
			a.logger.Info("Event subscriptions registered", zap.String("module", m.Name()))
		}
	}

	// 启动 Outbox Processor
	outboxProcessor := eventstore.NewOutboxProcessor(
		a.infra.DB,
		a.infra.TaskPublisher,
		a.logger,
	)
	go func() {
		if err := outboxProcessor.Start(context.Background()); err != nil {
			a.logger.Error("Outbox processor stopped with error", zap.Error(err))
		}
	}()
	a.logger.Info("Outbox processor started")
}

// SetupRouter 设置 HTTP 路由
func (a *Application) SetupRouter() *gin.Engine {
	router := gin.New()

	// 创建中间件工厂
	mwFactory := middleware.NewMiddlewareFactory(&middleware.MiddlewareConfig{
		Logger:      a.logger,
		ErrorMapper: a.infra.ErrorMapper,
	})

	// 应用全局中间件链
	for _, mw := range mwFactory.Chain() {
		switch v := mw.(type) {
		case gin.HandlerFunc:
			router.Use(v)
		case interface{ Handler() gin.HandlerFunc }:
			router.Use(v.Handler())
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		traceID := middleware.GetTraceID(c)
		c.JSON(200, gin.H{
			"status":    "healthy",
			"trace_id":  traceID,
			"timestamp": util.Now().Timestamp(),
		})
	})

	// Swagger UI (仅开发环境)
	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", middleware.Swagger())
	}

	// 注册模块路由
	api := router.Group("/api/v1")
	for _, m := range a.modules {
		if h, ok := m.(bootstrap.HTTPModule); ok {
			h.RegisterHTTP(api)
			a.logger.Info("HTTP routes registered", zap.String("module", m.Name()))
		}
	}

	return router
}

// Run 启动 HTTP 服务器
func (a *Application) Run(router *gin.Engine) {
	addr := ":" + a.config.Server.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
	}

	go func() {
		a.logger.Info("Server listening", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	// 优雅关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("Failed to shutdown server gracefully", zap.Error(err))
	}

	a.logger.Info("Server stopped")
}
