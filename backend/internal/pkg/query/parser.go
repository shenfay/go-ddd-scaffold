package query

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ParseQueryParams 从gin上下文中解析查询参数到query.QueryParams
func ParseQueryParams(c *gin.Context) *QueryParams {
	params := &QueryParams{
		Filter: make(FilterParams),
	}

	// 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err == nil && page >= 1 {
		params.Page = page
	} else {
		params.Page = 1 // 默认值
	}

	pageSizeStr := c.DefaultQuery("pageSize", "20")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err == nil && pageSize >= 1 && pageSize <= 100 { // 每页最多100条
		params.Size = pageSize
	} else {
		params.Size = 20 // 默认值
	}

	// 解析排序参数
	sortStr := c.Query("sort")
	if sortStr != "" {
		sortFields := strings.Split(sortStr, ",")
		for _, field := range sortFields {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}
			dir := SortAsc
			if strings.HasPrefix(field, "-") {
				dir = SortDesc
				field = strings.TrimPrefix(field, "-")
			} else if strings.HasPrefix(field, "+") {
				field = strings.TrimPrefix(field, "+")
			}
			params.Sort = append(params.Sort, SortParam{
				Field:     field,
				Direction: dir,
			})
		}
	}

	// 解析过滤器参数
	for key, values := range c.Request.URL.Query() {
		if !strings.HasPrefix(key, "filter[") || !strings.HasSuffix(key, "]") {
			continue // 不是过滤器参数
		}

		fieldName := strings.TrimPrefix(key, "filter[")
		fieldName = strings.TrimSuffix(fieldName, "]")

		for _, valueStr := range values {
			valueStr = strings.TrimSpace(valueStr)
			if valueStr == "" {
				continue
			}

			// 检查操作符如gt:, lt:, in:等
			operator := FilterEqual // 默认值
			opValue := valueStr       // 默认值

			switch {
			case strings.HasPrefix(valueStr, "gt:"):
				operator = FilterGreaterThan
				opValue = strings.TrimPrefix(valueStr, "gt:")
			case strings.HasPrefix(valueStr, ">"):
				operator = FilterGreaterThan
				opValue = strings.TrimPrefix(valueStr, ">")
			case strings.HasPrefix(valueStr, "lt:"):
				operator = FilterLessThan
				opValue = strings.TrimPrefix(valueStr, "lt:")
			case strings.HasPrefix(valueStr, "<"):
				operator = FilterLessThan
				opValue = strings.TrimPrefix(valueStr, "<")
			case strings.HasPrefix(valueStr, "gte:"):
				operator = FilterGreaterThanOrEqual
				opValue = strings.TrimPrefix(valueStr, "gte:")
			case strings.HasPrefix(valueStr, ">="):
				operator = FilterGreaterThanOrEqual
				opValue = strings.TrimPrefix(valueStr, ">=")
			case strings.HasPrefix(valueStr, "lte:"):
				operator = FilterLessThanOrEqual
				opValue = strings.TrimPrefix(valueStr, "lte:")
			case strings.HasPrefix(valueStr, "<="):
				operator = FilterLessThanOrEqual
				opValue = strings.TrimPrefix(valueStr, "<=")
			case strings.HasPrefix(valueStr, "not:"):
				operator = FilterNotEqual
				opValue = strings.TrimPrefix(valueStr, "not:")
			case strings.HasPrefix(valueStr, "!="):
				operator = FilterNotEqual
				opValue = strings.TrimPrefix(valueStr, "!=")
			case strings.HasPrefix(valueStr, "like:"):
				operator = FilterLike
				opValue = strings.TrimPrefix(valueStr, "like:")
			case strings.HasPrefix(valueStr, "in:"):
				operator = FilterIn
				opValue = strings.TrimPrefix(valueStr, "in:")
			case strings.HasPrefix(valueStr, "between:"):
				operator = FilterBetween
				opValue = strings.TrimPrefix(valueStr, "between:")
			case strings.HasPrefix(valueStr, "date:"):
				operator = FilterDate
				opValue = strings.TrimPrefix(valueStr, "date:")
			case strings.HasPrefix(valueStr, "date>:"):
				operator = FilterDateGreater
				opValue = strings.TrimPrefix(valueStr, "date>:")
			case strings.HasPrefix(valueStr, "date<:"):
				operator = FilterDateLess
				opValue = strings.TrimPrefix(valueStr, "date<:")
			case strings.HasPrefix(valueStr, "date>=:"):
				operator = FilterDateGreaterOrEqual
				opValue = strings.TrimPrefix(valueStr, "date>=:")
			case strings.HasPrefix(valueStr, "date<=:"):
				operator = FilterDateLessOrEqual
				opValue = strings.TrimPrefix(valueStr, "date<=:")
			default:
				// 如果没有操作符前缀，则为简单相等或用于IN的逗号分隔列表
				if strings.Contains(valueStr, ",") {
					operator = FilterIn
					individualValues := strings.Split(valueStr, ",")
					for _, v := range individualValues {
						v = strings.TrimSpace(v)
						if v != "" {
							params.Filter[fieldName] = append(params.Filter[fieldName], FilterCondition{
								Operator: FilterIn, // 对于逗号分隔，使用IN
								Value:    v,
							})
						}
					}
					continue // 移动到外层循环的下一个值，因为我们已将此值处理为多个条件
				}
				// 如果执行到这里，说明是简单相等
				opValue = valueStr
				operator = FilterEqual
			}

			// 处理暗示多个值的操作符的逗号分隔值（如IN、带范围的BETWEEN）
			finalOpValue := opValue
			if operator == FilterIn || operator == FilterBetween {
				// 按逗号分割IN/BETWEEN的值，但仅当上面未处理隐式IN时
				if operator == FilterIn && !strings.Contains(valueStr, ",") {
					// 此情况用于显式的'in:value1,value2'
					individualValues := strings.Split(opValue, ",")
					for _, v := range individualValues {
						v = strings.TrimSpace(v)
						if v != "" {
							params.Filter[fieldName] = append(params.Filter[fieldName], FilterCondition{
								Operator: FilterIn, // 对每个分割的值使用IN操作符
								Value:    v,
							})
						}
					}
					continue // 移动到外层循环的下一个值
				} else if operator == FilterBetween {
					// 期望"value1,value2"
					rangeValues := strings.Split(opValue, ",")
					if len(rangeValues) == 2 {
						finalOpValue = opValue // e.g., "3000,5000"
					}
					// 如果len != 2，忽略还是记录错误？目前按原样存储
				}
			}

			// 添加解析后的条件
			if fieldName != "" {
				params.Filter[fieldName] = append(params.Filter[fieldName], FilterCondition{
					Operator: operator,
					Value:    finalOpValue,
				})
			}
		}
	}

	return params
}