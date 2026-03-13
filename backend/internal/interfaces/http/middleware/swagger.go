package middleware

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerConfig Swagger 中间件配置
type SwaggerConfig struct {
	Enabled bool   // 是否启用
	URL     string // Swagger JSON 路径
}

// DefaultSwaggerConfig 默认配置
var DefaultSwaggerConfig = SwaggerConfig{
	Enabled: true,
	URL:     "/swagger/doc.json",
}

// Swagger Swagger UI 中间件
// 该中间件仅在开发环境生效，提供 API 文档浏览界面
// 使用方式：r.ginEngine.GET("/swagger/*any", middleware.Swagger())
// @Summary Swagger API 文档
// @Description 提供 API 文档浏览界面
func Swagger(config ...SwaggerConfig) gin.HandlerFunc {
	cfg := DefaultSwaggerConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL(cfg.URL),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.DeepLinking(true),
		ginSwagger.DocExpansion("none"),
	)
}
