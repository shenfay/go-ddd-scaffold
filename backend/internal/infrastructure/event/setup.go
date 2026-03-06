package event

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	user_event "go-ddd-scaffold/internal/domain/user/event"

	"gorm.io/gorm"
)

// InitializeEventHandlers 初始化事件处理器
// 在应用启动时调用，注册所有事件处理器
func InitializeEventHandlers(eventManager *EventManager, db *gorm.DB) {
	log.Println("🚀 开始初始化事件处理器...")

	// 用户注册事件处理器
	userCreatedHandler := func(ctx context.Context, event DomainEvent) error {
		registeredEvent, ok := event.(*user_event.UserRegisteredEvent)
		if !ok {
			return nil
		}

		log.Printf("✅ 处理用户注册事件：UserID=%s, Email=%s", registeredEvent.UserID, registeredEvent.Email)

		// TODO: 根据业务需要添加实际处理逻辑
		// 示例：发送欢迎邮件（需要实现邮件服务）
		// emailService.SendWelcomeEmail(registeredEvent.Email)

		// 示例：初始化用户数据（如创建默认配置、统计信息等）
		// initUserStatistics(db, registeredEvent.UserID)

		return nil
	}
	eventManager.RegisterHandler("UserRegistered", userCreatedHandler)

	// 用户登录事件处理器
	userLoggedInHandler := func(ctx context.Context, event DomainEvent) error {
		loginEvent, ok := event.(*user_event.UserLoggedInEvent)
		if !ok {
			return nil
		}

		log.Printf("✅ 处理用户登录事件：UserID=%s, IP=%s", loginEvent.UserID, loginEvent.IP)

		// 实际业务逻辑：记录登录日志
		if err := recordLoginLog(db, loginEvent); err != nil {
			log.Printf("⚠️ 记录登录日志失败：%v", err)
			// 不返回错误，避免影响主流程
		}

		return nil
	}
	eventManager.RegisterHandler("UserLoggedIn", userLoggedInHandler)

	log.Println("✅ 事件处理器初始化完成")
	log.Printf("📊 已注册事件类型：%v", eventManager.GetEventBus().ListRegisteredEvents())
}

// recordLoginLog 记录用户登录日志（使用 Map 方式插入，避免字段匹配问题）
func recordLoginLog(db *gorm.DB, event *user_event.UserLoggedInEvent) error {
	now := time.Now()

	data := map[string]interface{}{
		"id":             uuid.New().String(),
		"user_id":        event.UserID.String(),
		"ip_address":     event.IP,
		"device_type":    event.DeviceType,
		"os_info":        "",  // TODO: 需要从 UserLoggedInEvent 添加 OSInfo 字段
		"browser_info":   "",  // TODO: 需要从 UserLoggedInEvent 添加 BrowserInfo 字段
		"login_status":   event.LoginStatus,
		"failure_reason": event.FailureReason,
		"logged_at":      now,
		"created_at":     now,
	}

	return db.Table("login_logs").Create(data).Error
}
