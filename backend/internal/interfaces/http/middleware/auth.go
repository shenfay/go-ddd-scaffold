package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/auth"
	"go.uber.org/zap"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(tokenService auth.TokenService) gin.HandlerFunc {
	logger := zap.L().Named("auth.middleware")

	return func(c *gin.Context) {
		// 1. 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    common.CodeUnauthorized,
				"message": "缺少认证令牌",
			})
			return
		}

		// 2. 提取 Bearer Token
		token := extractBearerToken(authHeader)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    common.CodeUnauthorized,
				"message": "无效的认证格式",
			})
			return
		}

		// 3. 验证 Token
		claims, err := tokenService.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    common.CodeTokenInvalid,
				"message": "令牌无效或已过期",
			})
			return
		}

		// 4. 检查 Token 是否在黑名单中
		isBlacklisted, err := tokenService.IsTokenBlacklisted(token)
		if err != nil {
			// Redis 检查失败，记录错误但不阻止请求（降级处理）
			logger.Warn("failed to check token blacklist",
				zap.String("token", token[:20]+"..."),
				zap.Error(err),
			)
		} else if isBlacklisted {
			logger.Info("blacklisted token detected",
				zap.String("token", token[:20]+"..."),
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    common.CodeTokenInvalid,
				"message": "令牌已注销",
			})
			return
		}

		// 5. 将用户信息注入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// extractBearerToken 从 Authorization Header 中提取 Bearer Token
func extractBearerToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

// GetUserIDFromContext 从上下文获取用户 ID
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

// GetUsernameFromContext 从上下文获取用户名
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}
