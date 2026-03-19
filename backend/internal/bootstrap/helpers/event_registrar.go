package helpers

import (
	"github.com/shenfay/go-ddd-scaffold/internal/container"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/audit"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/loginlog"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	eventInterface "github.com/shenfay/go-ddd-scaffold/internal/interfaces/event"
	"go.uber.org/zap"
)

// EventHandlers 事件处理器集合
type EventHandlers struct {
	Audit    *audit.EventHandler
	LoginLog *loginlog.EventHandler
}

// RegisterEventHandlers 注册所有领域事件处理器
func RegisterEventHandlers(
	eventBus kernel.EventBus,
	container container.Container,
	userSideEffectHandler *user.SideEffectHandler,
	logger *zap.Logger,
) (*EventHandlers, error) {
	logger.Info("Registering event handlers...")

	// 创建事件订阅器
	subscriber := eventInterface.NewSubscriber(eventBus)

	// 创建 User-Agent 解析器适配器
	uaParser := NewUAParserAdapter()

	// 创建领域事件处理器
	auditHandler := audit.NewEventHandler(
		container.GetAuditLogRepo(),
		container.GetSnowflake(),
	)

	loginLogHandler := loginlog.NewEventHandler(
		container.GetLoginLogRepo(),
		container.GetSnowflake(),
		uaParser,
	)

	// 注册所有事件处理器
	subscriber.SubscribeAll(&eventInterface.Dependencies{
		AuditHandler:          auditHandler,
		LoginLogHandler:       loginLogHandler,
		UserSideEffectHandler: userSideEffectHandler,
	})

	logger.Info("Event handlers registered successfully")

	return &EventHandlers{
		Audit:    auditHandler,
		LoginLog: loginLogHandler,
	}, nil
}
