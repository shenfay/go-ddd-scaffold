package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/pkg/response"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtService entity.JWTService
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtService entity.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// HandlerFunc 返回 Gin 中间件处理函数
// 验证 JWT Token 并将用户信息注入 Context
func (m *AuthMiddleware) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 提取 Token（支持 Bearer 格式）
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Unauthorized(c.Request.Context(), "missing authorization token"))
			return
		}

		// 2. 验证 JWT
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Unauthorized(c.Request.Context(), "invalid or expired token"))
			return
		}

		// 3. 注入用户信息到 Context
		c.Set("userID", claims.UserID)
		c.Set("tenantID", claims.TenantID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequirePermission 创建权限检查中间件
// resource: 资源名称，如 "children", "progress", "self"
// action: 操作名称，如 "read", "write", "delete"
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Context 获取用户信息
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Unauthorized(c.Request.Context(), "user not authenticated"))
			return
		}

		tenantID, exists := c.Get("tenantID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.BusinessErr(c.Request.Context(), "InvalidParameter", "tenant not found"))
			return
		}

		// 转换为字符串
		userIDStr := userID.(uuid.UUID).String()
		tenantIDStr := tenantID.(uuid.UUID).String()

		// 超级管理员拥有所有权限
		if role, exists := c.Get("role"); exists {
			if role == entity.RoleSuperAdmin {
				c.Next()
				return
			}
		}

		// 通过 Casbin 检查权限（需要从 Context 获取 CasbinService）
		// 注意：这里需要通过 gin.Context 传递 CasbinService
		casbinService, exists := c.Get("casbinService")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ServerErr(c.Request.Context()))
			return
		}

		allowed, err := casbinService.(interface {
			Enforce(sub, dom, obj, act string) (bool, error)
		}).Enforce(userIDStr, tenantIDStr, resource, action)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ServerErr(c.Request.Context()))
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Forbidden(c.Request.Context(), "permission denied"))
			return
		}

		c.Next()
	}
}

// RequireRole 创建角色检查中间件
// 至少拥有指定角色之一才能访问
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Context 获取角色
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Unauthorized(c.Request.Context(), "user not authenticated"))
			return
		}

		userRole := string(role.(entity.UserRole))

		// 检查是否拥有所需角色
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, response.Forbidden(c.Request.Context(), "insufficient role privileges"))
	}
}

// GetUserID 从 Context 获取用户ID
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("user not authenticated")
	}
	return userID.(uuid.UUID), nil
}

// GetTenantID 从 Context 获取租户ID
func GetTenantID(c *gin.Context) (uuid.UUID, error) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		return uuid.Nil, errors.New("tenant not found")
	}
	return tenantID.(uuid.UUID), nil
}

// GetRole 从 Context 获取角色
func GetRole(c *gin.Context) (entity.UserRole, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role not found")
	}
	return role.(entity.UserRole), nil
}

// extractToken 从 Header 提取 Token
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// 支持 Bearer 格式
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return authHeader
}
