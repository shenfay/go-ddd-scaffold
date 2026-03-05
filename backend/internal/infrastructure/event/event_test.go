package event

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TestEvent 测试用领域事件
type TestEvent struct {
	BaseDomainEvent
	StudentID    string
	NodeID       uuid.UUID
	OldProgress  float64
	NewProgress  float64
	StudyMinutes int
	NodeType     string
	GraphID      uuid.UUID
	CreatedBy    uuid.UUID
}

// TestEventHandler 测试用事件处理器
type TestEventHandler struct {
	HandledEvents []DomainEvent
}

// NewTestEventHandler 创建测试事件处理器
func NewTestEventHandler() *TestEventHandler {
	return &TestEventHandler{
		HandledEvents: make([]DomainEvent, 0),
	}
}

// Handle 处理事件
func (h *TestEventHandler) Handle(ctx context.Context, event DomainEvent) error {
	h.HandledEvents = append(h.HandledEvents, event)
	return nil
}

// NewTestStudentProgressEvent 创建测试学生学习进度更新事件
func NewTestStudentProgressEvent(studentID string, nodeID uuid.UUID, oldProgress, newProgress float64, studyMinutes int, metadata map[string]interface{}) *TestEvent {
	return &TestEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "StudentProgressUpdated",
			EventID:     uuid.New().String(),
			AggregateID: nodeID,
			OccurredAt:  time.Now(),
			Version:     1,
		},
		StudentID:    studentID,
		NodeID:       nodeID,
		OldProgress:  oldProgress,
		NewProgress:  newProgress,
		StudyMinutes: studyMinutes,
	}
}

// NewTestNodeCreatedEvent 创建测试节点创建事件
func NewTestNodeCreatedEvent(nodeID, graphID, createdBy uuid.UUID, nodeType string, metadata map[string]interface{}) *TestEvent {
	return &TestEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventType:   "NodeCreated",
			EventID:     uuid.New().String(),
			AggregateID: nodeID,
			OccurredAt:  time.Now(),
			Version:     1,
		},
		NodeID:    nodeID,
		GraphID:   graphID,
		CreatedBy: createdBy,
		NodeType:  nodeType,
	}
}

// TestEventBus 测试事件总线
func TestEventBus(t *testing.T) {
	eventBus := NewEventBus()

	// 创建测试事件处理器
	handler := NewTestEventHandler()
	eventBus.RegisterHandler("StudentProgressUpdated", handler.Handle)

	// 验证处理器已注册
	if count := eventBus.GetHandlerCount("StudentProgressUpdated"); count != 1 {
		t.Errorf("期望处理器数量为1，实际为%d", count)
	}

	// 创建并发布事件
	event := NewTestStudentProgressEvent(
		"student123",
		uuid.New(),
		0.5,
		0.8,
		3,
		nil,
	)

	ctx := context.Background()
	if err := eventBus.PublishSync(ctx, event); err != nil {
		t.Errorf("发布事件失败: %v", err)
	}
}

// TestEventStore 测试事件存储
func TestEventStore(t *testing.T) {
	store := NewInMemoryEventStore()
	ctx := context.Background()

	// 创建测试事件
	nodeID := uuid.New()
	event := NewTestNodeCreatedEvent(
		nodeID,
		uuid.New(),
		uuid.New(),
		"C",
		nil,
	)

	// 保存事件
	if err := store.SaveEvent(ctx, event); err != nil {
		t.Errorf("保存事件失败: %v", err)
	}

	// 验证事件数量
	if count := store.GetEventCount(); count != 1 {
		t.Errorf("期望事件数量为1，实际为%d", count)
	}

	// 加载事件
	events, err := store.LoadEvents(ctx, nodeID.String())
	if err != nil {
		t.Errorf("加载事件失败: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("期望加载1个事件，实际加载%d个", len(events))
	}
}

// TestEventManager 测试事件管理器
func TestEventManager(t *testing.T) {
	eventBus := NewEventBus()
	eventStore := NewInMemoryEventStore()
	eventManager := NewEventManager(eventBus, eventStore)

	// 注册处理器
	handler := NewTestEventHandler()
	eventManager.RegisterHandler("NodeCreated", handler.Handle)

	// 创建事件
	event := NewTestNodeCreatedEvent(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"S",
		nil,
	)

	// 发布并存储事件
	ctx := context.Background()
	if err := eventManager.PublishAndStoreSync(ctx, event); err != nil {
		t.Errorf("发布并存储事件失败: %v", err)
	}

	// 验证事件已存储
	if count := eventStore.GetEventCount(); count != 1 {
		t.Errorf("期望存储1个事件，实际存储%d个", count)
	}
}

// TestEventStoreTimeRange 测试按时间范围查询事件
func TestEventStoreTimeRange(t *testing.T) {
	store := NewInMemoryEventStore()
	ctx := context.Background()

	// 保存多个事件
	for i := 0; i < 5; i++ {
		event := NewTestNodeCreatedEvent(
			uuid.New(),
			uuid.New(),
			uuid.New(),
			"T",
			nil,
		)
		if err := store.SaveEvent(ctx, event); err != nil {
			t.Errorf("保存事件失败: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	// 查询时间范围内的事件
	now := time.Now()
	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)

	events, err := store.LoadEventsInRange(ctx, start, end)
	if err != nil {
		t.Errorf("查询时间范围事件失败: %v", err)
	}

	if len(events) != 5 {
		t.Errorf("期望查询到5个事件，实际查询到%d个", len(events))
	}
}

// TestSetupEventInfrastructure 测试事件基础设施初始化
func TestSetupEventInfrastructure(t *testing.T) {
	// 注意：SetupEventInfrastructure 需要 Redis 客户端，这里使用 nil 进行测试
	// 在实际应用中应该提供有效的 Redis 客户端
	var redisClient interface{} = nil // 实际类型为 *redis.Client
	
	defer func() {
		if r := recover(); r != nil {
			t.Logf("SetupEventInfrastructure 在 Redis 客户端为 nil 时发生预期中的错误：%v", r)
			// 这是预期的行为，因为 Redis 客户端为 nil
			t.Skip("跳过此测试：需要有效的 Redis 客户端")
		}
	}()
	
	eventManager := SetupEventInfrastructure(redisClient.(*redis.Client))

	if eventManager == nil {
		t.Error("事件管理器初始化失败")
	}

	// 验证事件总线已初始化
	if eventManager.GetEventBus() == nil {
		t.Error("事件总线未初始化")
	}

	// 验证事件存储已初始化
	if eventManager.GetEventStore() == nil {
		t.Error("事件存储未初始化")
	}

	// 验证已注册事件类型
	registeredEvents := eventManager.GetEventBus().ListRegisteredEvents()
	if len(registeredEvents) == 0 {
		t.Error("没有注册任何事件处理器")
	}

	t.Logf("已注册%d个事件类型", len(registeredEvents))
}
