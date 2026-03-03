# 错误处理构建器技术参考

## 目录结构规范

```
.qoder/skills/error-handler-builder/
├── SKILL.md              # 主技能定义文件（必需，大写）
├── REFERENCE.md          # 技术参考文档（大写）
├── EXAMPLES.md           # 使用示例文档（大写）
├── QUICKSTART.md         # 快速入门指南
├── config.yaml           # 配置文件
├── scripts/              # 辅助脚本目录
│   └── helper.sh         # 主要辅助脚本
└── templates/            # 模板文件目录
    ├── error_types.tmpl  # 错误类型模板
    ├── handler.tmpl      # 错误处理器模板
    ├── middleware.tmpl   # 中间件模板
    └── logger.tmpl       # 日志器模板
```

## 核心组件详解

### 1. 错误分类体系

#### 错误类型定义
```go
type ErrorCategory string

const (
    BusinessError    ErrorCategory = "BUSINESS"
    SystemError      ErrorCategory = "SYSTEM"
    ValidationError  ErrorCategory = "VALIDATION"
    AuthenticationError ErrorCategory = "AUTHENTICATION"
    AuthorizationError ErrorCategory = "AUTHORIZATION"
    ExternalError    ErrorCategory = "EXTERNAL"
)

type ErrorCode string

const (
    // 业务错误码
    ErrCodeUserNotFound      ErrorCode = "BUSINESS_USER_NOT_FOUND"
    ErrCodeInsufficientFunds ErrorCode = "BUSINESS_INSUFFICIENT_FUNDS"
    
    // 系统错误码
    ErrCodeDatabaseFailure   ErrorCode = "SYSTEM_DATABASE_FAILURE"
    ErrCodeServiceTimeout    ErrorCode = "SYSTEM_SERVICE_TIMEOUT"
    
    // 验证错误码
    ErrCodeInvalidInput      ErrorCode = "VALIDATION_INVALID_INPUT"
    ErrCodeMissingRequired   ErrorCode = "VALIDATION_MISSING_REQUIRED"
)
```

### 2. 统一错误响应结构

#### 标准响应格式
```go
type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Code       string            `json:"code"`
    Message    string            `json:"message"`
    Details    []ErrorDetailItem `json:"details,omitempty"`
    RequestID  string            `json:"request_id,omitempty"`
    Timestamp  time.Time         `json:"timestamp"`
    TraceID    string            `json:"trace_id,omitempty"`
}

type ErrorDetailItem struct {
    Field   string `json:"field,omitempty"`
    Message string `json:"message"`
    Value   string `json:"value,omitempty"`
}
```

### 3. 错误处理中间件

#### 全局错误处理中间件
```go
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录panic错误
                logger.Error("Panic recovered", 
                    zap.Any("error", err),
                    zap.String("stack", string(debug.Stack())),
                )
                
                // 返回500错误
                c.JSON(http.StatusInternalServerError, ErrorResponse{
                    Error: ErrorDetail{
                        Code:      string(ErrCodeInternalServer),
                        Message:   "Internal server error",
                        Timestamp: time.Now(),
                    },
                })
                c.Abort()
            }
        }()
        
        c.Next()
        
        // 处理业务错误
        if len(c.Errors) > 0 {
            handleBusinessErrors(c)
        }
    }
}

func handleBusinessErrors(c *gin.Context) {
    var errors []ErrorDetailItem
    
    for _, err := range c.Errors {
        if bizErr, ok := err.Err.(*BusinessError); ok {
            errors = append(errors, ErrorDetailItem{
                Field:   bizErr.Field,
                Message: bizErr.Message,
                Value:   bizErr.Value,
            })
        }
    }
    
    if len(errors) > 0 {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: ErrorDetail{
                Code:      string(ErrCodeInvalidInput),
                Message:   "Request validation failed",
                Details:   errors,
                Timestamp: time.Now(),
            },
        })
        c.Abort()
    }
}
```

### 4. 日志记录系统

#### 结构化日志记录
```go
type Logger struct {
    logger *zap.Logger
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
    l.logger.Error(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
    l.logger.Warn(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
    l.logger.Info(msg, fields...)
}

// 错误日志记录辅助函数
func LogError(err error, ctx context.Context, additionalFields ...zap.Field) {
    fields := []zap.Field{
        zap.Error(err),
        zap.String("request_id", GetRequestID(ctx)),
        zap.String("trace_id", GetTraceID(ctx)),
    }
    fields = append(fields, additionalFields...)
    
    logger.Error("Error occurred", fields...)
}
```

## 配置管理系统

