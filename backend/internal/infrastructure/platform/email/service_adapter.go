package email

import (
	"context"

	ports_email "github.com/shenfay/go-ddd-scaffold/internal/application/ports/email"
)

// EmailServiceAdapter EmailService 端口适配器
type EmailServiceAdapter struct {
	service Service
}

// NewEmailServiceAdapter 创建邮件服务适配器
func NewEmailServiceAdapter(service Service) *EmailServiceAdapter {
	return &EmailServiceAdapter{
		service: service,
	}
}

// Send 发送邮件
func (a *EmailServiceAdapter) Send(ctx context.Context, to, subject, body string) error {
	// TODO: 实现通用邮件发送
	return nil
}

// SendHTML 发送 HTML 邮件
func (a *EmailServiceAdapter) SendHTML(ctx context.Context, to, subject, htmlBody string) error {
	// TODO: 实现 HTML 邮件发送
	return nil
}

// SendWithAttachment 发送带附件的邮件
func (a *EmailServiceAdapter) SendWithAttachment(ctx context.Context, to, subject, body string, attachments []string) error {
	// TODO: 实现带附件的邮件发送
	return nil
}

var _ ports_email.EmailService = (*EmailServiceAdapter)(nil)
