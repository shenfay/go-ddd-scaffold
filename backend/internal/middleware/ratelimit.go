package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

// visitor 访问者信息
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}
}

// allow 检查是否允许请求
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		v = &visitor{limiter: limiter, lastSeen: time.Now()}
		rl.visitors[ip] = v
	} else {
		v.lastSeen = time.Now()
	}

	// 清理长时间未访问的用户（可选）
	go rl.cleanup()

	return v.limiter.Allow()
}

// cleanup 清理长时间未访问的用户
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "TOO_MANY_REQUESTS",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimit 登录接口专用速率限制（更严格）
func LoginRateLimit() gin.HandlerFunc {
	// 5 次/分钟，突发 10 次
	rl := NewRateLimiter(rate.Every(time.Minute/5), 10)
	return RateLimitMiddleware(rl)
}

// GeneralRateLimit 通用速率限制
func GeneralRateLimit() gin.HandlerFunc {
	// 60 次/分钟，突发 100 次
	rl := NewRateLimiter(rate.Every(time.Minute/60), 100)
	return RateLimitMiddleware(rl)
}
