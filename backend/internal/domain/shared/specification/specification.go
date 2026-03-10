package specification

// Specification 规格接口（用于封装业务规则）
type Specification[T any] interface {
	// IsSatisfiedBy 判断候选对象是否满足规格
	IsSatisfiedBy(candidate T) bool
}

// AndSpec 与规格（组合多个规格）
type AndSpec[T any] struct {
	specs []Specification[T]
}

func (s *AndSpec[T]) IsSatisfiedBy(candidate T) bool {
	for _, spec := range s.specs {
		if !spec.IsSatisfiedBy(candidate) {
			return false
		}
	}
	return true
}

// OrSpec 或规格（组合多个规格）
type OrSpec[T any] struct {
	specs []Specification[T]
}

func (s *OrSpec[T]) IsSatisfiedBy(candidate T) bool {
	for _, spec := range s.specs {
		if spec.IsSatisfiedBy(candidate) {
			return true
		}
	}
	return false
}

// NotSpec 非规格（取反）
type NotSpec[T any] struct {
	spec Specification[T]
}

func (s *NotSpec[T]) IsSatisfiedBy(candidate T) bool {
	return !s.spec.IsSatisfiedBy(candidate)
}

// And 组合多个规格为与关系
func And[T any](specs ...Specification[T]) Specification[T] {
	return &AndSpec[T]{specs: specs}
}

// Or 组合多个规格为或关系
func Or[T any](specs ...Specification[T]) Specification[T] {
	return &OrSpec[T]{specs: specs}
}

// Not 对规格取反
func Not[T any](spec Specification[T]) Specification[T] {
	return &NotSpec[T]{spec: spec}
}
