// Package metrics Prometheus 监控指标
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics 监控指标集合
type Metrics struct {
	// Redis 相关指标
	redisRequestsTotal   *prometheus.CounterVec
	redisRequestDuration *prometheus.HistogramVec
	redisErrorsTotal     *prometheus.CounterVec
	redisPipelineSize    *prometheus.HistogramVec

	// Token 黑名单指标
	tokenBlacklistChecksTotal   *prometheus.CounterVec
	tokenBlacklistHitsTotal     *prometheus.CounterVec
	tokenBlacklistMissTotal     *prometheus.CounterVec
	tokenBlacklistCheckDuration *prometheus.HistogramVec

	// JWT 相关指标
	jwtIssuedTotal    *prometheus.CounterVec
	jwtValidatedTotal *prometheus.CounterVec
	jwtErrorsTotal    *prometheus.CounterVec
	jwtDuration       *prometheus.HistogramVec

	// 限流熔断指标
	rateLimitTriggeredTotal *prometheus.CounterVec
	circuitBreakerState     *prometheus.GaugeVec
	circuitBreakerTrips     *prometheus.CounterVec
}

// NewMetrics 创建监控指标集合
func NewMetrics(registry prometheus.Registerer) *Metrics {
	m := &Metrics{
		// Redis 指标
		redisRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "redis_requests_total",
				Help: "Redis 请求总数",
			},
			[]string{"operation"},
		),
		redisRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "redis_request_duration_seconds",
				Help:    "Redis 请求延迟（秒）",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
			},
			[]string{"operation"},
		),
		redisErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "redis_errors_total",
				Help: "Redis 错误总数",
			},
			[]string{"operation"},
		),
		redisPipelineSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "redis_pipeline_size",
				Help:    "Redis Pipeline 大小",
				Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"operation"},
		),

		// Token 黑名单指标
		tokenBlacklistChecksTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "token_blacklist_checks_total",
				Help: "Token 黑名单检查总次数",
			},
			[]string{"type"}, // single, batch
		),
		tokenBlacklistHitsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "token_blacklist_hits_total",
				Help: "Token 黑名单命中次数（在黑名单中）",
			},
			[]string{"type"},
		),
		tokenBlacklistMissTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "token_blacklist_miss_total",
				Help: "Token 黑名单未命中次数（不在黑名单中）",
			},
			[]string{"type"},
		),
		tokenBlacklistCheckDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "token_blacklist_check_duration_seconds",
				Help:    "Token 黑名单检查延迟（秒）",
				Buckets: []float64{0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1},
			},
			[]string{"type"},
		),

		// JWT 指标
		jwtIssuedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "jwt_issued_total",
				Help: "JWT 签发总数",
			},
			[]string{"type"}, // access_token, refresh_token
		),
		jwtValidatedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "jwt_validated_total",
				Help: "JWT 验证总次数",
			},
			[]string{"result"}, // success, failure
		),
		jwtErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "jwt_errors_total",
				Help: "JWT 错误总数",
			},
			[]string{"error_type"},
		),
		jwtDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "jwt_operation_duration_seconds",
				Help:    "JWT 操作延迟（秒）",
				Buckets: []float64{0.0001, 0.0005, 0.001, 0.0025, 0.005, 0.01},
			},
			[]string{"operation"},
		),

		// 限流熔断指标
		rateLimitTriggeredTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limit_triggered_total",
				Help: "限流触发总次数",
			},
			[]string{"resource"},
		),
		circuitBreakerState: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "circuit_breaker_state",
				Help: "熔断器状态（0=关闭，1=打开，2=半开）",
			},
			[]string{"resource"},
		),
		circuitBreakerTrips: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "circuit_breaker_trips_total",
				Help: "熔断器跳闸总次数",
			},
			[]string{"resource"},
		),
	}

	// 注册所有指标
	if registry != nil {
		registry.MustRegister(
			m.redisRequestsTotal,
			m.redisRequestDuration,
			m.redisErrorsTotal,
			m.redisPipelineSize,
			m.tokenBlacklistChecksTotal,
			m.tokenBlacklistHitsTotal,
			m.tokenBlacklistMissTotal,
			m.tokenBlacklistCheckDuration,
			m.jwtIssuedTotal,
			m.jwtValidatedTotal,
			m.jwtErrorsTotal,
			m.jwtDuration,
			m.rateLimitTriggeredTotal,
			m.circuitBreakerState,
			m.circuitBreakerTrips,
		)
	}

	return m
}

