# 错误处理构建器使用示例

## 基础使用示例

### 示例 1：标准错误处理生成

**场景**：为新项目生成完整的错误处理体系

```
/error-build-handler
```

**生成的文件结构：**
```
internal/
├── pkg/
│   └── errors/
│       ├── types.go          # 错误类型定义
│       ├── factory.go        # 错误工厂函数
│       └── codes.go          # 错误码枚举
└── middleware/
    └── error_handler.go      # 错误处理中间件
```

**生成的代码示例：**

**types.go**
```go
package errors

import (
    "fmt"
    "time"
)

// ErrorCategory 错误分类
type ErrorCategory string

const (
    BusinessError    ErrorCategory = "BUSINESS"
    SystemError      ErrorCategory = "SYSTEM"
    ValidationError  ErrorCategory = "VALIDATION"
    AuthenticationError ErrorCategory = "AUTHENTICATION"
)

// ErrorCode 错误码
type ErrorCode string

const (
    // 业务错误
    ErrCodeUserNotFound      ErrorCode = "BUSINESS_USER_NOT_FOUND"
    ErrCodeResourceNotFound  ErrorCode = "BUSINESS_RESOURCE_NOT_FOUND"
    ErrCodeInsufficientFunds ErrorCode = "BUSINESS_INSUFFICIENT_FUNDS"
    
    // 系统错误
    ErrCodeDatabaseFailure   ErrorCode = "SYSTEM_DATABASE_FAILURE"
    ErrCodeServiceTimeout    ErrorCode = "SYSTEM_SERVICE_TIMEOUT"
    ErrCodeInternalServer    ErrorCode = "SYSTEM_INTERNAL_SERVER"
    
    // 验证错误
    ErrCodeInvalidInput      ErrorCode = "VALIDATION_INVALID_INPUT"
    ErrCodeMissingRequired   ErrorCode = "VALIDATION_MISSING_REQUIRED"
)

// BusinessError 业务错误
type BusinessError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Field   string    `json:"field,omitempty"`
    Value   string    `json:"value,omitempty"`
    Time    time.Time `json:"time"`
}

func (e *BusinessError) Error() string {
    if e.Field != "" {
        return fmt.Sprintf("%s: %s (field: %s, value: %s)", e.Code, e.Message, e.Field, e.Value)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// SystemError 系统错误
type SystemError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Err     error     `json:"-"`
    Time    time.Time `json:"time"`
}

func (e *SystemError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *SystemError) Unwrap() error {
    return e.Err
}
```

**factory.go**
```go
package errors

import (
    "fmt"
    "time"
)

// 用户相关错误
func NewUserNotFoundError(userID string) *BusinessError {
    return &BusinessError{
        Code:    ErrCodeUserNotFound,
        Message: "用户不存在",
        Field:   "user_id",
        Value:   userID,
        Time:    time.Now(),
    }
}

func NewInsufficientFundsError(balance, required float64) *BusinessError {
    return &BusinessError{
        Code:    ErrCodeInsufficientFunds,
        Message: fmt.Sprintf("余额不足，当前余额: %.2f, 需要: %.2f", balance, required),
        Time:    time.Now(),
    }
}

// 系统相关错误
func NewDatabaseError(operation string, err error) *SystemError {
    return &SystemError{
        Code:    ErrCodeDatabaseFailure,
        Message: fmt.Sprintf("数据库操作失败: %s", operation),
        Err:     err,
        Time:    time.Now(),
    }
}

func NewServiceTimeoutError(service string, timeout string) *SystemError {
    return &SystemError{
        Code:    ErrCodeServiceTimeout,
        Message: fmt.Sprintf("服务调用超时: %s (timeout: %s)", service, timeout),
        Time:    time.Now(),
    }
}
```

### 示例 2：带监控的错误处理

**场景**：生成包含Prometheus监控的错误处理

```
/error-build-handler --with-monitoring --with-alerting
```

**生成的监控代码：**
```go
// metrics/error_metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    ErrorCounter = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mathfun_errors_total",
            Help: "Total number of errors by type and category",
        },
        []string{"type", "category", "code"},
    )
    
    ErrorDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "mathfun_error_duration_seconds",
            Help:    "Error handling duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"category"},
    )
)

// RecordError 记录错误指标
func RecordError(category string, code string) {
    ErrorCounter.WithLabelValues("error", category, code).Inc()
}

// RecordErrorDuration 记录错误处理耗时
func RecordErrorDuration(category string, duration float64) {
    ErrorDuration.WithLabelValues(category).Observe(duration)
}
```

## 高级使用示例

### 示例 3：多语言错误消息

**场景**：支持中英文错误消息

