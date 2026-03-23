package event

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
)

// Subscriber 事件订阅器
// 负责将领域事件处理器注册到事件总线
type Subscriber struct {
	bus kernel.EventBus
}

// NewSubscriber 创建事件订阅器
func NewSubscriber(bus kernel.EventBus) *Subscriber {
	return &Subscriber{bus: bus}
}

// Dependencies 事件订阅器依赖
type Dependencies struct {
	UserSideEffectHandler *userEvent.SideEffectHandler
}

// SubscribeAll 注册所有事件处理器
func (s *Subscriber) SubscribeAll(deps *Dependencies) {
	// ========== 同步事件处理器 ==========
	// 注意：以下事件仅在 EventBus 中同步处理（需要立即响应）
	// 审计日志和登录日志已移至 EventPublisher 异步处理，避免重复执行

	// 用户领域副作用处理器（需要同步执行的业务逻辑）
	if deps.UserSideEffectHandler != nil {
		s.bus.Subscribe("UserRegistered", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
		s.bus.Subscribe("UserActivated", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
		s.bus.Subscribe("UserDeactivated", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
		s.bus.Subscribe("UserLoggedIn", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
		s.bus.Subscribe("UserPasswordChanged", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
		s.bus.Subscribe("UserEmailChanged", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
		s.bus.Subscribe("UserLocked", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
		s.bus.Subscribe("UserUnlocked", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
		s.bus.Subscribe("UserProfileUpdated", func(ctx context.Context, event kernel.DomainEvent) error {
			return deps.UserSideEffectHandler.Handle(ctx, event)
		})
	}
}
