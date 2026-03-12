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

// ToUint 转换为 uint
func ToUint(v interface{}) uint {
	return cast.ToUint(v)
}

// ToUint64 转换为 uint64
func ToUint64(v interface{}) uint64 {
	return cast.ToUint64(v)
}

// ToFloat32 转换为 float32
func ToFloat32(v interface{}) float32 {
	return cast.ToFloat32(v)
}

// ToFloat64 转换为 float64
func ToFloat64(v interface{}) float64 {
	return cast.ToFloat64(v)
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