// RecordRedisRequest 记录 Redis 请求
func (m *Metrics) RecordRedisRequest(operation string, duration time.Duration, err error) {
	m.redisRequestsTotal.WithLabelValues(operation).Inc()
	m.redisRequestDuration.WithLabelValues(operation).Observe(duration.Seconds())
	if err != nil {
		m.redisErrorsTotal.WithLabelValues(operation).Inc()
	}
}

// RecordRedisPipeline 记录 Redis Pipeline 使用
func (m *Metrics) RecordRedisPipeline(operation string, size int) {
	m.redisPipelineSize.WithLabelValues(operation).Observe(float64(size))
}

// RecordTokenBlacklistCheck 记录 Token 黑名单检查
func (m *Metrics) RecordTokenBlacklistCheck(checkType string, inBlacklist bool, duration time.Duration) {
	m.tokenBlacklistChecksTotal.WithLabelValues(checkType).Inc()
	if inBlacklist {
		m.tokenBlacklistHitsTotal.WithLabelValues(checkType).Inc()
	} else {
		m.tokenBlacklistMissTotal.WithLabelValues(checkType).Inc()
	}
	m.tokenBlacklistCheckDuration.WithLabelValues(checkType).Observe(duration.Seconds())
}

// RecordJWTOperation 记录 JWT 操作
func (m *Metrics) RecordJWTOperation(operation string, duration time.Duration) {
	m.jwtDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordJWTIssued 记录 JWT 签发
func (m *Metrics) RecordJWTIssued(tokenType string) {
	m.jwtIssuedTotal.WithLabelValues(tokenType).Inc()
}

// RecordJWTValidated 记录 JWT 验证
func (m *Metrics) RecordJWTValidated(success bool) {
	if !success {
		m.jwtValidatedTotal.WithLabelValues("failure").Inc()
	} else {
		m.jwtValidatedTotal.WithLabelValues("success").Inc()
	}
}

// RecordJWTError 记录 JWT 错误
func (m *Metrics) RecordJWTError(errorType string) {
	m.jwtErrorsTotal.WithLabelValues(errorType).Inc()
}

// RecordRateLimit 记录限流触发
func (m *Metrics) RecordRateLimit(resource string) {
	m.rateLimitTriggeredTotal.WithLabelValues(resource).Inc()
}

// SetCircuitBreakerState 设置熔断器状态
func (m *Metrics) SetCircuitBreakerState(resource string, state int) {
	// state: 0=closed, 1=open, 2=half-open
	m.circuitBreakerState.WithLabelValues(resource).Set(float64(state))
}

// RecordCircuitBreakerTrip 记录熔断器跳闸
func (m *Metrics) RecordCircuitBreakerTrip(resource string) {
	m.circuitBreakerTrips.WithLabelValues(resource).Inc()
}

// Observer 用于记录操作延迟的工具
type Observer struct {
	startTime time.Time
	operation string
	metrics   *Metrics
	recordFn  func(string, time.Duration)
}

// NewObserver 创建观察者
func NewObserver(metrics *Metrics, operation string, recordFn func(string, time.Duration)) *Observer {
	return &Observer{
		startTime: time.Now(),
		operation: operation,
		metrics:   metrics,
		recordFn:  recordFn,
	}
}

// Observe 记录从创建到现在的延迟
func (o *Observer) Observe() {
	duration := time.Since(o.startTime)
	o.recordFn(o.operation, duration)
}

// ObserveErr 带错误的观察
func (o *Observer) ObserveErr(err error) {
	duration := time.Since(o.startTime)
	o.recordFn(o.operation, duration)
}
