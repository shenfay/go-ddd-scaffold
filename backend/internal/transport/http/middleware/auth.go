package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/app/authentication"
)

// JWTAuthConfig JWT 认证中间件配置
type JWTAuthConfig struct {
	TokenService authentication.TokenService
}

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware(config JWTAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := config.TokenService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_TOKEN",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息注入到上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}
