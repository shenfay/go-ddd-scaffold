package kernel

import (
	"encoding/json"
	"fmt"
)

// Identity 标识值对象接口
type Identity interface {
	ValueObject
	String() string
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
