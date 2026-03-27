package common

// ValueObject 值对象接口
type ValueObject interface {
	Equals(other ValueObject) bool
}
