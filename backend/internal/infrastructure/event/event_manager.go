package event

import (
	"context"
	"log"
)

// EventManager 事件管理器
// 协调事件总线和事件存储，提供统一的事件处理入口
type EventManager struct {
	eventBus   *EventBus
	eventStore EventStore
}

// NewEventManager 创建事件管理器
func NewEventManager(eventBus *EventBus, eventStore EventStore) *EventManager {
	return &EventManager{
		eventBus:   eventBus,
		eventStore: eventStore,
	}
}

// PublishAndStore 发布事件并存储
func (em *EventManager) PublishAndStore(ctx context.Context, domainEvent DomainEvent) error {
	// 先存储事件（保证事件不丢失）
	if err := em.eventStore.Store(ctx, domainEvent); err != nil {
		log.Printf("❌ 保存事件失败: %v", err)
		return err
	}

	// 再发布事件
	if err := em.eventBus.Publish(ctx, domainEvent); err != nil {
		log.Printf("⚠️ 发布事件失败: %v，但事件已存储", err)
		return err
	}

	log.Printf("✅ 事件已发布并存储: 类型=%s, ID=%s",
		domainEvent.GetEventType(), domainEvent.GetEventID())

	return nil
}

// PublishAndStoreSync 同步发布事件并存储
func (em *EventManager) PublishAndStoreSync(ctx context.Context, domainEvent DomainEvent) error {
	// 先存储事件
	if err := em.eventStore.Store(ctx, domainEvent); err != nil {
		log.Printf("❌ 保存事件失败: %v", err)
		return err
	}

	// 再同步发布事件
	if err := em.eventBus.PublishSync(ctx, domainEvent); err != nil {
		log.Printf("⚠️ 同步发布事件失败: %v，但事件已存储", err)
		return err
	}

	log.Printf("✅ 事件已同步发布并存储: 类型=%s, ID=%s",
		domainEvent.GetEventType(), domainEvent.GetEventID())

	return nil
}

// RegisterHandler 注册事件处理器
func (em *EventManager) RegisterHandler(eventType string, handler EventHandler) {
	em.eventBus.RegisterHandler(eventType, handler)
	log.Printf("📝 已注册事件处理器：事件类型=%s", eventType)
}

// GetEventBus 获取事件总线
func (em *EventManager) GetEventBus() *EventBus {
	return em.eventBus
}

// GetEventStore 获取事件存储
func (em *EventManager) GetEventStore() EventStore {
	return em.eventStore
}
