package query

// QueryParams 封装列表查询的排序和过滤参数
type QueryParams struct {
	Sort   []SortParam  `json:"sort,omitempty"`
	Filter FilterParams `json:"filter,omitempty"`
	Page   int          `json:"page,omitempty"`
	Size   int          `json:"size,omitempty"` // 内部从pageSize重命名为pageSize，但查询参数是pageSize
}

// SortParam 表示单个排序条件
type SortParam struct {
	Field     string // 要排序的字段名
	Direction SortDirection
}

// SortDirection 表示排序方向
type SortDirection int

const (
	SortAsc SortDirection = iota
	SortDesc
)

// FilterParams 表示过滤条件集合
// 键是字段名，值是过滤条件
type FilterParams map[string][]FilterCondition

// FilterCondition 表示单个过滤条件
type FilterCondition struct {
	Operator FilterOperator
	Value    string
}

// FilterOperator 表示过滤条件的操作符
type FilterOperator int

const (
	FilterEqual FilterOperator = iota
	FilterGreaterThan
	FilterLessThan
	FilterGreaterThanOrEqual
	FilterLessThanOrEqual
	FilterNotEqual
	FilterLike
	FilterIn
	FilterBetween
	FilterDate
	FilterDateGreater
	FilterDateLess
	FilterDateGreaterOrEqual
	FilterDateLessOrEqual
)
