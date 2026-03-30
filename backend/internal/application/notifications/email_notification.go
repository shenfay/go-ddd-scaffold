package notifications

import (
	"context"

	"go.uber.org/zap"
)

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
