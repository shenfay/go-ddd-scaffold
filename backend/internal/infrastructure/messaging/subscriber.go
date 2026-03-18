package messaging

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/audit"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/loginlog"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/messaging/handlers"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
	"github.com/shenfay/go-ddd-scaffold/shared/kernel"
)

// Subscriber 事件订阅器
type Subscriber struct {
	bus kernel.EventBus
}

// NewSubscriber 创建事件订阅器
func NewSubscriber(bus kernel.EventBus) *Subscriber {
	return &Subscriber{bus: bus}
}

// SubscribeHandlers 注册所有事件处理器
func (s *Subscriber) SubscribeHandlers(
	auditRepo audit.AuditLogRepository,
	loginLogRepo loginlog.LoginLogRepository,
	snowflakeNode *snowflake.Node,
) {
	// 审计日志处理器
	auditHandler := handlers.NewAuditLogHandler(auditRepo, snowflakeNode)
	s.bus.Subscribe("UserRegistered", func(ctx context.Context, event kernel.DomainEvent) error {
		return auditHandler.Handle(ctx, event)
	})
	s.bus.Subscribe("UserLoggedIn", func(ctx context.Context, event kernel.DomainEvent) error {
		return auditHandler.Handle(ctx, event)
	})

	// 登录日志处理器
	loginHandler := handlers.NewLoginLogHandler(loginLogRepo, snowflakeNode)
	s.bus.Subscribe("UserLoggedIn", func(ctx context.Context, event kernel.DomainEvent) error {
		return loginHandler.Handle(ctx, event)
	})
}
