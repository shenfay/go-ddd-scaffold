package subscriber

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
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

// StatisticsRepository 统计仓库接口
type StatisticsRepository interface {
	InitializeUserStats(userID int64) error
}

// UserEventSubscriber 用户事件订阅器
// 职责：处理用户领域事件产生的副作用（跨领域操作）
// 位置：应用层（协调基础设施服务）
// 注意：不再负责写入 ActivityLog，ActivityLog 由 UseCase 在事务内直接写入
type UserEventSubscriber struct {
	logger       *zap.Logger
	emailService EmailService
	statsRepo    StatisticsRepository
}

// NewUserEventSubscriber 创建用户事件订阅器
func NewUserEventSubscriber(
	logger *zap.Logger,
	emailService EmailService,
	statsRepo StatisticsRepository,
) *UserEventSubscriber {
	return &UserEventSubscriber{
		logger:       logger,
		emailService: emailService,
		statsRepo:    statsRepo,
	}
}

// CanHandle 检查是否可以处理该事件类型
// 实现 worker.Handler 接口，支持在 Worker 中使用
func (s *UserEventSubscriber) CanHandle(eventType string) bool {
	switch eventType {
	case "UserRegistered", "UserActivated", "UserDeactivated", "UserLoggedIn",
		"UserPasswordChanged", "UserEmailChanged", "UserLocked", "UserUnlocked", "UserProfileUpdated":
		return true
	default:
		return false
	}
}

// Handle 处理领域事件
// 实现 worker.Handler 接口，支持在 Worker 中使用
func (s *UserEventSubscriber) Handle(ctx context.Context, event kernel.DomainEvent) error {
	switch e := event.(type) {
	case *userEvent.UserRegisteredEvent:
		return s.handleUserRegistered(ctx, e)
	case *userEvent.UserActivatedEvent:
		return s.handleUserActivated(ctx, e)
	case *userEvent.UserDeactivatedEvent:
		return s.handleUserDeactivated(ctx, e)
	case *userEvent.UserLoggedInEvent:
		return s.handleUserLoggedIn(ctx, e)
	case *userEvent.UserPasswordChangedEvent:
		return s.handlePasswordChanged(ctx, e)
	case *userEvent.UserEmailChangedEvent:
		return s.handleEmailChanged(ctx, e)
	case *userEvent.UserLockedEvent:
		return s.handleUserLocked(ctx, e)
	case *userEvent.UserUnlockedEvent:
		return s.handleUserUnlocked(ctx, e)
	case *userEvent.UserProfileUpdatedEvent:
		return s.handleUserProfileUpdated(ctx, e)
	default:
		// 忽略未知事件类型
		s.logger.Debug("Unknown event type", zap.String("type", event.EventName()))
		return nil
	}
}

// Subscribe 注册事件处理器到事件总线
func (s *UserEventSubscriber) Subscribe(bus kernel.EventBus) {
	bus.Subscribe("UserRegistered", s.handleUserRegistered)
	bus.Subscribe("UserActivated", s.handleUserActivated)
	bus.Subscribe("UserDeactivated", s.handleUserDeactivated)
	bus.Subscribe("UserLoggedIn", s.handleUserLoggedIn)
	bus.Subscribe("UserPasswordChanged", s.handlePasswordChanged)
	bus.Subscribe("UserEmailChanged", s.handleEmailChanged)
	bus.Subscribe("UserLocked", s.handleUserLocked)
	bus.Subscribe("UserUnlocked", s.handleUserUnlocked)
	bus.Subscribe("UserProfileUpdated", s.handleUserProfileUpdated)
}

// handleUserRegistered 处理用户注册事件
// 副作用：发送欢迎邮件、初始化用户统计
func (s *UserEventSubscriber) handleUserRegistered(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserRegisteredEvent)
	if !ok {
		s.logger.Error("Invalid event type for UserRegistered", zap.String("type", event.EventName()))
		return nil
	}

	s.logger.Info("Handling UserRegistered event",
		zap.String("user_id", e.UserID.String()),
		zap.String("username", e.Username),
		zap.String("email", e.Email),
	)

	// 1. 发送欢迎邮件（异步）
	if s.emailService != nil {
		go func() {
			if err := s.emailService.SendWelcomeEmail(e.Email, e.Username); err != nil {
				s.logger.Error("发送欢迎邮件失败",
					zap.String("email", e.Email),
					zap.Error(err),
				)
			}
		}()
	}

	// 2. 初始化用户统计信息（异步）
	if s.statsRepo != nil {
		go func() {
			if err := s.statsRepo.InitializeUserStats(e.UserID.Int64()); err != nil {
				s.logger.Error("初始化用户统计失败", zap.Error(err))
			}
		}()
	}

	return nil
}

