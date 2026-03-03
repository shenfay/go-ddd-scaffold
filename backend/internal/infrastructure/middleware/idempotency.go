package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IdempotencyConfig 幂等性配置
type IdempotencyConfig struct {
	// 内存存储(开发环境使用)
	inMemoryStore sync.Map
	// 幂等Key前缀
	KeyPrefix string
	// 过期时间(秒)，默认24小时
	ExpireSeconds int
	// 需要幂等处理的HTTP方法
	IdempotentMethods []string
}

// DefaultIdempotencyConfig 默认配置
func DefaultIdempotencyConfig() *IdempotencyConfig {
	return &IdempotencyConfig{
		inMemoryStore:     sync.Map{},
		KeyPrefix:         "idempotency:",
		ExpireSeconds:     24 * 3600, // 24小时
		IdempotentMethods: []string{"POST", "PUT", "DELETE", "PATCH"},
	}
}

// Idempotency 幂等性中间件
// 防止重复提交造成的数据不一致问题
func Idempotency(config *IdempotencyConfig) gin.HandlerFunc {
	// 启动清理过期数据的goroutine
	go cleanupExpiredEntries(config)

	return func(c *gin.Context) {
		// 1. 检查是否需要幂等处理
		if !isIdempotentMethod(c.Request.Method, config.IdempotentMethods) {
			c.Next()
			return
		}

		// 2. 获取幂等Key
		idempotencyKey := getIdempotencyKey(c)
		if idempotencyKey == "" {
			// 没有提供幂等Key，继续正常处理
			c.Next()
			return
		}

		fullKey := config.KeyPrefix + idempotencyKey

		// 3. 检查是否已存在处理记录
		if cachedResp, ok := config.inMemoryStore.Load(fullKey); ok {
			storedResp := cachedResp.(IdempotencyResponse)

			// 检查是否过期
			if time.Since(storedResp.TimestampTime) > time.Duration(config.ExpireSeconds)*time.Second {
				config.inMemoryStore.Delete(fullKey)
			} else {
				// 返回缓存的响应
				c.Data(storedResp.StatusCode, "application/json", []byte(storedResp.Body))
				c.Abort()
				return
			}
		}

		// 4. 首次请求，记录开始时间
		startTime := time.Now()

		// 5. 捕获响应
		responseCapture := &responseCaptureWriter{ResponseWriter: c.Writer}
		c.Writer = responseCapture

		// 6. 继续处理请求
		c.Next()

		// 7. 如果请求成功，存储响应结果
		if responseCapture.statusCode >= 200 && responseCapture.statusCode < 300 {
			resp := IdempotencyResponse{
				StatusCode:    responseCapture.statusCode,
				Body:          string(responseCapture.body),
				Timestamp:     startTime.Format(time.RFC3339),
				TimestampTime: startTime,
			}

			config.inMemoryStore.Store(fullKey, resp)
		}
	}
}

// IdempotencyResponse 存储的幂等响应结构
type IdempotencyResponse struct {
	StatusCode    int       `json:"status_code"`
	Body          string    `json:"body"`
	Timestamp     string    `json:"timestamp"`
	TimestampTime time.Time `json:"-"`
}

// responseCaptureWriter 捕获响应的Writer
type responseCaptureWriter struct {
	gin.ResponseWriter
	statusCode int
	body       []byte
}

func (w *responseCaptureWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return w.ResponseWriter.Write(data)
}

// getIdempotencyKey 从请求中获取幂等Key
func getIdempotencyKey(c *gin.Context) string {
	// 优先从Header获取
	key := c.GetHeader("Idempotency-Key")
	if key != "" {
		return key
	}

	// 从请求参数获取
	key = c.Query("idempotency_key")
	if key != "" {
		return key
	}

	// 自动生成基于请求内容的Key
	return generateIdempotencyKey(c)
}

// generateIdempotencyKey 自动生成幂等Key
// 基于: 方法 + 路径 + 请求体内容的MD5
func generateIdempotencyKey(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path

	// 获取请求体内容(注意: 这会消耗body，需要重新设置)
	bodyBytes, _ := c.GetRawData()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// 构造唯一标识字符串
	uniqueStr := fmt.Sprintf("%s:%s:%s", method, path, string(bodyBytes))

	// MD5哈希确保唯一性
	hash := md5.Sum([]byte(uniqueStr))
	return hex.EncodeToString(hash[:])
}

// isIdempotentMethod 判断HTTP方法是否需要幂等处理
func isIdempotentMethod(method string, idempotentMethods []string) bool {
	for _, m := range idempotentMethods {
		if strings.ToUpper(method) == strings.ToUpper(m) {
			return true
		}
	}
	return false
}

// GenerateIdempotencyKey 手动生成幂等Key的工具函数
// 供客户端在发起请求前调用
func GenerateIdempotencyKey(method, path, requestBody string) string {
	uniqueStr := fmt.Sprintf("%s:%s:%s",
		strings.ToUpper(method),
		strings.ToLower(path),
		requestBody)

	hash := md5.Sum([]byte(uniqueStr))
	return hex.EncodeToString(hash[:])
}

// GenerateUUIDIdempotencyKey 生成UUID格式的幂等Key
func GenerateUUIDIdempotencyKey() string {
	return uuid.New().String()
}

// cleanupExpiredEntries 定期清理过期的幂等记录
func cleanupExpiredEntries(config *IdempotencyConfig) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		config.inMemoryStore.Range(func(key, value interface{}) bool {
			resp := value.(IdempotencyResponse)
			if now.Sub(resp.TimestampTime) > time.Duration(config.ExpireSeconds)*time.Second {
				config.inMemoryStore.Delete(key)
			}
			return true
		})
	}
}
