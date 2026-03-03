// Package converter 提供转换器工厂
package converter

// NewConverter 创建一个新的类型转换器实例
func NewConverter() Converter {
	return NewDefaultConverter()
}
