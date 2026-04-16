package shared

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewBaseEvent 测试基础事件创建
func TestNewBaseEvent(t *testing.T) {
	t.Run("should create event with valid timestamp", func(t *testing.T) {
		before := time.Now()
		event := NewBaseEvent()
		after := time.Now()

		// Timestamp 应该在合理范围内
		assert.True(t, event.Timestamp.After(before) || event.Timestamp.Equal(before))
		assert.True(t, event.Timestamp.Before(after) || event.Timestamp.Equal(after))
	})

	t.Run("should have increasing timestamps", func(t *testing.T) {
		event1 := NewBaseEvent()
		time.Sleep(1 * time.Millisecond)
		event2 := NewBaseEvent()

		// 后创建的事件时间戳应该不小于先创建的
		assert.True(t, event2.Timestamp.After(event1.Timestamp) || event2.Timestamp.Equal(event1.Timestamp))
	})
}

// TestBaseEvent_Fields 测试基础事件字段
func TestBaseEvent_Fields(t *testing.T) {
	t.Run("should have timestamp populated", func(t *testing.T) {
		event := NewBaseEvent()

		assert.NotZero(t, event.Timestamp)
		assert.IsType(t, time.Time{}, event.Timestamp)
	})
}

// TestDomainErrors 测试领域错误定义
func TestDomainErrors(t *testing.T) {
	t.Run("ErrInvalidArgument should be defined", func(t *testing.T) {
		assert.NotNil(t, ErrInvalidArgument)
		assert.Equal(t, "invalid argument", ErrInvalidArgument.Error())
	})

	t.Run("ErrNotFound should be defined", func(t *testing.T) {
		assert.NotNil(t, ErrNotFound)
		assert.Equal(t, "not found", ErrNotFound.Error())
	})

	t.Run("errors should be distinct", func(t *testing.T) {
		assert.NotEqual(t, ErrInvalidArgument, ErrNotFound)
	})
}
