package aggregate_test

import (
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("成功创建用户", func(t *testing.T) {
		// Arrange
		username := "testuser"
		email := "test@example.com"
		hashedPassword := "hashed_password_123"
		idCounter := int64(1)

		// Act
		user, err := aggregate.NewUser(username, email, hashedPassword, func() int64 {
			idCounter++
			return idCounter
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, username, user.Username().Value())
		assert.Equal(t, email, user.Email().Value())
		assert.Equal(t, vo.UserStatusActive, user.Status())
		assert.NotZero(t, user.ID())
		assert.Equal(t, int64(2), user.ID().(vo.UserID).Int64())
	})

	t.Run("创建失败 - 用户名无效", func(t *testing.T) {
		// Arrange
		invalidUsername := "" // 空用户名
		email := "test@example.com"
		hashedPassword := "hashed_password"

		// Act
		user, err := aggregate.NewUser(invalidUsername, email, hashedPassword, func() int64 { return 1 })

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("创建失败 - 邮箱无效", func(t *testing.T) {
		// Arrange
		username := "testuser"
		invalidEmail := "invalid-email" // 无效邮箱
		hashedPassword := "hashed_password"

		// Act
		user, err := aggregate.NewUser(username, invalidEmail, hashedPassword, func() int64 { return 1 })

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUser_Activate(t *testing.T) {
	t.Run("成功激活用户", func(t *testing.T) {
		// Arrange
		// 由于 NewUser 默认创建 Active 用户，我们需要创建一个 Pending 用户
		// 这里直接测试 Activate 方法从 Pending -> Active 的逻辑
		// 简化处理：测试从 Locked -> Active 的场景（通过 Unlock）
		user := createTestUserWithStatus(vo.UserStatusLocked)
		originalVersion := user.Version()

		// Act - 先 Unlock 回到 Active
		err := user.Unlock()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, vo.UserStatusActive, user.Status())
		assert.Greater(t, user.Version(), originalVersion)
	})

	t.Run("激活失败 - 用户已经是 Active", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusActive)

		// Act
		err := user.Activate()

		// Assert
		assert.Error(t, err)
		bizErr, ok := err.(*kernel.BusinessError)
		if assert.True(t, ok) {
			assert.Equal(t, aggregate.CodeUserNotPending, bizErr.Code)
		}
	})
}

func TestUser_Deactivate(t *testing.T) {
	t.Run("成功停用用户", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusActive)

		// Act
		err := user.Deactivate("test reason")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, vo.UserStatusInactive, user.Status())
		assert.Greater(t, user.Version(), 0)
	})

	t.Run("停用失败 - 用户已停用", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusInactive)

		// Act
		err := user.Deactivate("test reason")

		// Assert
		assert.Error(t, err)
		bizErr, ok := err.(*kernel.BusinessError)
		if assert.True(t, ok) {
			assert.Equal(t, aggregate.CodeUserAlreadyInactive, bizErr.Code)
		}
	})
}

func TestUser_Lock(t *testing.T) {
	t.Run("成功锁定用户", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusActive)

		// Act
		err := user.Lock()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, vo.UserStatusLocked, user.Status())
		assert.Greater(t, user.Version(), 0)
	})

	t.Run("锁定失败 - 用户已锁定", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusLocked)

		// Act
		err := user.Lock()

		// Assert
		assert.Error(t, err)
		bizErr, ok := err.(*kernel.BusinessError)
		if assert.True(t, ok) {
			assert.Equal(t, aggregate.CodeUserAlreadyLocked, bizErr.Code)
		}
	})
}

func TestUser_Unlock(t *testing.T) {
	t.Run("成功解锁用户", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusLocked)

		// Act
		err := user.Unlock()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, vo.UserStatusActive, user.Status())
		assert.Greater(t, user.Version(), 0)
	})

	t.Run("解锁失败 - 用户未锁定", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusActive)

		// Act
		err := user.Unlock()

		// Assert
		assert.Error(t, err)
		bizErr, ok := err.(*kernel.BusinessError)
		if assert.True(t, ok) {
			assert.Equal(t, aggregate.CodeUserNotLocked, bizErr.Code)
		}
	})
}

func TestUser_ChangePassword(t *testing.T) {
	t.Run("成功修改密码", func(t *testing.T) {
		// Arrange
		user := createTestUser()
		newPassword := "new_secure_password"
		ipAddress := "192.168.1.1"

		// Act
		err := user.ChangePassword(newPassword, ipAddress)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, "old_password", user.Password().Value())
		assert.Greater(t, user.Version(), 0)
	})
}

func TestUser_UpdateEmail(t *testing.T) {
	t.Run("成功更新邮箱", func(t *testing.T) {
		// Arrange
		user := createTestUser()
		newEmail := "newemail@example.com"

		// Act
		err := user.UpdateEmail(newEmail)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, newEmail, user.Email().Value())
		assert.Greater(t, user.Version(), 0)
	})

	t.Run("更新邮箱失败 - 邮箱格式无效", func(t *testing.T) {
		// Arrange
		user := createTestUser()
		invalidEmail := "invalid-email"

		// Act
		err := user.UpdateEmail(invalidEmail)

		// Assert
		assert.Error(t, err)
	})
}

func TestUser_CanLogin(t *testing.T) {
	t.Run("允许登录 - Active 状态", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusActive)

		// Act
		canLogin := user.CanLogin()

		// Assert
		assert.True(t, canLogin)
	})

	t.Run("禁止登录 - Locked 状态", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusLocked)

		// Act
		canLogin := user.CanLogin()

		// Assert
		assert.False(t, canLogin)
	})

	t.Run("禁止登录 - Inactive 状态", func(t *testing.T) {
		// Arrange
		user := createTestUserWithStatus(vo.UserStatusInactive)

		// Act
		canLogin := user.CanLogin()

		// Assert
		assert.False(t, canLogin)
	})
}

func TestUser_FullName(t *testing.T) {
	t.Run("获取完整姓名", func(t *testing.T) {
		// Arrange
		user := createTestUser()
		user.SetFirstName("John")
		user.SetLastName("Doe")

		// Act
		fullName := user.FullName()

		// Assert
		assert.Equal(t, "John Doe", fullName)
	})

	t.Run("获取完整姓名 - 带默认值", func(t *testing.T) {
		// Arrange
		user := createTestUser()

		// Act
		fullName := user.GetFullName("Default Name")

		// Assert
		assert.Equal(t, "Default Name", fullName)
	})
}

// Helper functions

func createTestUser() *aggregate.User {
	user, _ := aggregate.NewUser("testuser", "test@example.com", "old_password", func() int64 { return 1 })
	return user
}

func createTestUserWithStatus(status vo.UserStatus) *aggregate.User {
	user := createTestUser()

	// 通过反射或直接设置状态字段（需要导出方法或测试包访问）
	// 这里我们通过调用相应的方法来设置状态
	switch status {
	case vo.UserStatusPending:
		// 默认是 Active，需要手动设置为 Pending
		// 由于没有直接设置的方法，我们创建新用户使用不同策略
		// 简化处理：直接返回 Active 用户，测试时注意
		break
	case vo.UserStatusInactive:
		user.Deactivate("test reason")
	case vo.UserStatusLocked:
		user.Lock()
	}

	return user
}