// handleUserActivated 处理用户激活事件
func (s *UserEventSubscriber) handleUserActivated(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserActivatedEvent)
	if !ok {
		return nil
	}

	s.logger.Info("Handling UserActivated event",
		zap.String("user_id", e.UserID.String()),
	)

	// TODO: 发送激活确认邮件
	// if s.emailService != nil {
	//     // 需要查询用户邮箱
	// }

	return nil
}

// handleUserDeactivated 处理用户禁用事件
func (s *UserEventSubscriber) handleUserDeactivated(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserDeactivatedEvent)
	if !ok {
		return nil
	}

	s.logger.Info("Handling UserDeactivated event",
		zap.String("user_id", e.UserID.String()),
		zap.String("reason", e.Reason),
	)

	// TODO: 发送账户禁用通知
	// if s.emailService != nil {
	//     // 需要查询用户邮箱
	// }

	return nil
}

// handleUserLoggedIn 处理用户登录事件
func (s *UserEventSubscriber) handleUserLoggedIn(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserLoggedInEvent)
	if !ok {
		return nil
	}

	s.logger.Info("Handling UserLoggedIn event",
		zap.String("user_id", e.UserID.String()),
		zap.String("ip_address", e.IPAddress),
		zap.Bool("success", e.Success),
	)

	// TODO: 更新用户统计（登录次数、最后登录时间）
	// TODO: 检测异常登录行为

	return nil
}

// handlePasswordChanged 处理用户修改密码事件
func (s *UserEventSubscriber) handlePasswordChanged(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserPasswordChangedEvent)
	if !ok {
		return nil
	}

	s.logger.Info("Handling UserPasswordChanged event",
		zap.String("user_id", e.UserID.String()),
		zap.String("ip_address", e.IPAddress),
	)

	// 发送密码修改通知邮件（异步）
	// 注意：需要从数据库查询用户邮箱
	// if s.emailService != nil {
	//     go func() {
	//         // 查询用户邮箱后发送通知
	//     }()
	// }

	return nil
}

// handleEmailChanged 处理用户修改邮箱事件
func (s *UserEventSubscriber) handleEmailChanged(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserEmailChangedEvent)
	if !ok {
		return nil
	}

	s.logger.Info("Handling UserEmailChanged event",
		zap.String("user_id", e.UserID.String()),
		zap.String("old_email", e.OldEmail),
		zap.String("new_email", e.NewEmail),
	)

	// 向旧邮箱发送变更通知（异步）
	if s.emailService != nil {
		go func() {
			if err := s.emailService.SendEmailChangedEmail(e.OldEmail, "", e.OldEmail, e.NewEmail); err != nil {
				s.logger.Error("发送邮箱变更通知失败", zap.Error(err))
			}
		}()
	}

	return nil
}

// handleUserLocked 处理用户锁定事件
func (s *UserEventSubscriber) handleUserLocked(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserLockedEvent)
	if !ok {
		return nil
	}

	s.logger.Info("Handling UserLocked event",
		zap.String("user_id", e.UserID.String()),
		zap.String("reason", e.Reason),
	)

	// TODO: 发送账户锁定通知（异步）
	// 需要从数据库查询用户邮箱

	return nil
}

// handleUserUnlocked 处理用户解锁事件
func (s *UserEventSubscriber) handleUserUnlocked(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserUnlockedEvent)
	if !ok {
		return nil
	}

	s.logger.Info("Handling UserUnlocked event",
		zap.String("user_id", e.UserID.String()),
	)

	// TODO: 发送账户解锁通知（异步）

	return nil
}

// handleUserProfileUpdated 处理用户资料更新事件
func (s *UserEventSubscriber) handleUserProfileUpdated(ctx context.Context, event kernel.DomainEvent) error {
	e, ok := event.(*userEvent.UserProfileUpdatedEvent)
	if !ok {
		return nil
	}

	s.logger.Info("Handling UserProfileUpdated event",
		zap.String("user_id", e.UserID.String()),
		zap.Strings("updated_fields", e.UpdatedFields),
	)

	// TODO: 如果关键信息变更（如手机号），发送确认通知

	return nil
}
