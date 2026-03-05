// Package cast 提供基础类型转换函数
package cast

import "github.com/google/uuid"

// ToUUID 将字符串转换为 UUID
func ToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// ToUUIDPtr 将字符串转换为 UUID 指针（空字符串返回 nil）
func ToUUIDPtr(s string) (*uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ToStringPtr 将字符串转换为字符串指针（空字符串返回 nil）
func ToStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
