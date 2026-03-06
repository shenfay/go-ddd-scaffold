package entity_test

import (
	"testing"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_Lock(t *testing.T) {
	t.Run("成功锁定活跃用户", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)

		// Act
		err := user.Lock()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.StatusLocked, user.Status)
		assert.Len(t, user.Events(), 1)
	})

	t.Run("锁定已锁定的用户返回错误", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusLocked)

		// Act
		err := user.Lock()

		// Assert
		assert.Error(t, err)
		assert.IsType(t, entity.ErrAlreadyLocked(""), err)
	})

	t.Run("锁定非活跃用户", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusInactive)

		// Act
		err := user.Lock()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.StatusLocked, user.Status)
	})
}

func TestUser_Activate(t *testing.T) {
	t.Run("激活锁定用户", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusLocked)

		// Act
		err := user.Activate()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.StatusActive, user.Status)
		assert.Len(t, user.Events(), 1)
	})

	t.Run("激活已活跃用户", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)

		// Act
		err := user.Activate()

		// Assert
		assert.NoError(t, err) // 应该不报错（幂等性）
		assert.Equal(t, entity.StatusActive, user.Status)
	})

	t.Run("激活非活跃用户", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusInactive)

		// Act
		err := user.Activate()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.StatusActive, user.Status)
	})
}

func TestUser_UpdateProfile(t *testing.T) {
	t.Run("成功更新个人资料", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)
		nickname, _ := valueobject.NewNickname("New Nickname")
		phone := "13800138000"
		bio := "Test bio"

		// Act
		err := user.UpdateProfile(nickname, &phone, &bio)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, nickname, user.Nickname)
		assert.Equal(t, phone, *user.Phone)
		assert.Equal(t, bio, *user.Bio)
		assert.Len(t, user.Events(), 1)
	})

	t.Run("更新部分字段", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)
		nickname, _ := valueobject.NewNickname("Updated Nickname")

		// Act
		err := user.UpdateProfile(nickname, nil, nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, nickname, user.Nickname)
		assert.Nil(t, user.Phone)
		assert.Nil(t, user.Bio)
	})
}

func TestUser_UpdateEmail(t *testing.T) {
	t.Run("成功更新邮箱", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)
		oldEmail := user.Email
		newEmail, _ := valueobject.NewEmail("newemail@example.com")

		// Act
		err := user.UpdateEmail(newEmail)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, newEmail, user.Email)
		assert.NotEqual(t, oldEmail, user.Email)
		assert.Len(t, user.Events(), 1)
	})

	t.Run("更新为相同邮箱", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)
		sameEmail := user.Email

		// Act
		err := user.UpdateEmail(sameEmail)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, user.Events(), 0) // 不应该产生事件
	})
}

func TestUser_AddDomainEvent(t *testing.T) {
	t.Run("添加领域事件", func(t *testing.T) {
		// Act & Assert - Events() 是内部方法，外部无法直接添加
		// 这个测试跳过，因为 addEvent 是内部方法
		t.Skip("addEvent 是内部方法，不对外暴露")
	})

	t.Run("获取事件列表", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)

		// Act
		events := user.Events()

		// Assert
		assert.Empty(t, events) // 初始为空
		// 注意：Events() 返回的是当前事件列表
	})
}

func TestUser_StatusTransitions(t *testing.T) {
	t.Run("完整的状态流转", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)

		// Act & Assert - Active -> Locked
		err := user.Lock()
		assert.NoError(t, err)
		assert.Equal(t, entity.StatusLocked, user.Status)

		// Locked -> Active
		err = user.Activate()
		assert.NoError(t, err)
		assert.Equal(t, entity.StatusActive, user.Status)

		// Active -> Inactive (假设有这个方法)
		// user.Deactivate()
		// assert.Equal(t, entity.StatusInactive, user.Status)

		// Inactive -> Active
		// err = user.Activate()
		// assert.NoError(t, err)
		// assert.Equal(t, entity.StatusActive, user.Status)
	})
}

func TestUser_DomainEvents(t *testing.T) {
	t.Run("Lock 产生 UserLockedEvent", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)

		// Act
		user.Lock()

		// Assert
		events := user.Events()
		assert.Len(t, events, 1)
		// 这里可以进一步检查事件类型和内容
	})

	t.Run("Activate 产生 UserActivatedEvent", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusLocked)

		// Act
		user.Activate()

		// Assert
		events := user.Events()
		assert.Len(t, events, 1)
	})

	t.Run("UpdateProfile 产生 UserProfileUpdatedEvent", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)
		nickname, _ := valueobject.NewNickname("Test Nickname")

		// Act
		user.UpdateProfile(nickname, nil, nil)

		// Assert
		events := user.Events()
		assert.Len(t, events, 1)
	})

	t.Run("UpdateEmail 产生 UserEmailChangedEvent", func(t *testing.T) {
		// Arrange
		user := createTestUser(entity.StatusActive)
		newEmail, _ := valueobject.NewEmail("new@example.com")

		// Act
		user.UpdateEmail(newEmail)

		// Assert
		events := user.Events()
		assert.Len(t, events, 1)
	})
}

// Helper functions
func createTestUser(status entity.UserStatus) *entity.User {
	email, _ := valueobject.NewEmail("test@example.com")
	nickname, _ := valueobject.NewNickname("Test User")

	return &entity.User{
		ID:       uuid.New(),
		Email:    email,
		Nickname: nickname,
		Status:   status,
		Password: entity.HashedPassword("hashed_password"),
	}
}

// Test event implementation
type testEvent struct{}

func (e *testEvent) GetEventID() string          { return "test-event-id" }
func (e *testEvent) GetEventType() string        { return "TestEvent" }
func (e *testEvent) GetAggregateID() uuid.UUID   { return uuid.Nil }
func (e *testEvent) GetOccurredAt() interface{}  { return nil }
func (e *testEvent) GetVersion() int             { return 1 }
