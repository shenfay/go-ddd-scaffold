package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

// CompressionConfig 压缩配置
type CompressionConfig struct {
	// 最小响应大小才启用压缩（字节），默认1KB
	MinSize int
	// 压缩级别：1-9，默认5
	Level int
	// 排除的路径
	ExcludedPaths []string
}

// DefaultCompressionConfig 默认压缩配置
var DefaultCompressionConfig = CompressionConfig{
	MinSize:       1024, // 1KB
	Level:         5,
	ExcludedPaths: []string{"/health", "/metrics"},
}

// GzipMiddleware gzip压缩中间件
// 对响应体进行gzip压缩，减小传输大小
func GzipMiddleware(config CompressionConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否排除该路径
		if isExcludedPath(c.Request.URL.Path, config.ExcludedPaths) {
			c.Next()
			return
		}

		// 检查客户端是否支持gzip
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next()
			return
		}

		// 创建gzip响应写入器
		gzipWriter, err := gzip.NewWriterLevel(c.Writer, config.Level)
		if err != nil {
			c.Next()
			return
		}
		defer gzipWriter.Close()

		// 设置响应头
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// 包装响应写入器
		c.Writer = &gzipResponseWriter{
			ResponseWriter: c.Writer,
			Writer:         gzipWriter,
		}

		c.Next()
	}
}

// gzipResponseWriter 包装gzip写入器的响应
type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

// isExcludedPath 检查路径是否排除
func isExcludedPath(path string, excludedPaths []string) bool {
	for _, p := range excludedPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
