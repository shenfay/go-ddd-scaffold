package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 默认 CORS 配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           3600,
	}
}

// CORSMiddleware CORS 中间件
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 检查来源是否允许
		allowed := false
		for _, o := range config.AllowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{
				"code":    "CORS_ERROR",
				"message": "Origin not allowed",
			})
			return
		}

		// 设置 CORS 响应头
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowedMethods))
		c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowedHeaders))
		c.Header("Access-Control-Allow-Credentials", boolToString(config.AllowCredentials))
		c.Header("Access-Control-Max-Age", intToString(config.MaxAge))

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatusJSON(204, nil)
			return
		}

		c.Next()
	}
}

// joinStrings 连接字符串数组
func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

// boolToString 布尔转字符串
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// intToString 整数转字符串
func intToString(i int) string {
	return string(rune(i))
}
