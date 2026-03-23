package task_queue

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// EventHandlerFunc 事件处理函数类型
type EventHandlerFunc func(ctx context.Context, event kernel.DomainEvent) error

// GenericHandler 通用事件处理器适配器
// 用于将任意事件处理函数适配为 Handler 接口
type GenericHandler struct {
	eventTypes []string
	handler    EventHandlerFunc
}

// NewGenericHandler 创建通用事件处理器
func NewGenericHandler(eventTypes []string, handler EventHandlerFunc) *GenericHandler {
	return &GenericHandler{
		eventTypes: eventTypes,
		handler:    handler,
	}
}

// CanHandle 检查是否可以处理该事件类型
func (h *GenericHandler) CanHandle(eventType string) bool {
	for _, et := range h.eventTypes {
		if et == eventType {
			return true
		}
	}
	return false
}

// Handle 处理事件
func (h *GenericHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
	return h.handler(ctx, event)
}

// AuditLogHandler 审计日志处理器接口（用于适配）
type AuditLogHandler interface {
	Handle(ctx context.Context, event kernel.DomainEvent) error
}

// NewAuditLogHandlerAdapter 创建审计日志处理器适配器
func NewAuditLogHandlerAdapter(handler AuditLogHandler) *GenericHandler {
	return NewGenericHandler(
		[]string{"UserRegistered", "UserLoggedIn"},
		handler.Handle,
	)
}

// LoginLogHandler 登录日志处理器接口（用于适配）
type LoginLogHandler interface {
	Handle(ctx context.Context, event kernel.DomainEvent) error
}

// NewLoginLogHandlerAdapter 创建登录日志处理器适配器
func NewLoginLogHandlerAdapter(handler LoginLogHandler) *GenericHandler {
	return NewGenericHandler(
		[]string{"UserLoggedIn"},
		handler.Handle,
	)
}
