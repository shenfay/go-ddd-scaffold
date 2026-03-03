package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	MaxAge       int  // 缓存最大时间（秒），0 表示不缓存
	EnableETag   bool // 是否启用 ETag
	Enable304    bool // 是否启用 304 响应
	PrivateCache bool // 是否为私有缓存
}

// DefaultCacheConfig 默认缓存配置
var DefaultCacheConfig = CacheConfig{
	MaxAge:       300, // 默认 5 分钟
	EnableETag:   true,
	Enable304:    true,
	PrivateCache: false,
}

// CacheMiddleware 创建缓存中间件
// 基于 ETag 和 Last-Modified 实现 HTTP 缓存
func CacheMiddleware(config CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只处理 GET 请求
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// 检查 If-None-Match
		if config.EnableETag {
			ifNoneMatch := c.GetHeader("If-None-Match")
			if ifNoneMatch != "" {
				currentETag := c.GetString("etag")
				if currentETag != "" && ifNoneMatch == currentETag {
					c.Status(http.StatusNotModified)
					c.Abort()
					return
				}
			}
		}

		// 检查 If-Modified-Since
		if config.Enable304 {
			ifModifiedSince := c.GetHeader("If-Modified-Since")
			if ifModifiedSince != "" {
				if t, err := time.Parse(http.TimeFormat, ifModifiedSince); err == nil {
					lastModified := c.GetTime("lastModified")
					if !lastModified.IsZero() && !lastModified.After(t) {
						c.Status(http.StatusNotModified)
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()

		// 处理响应后设置缓存头
		if c.Writer.Status() == http.StatusOK {
			setCacheHeaders(c, config)
		}
	}
}

// setCacheHeaders 设置缓存响应头
func setCacheHeaders(c *gin.Context, config CacheConfig) {
	// 设置 Cache-Control
	cacheControl := "public"
	if config.PrivateCache {
		cacheControl = "private"
	}
	if config.MaxAge > 0 {
		cacheControl += fmt.Sprintf(", max-age=%d", config.MaxAge)
	} else {
		cacheControl += ", no-cache, must-revalidate"
	}
	c.Header("Cache-Control", cacheControl)

	// 设置 ETag
	if config.EnableETag {
		etag := c.GetString("etag")
		if etag != "" {
			c.Header("ETag", etag)
		}
	}

	// 设置 Last-Modified
	lastModified := c.GetTime("lastModified")
	if !lastModified.IsZero() {
		c.Header("Last-Modified", lastModified.Format(http.TimeFormat))
	}

	// 设置Expires
	if config.MaxAge > 0 {
		expires := time.Now().Add(time.Duration(config.MaxAge) * time.Second)
		c.Header("Expires", expires.Format(http.TimeFormat))
	}
}

// GenerateETag 生成 ETag
// 基于内容生成，用于缓存验证
func GenerateETag(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf(`"%s"`, hex.EncodeToString(hash[:]))
}

// GenerateETagFromBytes 从字节数据生成 ETag
func GenerateETagFromBytes(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf(`"%s"`, hex.EncodeToString(hash[:]))
}

// GenerateETagFromInt 从整数值生成 ETag
func GenerateETagFromInt(value int) string {
	return GenerateETag(strconv.Itoa(value))
}

// SetETag 设置 ETag 到上下文
func SetETag(c *gin.Context, etag string) {
	c.Set("etag", etag)
}

// SetLastModified 设置最后修改时间
func SetLastModified(c *gin.Context, t time.Time) {
	c.Set("lastModified", t)
}

// CachedHandler 创建带缓存的处理器
// 自动为响应生成 ETag 并处理缓存
func CachedHandler(config CacheConfig, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 执行处理函数
		handler(c)

		// 如果响应成功，设置缓存头
		if c.Writer.Status() == http.StatusOK {
			// 如果还没有 ETag，基于响应体生成
			if config.EnableETag && c.GetString("etag") == "" {
				// 这里简化处理，实际可以从响应数据生成
				etag := GenerateETag(fmt.Sprintf("%d-%s", time.Now().Unix(), c.Request.URL.Path))
				c.Header("ETag", etag)
			}

			setCacheHeaders(c, config)
		}
	}
}

// NoCacheMiddleware 禁用缓存中间件
func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}

// StaticFileCacheMiddleware 静态文件缓存中间件
// 适用于静态资源的长缓存
func StaticFileCacheMiddleware(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置长期缓存
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%.0f", maxAge.Seconds()))

		// 生成 ETag
		etag := GenerateETag(fmt.Sprintf("%d-%s", time.Now().Unix(), c.Request.URL.Path))
		c.Header("ETag", etag)

		// 检查 ETag
		ifNoneMatch := c.GetHeader("If-None-Match")
		if ifNoneMatch == etag {
			c.Status(http.StatusNotModified)
			c.Abort()
			return
		}

		c.Next()
	}
}

// VaryHeaderMiddleware 设置 Vary 头
// 告诉缓存代理需要考虑哪些请求头进行缓存
func VaryHeaderMiddleware(varyHeaders ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(varyHeaders) > 0 {
			c.Header("Vary", strings.Join(varyHeaders, ", "))
		} else {
			c.Header("Vary", "Accept-Encoding, Accept-Language")
		}
		c.Next()
	}
}
