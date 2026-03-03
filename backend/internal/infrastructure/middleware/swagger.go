// Package middleware 提供Swagger相关中间件
package middleware

import (
	"net/http"

	"go-ddd-scaffold/internal/pkg/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

// SwaggerConfig Swagger中间件配置
type SwaggerConfig struct {
	// 是否启用Swagger UI
	EnableUI bool
	// 是否启用JSON文档
	EnableJSON bool
	// UI路径前缀，默认为"/swagger"
	UIPath string
	// JSON路径，默认为"/swagger/doc.json"
	JSONPath string
}

// DefaultSwaggerConfig 返回默认的Swagger配置
func DefaultSwaggerConfig() SwaggerConfig {
	return SwaggerConfig{
		EnableUI:   true,
		EnableJSON: true,
		UIPath:     "/swagger",
		JSONPath:   "/swagger/doc.json",
	}
}

// Swagger 返回Swagger中间件
func Swagger() gin.HandlerFunc {
	return SwaggerWithConfig(DefaultSwaggerConfig())
}

// SwaggerWithConfig 根据配置返回Swagger中间件
func SwaggerWithConfig(config SwaggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 处理 /swagger 或 /swagger/ 重定向到 /swagger/
		if config.EnableUI && (path == config.UIPath || path == config.UIPath+"/") {
			c.Redirect(http.StatusFound, config.UIPath+"/index.html")
			c.Abort()
			return
		}

		// 处理Swagger UI静态文件请求
		if config.EnableUI && len(path) > len(config.UIPath) && path[:len(config.UIPath)+1] == config.UIPath+"/" {
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
			return
		}

		// 处理Swagger JSON文档请求
		if config.EnableJSON && path == config.JSONPath {
			doc, err := swag.ReadDoc()
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ServerErr(c.Request.Context()))
				return
			}
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, doc)
			return
		}

		// 如果不是Swagger相关请求，继续处理
		c.Next()
	}
}

// SwaggerUIOnly 只提供Swagger UI的中间件
func SwaggerUIOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 处理 /swagger 或 /swagger/ 重定向
		if path == "/swagger" || path == "/swagger/" {
			c.Redirect(http.StatusFound, "/swagger/index.html")
			c.Abort()
			return
		}
		// 处理静态文件
		if len(path) > 10 && path[:10] == "/swagger/" {
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
			return
		}
		c.Next()
	}
}

// SwaggerJSONOnly 只提供Swagger JSON文档的中间件
func SwaggerJSONOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/swagger/doc.json" {
			doc, err := swag.ReadDoc()
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ServerErr(c.Request.Context()))
				return
			}
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, doc)
			return
		}
		c.Next()
	}
}
