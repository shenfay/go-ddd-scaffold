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

// EmailNotification 邮件通知示例
type EmailNotification struct {
	logger *zap.Logger
}

// NewEmailNotification 创建邮件通知处理器
func NewEmailNotification(logger *zap.Logger) *EmailNotification {
	return &EmailNotification{
		logger: logger,
	}
}

// HandleWelcomeEmail 发送欢迎邮件
func (n *EmailNotification) HandleWelcomeEmail(ctx context.Context, email string, username string) error {
	n.logger.Debug("Sending welcome email", zap.String("email", email), zap.String("username", username))

	// TODO: 实现具体的邮件发送逻辑
	// 1. 加载邮件模板
	// 2. 渲染邮件内容
	// 3. 调用邮件服务发送

	n.logger.Info("Welcome email sent", zap.String("email", email))
	return nil
}

// Handle 处理通知（通用接口实现）
func (n *EmailNotification) Handle(ctx context.Context, payload map[string]interface{}) error {
	email, _ := payload["email"].(string)
	username, _ := payload["username"].(string)

	return n.HandleWelcomeEmail(ctx, email, username)
}

// Queue 返回队列名称
func (n *EmailNotification) Queue() string {
	return "notifications_default"
}
