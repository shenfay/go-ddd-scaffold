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

// Send 发送邮件 (通用邮件)
func (a *EmailServiceAdapter) Send(ctx context.Context, to, subject, body string) error {
	return a.service.SendGenericEmail(ctx, to, subject, body)
}

// SendHTML 发送 HTML 邮件
func (a *EmailServiceAdapter) SendHTML(ctx context.Context, to, subject, htmlBody string) error {
	// 使用与 Send 相同的实现，因为默认就是发送 HTML 邮件
	return a.Send(ctx, to, subject, htmlBody)
}

// SendWithAttachment 发送带附件的邮件
func (a *EmailServiceAdapter) SendWithAttachment(ctx context.Context, to, subject, body string, attachments []string) error {
	return a.service.SendGenericEmailWithAttachment(ctx, to, subject, body, attachments)
}

var _ ports_email.EmailService = (*EmailServiceAdapter)(nil)
