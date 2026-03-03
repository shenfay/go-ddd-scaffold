// Package converter 提供默认类型转换器实现
package converter

import (
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

// DefaultConverter 默认转换器实现
type DefaultConverter struct{}

// NewDefaultConverter 创建默认转换器实例
func NewDefaultConverter() *DefaultConverter {
	return &DefaultConverter{}
}

// ToUUID 将字符串转换为UUID
func (dc *DefaultConverter) ToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// ToUUIDPtr 将字符串转换为UUID指针
func (dc *DefaultConverter) ToUUIDPtr(s string) (*uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ToUUIDSlice 将字符串切片转换为UUID切片
func (dc *DefaultConverter) ToUUIDSlice(ss []string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, len(ss))
	for i, s := range ss {
		u, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		result[i] = u
	}
	return result, nil
}

// ToUUIDPtrSlice 将字符串切片转换为UUID指针切片
func (dc *DefaultConverter) ToUUIDPtrSlice(ss []string) ([]*uuid.UUID, error) {
	result := make([]*uuid.UUID, len(ss))
	for i, s := range ss {
		if s == "" {
			result[i] = nil
			continue
		}
		u, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		result[i] = &u
	}
	return result, nil
}

// ToUUIDOrDefault 将字符串转换为UUID，转换失败时返回默认值
func (dc *DefaultConverter) ToUUIDOrDefault(s string, defaultValue uuid.UUID) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return defaultValue
	}
	return u
}

// ToUUIDPtrOrDefault 将字符串转换为UUID指针，转换失败时返回默认值
func (dc *DefaultConverter) ToUUIDPtrOrDefault(s string, defaultValue *uuid.UUID) *uuid.UUID {
	if s == "" {
		return nil
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return defaultValue
	}
	return &u
}

// ToStringPtr 将字符串转换为字符串指针
func (dc *DefaultConverter) ToStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ToInt64Slice 将字符串切片转换为int64切片
func (dc *DefaultConverter) ToInt64Slice(ss []string) ([]int64, error) {
	return cast.ToInt64SliceE(ss)
}

// ToFloat64Slice 将字符串切片转换为float64切片
func (dc *DefaultConverter) ToFloat64Slice(ss []string) ([]float64, error) {
	return cast.ToFloat64SliceE(ss)
}

// ToBoolSlice 将字符串切片转换为bool切片
func (dc *DefaultConverter) ToBoolSlice(ss []string) ([]bool, error) {
	return cast.ToBoolSliceE(ss)
}

// ToInt 将接口转换为int
func (dc *DefaultConverter) ToInt(s interface{}) (int, error) {
	return cast.ToIntE(s)
}

// ToInt32 将接口转换为int32
func (dc *DefaultConverter) ToInt32(s interface{}) (int32, error) {
	return cast.ToInt32E(s)
}

// ToInt64 将接口转换为int64
func (dc *DefaultConverter) ToInt64(s interface{}) (int64, error) {
	return cast.ToInt64E(s)
}

// ToFloat32 将接口转换为float32
func (dc *DefaultConverter) ToFloat32(s interface{}) (float32, error) {
	return cast.ToFloat32E(s)
}

// ToFloat64 将接口转换为float64
func (dc *DefaultConverter) ToFloat64(s interface{}) (float64, error) {
	return cast.ToFloat64E(s)
}

// ToBool 将接口转换为bool
func (dc *DefaultConverter) ToBool(s interface{}) (bool, error) {
	return cast.ToBoolE(s)
}

// ToString 将接口转换为string
func (dc *DefaultConverter) ToString(s interface{}) (string, error) {
	return cast.ToStringE(s)
}

// ToSlice 将接口转换为切片
func (dc *DefaultConverter) ToSlice(s interface{}) ([]interface{}, error) {
	return cast.ToSliceE(s)
}

// ToMap 将接口转换为map
func (dc *DefaultConverter) ToMap(s interface{}) (map[string]interface{}, error) {
	return cast.ToStringMapE(s)
}
