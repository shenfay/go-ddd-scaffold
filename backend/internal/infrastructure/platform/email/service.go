package email

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"go.uber.org/zap"
)

// Service 邮件服务接口
type Service interface {
	SendWelcomeEmail(ctx context.Context, to, username string) error
	SendPasswordChangedEmail(ctx context.Context, to, username string) error
	SendEmailChangedEmail(ctx context.Context, to, username, oldEmail, newEmail string) error
	SendAccountLockedEmail(ctx context.Context, to, username, reason string) error
	SendAccountUnlockedEmail(ctx context.Context, to, username string) error
}

// SMTPService SMTP 邮件服务实现
type SMTPService struct {
	config config.EmailConfig
	logger *zap.Logger
}

// NewSMTPService 创建 SMTP 邮件服务
func NewSMTPService(cfg config.EmailConfig, logger *zap.Logger) *SMTPService {
	return &SMTPService{
		config: cfg,
		logger: logger,
	}
}

// SendWelcomeEmail 发送欢迎邮件
func (s *SMTPService) SendWelcomeEmail(ctx context.Context, to, username string) error {
	subject := "欢迎加入我们的平台"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>欢迎加入</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #4CAF50;">欢迎加入我们的平台！</h2>
        <p>亲爱的 %s，</p>
        <p>感谢您注册我们的平台。您的账户已成功创建，现在可以开始使用我们的服务了。</p>
        <p>如果您有任何问题，请随时联系我们的客服团队。</p>
        <p style="margin-top: 30px;">祝您使用愉快！</p>
        <p style="color: #666; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
    </div>
</body>
</html>
`, username)

	return s.sendEmail(to, subject, body)
}

// SendPasswordChangedEmail 发送密码修改通知邮件
func (s *SMTPService) SendPasswordChangedEmail(ctx context.Context, to, username string) error {
	subject := "密码修改通知"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>密码修改通知</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2196F3;">密码修改通知</h2>
        <p>亲爱的 %s，</p>
        <p>您的账户密码已于刚刚修改。如果这是您本人的操作，请忽略此邮件。</p>
        <p style="color: #f44336; font-weight: bold;">如果您没有进行此操作，请立即联系客服或重置密码。</p>
        <p style="margin-top: 30px;">谢谢！</p>
        <p style="color: #666; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
    </div>
</body>
</html>
`, username)

	return s.sendEmail(to, subject, body)
}

// SendEmailChangedEmail 发送邮箱变更通知邮件
func (s *SMTPService) SendEmailChangedEmail(ctx context.Context, to, username, oldEmail, newEmail string) error {
	subject := "邮箱变更通知"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>邮箱变更通知</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2196F3;">邮箱变更通知</h2>
        <p>亲爱的 %s，</p>
        <p>您的账户邮箱已变更：</p>
        <ul>
            <li>原邮箱：%s</li>
            <li>新邮箱：%s</li>
        </ul>
        <p style="color: #f44336; font-weight: bold;">如果这不是您本人的操作，请立即联系客服。</p>
        <p style="margin-top: 30px;">谢谢！</p>
        <p style="color: #666; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
    </div>
</body>
</html>
`, username, oldEmail, newEmail)

	return s.sendEmail(to, subject, body)
}

// SendAccountLockedEmail 发送账户锁定通知邮件
func (s *SMTPService) SendAccountLockedEmail(ctx context.Context, to, username, reason string) error {
	subject := "账户锁定通知"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>账户锁定通知</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #f44336;">账户锁定通知</h2>
        <p>亲爱的 %s，</p>
        <p>您的账户已被锁定。</p>
        <p><strong>原因：</strong>%s</p>
        <p>如需解锁账户，请联系客服团队。</p>
        <p style="margin-top: 30px;">谢谢！</p>
        <p style="color: #666; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
    </div>
</body>
</html>
`, username, reason)

	return s.sendEmail(to, subject, body)
}

// SendAccountUnlockedEmail 发送账户解锁通知邮件
func (s *SMTPService) SendAccountUnlockedEmail(ctx context.Context, to, username string) error {
	subject := "账户解锁通知"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>账户解锁通知</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #4CAF50;">账户解锁通知</h2>
        <p>亲爱的 %s，</p>
        <p>您的账户已成功解锁，现在可以正常登录使用了。</p>
        <p style="margin-top: 30px;">谢谢！</p>
        <p style="color: #666; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
    </div>
</body>
</html>
`, username)

	return s.sendEmail(to, subject, body)
}

// sendEmail 发送邮件
func (s *SMTPService) sendEmail(to, subject, body string) error {
	if s.config.SMTPHost == "" {
		s.logger.Warn("SMTP 未配置，跳过发送邮件",
			zap.String("to", to),
			zap.String("subject", subject),
		)
		return nil
	}

	from := s.config.From
	if from == "" {
		from = s.config.Username
	}

	fromName := s.config.FromName
	if fromName == "" {
		fromName = "系统通知"
	}

	// 构建邮件内容
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg.Bytes())
	if err != nil {
		s.logger.Error("发送邮件失败",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("邮件发送成功",
		zap.String("to", to),
		zap.String("subject", subject),
	)

	return nil
}

// NoOpService 空邮件服务（用于测试或未配置邮件时）
type NoOpService struct {
	logger *zap.Logger
}

// NewNoOpService 创建空邮件服务
func NewNoOpService(logger *zap.Logger) *NoOpService {
	return &NoOpService{logger: logger}
}

func (s *NoOpService) SendWelcomeEmail(ctx context.Context, to, username string) error {
	s.logger.Info("[NoOp] 发送欢迎邮件", zap.String("to", to), zap.String("username", username))
	return nil
}

func (s *NoOpService) SendPasswordChangedEmail(ctx context.Context, to, username string) error {
	s.logger.Info("[NoOp] 发送密码修改通知", zap.String("to", to), zap.String("username", username))
	return nil
}

func (s *NoOpService) SendEmailChangedEmail(ctx context.Context, to, username, oldEmail, newEmail string) error {
	s.logger.Info("[NoOp] 发送邮箱变更通知", zap.String("to", to), zap.String("username", username))
	return nil
}

func (s *NoOpService) SendAccountLockedEmail(ctx context.Context, to, username, reason string) error {
	s.logger.Info("[NoOp] 发送账户锁定通知", zap.String("to", to), zap.String("username", username))
	return nil
}

func (s *NoOpService) SendAccountUnlockedEmail(ctx context.Context, to, username string) error {
	s.logger.Info("[NoOp] 发送账户解锁通知", zap.String("to", to), zap.String("username", username))
	return nil
}
