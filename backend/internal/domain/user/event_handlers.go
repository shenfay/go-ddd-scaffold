package user

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"go.uber.org/zap"
)

// SideEffectHandler 用户领域副作用处理器
// 处理用户领域事件产生的跨领域副作用（如发送邮件、初始化统计等）
type SideEffectHandler struct {
	logger *zap.Logger
	// TODO: 可以注入其他服务来处理副作用
	// emailService EmailService
	// statsRepo    StatisticsRepository
	// auditLogger  AuditLogger
}

// NewSideEffectHandler 创建用户领域副作用处理器
func NewSideEffectHandler(logger *zap.Logger) *SideEffectHandler {
	return &SideEffectHandler{
		logger: logger,
	}
}

// Handle 处理领域事件
func (h *SideEffectHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
	switch e := event.(type) {
	case *UserRegisteredEvent:
		return h.handleUserRegistered(ctx, e)
	case *UserActivatedEvent:
		return h.handleUserActivated(ctx, e)
	case *UserDeactivatedEvent:
		return h.handleUserDeactivated(ctx, e)
	case *UserLoggedInEvent:
		return h.handleUserLoggedIn(ctx, e)
	case *UserPasswordChangedEvent:
		return h.handleUserPasswordChanged(ctx, e)
	case *UserEmailChangedEvent:
		return h.handleUserEmailChanged(ctx, e)
	case *UserLockedEvent:
		return h.handleUserLocked(ctx, e)
	case *UserUnlockedEvent:
		return h.handleUserUnlocked(ctx, e)
	case *UserProfileUpdatedEvent:
		return h.handleUserProfileUpdated(ctx, e)
	default:
		// 忽略未知事件类型
		h.logger.Debug("Unknown event type", zap.String("type", event.EventName()))
		return nil
	}
}

// handleUserRegistered 处理用户注册事件
func (h *SideEffectHandler) handleUserRegistered(ctx context.Context, event *UserRegisteredEvent) error {
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

	return nil
}

// handleUserActivated 处理用户激活事件
func (h *SideEffectHandler) handleUserActivated(ctx context.Context, event *UserActivatedEvent) error {
	h.logger.Info("Handling UserActivated event",
		zap.String("user_id", event.UserID.String()),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送激活确认邮件

	return nil
}

// handleUserDeactivated 处理用户禁用事件
func (h *SideEffectHandler) handleUserDeactivated(ctx context.Context, event *UserDeactivatedEvent) error {
	h.logger.Info("Handling UserDeactivated event",
		zap.String("user_id", event.UserID.String()),
		zap.String("reason", event.Reason),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送账户禁用通知

	return nil
}

// handleUserLoggedIn 处理用户登录事件
func (h *SideEffectHandler) handleUserLoggedIn(ctx context.Context, event *UserLoggedInEvent) error {
	h.logger.Info("Handling UserLoggedIn event",
		zap.String("user_id", event.UserID.String()),
		zap.String("ip_address", event.IPAddress),
	)

	// TODO: 实现副作用逻辑
	// 1. 更新用户统计（登录次数、最后登录时间）
	// 2. 检测异常登录行为

	return nil
}

// handleUserPasswordChanged 处理用户修改密码事件
func (h *SideEffectHandler) handleUserPasswordChanged(ctx context.Context, event *UserPasswordChangedEvent) error {
	h.logger.Info("Handling UserPasswordChanged event",
		zap.String("user_id", event.UserID.String()),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送密码修改通知邮件
	// 2. 如果是异常操作，触发安全告警

	return nil
}

// handleUserEmailChanged 处理用户修改邮箱事件
func (h *SideEffectHandler) handleUserEmailChanged(ctx context.Context, event *UserEmailChangedEvent) error {
	h.logger.Info("Handling UserEmailChanged event",
		zap.String("user_id", event.UserID.String()),
		zap.String("new_email", event.NewEmail),
	)

	// TODO: 实现副作用逻辑
	// 1. 向旧邮箱发送变更通知
	// 2. 向新邮箱发送验证邮件

	return nil
}

// handleUserLocked 处理用户锁定事件
func (h *SideEffectHandler) handleUserLocked(ctx context.Context, event *UserLockedEvent) error {
	h.logger.Info("Handling UserLocked event",
		zap.String("user_id", event.UserID.String()),
		zap.String("reason", event.Reason),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送账户锁定通知

	return nil
}

// handleUserUnlocked 处理用户解锁事件
func (h *SideEffectHandler) handleUserUnlocked(ctx context.Context, event *UserUnlockedEvent) error {
	h.logger.Info("Handling UserUnlocked event",
		zap.String("user_id", event.UserID.String()),
	)

	// TODO: 实现副作用逻辑
	// 1. 发送账户解锁通知

	return nil
}

// handleUserProfileUpdated 处理用户资料更新事件
func (h *SideEffectHandler) handleUserProfileUpdated(ctx context.Context, event *UserProfileUpdatedEvent) error {
	h.logger.Info("Handling UserProfileUpdated event",
		zap.String("user_id", event.UserID.String()),
	)

	// TODO: 实现副作用逻辑
	// 1. 如果关键信息变更（如手机号），发送确认通知

	return nil
}
