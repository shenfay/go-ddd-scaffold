package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/app"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/eventstore"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/internal/module"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
	"go.uber.org/zap"
)

// Server HTTP 服务器结构
type Server struct {
	infra   *app.Infrastructure
	logger  *zap.Logger
	modules []app.Module
}

// NewServer 创建 HTTP 服务器
func NewServer(infra *app.Infrastructure, logger *zap.Logger) *Server {
	return &Server{
		infra:  infra,
		logger: logger,
	}
}

// Run 运行服务器（包含完整的启动流程）
func (s *Server) Run() {
	s.createModules()
	router := s.setupRouter()
	s.startServer(router)
}

// createModules 创建并注册模块
func (s *Server) createModules() {
	// 创建用户模块
	userMod := module.NewUserModule(s.infra)
	s.logger.Info("User module created", zap.String("module", userMod.Name()))

	// 创建认证模块
	authMod := module.NewAuthModule(s.infra)
	s.logger.Info("Auth module created", zap.String("module", authMod.Name()))

	s.modules = []app.Module{authMod, userMod}

	// 注册事件订阅
	for _, m := range s.modules {
		if em, ok := m.(app.EventModule); ok {
			em.RegisterSubscriptions(s.infra.EventBus)
		}
	}

	// 启动 Outbox Processor
	s.startOutboxProcessor()
}

// startOutboxProcessor 启动领域事件 Outbox 处理器
func (s *Server) startOutboxProcessor() {
	outboxProcessor := eventstore.NewOutboxProcessor(
		s.infra.DB,
		s.infra.TaskPublisher,
		s.logger,
	)
	go func() {
		if err := outboxProcessor.Start(context.Background()); err != nil {
			s.logger.Error("Outbox processor stopped", zap.Error(err))
		}
	}()
	s.logger.Info("Outbox processor started")
}

// setupRouter 设置 HTTP 路由
func (s *Server) setupRouter() *gin.Engine {
	router := gin.New()

	// 中间件
	mwFactory := middleware.NewMiddlewareFactory(&middleware.MiddlewareConfig{
		Logger:      s.logger,
		ErrorMapper: s.infra.ErrorMapper,
	})

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
	for _, m := range s.modules {
		if h, ok := m.(app.HTTPModule); ok {
			h.RegisterHTTP(api)
			s.logger.Info("HTTP routes registered", zap.String("module", m.Name()))
		}
	}

	return router
}

// startServer 启动 HTTP 服务器（包含优雅关闭）
func (s *Server) startServer(router *gin.Engine) {
	addr := ":" + s.infra.Config.Server.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  s.infra.Config.Server.ReadTimeout,
		WriteTimeout: s.infra.Config.Server.WriteTimeout,
	}

	go func() {
		s.logger.Info("Server listening", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("Failed to shutdown server", zap.Error(err))
	}

	s.logger.Info("Server stopped")
}
