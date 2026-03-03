// Package middleware 提供CORS中间件配置
package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// DefaultCORSConfig 返回默认的CORS配置
func DefaultCORSConfig() cors.Config {
	return cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
			"X-Requested-With", "X-Request-Id", "X-Timestamp",
		},
		MaxAge: 12 * time.Hour,
	}
}

// CORS 返回CORS中间件处理器
func CORS() gin.HandlerFunc {
	return cors.New(DefaultCORSConfig())
}

// CORSWithConfig 根据配置返回CORS中间件
func CORSWithConfig(config cors.Config) gin.HandlerFunc {
	return cors.New(config)
}
