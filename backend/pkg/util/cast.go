package util

import (
	"time"

	"github.com/spf13/cast"
)

// ToString 转换为字符串
func ToString(v interface{}) string {
	return cast.ToString(v)
}

// ToInt 转换为 int
func ToInt(v interface{}) int {
	return cast.ToInt(v)
}

// ToInt64 转换为 int64
func ToInt64(v interface{}) int64 {
	return cast.ToInt64(v)
}

// ToInt32 转换为 int32
func ToInt32(v interface{}) int32 {
	return cast.ToInt32(v)
}

// ToInt16 转换为 int16
func ToInt16(v interface{}) int16 {
	return cast.ToInt16(v)
}

// ToBool 转换为 bool
func ToBool(v interface{}) bool {
	return cast.ToBool(v)
}

// ToTime 转换为 time.Time
func ToTime(v interface{}) time.Time {
	return cast.ToTime(v)
}

// ToStringSlice 转换为 []string
func ToStringSlice(v interface{}) []string {
	return cast.ToStringSlice(v)
}

// ToIntSlice 转换为 []int
func ToIntSlice(v interface{}) []int {
	return cast.ToIntSlice(v)
}

// ToInt64Slice 转换为 []int64
func ToInt64Slice(v interface{}) []int64 {
	return cast.ToInt64Slice(v)
}

// ToStringMap 转换为 map[string]interface{}
func ToStringMap(v interface{}) map[string]interface{} {
	return cast.ToStringMap(v)
}

// ToStringMapString 转换为 map[string]string
func ToStringMapString(v interface{}) map[string]string {
	return cast.ToStringMapString(v)
}

// ToStringMapStringSlice 转换为 map[string][]string
func ToStringMapStringSlice(v interface{}) map[string][]string {
	return cast.ToStringMapStringSlice(v)
}

// StringPtr 创建 *string
func StringPtr(s string) *string {
	return &s
}

// BoolPtr 创建 *bool
func BoolPtr(b bool) *bool {
	return &b
}

// Int32Ptr 创建 *int32
func Int32Ptr(i int) *int32 {
	v := int32(i)
	return &v
}

// Int64Ptr 创建 *int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// Int16Ptr 创建 *int16
func Int16Ptr(i int) *int16 {
	v := int16(i)
	return &v
}

// TimePtr 创建 *time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// StringValue 获取 *string 的值
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BoolValue 获取 *bool 的值
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Int32Value 获取 *int32 的值
func Int32Value(i *int32) int {
	if i == nil {
		return 0
	}
	return int(*i)
}

// Int64Value 获取 *int64 的值
func Int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// Int16Value 获取 *int16 的值
func Int16Value(i *int16) int16 {
	if i == nil {
		return 0
	}
	return *i
}

// TimeValue 获取 *time.Time 的值
func TimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
