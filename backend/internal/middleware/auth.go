package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// UserClaims JWT 用户声明
type UserClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"` // access 或 refresh
	jwt.RegisteredClaims
}

// TokenValidator Token 验证接口
type TokenValidator interface {
	ValidateAccessToken(tokenString string) (*UserClaims, error)
}

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware(tokenService TokenValidator) gin.HandlerFunc {
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

		// 验证 Token
		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrorCodeInvalidToken,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("token_type", claims.TokenType)

		c.Next()
	}
}

// CurrentUser 从上下文中获取当前用户信息
func CurrentUser(c *gin.Context) (userID string, email string, ok bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return "", "", false
	}

	emailVal, exists := c.Get("user_email")
	if !exists {
		return "", "", false
	}

	userID, ok = userIDVal.(string)
	if !ok {
		return "", "", false
	}

	email, ok = emailVal.(string)
	if !ok {
		return "", "", false
	}

	return userID, email, true
}

// RequireAuth 需要认证的辅助函数（返回 401）
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, _, ok := CurrentUser(c); !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrorCodeUnauthorized,
				"message": "Authentication required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
