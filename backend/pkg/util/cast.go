package util

import (
	"time"

	"github.com/spf13/cast"
)

// ========== 1. 类型转换（To 前缀，处理任意输入） ==========

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

// ========== 2. 创建指针（类型名即函数名，类型必须匹配） ==========

// String 创建 *string
func String(s string) *string {
	return &s
}

// Bool 创建 *bool
func Bool(b bool) *bool {
	return &b
}

// Int32 创建 *int32
func Int32(i int32) *int32 {
	return &i
}

// Int64 创建 *int64
func Int64(i int64) *int64 {
	return &i
}

// Int16 创建 *int16
func Int16(i int16) *int16 {
	return &i
}

// Time 创建 *time.Time
func Time(t time.Time) *time.Time {
	return &t
}

// ========== 3. 获取值（Value 后缀，安全防护） ==========

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
func Int32Value(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
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

// ========== 4. 智能转换（根据值决定是否返回 nil） ==========

// StringPtrNilIfEmpty 空字符串返回 nil，否则返回指针
func StringPtrNilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Int16PtrNilIfZero 零值返回 nil，否则返回指针
func Int16PtrNilIfZero(i int16) *int16 {
	if i == 0 {
		return nil
	}
	return &i
}

// Int32PtrNilIfZero 零值返回 nil，否则返回指针
func Int32PtrNilIfZero(i int32) *int32 {
	if i == 0 {
		return nil
	}
	return &i
}

// Int64PtrNilIfZero 零值返回 nil，否则返回指针
func Int64PtrNilIfZero(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}
