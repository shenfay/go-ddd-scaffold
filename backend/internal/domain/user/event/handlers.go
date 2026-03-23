package event

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"go.uber.org/zap"
)

// EmailService 邮件服务接口
type EmailService interface {
	SendWelcomeEmail(to, username string) error
	SendPasswordChangedEmail(to, username string) error
	SendEmailChangedEmail(to, username, oldEmail, newEmail string) error
	SendAccountLockedEmail(to, username, reason string) error
	SendAccountUnlockedEmail(to, username string) error
}

// StatisticsRepository 统计仓库接口（预留）
type StatisticsRepository interface {
	InitializeUserStats(userID int64) error
}

// AuditLogger 审计日志接口（预留）
type AuditLogger interface {
	Log(ctx context.Context, action string, userID int64, details map[string]interface{}) error
}

// SideEffectHandler 用户领域副作用处理器
// 处理用户领域事件产生的跨领域副作用（如发送邮件、初始化统计等）
type SideEffectHandler struct {
	logger *zap.Logger
	// TODO: 可以注入其他服务来处理副作用
	emailService EmailService
	statsRepo    StatisticsRepository
	auditLogger  AuditLogger
}

// NewSideEffectHandler 创建用户领域副作用处理器
func NewSideEffectHandler(logger *zap.Logger, emailService EmailService) *SideEffectHandler {
	return &SideEffectHandler{
		logger:       logger,
		emailService: emailService,
	}
}

// CanHandle 检查是否可以处理该事件类型
func (h *SideEffectHandler) CanHandle(eventType string) bool {
	switch eventType {
	case "UserRegistered", "UserActivated", "UserDeactivated", "UserLoggedIn",
		"UserPasswordChanged", "UserEmailChanged", "UserLocked", "UserUnlocked", "UserProfileUpdated":
		return true
	default:
		return false
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

	// 1. 发送欢迎邮件
	if h.emailService != nil {
		if err := h.emailService.SendWelcomeEmail(event.Email, event.Username); err != nil {
			h.logger.Error("发送欢迎邮件失败",
				zap.String("email", event.Email),
				zap.Error(err),
			)
			// 邮件发送失败不阻断主流程
		}
	}

	// 2. 初始化用户统计信息（预留）
	// if h.statsRepo != nil {
	//     if err := h.statsRepo.InitializeUserStats(event.UserID.Int64()); err != nil {
	//         h.logger.Error("初始化用户统计失败", zap.Error(err))
	//     }
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

	// 发送密码修改通知邮件
	// 注意：这里需要获取用户邮箱，但事件中只有 UserID
	// 实际实现中可能需要从数据库查询用户邮箱
	// if h.emailService != nil {
	//     // 需要查询用户邮箱
	//     // h.emailService.SendPasswordChangedEmail(userEmail, username)
	// }

	return nil
}

// handleUserEmailChanged 处理用户修改邮箱事件
func (h *SideEffectHandler) handleUserEmailChanged(ctx context.Context, event *UserEmailChangedEvent) error {
	h.logger.Info("Handling UserEmailChanged event",
		zap.String("user_id", event.UserID.String()),
		zap.String("new_email", event.NewEmail),
	)

	// 向旧邮箱发送变更通知
	if h.emailService != nil {
		if err := h.emailService.SendEmailChangedEmail(event.OldEmail, "", event.OldEmail, event.NewEmail); err != nil {
			h.logger.Error("发送邮箱变更通知失败", zap.Error(err))
		}
	}

	return nil
}

// handleUserLocked 处理用户锁定事件
func (h *SideEffectHandler) handleUserLocked(ctx context.Context, event *UserLockedEvent) error {
	h.logger.Info("Handling UserLocked event",
		zap.String("user_id", event.UserID.String()),
		zap.String("reason", event.Reason),
	)

	// 发送账户锁定通知
	// 注意：需要获取用户邮箱
	// if h.emailService != nil {
	//     // 查询用户邮箱后发送通知
	// }

	return nil
}

// handleUserUnlocked 处理用户解锁事件
func (h *SideEffectHandler) handleUserUnlocked(ctx context.Context, event *UserUnlockedEvent) error {
	h.logger.Info("Handling UserUnlocked event",
		zap.String("user_id", event.UserID.String()),
	)

	// 发送账户解锁通知
	// 注意：需要获取用户邮箱
	// if h.emailService != nil {
	//     // 查询用户邮箱后发送通知
	// }

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
