package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware"
	"github.com/shenfay/go-ddd-scaffold/pkg/metrics"
)

// Middlewares 注册全局中间件
func Middlewares(engine *gin.Engine, m *metrics.Metrics) {
	// CORS 中间件
	engine.Use(middleware.CORSMiddleware(middleware.CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	// Prometheus 监控中间件
	engine.Use(middleware.PrometheusMiddleware(m))

	// 请求日志中间件（Gin默认）
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))

	// Recovery 中间件（必须）
	engine.Use(gin.Recovery())
}
