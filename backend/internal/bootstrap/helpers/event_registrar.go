package helpers

import (
	"github.com/shenfay/go-ddd-scaffold/internal/container"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	eventInterface "github.com/shenfay/go-ddd-scaffold/internal/interfaces/event"
	"github.com/shenfay/go-ddd-scaffold/pkg/useragent"
	"go.uber.org/zap"
)

// EventHandlers 事件处理器集合
type EventHandlers struct {
	Audit    *eventInterface.AuditSubscriber
	LoginLog *eventInterface.LoginLogSubscriber
}

// uaParserAdapter User-Agent 解析器适配器
type uaParserAdapter struct {
	parser *useragent.Parser
}

// newUAParserAdapter 创建 User-Agent 解析器适配器
func newUAParserAdapter() *uaParserAdapter {
	return &uaParserAdapter{
		parser: useragent.NewParser(),
	}
}

// Parse 解析 User-Agent 字符串
func (a *uaParserAdapter) Parse(ua string) eventInterface.DeviceInfo {
	info := a.parser.Parse(ua)
	return eventInterface.DeviceInfo{
		DeviceType: info.DeviceType,
		OS:         info.OS,
		Browser:    info.Browser,
	}
}

// RegisterEventHandlers 注册所有领域事件处理器
func RegisterEventHandlers(
	eventBus kernel.EventBus,
	container container.Container,
	userSideEffectHandler *event.SideEffectHandler,
	logger *zap.Logger,
) (*EventHandlers, error) {
	logger.Info("Registering event handlers...")

	// 创建事件订阅器
	subscriber := eventInterface.NewSubscriber(eventBus)

	// 创建 User-Agent 解析器适配器
	uaParser := newUAParserAdapter()

	// 创建领域事件订阅者
	auditSubscriber := eventInterface.NewAuditSubscriber(
		container.GetAuditLogRepo(),
		container.GetSnowflake(),
	)

	loginLogSubscriber := eventInterface.NewLoginLogSubscriber(
		container.GetLoginLogRepo(),
		container.GetSnowflake(),
		uaParser,
	)

	// 注册所有事件处理器
	subscriber.SubscribeAll(&eventInterface.Dependencies{
		AuditSubscriber:       auditSubscriber,
		LoginLogSubscriber:    loginLogSubscriber,
		UserSideEffectHandler: userSideEffectHandler,
	})

	logger.Info("Event handlers registered successfully")

	return &EventHandlers{
		Audit:    auditSubscriber,
		LoginLog: loginLogSubscriber,
	}, nil
}
