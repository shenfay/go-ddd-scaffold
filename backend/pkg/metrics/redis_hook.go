package metrics

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// MetricsRedisHook Redis 指标收集 Hook
type MetricsRedisHook struct {
	metrics        *Metrics
	metricsEnabled bool
}

// NewMetricsRedisHook 创建 Redis 指标收集 Hook
func NewMetricsRedisHook(m *Metrics, metricsEnabled bool) *MetricsRedisHook {
	return &MetricsRedisHook{
		metrics:        m,
		metricsEnabled: metricsEnabled,
	}
}

// DialHook implements redis.DialHook
func (h *MetricsRedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// ProcessHook implements redis.ProcessHook
func (h *MetricsRedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if !h.metricsEnabled || h.metrics == nil {
			return next(ctx, cmd)
		}

		start := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(start).Seconds()

		// 提取命令名称
		cmdName := cmd.Name()
		if cmdName != "" {
			h.metrics.IncRedisCommand(cmdName)
			h.metrics.ObserveRedisCommandDuration(cmdName, duration)
		}

		return err
	}
}

// ProcessPipelineHook implements redis.ProcessPipelineHook
func (h *MetricsRedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if !h.metricsEnabled || h.metrics == nil {
			return next(ctx, cmds)
		}

		start := time.Now()
		err := next(ctx, cmds)
		duration := time.Since(start).Seconds()

		// 统计 pipeline 中的命令
		cmdCount := make(map[string]int)
		for _, cmd := range cmds {
			if cmdName := cmd.Name(); cmdName != "" {
				cmdCount[cmdName]++
			}
		}

		for cmdName, count := range cmdCount {
			for i := 0; i < count; i++ {
				h.metrics.IncRedisCommand(cmdName)
			}
		}

		// 记录平均耗时
		if len(cmds) > 0 {
			h.metrics.ObserveRedisCommandDuration("pipeline", duration/float64(len(cmds)))
		}

		return err
	}
}