```yaml
# error-localization.yaml
localization:
  enabled: true
  default_language: "zh"
  supported_languages: ["zh", "en"]
  messages:
    BUSINESS_USER_NOT_FOUND:
      zh: "用户不存在"
      en: "User not found"
    VALIDATION_INVALID_INPUT:
      zh: "输入参数无效"
      en: "Invalid input parameters"
```

```
/error-build-handler --config error-localization.yaml
```

**生成的多语言支持代码：**
```go
// localization/messages.go
package localization

import "golang.org/x/text/language"

type MessageLookup map[language.Tag]map[string]string

var Messages = MessageLookup{
    language.Chinese: {
        "BUSINESS_USER_NOT_FOUND":      "用户不存在",
        "VALIDATION_INVALID_INPUT":     "输入参数无效",
        "SYSTEM_DATABASE_FAILURE":      "数据库操作失败",
    },
    language.English: {
        "BUSINESS_USER_NOT_FOUND":      "User not found",
        "VALIDATION_INVALID_INPUT":     "Invalid input parameters",
        "SYSTEM_DATABASE_FAILURE":      "Database operation failed",
    },
}

func GetMessage(code string, lang language.Tag) string {
    if msgs, ok := Messages[lang]; ok {
        if msg, exists := msgs[code]; exists {
            return msg
        }
    }
    // 回退到默认语言
    if msgs, ok := Messages[language.Chinese]; ok {
        if msg, exists := msgs[code]; exists {
            return msg
        }
    }
    return code
}
```

### 示例 4：与API端点生成器集成

**场景**：为生成的API端点自动集成错误处理

```
# 1. 生成错误处理代码
/error-build-handler --output ./backend

# 2. 生成API端点（自动集成错误处理）
/api-generate-endpoints --aggregate UserAggregate --integrate error-handler
```

**自动生成的集成代码：**
```go
// api/user/handler.go
package user

import (
    "mathfun/internal/pkg/errors"
    "mathfun/internal/middleware"
)

type UserHandler struct {
    userService *UserService
    errorHandler *middleware.ErrorHandler
}

func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    user, err := h.userService.GetUser(userID)
    if err != nil {
        // 自动使用标准错误处理
        h.errorHandler.HandleError(c, err)
        return
    }
    
    c.JSON(http.StatusOK, UserResponse{User: user})
}
```

## 团队协作示例

### 示例 5：团队错误规范统一

**场景**：建立团队统一的错误处理规范

```yaml
# team-error-standard.yaml
error_handling:
  # 团队统一的错误码前缀
  team_prefix: "MATHFUN_"
  
  # 统一的错误分类
  categories:
    business:
      prefix: "BIZ_"
      log_level: "warn"
    system:
      prefix: "SYS_"
      log_level: "error"
    validation:
      prefix: "VAL_"
      log_level: "info"
  
  # 统一的日志格式
  logging:
    format: "json"
    include_fields: ["request_id", "user_id", "trace_id"]
    redact_sensitive: true
```

```
/error-build-handler --config team-error-standard.yaml --team-mode
```

### 示例 6：错误处理重构

**场景**：重构现有系统的错误处理

```
# 1. 分析现有错误使用情况
/error-analyze-existing-errors --source ./backend

# 2. 生成迁移方案
/error-generate-migration-plan --analysis-result analysis.json

# 3. 执行重构
/error-refactor-existing --migration-plan migration.yaml
```

## 故障排除示例

### 示例 7：处理错误码冲突

**场景**：解决不同模块间错误码命名冲突

```
用户: "/error-build-handler"
Skills: "检测到以下错误码冲突：
1. USER_NOT_FOUND 在 user 和 admin 模块中都定义
2. DATABASE_ERROR 在多个服务中重复定义

建议解决方案：
1. 添加模块前缀：USER_NOT_FOUND → USER_SERVICE_USER_NOT_FOUND
2. 使用统一的系统错误码
3. 建立错误码注册机制"

用户: "采用方案1"
Skills: "正在重新生成带前缀的错误码..."
[生成更新后的错误码]
```

### 示例 8：性能优化配置

**场景**：优化高并发场景下的错误处理性能

```yaml
# performance-optimized.yaml
performance:
  # 错误对象池化
  object_pooling:
    enabled: true
    pool_size: 1000
    
  # 异步日志记录
  async_logging:
    enabled: true
    buffer_size: 10000
    worker_count: 4
    
  # 错误统计采样
  sampling:
    enabled: true
    rate: 0.1  # 10%采样率
```

```
/error-build-handler --config performance-optimized.yaml
```

这些示例展示了错误处理构建器在各种实际场景中的灵活应用，从基础错误处理到复杂的企业级需求都能有效支持。