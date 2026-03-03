package event

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domain_event "go-ddd-scaffold/internal/domain/knowledge/event"
)

// TestEventBus 测试事件总线
func TestEventBus(t *testing.T) {
	eventBus := NewEventBus()

	// 创建测试事件处理器
	handler := domain_event.NewLearningProgressEventHandler()
	eventBus.RegisterHandler("StudentProgressUpdated", handler)

	// 验证处理器已注册
	if count := eventBus.GetHandlerCount("StudentProgressUpdated"); count != 1 {
		t.Errorf("期望处理器数量为1，实际为%d", count)
	}

	// 创建并发布事件
	event := domain_event.NewStudentProgressUpdatedEvent(
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
	event := domain_event.NewNodeCreatedEvent(
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
	handler := domain_event.NewAnalyticsEventHandler()
	eventManager.RegisterHandler("NodeCreated", handler)

	// 创建事件
	event := domain_event.NewNodeCreatedEvent(
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
		event := domain_event.NewNodeCreatedEvent(
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
	eventManager := SetupEventInfrastructure()

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
