# 工具包使用指南

本文档介绍项目中提供的公共工具包及其使用方法。

## 概述

项目提供了以下公共包，位于 `pkg/` 目录下：

- **pkg/response/** - HTTP 响应公共包（结构体定义 + 构造函数）
- **pkg/util/cast.go** - 基于 `github.com/spf13/cast` 的类型转换工具集
- **pkg/util/time.go** - 基于 `github.com/dromara/carbon/v2` 的时间处理工具集

## HTTP 响应包 (pkg/response)

### 导入方式
```go
import "github.com/shenfay/go-ddd-scaffold/pkg/response"
```

### 响应结构体

#### Response - 成功响应
```go
type Response struct {
    Code      int         `json:"code"`               // 业务错误码
    Message   string      `json:"message"`            // 错误消息
    Data      interface{} `json:"data,omitempty"`     // 响应数据
    Details   interface{} `json:"details,omitempty"`  // 详细错误信息
    TraceID   string      `json:"trace_id,omitempty"` // 请求追踪 ID
    Timestamp int64       `json:"timestamp"`          // 时间戳
}
```

#### ErrorResponse - 错误响应
```go
type ErrorResponse struct {
    Code      int         `json:"code"`               // 错误码
    Message   string      `json:"message"`            // 错误消息
    Details   interface{} `json:"details,omitempty"`  // 详细错误信息
    TraceID   string      `json:"trace_id,omitempty"` // 请求追踪 ID
    Timestamp int64       `json:"timestamp"`          // 时间戳
}
```

#### PageData - 分页数据
```go
type PageData struct {
    Items     interface{} `json:"items"`      // 数据列表
    Total     int64       `json:"total"`      // 总数
    Page      int         `json:"page"`       // 当前页码
    PageSize  int         `json:"page_size"`  // 每页数量
    TotalPage int         `json:"total_page"` // 总页数
}
```

### 构造函数

#### NewResponse - 创建成功响应
```go
resp := response.NewResponse(data)
// 返回：&Response{Code: 0, Message: "success", Data: data, Timestamp: now}
```

#### NewErrorResponse - 创建错误响应
```go
errResp := response.NewErrorResponse(code, message, details)
// 返回：&ErrorResponse{Code: code, Message: message, Details: details, Timestamp: now}
```

#### NewPageResponse - 创建分页响应
```go
pageResp := response.NewPageResponse(items, total, page, pageSize)
// 自动计算 TotalPage
// 返回：&Response{Data: PageData{...}}
```

### 链式调用方法

```go
// 添加追踪 ID
resp := response.NewResponse(data).
    WithTraceID(traceID).
    WithMessage("custom message")

// 添加详细信息
errResp := response.NewErrorResponse(code, msg, nil).
    WithDetails(map[string]interface{}{
        "field": "email",
        "rule": "required",
    })
```

### 实际应用场景

#### 1. HTTP Handler 中的使用

```go
func GetUserHandler(c *gin.Context) {
    userID := util.ToInt64(c.Param("id"))
    
    user, err := userService.GetUser(userID)
    if err != nil {
        // 错误响应
        c.JSON(http.StatusOK, response.NewErrorResponse(
            kernel.CodeUserNotFound,
            "user not found",
            nil,
        ))
        return
    }
    
    // 成功响应
    c.JSON(http.StatusOK, response.NewResponse(user).WithTraceID(traceID))
}
```

#### 2. 分页查询响应

```go
func ListUsersHandler(c *gin.Context) {
    page := util.ToInt(c.Query("page"))
    pageSize := util.ToInt(c.Query("page_size"))
    
    users, total, err := userService.ListUsers(page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, 
            response.NewErrorResponse(kernel.CodeInternalError, err.Error(), nil))
        return
    }
    
    c.JSON(http.StatusOK, response.NewPageResponse(users, total, page, pageSize))
}
```

#### 3. 与 ErrorMapper 配合使用

```go
// internal/domain/shared/kernel/response.go
mapper := kernel.NewErrorMapper()
httpStatus, code, message, details := mapper.Map(err)

// pkg/response/response.go
c.JSON(httpStatus, response.NewErrorResponse(code, message, details))
```

---

## 类型转换工具 (util 包)

### 基本类型转换

```go
// 转换为字符串
str := util.ToString(123)           // "123"
str := util.ToString(12.34)         // "12.34"
str := util.ToString(true)          // "true"

// 转换为整数
num := util.ToInt("123")            // 123
num := util.ToInt64("1234567890")   // 1234567890
num := util.ToUint("100")           // 100
num := util.ToUint64("100")         // 100

// 转换为浮点数
f := util.ToFloat64("12.34")        // 12.34

// 转换为布尔值
b := util.ToBool("true")            // true
b := util.ToBool("false")           // false

// 转换为时间
t := util.ToTime("2024-01-01")      // time.Time
```

### 实际应用场景

#### 1. HTTP 参数解析

```go
// 旧代码
userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
if err != nil {
    return err
}

// 新代码
userID := util.ToInt64(c.Param("id"))
if userID == 0 {
    return errors.New("invalid user id")
}
```

#### 2. SQL 参数构建

```go
// 旧代码
func itoa(i int) string {
    if i == 0 {
        return "0"
    }
    result := ""
    for i > 0 {
        result = string('0'+i%10) + result
        i /= 10
    }
    return result
}

query += " LIMIT $" + itoa(argPos)

// 新代码
query += " LIMIT $" + util.ToString(argPos)
```

## 时间处理工具 (util 包)

### 创建时间

```go
// 当前时间
now := util.Now()

// 昨天/明天
yesterday := util.Yesterday()
tomorrow := util.Tomorrow()

// 解析时间字符串
carb := util.Parse("2024-01-01 12:00:00")
carb := util.ParseByFormat("2024/01/01", "2006/01/02")
carb := util.ParseByLayout("2024-01-01", "2006-01-02")

// 从时间戳创建（秒级）
carb := util.FromTimestamp(1704067200)

// 从时间戳创建（毫秒级）
carb := util.FromTimestampMilli(1704067200000)

// 从日期时间创建
carb := util.FromDateTime(2024, 1, 1, 12, 0, 0)
carb := util.FromDate(2024, 1, 1)
```

### 格式化时间

```go
carb := util.Now()

// 获取 Unix 时间戳（秒级）
ts := carb.Timestamp()              // 1704067200

// 获取 Unix 时间戳（毫秒级）
tsMilli := carb.TimestampMilli()    // 1704067200000

// 格式化输出
str := util.Format(carb, "2006-01-02 15:04:05")
str := util.ToDateTimeString(carb)   // "2024-01-01 12:00:00"
str := util.ToDateString(carb)       // "2024-01-01"
str := util.ToTimeString(carb)       // "12:00:00"
```

### 时间计算

```go
carb := util.Now()

// 增加时间
future := util.AddSeconds(carb, 60)
future := util.AddMinutes(carb, 30)
future := util.AddHours(carb, 2)
future := util.AddDays(carb, 7)
future := util.AddMonths(carb, 1)
future := util.AddYears(carb, 1)

// 减少时间
past := util.SubSeconds(carb, 60)
past := util.SubMinutes(carb, 30)
past := util.SubHours(carb, 2)
past := util.SubDays(carb, 7)
```

### 时间比较

```go
carb1 := util.Now()
carb2 := util.Yesterday()

// 计算差异
secs := util.DiffInSeconds(carb1, carb2)
mins := util.DiffInMinutes(carb1, carb2)
hours := util.DiffInHours(carb1, carb2)
days := util.DiffInDays(carb1, carb2)

// 判断
isToday := util.IsToday(carb1)        // true
isYesterday := util.IsYesterday(carb2) // true
isFuture := util.IsFuture(future)     // true
isPast := util.IsPast(past)           // true
```

### 特殊时间点

```go
carb := util.Now()

// 一天的开始和结束
startOfDay := util.StartOfDay(carb)    // 2024-01-01 00:00:00
endOfDay := util.EndOfDay(carb)        // 2024-01-01 23:59:59

// 一月的开始和结束
startOfMonth := util.StartOfMonth(carb)
endOfMonth := util.EndOfMonth(carb)

// 一年的开始和结束
startOfYear := util.StartOfYear(carb)
endOfYear := util.EndOfYear(carb)
```

### 实际应用场景

#### 1. 替换 time.Now()

```go
// 旧代码
createdAt := time.Now()

// 新代码
createdAt := util.Now()
```

#### 2. 替换 time.Now().Unix()

```go
// 旧代码
timestamp := time.Now().Unix()

// 新代码
timestamp := util.Now().Timestamp()
```

#### 3. 替换 time.Now().UnixMilli()

```go
// 旧代码
now := time.Now().UnixMilli()

// 新代码
now := util.Now().TimestampMilli()
```

#### 4. 时间格式化

```go
// 旧代码
today := time.Now().Format("2006-01-02")

// 新代码
today := util.Now().ToDateString()
```

#### 5. 时间计算

```go
// 旧代码
lockUntil := time.Now().Add(24 * time.Hour)

// 新代码
lockUntil := util.AddHours(util.Now(), 24)
```

## 最佳实践

### 1. 优先使用工具函数

- ✅ 使用 `util.ToInt64()` 而不是 `strconv.ParseInt()`
- ✅ 使用 `util.Now()` 而不是 `time.Now()`
- ✅ 使用 `util.Now().Timestamp()` 而不是 `time.Now().Unix()`

### 2. 错误处理

```go
// 类型转换失败时，cast 包会返回零值
userID := util.ToInt64(c.Param("id"))
if userID == 0 {
    // 处理无效 ID
    return errors.New("invalid user id")
}
```

### 3. 性能考虑

- `cast` 包内部做了优化，性能优于手动实现
- `carbon` 包提供了丰富的链式调用，代码更简洁

### 4. 代码一致性

在整个项目中保持使用统一的工具函数：

```go
// ❌ 不推荐：混用不同方式
userID, _ := strconv.ParseInt(id, 10, 64)
now := time.Now()

// ✅ 推荐：统一使用工具函数
userID := util.ToInt64(id)
now := util.Now()
```

## 迁移指南

### 逐步迁移策略

1. **基础设施层优先**: 先迁移 `infrastructure` 层的代码
2. **接口层跟进**: 然后迁移 `interfaces` 层的代码
3. **应用层可选**: 应用层可以根据需要选择性迁移
4. **领域层保持**: 领域层建议保持使用标准库，避免外部依赖

### 已迁移的文件

- ✅ `internal/infrastructure/persistence/user_projector.go`
- ✅ `internal/interfaces/http/user/handler.go`
- ✅ `internal/interfaces/http/middleware/error.go`
- ✅ `internal/interfaces/http/router.go`
- ✅ `internal/infrastructure/snowflake/snowflake.go`

## 依赖版本

```toml
github.com/spf13/cast v1.10.0
github.com/dromara/carbon/v2 v2.6.16
```

## 总结

通过使用 `cast` 和 `carbon` 这两个成熟的第三方库，我们可以：

1. **简化代码**: 减少样板代码，提高可读性
2. **提高可靠性**: 使用经过广泛测试的成熟库
3. **统一风格**: 全项目使用一致的工具类
4. **便于维护**: 集中管理类型转换和时间处理逻辑
