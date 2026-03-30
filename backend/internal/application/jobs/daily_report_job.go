package jobs

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// DailyReportJob 每日报表作业示例
type DailyReportJob struct {
	logger *zap.Logger
}

// NewDailyReportJob 创建每日报表作业
func NewDailyReportJob(logger *zap.Logger) *DailyReportJob {
	return &DailyReportJob{
		logger: logger,
	}
}

// Execute 执行每日报表生成任务
func (j *DailyReportJob) Execute(ctx context.Context, payload map[string]interface{}) error {
	j.logger.Info("Starting daily report generation", zap.Any("payload", payload))

	// TODO: 实现具体的报表生成逻辑
	// 1. 查询昨天的数据
	// 2. 生成统计报表
	// 3. 发送通知（可选）

	j.logger.Info("Daily report completed")
	return nil
}

// Queue 返回队列名称
func (j *DailyReportJob) Queue() string {
	return "jobs_default"
}

// MaxRetry 返回最大重试次数
func (j *DailyReportJob) MaxRetry() int {
	return 3
}

// Timeout 返回超时时间
func (j *DailyReportJob) Timeout() time.Duration {
	return 5 * time.Minute
}

// Schedule 返回 Cron 表达式（每天凌晨 2 点执行）
func (j *DailyReportJob) Schedule() string {
	return "0 2 * * *"
}
