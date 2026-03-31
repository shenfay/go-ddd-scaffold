package common

// ValueObject 值对象标记接口
// 值对象应该是不可变的，通过属性判断相等性
// 注意：这是一个标记接口，具体的 Equals 方法由各值对象自行实现
type ValueObject interface {
}
