package common

import (
	"context"
	"fmt"
	"sync"
)

// EventType 事件类型枚举
type EventType string

const (
	// 用户生命周期事件
	UserRegistered  EventType = "UserRegistered"
	UserActivated   EventType = "UserActivated"
	UserDeactivated EventType = "UserDeactivated"
	UserLocked      EventType = "UserLocked"
	UserUnlocked    EventType = "UserUnlocked"

	// 认证相关事件
	UserLoggedIn        EventType = "UserLoggedIn"
	UserPasswordChanged EventType = "UserPasswordChanged"
	UserLoggedOut       EventType = "UserLoggedOut"

	// 个人信息事件
	UserEmailChanged   EventType = "UserEmailChanged"
	UserProfileUpdated EventType = "UserProfileUpdated"
	UserPhoneChanged   EventType = "UserPhoneChanged"

	// 安全相关事件
	UserTwoFactorEnabled  EventType = "UserTwoFactorEnabled"
	UserTwoFactorDisabled EventType = "UserTwoFactorDisabled"
	UserSessionRevoked    EventType = "UserSessionRevoked"
)

// EventRegistry 事件注册中心
type EventRegistry struct {
	creators   map[EventType]EventCreator
	handlers   map[EventType][]EventHandler
	validators map[EventType]EventValidator
	mu         sync.RWMutex
}

// EventCreator 事件创建器函数类型
type EventCreator func() DomainEvent

// EventHandler 事件处理器函数类型
type EventHandler func(context.Context, DomainEvent) error

// EventValidator 事件验证器函数类型
type EventValidator func(DomainEvent) error

// NewEventRegistry 创建事件注册中心
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		creators:   make(map[EventType]EventCreator),
		handlers:   make(map[EventType][]EventHandler),
		validators: make(map[EventType]EventValidator),
	}
}

// RegisterEvent 注册事件类型及其创建器
func (r *EventRegistry) RegisterEvent(eventType EventType, creator EventCreator) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if creator == nil {
		return fmt.Errorf("event creator cannot be nil for event type: %s", eventType)
	}

	if _, exists := r.creators[eventType]; exists {
		return fmt.Errorf("event type already registered: %s", eventType)
	}

	r.creators[eventType] = creator
	return nil
}

// RegisterHandler 注册事件处理器
func (r *EventRegistry) RegisterHandler(eventType EventType, handler EventHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if handler == nil {
		return fmt.Errorf("event handler cannot be nil for event type: %s", eventType)
	}

	r.handlers[eventType] = append(r.handlers[eventType], handler)
	return nil
}

// RegisterValidator 注册事件验证器
func (r *EventRegistry) RegisterValidator(eventType EventType, validator EventValidator) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if validator == nil {
		return fmt.Errorf("event validator cannot be nil for event type: %s", eventType)
	}

	if _, exists := r.validators[eventType]; exists {
		return fmt.Errorf("validator already registered for event type: %s", eventType)
	}

	r.validators[eventType] = validator
	return nil
}

// CreateEvent 创建事件实例
func (r *EventRegistry) CreateEvent(eventType EventType) (DomainEvent, error) {
	r.mu.RLock()
	creator, exists := r.creators[eventType]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no creator registered for event type: %s", eventType)
	}

	return creator(), nil
}

// HandleEvent 处理事件
func (r *EventRegistry) HandleEvent(ctx context.Context, event DomainEvent) error {
	eventType := EventType(event.EventName())

	r.mu.RLock()
	handlers, exists := r.handlers[eventType]
	r.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		// 没有处理器时不报错，只是记录
		return nil
	}

	// 验证事件（如果有的话）
	if validator, hasValidator := r.validators[eventType]; hasValidator {
		if err := validator(event); err != nil {
			return fmt.Errorf("event validation failed: %w", err)
		}
	}

	// 执行所有处理器
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			// 可以选择继续执行其他处理器或立即返回错误
			return fmt.Errorf("handler failed for event %s: %w", eventType, err)
		}
	}

	return nil
}

// IsValidEventType 检查事件类型是否有效
func (r *EventRegistry) IsValidEventType(eventType EventType) bool {
	r.mu.RLock()
	_, exists := r.creators[eventType]
	r.mu.RUnlock()
	return exists
}

// GetRegisteredEvents 获取所有已注册的事件类型
func (r *EventRegistry) GetRegisteredEvents() []EventType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]EventType, 0, len(r.creators))
	for eventType := range r.creators {
		events = append(events, eventType)
	}

	return events
}

// GetHandlersCount 获取指定事件类型的处理器数量
func (r *EventRegistry) GetHandlersCount(eventType EventType) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.handlers[eventType])
}
