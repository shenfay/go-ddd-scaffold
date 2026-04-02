package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware(tokenService interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrorCodeUnauthorized,
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		// 提取 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrorCodeUnauthorized,
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// TODO: 验证 Token（需要 tokenService 实现）
		_ = tokenString
		// claims, err := tokenService.ValidateAccessToken(tokenString)
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"code":    errors.ErrorCodeInvalidToken,
		// 		"message": err.Error(),
		// 	})
		// 	c.Abort()
		// 	return
		// }

		// // 将用户信息存入上下文
		// c.Set("user_id", claims.UserID)
		// c.Set("user_email", claims.Email)

		c.Next()
	}
}
