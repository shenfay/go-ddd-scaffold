// Package event 事件模块错误定义
package event

import "go-ddd-scaffold/internal/pkg/errors"

var (
	// ErrNilEvent 事件为空
	ErrNilEvent = errors.New("EVENT_NIL_EVENT", "事件对象不能为空")

	// ErrSerializationFailed 序列化失败
	ErrSerializationFailed = errors.New("EVENT_SERIALIZATION_FAILED", "事件序列化失败")

	// ErrDeserializationFailed 反序列化失败
	ErrDeserializationFailed = errors.New("EVENT_DESERIALIZATION_FAILED", "事件反序列化失败")

	// ErrDatabaseOperation 数据库操作失败
	ErrDatabaseOperation = errors.New("EVENT_DATABASE_OPERATION_FAILED", "数据库操作失败")

	// ErrEventNotFound 事件未找到
	ErrEventNotFound = errors.New("EVENT_NOT_FOUND", "事件未找到")

	// ErrInvalidEventType 无效的事件类型
	ErrInvalidEventType = errors.New("EVENT_INVALID_TYPE", "无效的事件类型")

	// ErrPublishFailed 事件发布失败
	ErrPublishFailed = errors.New("EVENT_PUBLISH_FAILED", "事件发布失败")

	// ErrHandlerExecution 处理器执行失败
	ErrHandlerExecution = errors.New("EVENT_HANDLER_EXECUTION_FAILED", "事件处理器执行失败")
)
