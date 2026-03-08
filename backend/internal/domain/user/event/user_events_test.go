package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-ddd-scaffold/internal/domain/user/event"
)

func TestUserRegisteredEvent(t *testing.T) {
	t.Run("创建用户注册事件", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"

		eventObj := event.NewUserRegisteredEvent(userID, email)

		// 验证基本字段
		assert.Equal(t, userID, eventObj.UserID)
		assert.Equal(t, email, eventObj.Email)
		assert.Equal(t, "UserRegistered", eventObj.EventType)
		assert.Equal(t, userID, eventObj.AggregateID)
		
		// 验证事件 ID 不为空
		assert.NotEmpty(t, eventObj.EventID)
		
		// 验证时间在合理范围内（前后 1 秒）
		now := time.Now()
		assert.WithinDuration(t, now, eventObj.OccurredAt, 1*time.Second)
		
		// 验证版本号
		assert.Equal(t, 1, eventObj.Version)
	})

	t.Run("事件接口方法", func(t *testing.T) {
		userID := uuid.New()
		eventObj := event.NewUserRegisteredEvent(userID, "test@example.com")

		assert.Equal(t, "UserRegistered", eventObj.GetEventType())
		assert.NotEmpty(t, eventObj.GetEventID())
		assert.Equal(t, userID, eventObj.GetAggregateID())
		assert.NotNil(t, eventObj.GetOccurredAt())
		assert.Equal(t, 1, eventObj.GetVersion())
	})
}

func TestUserLoggedInEvent(t *testing.T) {
	t.Run("创建用户登录事件 - 成功", func(t *testing.T) {
		userID := uuid.New()
		ip := "192.168.1.100"
		userAgent := "Mozilla/5.0"
		deviceType := "web"
		loginStatus := "success"

		eventObj := event.NewUserLoggedInEvent(
			userID, ip, userAgent, deviceType, loginStatus, nil,
		)

		assert.Equal(t, userID, eventObj.UserID)
		assert.Equal(t, ip, eventObj.IP)
		assert.Equal(t, userAgent, eventObj.UserAgent)
		assert.Equal(t, deviceType, eventObj.DeviceType)
		assert.Equal(t, loginStatus, eventObj.LoginStatus)
		assert.Nil(t, eventObj.FailureReason)
		assert.Equal(t, "UserLoggedIn", eventObj.EventType)
	})

	t.Run("创建用户登录事件 - 失败", func(t *testing.T) {
		userID := uuid.New()
		ip := "192.168.1.100"
		userAgent := "Mozilla/5.0"
		deviceType := "mobile"
		loginStatus := "failed"
		failureReason := "密码错误"

		eventObj := event.NewUserLoggedInEvent(
			userID, ip, userAgent, deviceType, loginStatus, &failureReason,
		)

		assert.Equal(t, userID, eventObj.UserID)
		assert.Equal(t, loginStatus, eventObj.LoginStatus)
		assert.NotNil(t, eventObj.FailureReason)
		assert.Equal(t, failureReason, *eventObj.FailureReason)
	})
}

func TestUserLockedEvent(t *testing.T) {
	t.Run("创建用户锁定事件", func(t *testing.T) {
		userID := uuid.New()

		eventObj := event.NewUserLockedEvent(userID)

		assert.Equal(t, userID, eventObj.UserID)
		assert.Equal(t, "UserLocked", eventObj.EventType)
		assert.Equal(t, userID, eventObj.AggregateID)
		assert.NotEmpty(t, eventObj.EventID)
		assert.NotNil(t, eventObj.OccurredAt)
		assert.Equal(t, 1, eventObj.Version)
	})
}

func TestUserActivatedEvent(t *testing.T) {
	t.Run("创建用户激活事件", func(t *testing.T) {
		userID := uuid.New()

		eventObj := event.NewUserActivatedEvent(userID)

		assert.Equal(t, userID, eventObj.UserID)
		assert.Equal(t, "UserActivated", eventObj.EventType)
		assert.Equal(t, userID, eventObj.AggregateID)
		assert.NotEmpty(t, eventObj.EventID)
		assert.NotNil(t, eventObj.OccurredAt)
		assert.Equal(t, 1, eventObj.Version)
	})
}

func TestUserProfileUpdatedEvent(t *testing.T) {
	t.Run("创建用户资料更新事件 - 完整信息", func(t *testing.T) {
		userID := uuid.New()
		nickname := "新昵称"
		phone := "13800138000"
		bio := "个人简介"

		eventObj := event.NewUserProfileUpdatedEvent(
			userID, nickname, &phone, &bio,
		)

		assert.Equal(t, userID, eventObj.UserID)
		assert.Equal(t, nickname, eventObj.Nickname)
		assert.Equal(t, phone, *eventObj.Phone)
		assert.Equal(t, bio, *eventObj.Bio)
		assert.Equal(t, "UserProfileUpdated", eventObj.EventType)
	})

	t.Run("创建用户资料更新事件 - 部分信息", func(t *testing.T) {
		userID := uuid.New()
		nickname := "新昵称"

		eventObj := event.NewUserProfileUpdatedEvent(
			userID, nickname, nil, nil,
		)

		assert.Equal(t, userID, eventObj.UserID)
		assert.Equal(t, nickname, eventObj.Nickname)
		assert.Nil(t, eventObj.Phone)
		assert.Nil(t, eventObj.Bio)
	})
}

func TestUserEmailChangedEvent(t *testing.T) {
	t.Run("创建用户邮箱变更事件", func(t *testing.T) {
		userID := uuid.New()
		oldEmail := "old@example.com"
		newEmail := "new@example.com"

		eventObj := event.NewUserEmailChangedEvent(userID, oldEmail, newEmail)

		assert.Equal(t, userID, eventObj.UserID)
		assert.Equal(t, oldEmail, eventObj.OldEmail)
		assert.Equal(t, newEmail, eventObj.NewEmail)
		assert.Equal(t, "UserEmailChanged", eventObj.EventType)
		assert.NotEmpty(t, eventObj.EventID)
		assert.NotNil(t, eventObj.OccurredAt)
		assert.Equal(t, 1, eventObj.Version)
	})
}
