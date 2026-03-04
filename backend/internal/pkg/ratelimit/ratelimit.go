// Package ratelimit 限流和熔断器实现
package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"

	"go-ddd-scaffold/internal/pkg/metrics"
)

var (
	// ErrRateLimited 限流错误
	ErrRateLimited = errors.New("请求过于频繁，已被限流")
	// ErrCircuitBreakerOpen 熔断器打开错误
	ErrCircuitBreakerOpen = errors.New("熔断器已打开，服务暂时不可用")
)

// CircuitBreakerState 熔断器状态
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota // 关闭（正常）
	StateOpen                               // 打开（熔断）
	StateHalfOpen                           // 半开（尝试恢复）
)

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	MaxFailures     int           // 最大失败次数触发熔断
	ResetTimeout    time.Duration // 熔断后自动恢复时间
	HalfOpenMaxCall int           // 半开状态允许的最大调用数
}

// DefaultCircuitBreakerConfig 默认熔断器配置
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxFailures:     5,               // 5 次失败触发熔断
		ResetTimeout:    30 * time.Second, // 30 秒后尝试恢复
		HalfOpenMaxCall: 3,               // 半开状态允许 3 次调用
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	name           string
	config         CircuitBreakerConfig
	state          CircuitBreakerState
	failures       int
	lastFailureAt  time.Time
	halfOpenCalls  int
	mutex          sync.RWMutex
	metrics        *metrics.Metrics
	onStateChange  func(CircuitBreakerState) // 状态变化回调
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(name string, config CircuitBreakerConfig, metrics *metrics.Metrics) *CircuitBreaker {
	cb := &CircuitBreaker{
		name:    name,
		config:  config,
		state:   StateClosed,
		metrics: metrics,
	}
	
	if metrics != nil {
		metrics.SetCircuitBreakerState(name, int(cb.state))
	}
	
	return cb
}

// Execute 执行操作（带熔断保护）
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	cb.mutex.Lock()
	
	// 检查是否需要改变状态
	now := time.Now()
	switch cb.state {
	case StateOpen:
		// 熔断状态：检查是否到了恢复时间
		if now.Sub(cb.lastFailureAt) > cb.config.ResetTimeout {
			cb.state = StateHalfOpen
			cb.halfOpenCalls = 0
			cb.mutex.Unlock()
			
			if cb.metrics != nil {
				cb.metrics.SetCircuitBreakerState(cb.name, int(cb.state))
			}
		} else {
			cb.mutex.Unlock()
			
			if cb.metrics != nil {
				cb.metrics.RecordCircuitBreakerTrip(cb.name)
			}
			
			return ErrCircuitBreakerOpen
		}
		
	case StateHalfOpen:
		// 半开状态：限制调用次数
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCall {
			cb.mutex.Unlock()
			return ErrCircuitBreakerOpen
		}
		cb.halfOpenCalls++
		cb.mutex.Unlock()
		
	default:
		cb.mutex.Unlock()
	}
	
	// 执行实际调用
	err := fn()
	
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	if err != nil {
		cb.failures++
		cb.lastFailureAt = now
		
		// 达到最大失败次数，触发熔断
		if cb.failures >= cb.config.MaxFailures {
			cb.state = StateOpen
			
			if cb.metrics != nil {
				cb.metrics.SetCircuitBreakerState(cb.name, int(cb.state))
				cb.metrics.RecordCircuitBreakerTrip(cb.name)
			}
			
			if cb.onStateChange != nil {
				go cb.onStateChange(StateOpen)
			}
		}
	} else {
		// 调用成功
		if cb.state == StateHalfOpen {
			// 半开状态成功，恢复正常
			cb.state = StateClosed
			cb.failures = 0
			cb.halfOpenCalls = 0
			
			if cb.metrics != nil {
				cb.metrics.SetCircuitBreakerState(cb.name, int(cb.state))
			}
			
			if cb.onStateChange != nil {
				go cb.onStateChange(StateClosed)
			}
		} else if cb.state == StateClosed {
			// 连续成功，减少失败计数
			if cb.failures > 0 {
				cb.failures--
			}
		}
	}
	
	return err
}

// State 获取当前状态
func (cb *CircuitBreaker) State() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Reset 手动重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.state = StateClosed
	cb.failures = 0
	cb.halfOpenCalls = 0
	
	if cb.metrics != nil {
		cb.metrics.SetCircuitBreakerState(cb.name, int(cb.state))
	}
	
	if cb.onStateChange != nil {
		go cb.onStateChange(StateClosed)
	}
}

// OnStateChange 设置状态变化回调
func (cb *CircuitBreaker) OnStateChange(fn func(CircuitBreakerState)) {
	cb.onStateChange = fn
}

// RateLimiter 简单的令牌桶限流器
type RateLimiter struct {
	rate       int           // 每秒允许的请求数
	burst      int           // 突发容量
	tokens     float64       // 当前令牌数
	lastRefill time.Time     // 上次补充令牌时间
	mutex      sync.Mutex
	metrics    *metrics.Metrics
	resource   string
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rate, burst int, resource string, metrics *metrics.Metrics) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		resource:   resource,
		metrics:    metrics,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += elapsed * float64(rl.rate)
	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}
	rl.lastRefill = now
	
	// 消耗令牌
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	
	// 限流触发
	if rl.metrics != nil {
		rl.metrics.RecordRateLimit(rl.resource)
	}
	
	return false
}

// Wait 等待直到允许请求（带超时）
func (rl *RateLimiter) Wait(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if rl.Allow() {
				return nil
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ProtectedExecutor 受保护的执行器（限流 + 熔断）
type ProtectedExecutor struct {
	rateLimiter    *RateLimiter
	circuitBreaker *CircuitBreaker
}

// NewProtectedExecutor 创建受保护的执行器
func NewProtectedExecutor(
	rateLimiter *RateLimiter,
	circuitBreaker *CircuitBreaker,
) *ProtectedExecutor {
	return &ProtectedExecutor{
		rateLimiter:    rateLimiter,
		circuitBreaker: circuitBreaker,
	}
}

// Execute 执行操作（限流 + 熔断双重保护）
func (pe *ProtectedExecutor) Execute(ctx context.Context, fn func() error) error {
	// 1. 限流检查
	if !pe.rateLimiter.Allow() {
		return ErrRateLimited
	}
	
	// 2. 熔断器保护
	return pe.circuitBreaker.Execute(ctx, fn)
}
