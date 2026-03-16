package ddd

import (
	"encoding/json"
	"fmt"
)

// ValueObject 值对象接口
type ValueObject interface {
	Equals(other ValueObject) bool
}

// Identity 值对象标识接口
type Identity interface {
	ValueObject
	String() string
}

// StringIdentity 字符串标识值对象
type StringIdentity struct {
	value string
}

// NewStringIdentity 创建字符串标识
func NewStringIdentity(value string) StringIdentity {
	return StringIdentity{value: value}
}

// String 返回标识字符串值
func (si StringIdentity) String() string {
	return si.value
}

// Equals 比较两个标识是否相等
func (si StringIdentity) Equals(other ValueObject) bool {
	if otherSi, ok := other.(StringIdentity); ok {
		return si.value == otherSi.value
	}
	return false
}

// Int64Identity 64 位整数标识值对象
type Int64Identity struct {
	value int64
}

// NewInt64Identity 创建 64 位整数标识
func NewInt64Identity(value int64) Int64Identity {
	return Int64Identity{value: value}
}

// Int64 返回标识的 64 位整数值
func (ii Int64Identity) Int64() int64 {
	return ii.value
}

// String 返回标识的字符串表示
func (ii Int64Identity) String() string {
	return fmt.Sprintf("%d", ii.value)
}

// MarshalJSON 实现 json.Marshaler 接口
func (ii Int64Identity) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", ii.value)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (ii *Int64Identity) UnmarshalJSON(data []byte) error {
	var value int64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	ii.value = value
	return nil
}

// Equals 比较两个标识是否相等
func (ii Int64Identity) Equals(other ValueObject) bool {
	if otherIi, ok := other.(Int64Identity); ok {
		return ii.value == otherIi.value
	}
	return false
}

// UInt64Identity 无符号 64 位整数标识值对象
type UInt64Identity struct {
	value uint64
}

// NewUInt64Identity 创建无符号 64 位整数标识
func NewUInt64Identity(value uint64) UInt64Identity {
	return UInt64Identity{value: value}
}

// UInt64 返回标识的无符号 64 位整数值
func (ui UInt64Identity) UInt64() uint64 {
	return ui.value
}

// String 返回标识的字符串表示
func (ui UInt64Identity) String() string {
	return fmt.Sprintf("%d", ui.value)
}

// MarshalJSON 实现 json.Marshaler 接口
func (ui UInt64Identity) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", ui.value)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (ui *UInt64Identity) UnmarshalJSON(data []byte) error {
	var value uint64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	ui.value = value
	return nil
}

// Equals 比较两个标识是否相等
func (ui UInt64Identity) Equals(other ValueObject) bool {
	if otherUi, ok := other.(UInt64Identity); ok {
		return ui.value == otherUi.value
	}
	return false
}
