// Package converter 提供类型转换功能的接口定义
package converter

import "github.com/google/uuid"

// Converter 类型转换器接口
type Converter interface {
	// UUID 相关转换
	ToUUID(s string) (uuid.UUID, error)
	ToUUIDPtr(s string) (*uuid.UUID, error)
	ToUUIDSlice(ss []string) ([]uuid.UUID, error)
	ToUUIDPtrSlice(ss []string) ([]*uuid.UUID, error)
	ToUUIDOrDefault(s string, defaultValue uuid.UUID) uuid.UUID
	ToUUIDPtrOrDefault(s string, defaultValue *uuid.UUID) *uuid.UUID

	// 基础类型转换
	ToStringPtr(s string) *string

	ToInt64Slice(ss []string) ([]int64, error)
	ToFloat64Slice(ss []string) ([]float64, error)
	ToBoolSlice(ss []string) ([]bool, error)

	// 其他基础类型转换
	ToInt(s interface{}) (int, error)
	ToInt32(s interface{}) (int32, error)
	ToInt64(s interface{}) (int64, error)
	ToFloat32(s interface{}) (float32, error)
	ToFloat64(s interface{}) (float64, error)
	ToBool(s interface{}) (bool, error)
	ToString(s interface{}) (string, error)
	ToSlice(s interface{}) ([]interface{}, error)
	ToMap(s interface{}) (map[string]interface{}, error)
}
