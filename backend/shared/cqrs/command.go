package cqrs

import (
	"context"
	"reflect"
)

// Command 命令接口
type Command interface {
	CommandName() string
	Validate() error
}

// CommandHandler 命令处理器接口
type CommandHandler interface {
	Handle(ctx context.Context, command Command) (interface{}, error)
}

// CommandBus 命令总线接口
type CommandBus interface {
	Dispatch(ctx context.Context, command Command) (interface{}, error)
	RegisterHandler(commandType reflect.Type, handler CommandHandler) error
}

// SimpleCommandBus 简单命令总线实现
type SimpleCommandBus struct {
	handlers map[reflect.Type]CommandHandler
}

// NewSimpleCommandBus 创建简单命令总线
func NewSimpleCommandBus() *SimpleCommandBus {
	return &SimpleCommandBus{
		handlers: make(map[reflect.Type]CommandHandler),
	}
}

// Dispatch 分发命令
func (cb *SimpleCommandBus) Dispatch(ctx context.Context, command Command) (interface{}, error) {
	commandType := reflect.TypeOf(command)
	handler, exists := cb.handlers[commandType]
	if !exists {
		return nil, &CommandHandlerNotFoundError{CommandType: commandType.Name()}
	}
	
	return handler.Handle(ctx, command)
}

// RegisterHandler 注册命令处理器
func (cb *SimpleCommandBus) RegisterHandler(commandType reflect.Type, handler CommandHandler) error {
	cb.handlers[commandType] = handler
	return nil
}

// CommandHandlerNotFoundError 命令处理器未找到错误
type CommandHandlerNotFoundError struct {
	CommandType string
}

func (e *CommandHandlerNotFoundError) Error() string {
	return "no handler found for command: " + e.CommandType
}

// Query 查询接口
type Query interface {
	QueryName() string
}

// QueryHandler 查询处理器接口
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (interface{}, error)
}

// QueryBus 查询总线接口
type QueryBus interface {
	Ask(ctx context.Context, query Query) (interface{}, error)
	RegisterHandler(queryType reflect.Type, handler QueryHandler) error
}

// SimpleQueryBus 简单查询总线实现
type SimpleQueryBus struct {
	handlers map[reflect.Type]QueryHandler
}

// NewSimpleQueryBus 创建简单查询总线
func NewSimpleQueryBus() *SimpleQueryBus {
	return &SimpleQueryBus{
		handlers: make(map[reflect.Type]QueryHandler),
	}
}

// Ask 执行查询
func (qb *SimpleQueryBus) Ask(ctx context.Context, query Query) (interface{}, error) {
	queryType := reflect.TypeOf(query)
	handler, exists := qb.handlers[queryType]
	if !exists {
		return nil, &QueryHandlerNotFoundError{QueryType: queryType.Name()}
	}
	
	return handler.Handle(ctx, query)
}

// RegisterHandler 注册查询处理器
func (qb *SimpleQueryBus) RegisterHandler(queryType reflect.Type, handler QueryHandler) error {
	qb.handlers[queryType] = handler
	return nil
}

// QueryHandlerNotFoundError 查询处理器未找到错误
type QueryHandlerNotFoundError struct {
	QueryType string
}

func (e *QueryHandlerNotFoundError) Error() string {
	return "no handler found for query: " + e.QueryType
}

// Result 分页查询结果
type Result[T any] struct {
	Data       T     `json:"data"`
	TotalCount int64 `json:"total_count,omitempty"`
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// NewResult 创建查询结果
func NewResult[T any](data T) *Result[T] {
	return &Result[T]{Data: data}
}

// NewPaginatedResult 创建分页查询结果
func NewPaginatedResult[T any](data T, totalCount int64, page, pageSize int) *Result[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int(totalCount) / pageSize
		if int(totalCount)%pageSize > 0 {
			totalPages++
		}
	}
	
	return &Result[T]{
		Data:       data,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}