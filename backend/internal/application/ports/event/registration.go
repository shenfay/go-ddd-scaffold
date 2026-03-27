package event

import (
	"context"
	"fmt"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
)

// RegisterAllEvents 注册所有事件类型和处理器
func RegisterAllEvents(registry *common.EventRegistry) error {
	// 注册用户相关事件
	events := map[common.EventType]common.EventCreator{
		common.UserRegistered:      func() common.DomainEvent { return &userEvent.UserRegisteredEvent{} },
		common.UserActivated:       func() common.DomainEvent { return &userEvent.UserActivatedEvent{} },
		common.UserDeactivated:     func() common.DomainEvent { return &userEvent.UserDeactivatedEvent{} },
		common.UserLoggedIn:        func() common.DomainEvent { return &userEvent.UserLoggedInEvent{} },
		common.UserPasswordChanged: func() common.DomainEvent { return &userEvent.UserPasswordChangedEvent{} },
		common.UserEmailChanged:    func() common.DomainEvent { return &userEvent.UserEmailChangedEvent{} },
		common.UserLocked:          func() common.DomainEvent { return &userEvent.UserLockedEvent{} },
		common.UserUnlocked:        func() common.DomainEvent { return &userEvent.UserUnlockedEvent{} },
		common.UserProfileUpdated:  func() common.DomainEvent { return &userEvent.UserProfileUpdatedEvent{} },
	}

	// 批量注册事件创建器
	for eventType, creator := range events {
		if err := registry.RegisterEvent(eventType, creator); err != nil {
			return fmt.Errorf("failed to register event %s: %w", eventType, err)
		}
	}

	return nil
}

// RegisterDefaultHandlers 注册默认事件处理器
func RegisterDefaultHandlers(registry *common.EventRegistry) error {
	// 示例：注册通用的日志处理器
	logHandler := func(ctx context.Context, event common.DomainEvent) error {
		// 这里可以实现通用的日志记录逻辑
		fmt.Printf("Event processed: %s, AggregateID: %v\n", event.EventName(), event.AggregateID())
		return nil
	}

	// 为所有事件类型注册日志处理器
	eventTypes := []common.EventType{
		common.UserRegistered,
		common.UserActivated,
		common.UserDeactivated,
		common.UserLoggedIn,
		common.UserPasswordChanged,
		common.UserEmailChanged,
		common.UserLocked,
		common.UserUnlocked,
		common.UserProfileUpdated,
	}

	for _, eventType := range eventTypes {
		if err := registry.RegisterHandler(eventType, logHandler); err != nil {
			return fmt.Errorf("failed to register handler for %s: %w", eventType, err)
		}
	}

	return nil
}

// ValidateEvent 验证事件的基本规则
func ValidateEvent(event common.DomainEvent) error {
	if event.EventName() == "" {
		return fmt.Errorf("event name cannot be empty")
	}

	if event.AggregateID() == nil {
		return fmt.Errorf("aggregate ID cannot be nil")
	}

	return nil
}
