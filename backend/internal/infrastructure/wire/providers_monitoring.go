// Package wire 监控和限流熔断依赖注入
package wire

import (
	"go-ddd-scaffold/internal/pkg/metrics"
	"go-ddd-scaffold/internal/pkg/ratelimit"

	"github.com/prometheus/client_golang/prometheus"
)

// InitializeMetrics 初始化 Prometheus 监控指标
func InitializeMetrics() *metrics.Metrics {
	registry := prometheus.DefaultRegisterer
	return metrics.NewMetrics(registry)
}

// InitializeRateLimiter 初始化限流器（Token 黑名单专用）
func InitializeRateLimiter(metrics *metrics.Metrics) *ratelimit.RateLimiter {
	// 配置：每秒 100 次请求，突发容量 200
	return ratelimit.NewRateLimiter(100, 200, "token_blacklist", metrics)
}

// InitializeCircuitBreaker 初始化熔断器（Redis 专用）
func InitializeCircuitBreaker(metrics *metrics.Metrics) *ratelimit.CircuitBreaker {
	config := ratelimit.DefaultCircuitBreakerConfig()
	config.MaxFailures = 5              // 5 次失败触发熔断
	config.ResetTimeout = 30            // 30 秒后尝试恢复
	config.HalfOpenMaxCall = 3          // 半开状态允许 3 次调用
	
	cb := ratelimit.NewCircuitBreaker("redis", config, metrics)
	
	// 设置状态变化回调
	cb.OnStateChange(func(state ratelimit.CircuitBreakerState) {
		// 可以在这里添加日志或告警
		switch state {
		case ratelimit.StateClosed:
			// 恢复正常
		case ratelimit.StateOpen:
			// 熔断打开，记录日志
		case ratelimit.StateHalfOpen:
			// 正在恢复
		}
	})
	
	return cb
}
