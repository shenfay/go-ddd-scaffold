package email

import "context"

// EmailService 邮件服务端口
type EmailService interface {
	// Send 发送邮件
	Send(ctx context.Context, to, subject, body string) error

	// SendHTML 发送 HTML 邮件
	SendHTML(ctx context.Context, to, subject, htmlBody string) error

	// SendWithAttachment 发送带附件的邮件
	SendWithAttachment(ctx context.Context, to, subject, body string, attachments []string) error
}
