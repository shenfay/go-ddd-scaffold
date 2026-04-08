package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/internal/infra/messaging"
	"github.com/shenfay/go-ddd-scaffold/pkg/event"
	"github.com/shenfay/go-ddd-scaffold/pkg/logger"
)

// Server Worker 服务器
type Server struct {
	srv      *asynq.Server
	eventBus *messaging.AsynqEventBus
}

// NewServer 创建 Worker 服务器
func NewServer(redisAddr, redisPwd string, db int, concurrency int) *Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPwd,
			DB:       db,
		},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			StrictPriority: true,
		},
	)

	return &Server{srv: srv}
}

// SetEventBus 设置事件总线
func (s *Server) SetEventBus(eventBus *messaging.AsynqEventBus) {
	s.eventBus = eventBus
}

// RegisterDomainEvents 自动注册所有领域事件处理器
func (s *Server) RegisterDomainEvents(mux *asynq.ServeMux) {
	if s.eventBus == nil {
		return
	}

	handlers := s.eventBus.GetSubscriptions()
	for eventType := range handlers {
		evtType := eventType // 闭包捕获
		mux.HandleFunc(evtType, func(ctx context.Context, task *asynq.Task) error {
			return s.processDomainEvent(ctx, evtType, task)
		})
		logger.Info("✓ Registered handler for: ", evtType)
	}
}

// RegisterHandler 注册自定义任务处理器
func (s *Server) RegisterHandler(mux *asynq.ServeMux, taskType string, handler func(ctx context.Context, task *asynq.Task) error) {
	mux.HandleFunc(taskType, handler)
	logger.Info("✓ Registered handler for type: ", taskType)
}

// Start 启动 Worker
func (s *Server) Start(mux *asynq.ServeMux) error {
	logger.Info("🎯 Starting Asynq Worker processor...")
	return s.srv.Run(mux)
}

// Shutdown 优雅关闭
func (s *Server) Shutdown() {
	s.srv.Shutdown()
	logger.Info("✅ Worker stopped gracefully")
}

// processDomainEvent 处理领域事件（调用订阅的 Listener）
func (s *Server) processDomainEvent(ctx context.Context, eventType string, task *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	// 重构建事件对象
	evt := &GenericEvent{
		Type:    eventType,
		Payload: payload,
	}

	// 调用所有订阅的处理器
	return s.eventBus.DispatchEvent(ctx, evt)
}

// GenericEvent 通用事件（用于反序列化）
type GenericEvent struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// GetType 获取事件类型
func (e *GenericEvent) GetType() string {
	return e.Type
}

// GetPayload 获取事件载荷
func (e *GenericEvent) GetPayload() interface{} {
	return e.Payload
}

// 确保实现 event.Event 接口
var _ event.Event = (*GenericEvent)(nil)
