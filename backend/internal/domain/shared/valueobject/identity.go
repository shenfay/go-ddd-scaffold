// Package valueobject 提供领域驱动设计中的值对象实现
package valueobject

import (
	"encoding/json"

	"github.com/google/uuid"
)

// UserID 用户 ID 值对象
type UserID struct {
	value uuid.UUID
}

// NewUserID 创建新的 UserID 值对象
func NewUserID(id uuid.UUID) UserID {
	return UserID{value: id}
}

// ParseUserID 从字符串解析 UserID
func ParseUserID(idStr string) (UserID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return UserID{}, err
	}
	return UserID{value: id}, nil
}

// Value 获取底层 UUID 值
func (id UserID) Value() uuid.UUID {
	return id.value
}

// String 转换为字符串
func (id UserID) String() string {
	return id.value.String()
}

// Equals 判断是否相等
func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}

// MarshalJSON 实现 json.Marshaler 接口
func (id UserID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.value.String())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (id *UserID) UnmarshalJSON(data []byte) error {
	var idStr string
	if err := json.Unmarshal(data, &idStr); err != nil {
		return err
	}
	uuid, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	id.value = uuid
	return nil
}

// TenantID 租户 ID 值对象
type TenantID struct {
	value uuid.UUID
}

// NewTenantID 创建新的 TenantID 值对象
func NewTenantID(id uuid.UUID) TenantID {
	return TenantID{value: id}
}

// ParseTenantID 从字符串解析 TenantID
func ParseTenantID(idStr string) (TenantID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return TenantID{}, err
	}
	return TenantID{value: id}, nil
}

// Value 获取底层 UUID 值
func (id TenantID) Value() uuid.UUID {
	return id.value
}

// String 转换为字符串
func (id TenantID) String() string {
	return id.value.String()
}

// Equals 判断是否相等
func (id TenantID) Equals(other TenantID) bool {
	return id.value == other.value
}

// MarshalJSON 实现 json.Marshaler 接口
func (id TenantID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.value.String())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (id *TenantID) UnmarshalJSON(data []byte) error {
	var idStr string
	if err := json.Unmarshal(data, &idStr); err != nil {
		return err
	}
	uuid, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	id.value = uuid
	return nil
}
