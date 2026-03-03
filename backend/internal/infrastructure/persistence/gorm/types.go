package gorm

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONMap 是用于在PostgreSQL中存储map[string]interface{}作为JSONB的自定义类型
type JSONMap map[string]interface{}

// Scan 实现GORM的Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("扫描JSONMap失败：不支持的类型")
	}

	return json.Unmarshal(bytes, j)
}

// Value 实现GORM的Valuer接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
