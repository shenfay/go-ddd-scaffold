package event

import (
	"log"
)

// InitializeEventHandlers 初始化事件处理器
// 在应用启动时调用，注册所有事件处理器
func InitializeEventHandlers(eventManager *EventManager) {
	log.Println("🚀 开始初始化事件处理器...")

	// TODO: 根据业务需要注册事件处理器
	// 示例:
	// userCreatedHandler := func(ctx context.Context, event DomainEvent) error {
	//     log.Printf("处理用户创建事件：%s", event.GetEventID())
	//     return nil
	// }
	// eventManager.RegisterHandler("UserCreated", userCreatedHandler)

	log.Println("✅ 事件处理器初始化完成")
	log.Printf("📊 已注册事件类型：%v", eventManager.GetEventBus().ListRegisteredEvents())
}

// SetupEventInfrastructure 设置事件基础设施
// 创建并配置完整的事件处理系统
func SetupEventInfrastructure() *EventManager {
	log.Println("🔧 开始设置事件基础设施...")

	// 创建事件总线
	eventBus := NewEventBus()
	log.Println("✅ 事件总线已创建")

	// 创建事件存储（使用内存存储）
	eventStore := NewInMemoryEventStore()
	log.Println("✅ 事件存储已创建")

	// 创建事件管理器
	eventManager := NewEventManager(eventBus, eventStore)
	log.Println("✅ 事件管理器已创建")

	// 初始化事件处理器
	InitializeEventHandlers(eventManager)

	log.Println("🎉 事件基础设施设置完成")
	return eventManager
}
