package user

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
	"go.uber.org/zap"
)

// UserEventHandler 用户领域事件处理器
type UserEventHandler struct {
	logger *zap.Logger
	// TODO: 可以注入其他服务来处理副作用
	// emailService EmailService
	// statsRepo    StatisticsRepository
	// auditLogger  AuditLogger
}

// NewUserEventHandler 创建用户领域事件处理器
func NewUserEventHandler(logger *zap.Logger) *UserEventHandler {
	return &UserEventHandler{
		logger: logger,
	}
}

// Handle 处理领域事件
func (h *UserEventHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	switch e := event.(type) {
	case *user.UserRegisteredEvent:
		return h.handleUserRegistered(ctx, e)
	case *user.UserActivatedEvent:
		return h.handleUserActivated(ctx, e)
	case *user.UserDeactivatedEvent:
		return h.handleUserDeactivated(ctx, e)
	case *user.UserLoggedInEvent:
		return h.handleUserLoggedIn(ctx, e)
	case *user.UserPasswordChangedEvent:
		return h.handleUserPasswordChanged(ctx, e)
	case *user.UserEmailChangedEvent:
		return h.handleUserEmailChanged(ctx, e)
	case *user.UserLockedEvent:
		return h.handleUserLocked(ctx, e)
	case *user.UserUnlockedEvent:
		return h.handleUserUnlocked(ctx, e)
	case *user.UserProfileUpdatedEvent:
		return h.handleUserProfileUpdated(ctx, e)
	default:
		// 忽略未知事件类型
		h.logger.Debug("Unknown event type", zap.String("type", event.EventName()))
		return nil
	}
}

// handleUserRegistered 处理用户注册事件
func (h *UserEventHandler) handleUserRegistered(ctx context.Context, event *user.UserRegisteredEvent) error {
	h.logger.Info("Handling UserRegistered event",
		zap.String("user_id", event.UserID.String()),
		zap.String("username", event.Username),
		zap.String("email", event.Email),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送欢迎邮件
	// if err := h.emailService.SendWelcomeEmail(event.Email, event.Username); err != nil {
	//     return err
	// }

	// 2. 初始化用户统计信息
	// if err := h.statsRepo.InitializeUserStats(event.UserID.Int64()); err != nil {
	//     return err
	// }

	// 3. 记录审计日志
	// h.auditLogger.Info("User registered", zap.String("user_id", event.UserID.String()))

	return nil
}

// handleUserActivated 处理用户激活事件
func (h *UserEventHandler) handleUserActivated(ctx context.Context, event *user.UserActivatedEvent) error {
	h.logger.Info("Handling UserActivated event",
		zap.String("user_id", event.UserID.String()),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送激活确认邮件
	// 2. 记录审计日志

	return nil
}

// handleUserDeactivated 处理用户禁用事件
func (h *UserEventHandler) handleUserDeactivated(ctx context.Context, event *user.UserDeactivatedEvent) error {
	h.logger.Info("Handling UserDeactivated event",
		zap.String("user_id", event.UserID.String()),
		zap.String("reason", event.Reason),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送账户禁用通知
	// 2. 记录审计日志

	return nil
}

// handleUserLoggedIn 处理用户登录事件
func (h *UserEventHandler) handleUserLoggedIn(ctx context.Context, event *user.UserLoggedInEvent) error {
	h.logger.Info("Handling UserLoggedIn event",
		zap.String("user_id", event.UserID.String()),
		zap.String("ip_address", event.IPAddress),
	)

	// TODO: 实现副作用逻辑
	// 1. 记录登录日志到 login_logs 表
	// 2. 更新用户统计（登录次数、最后登录时间）
	// 3. 检测异常登录行为

	return nil
}

// handleUserPasswordChanged 处理用户修改密码事件
func (h *UserEventHandler) handleUserPasswordChanged(ctx context.Context, event *user.UserPasswordChangedEvent) error {
	h.logger.Info("Handling UserPasswordChanged event",
		zap.String("user_id", event.UserID.String()),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送密码修改通知邮件
	// 2. 记录审计日志
	// 3. 如果是异常操作，触发安全告警

	return nil
}

// handleUserEmailChanged 处理用户修改邮箱事件
func (h *UserEventHandler) handleUserEmailChanged(ctx context.Context, event *user.UserEmailChangedEvent) error {
	h.logger.Info("Handling UserEmailChanged event",
		zap.String("user_id", event.UserID.String()),
		zap.String("new_email", event.NewEmail),
	)

	// TODO: 实现副作用逻辑
	// 1. 向旧邮箱发送变更通知
	// 2. 向新邮箱发送验证邮件
	// 3. 记录审计日志

	return nil
}

// handleUserLocked 处理用户锁定事件
func (h *UserEventHandler) handleUserLocked(ctx context.Context, event *user.UserLockedEvent) error {
	h.logger.Info("Handling UserLocked event",
		zap.String("user_id", event.UserID.String()),
		zap.String("reason", event.Reason),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送账户锁定通知
	// 2. 记录安全日志

	return nil
}

// handleUserUnlocked 处理用户解锁事件
func (h *UserEventHandler) handleUserUnlocked(ctx context.Context, event *user.UserUnlockedEvent) error {
	h.logger.Info("Handling UserUnlocked event",
		zap.String("user_id", event.UserID.String()),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送账户解锁通知
	// 2. 记录审计日志

	return nil
}

// handleUserProfileUpdated 处理用户资料更新事件
func (h *UserEventHandler) handleUserProfileUpdated(ctx context.Context, event *user.UserProfileUpdatedEvent) error {
	h.logger.Info("Handling UserProfileUpdated event",
		zap.String("user_id", event.UserID.String()),
	)

	// TODO: 实现副作用逻辑
	// 1. 如果关键信息变更（如手机号），发送确认通知
	// 2. 记录审计日志

	return nil
}