### 完整配置结构
```yaml
error_handling:
  # 响应格式配置
  response_format:
    standard: true
    include_stack_trace: false
    include_error_id: true
    mask_sensitive_data: true
    sensitive_field_patterns:
      - "password"
      - "secret"
      - "token"
      - "key"
  
  # 错误分类配置
  error_categories:
    business_errors:
      prefix: "BUSINESS_"
      log_level: "warn"
      http_status: 400
      
    system_errors:
      prefix: "SYSTEM_"
      log_level: "error"
      http_status: 500
      
    validation_errors:
      prefix: "VALIDATION_"
      log_level: "info"
      http_status: 400
      
    authentication_errors:
      prefix: "AUTH_"
      log_level: "warn"
      http_status: 401
      
    authorization_errors:
      prefix: "PERMISSION_"
      log_level: "warn"
      http_status: 403

  # 日志配置
  logging:
    enabled: true
    format: "json"
    levels: ["debug", "info", "warn", "error"]
    include_context: true
    context_fields:
      - "request_id"
      - "user_id"
      - "trace_id"
      - "span_id"

# 监控集成配置
monitoring:
  prometheus:
    enabled: true
    metrics_prefix: "mathfun_"
    error_counter_name: "errors_total"
    error_histogram_name: "error_duration_seconds"
  
  alerting:
    enabled: true
    rules:
      - name: "high_error_rate"
        condition: "rate(errors_total[5m]) > 10"
        severity: "warning"
        description: "High error rate detected"
      
      - name: "critical_system_error"
        condition: "errors_total{category='SYSTEM'} > 0"
        severity: "critical"
        description: "Critical system error occurred"

# 多语言支持
localization:
  enabled: true
  default_language: "zh"
  supported_languages: ["zh", "en"]
  message_sources:
    zh: "locales/zh/errors.yaml"
    en: "locales/en/errors.yaml"
```

## API参考

### 核心命令

#### /error-build-handler
生成错误处理代码
```
/error-build-handler [--config {config_file}] [--output {path}]
```

#### /error-list-codes
列出所有错误码
```
/error-list-codes [--category {category}] [--filter {pattern}]
```

#### /error-validate-config
验证错误处理配置
```
/error-validate-config [--config {config_file}]
```

#### /error-generate-tests
生成错误处理测试代码
```
/error-generate-tests [--coverage-target {percentage}]
```

## 生成代码结构

### 1. 错误类型定义
```go
// errors/types.go
package errors

import "fmt"

type BusinessError struct {
    Code    ErrorCode
    Message string
    Field   string
    Value   string
}

func (e *BusinessError) Error() string {
    if e.Field != "" {
        return fmt.Sprintf("%s: %s (field: %s)", e.Code, e.Message, e.Field)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

type SystemError struct {
    Code    ErrorCode
    Message string
    Err     error
}

func (e *SystemError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *SystemError) Unwrap() error {
    return e.Err
}
```

### 2. 错误工厂函数
```go
// errors/factory.go
package errors

func NewUserNotFoundError(userID string) *BusinessError {
    return &BusinessError{
        Code:    ErrCodeUserNotFound,
        Message: "用户不存在",
        Field:   "user_id",
        Value:   userID,
    }
}

func NewDatabaseError(operation string, err error) *SystemError {
    return &SystemError{
        Code:    ErrCodeDatabaseFailure,
        Message: fmt.Sprintf("数据库操作失败: %s", operation),
        Err:     err,
    }
}
```

### 3. 全局错误处理器
```go
// middleware/error_handler.go
package middleware

func GlobalErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // 处理累积的错误
        if len(c.Errors) > 0 {
            processErrors(c)
        }
    }
}

func processErrors(c *gin.Context) {
    var errorDetails []ErrorDetailItem
    
    for _, ginErr := range c.Errors {
        switch err := ginErr.Err.(type) {
        case *BusinessError:
            errorDetails = append(errorDetails, ErrorDetailItem{
                Field:   err.Field,
                Message: err.Message,
                Value:   err.Value,
            })
        case *SystemError:
            logger.Error("System error", 
                zap.String("code", string(err.Code)),
                zap.Error(err.Err),
            )
        }
    }
    
    if len(errorDetails) > 0 {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: ErrorDetail{
                Code:      string(ErrCodeInvalidInput),
                Message:   "请求验证失败",
                Details:   errorDetails,
                RequestID: getRequestID(c),
                Timestamp: time.Now(),
            },
        })
    }
}
```

## 最佳实践

### 错误处理原则
1. **早发现早处理** - 在最接近问题源头的地方处理错误
2. **错误不丢失** - 确保所有错误都被适当记录和处理
3. **用户友好** - 向用户显示清晰、有用的错误信息
4. **安全第一** - 不在错误响应中泄露敏感信息
5. **可追溯性** - 通过request_id等标识关联错误日志

### 性能优化
- 错误对象池化减少GC压力
- 异步日志记录避免阻塞
- 错误统计采样减少监控开销
- 缓存频繁的错误消息

---
*本文档遵循Qoder Skills技术规范，定期更新以反映最新功能和最佳实践*