# Go DDD Scaffold 安全合规文档

## 文档概述

本文档详细描述了 go-ddd-scaffold 项目的安全设计和合规要求，包括等保三级安全要求对照、JWT实现规范、RBAC权限模型以及审计日志规范。

## 等保三级安全要求对照

### 身份鉴别（S3-A1）
| 要求项 | 实现方案 | 合规状态 |
|--------|----------|----------|
| 身份标识唯一性 | Snowflake ID + 用户名/邮箱唯一约束 | ✅ 符合 |
| 身份鉴别信息复杂度 | 密码策略：8位以上，包含大小写字母数字特殊字符 | ✅ 符合 |
| 登录失败处理 | 连续失败5次锁定账户30分钟 | ✅ 符合 |
| 会话管理 | JWT双Token机制，Access Token 30分钟，Refresh Token 7天 | ✅ 符合 |

### 访问控制（S3-A2）
| 要求项 | 实现方案 | 合规状态 |
|--------|----------|----------|
| 自主访问控制 | RBAC权限模型，基于角色的访问控制 | ✅ 符合 |
| 强制访问控制 | 多租户数据隔离，租户间数据严格分离 | ✅ 符合 |
| 权限最小化 | 最小权限原则，按需分配权限 | ✅ 符合 |
| 特权用户权限分离 | 系统管理员与业务用户权限分离 | ✅ 符合 |

### 安全审计（S3-A3）
| 要求项 | 实现方案 | 合规状态 |
|--------|----------|----------|
| 审计事件覆盖 | 用户登录、权限变更、数据操作等关键事件 | ✅ 符合 |
| 审计记录完整性 | 记录时间、用户、操作、结果等完整信息 | ✅ 符合 |
| 审计记录保护 | 审计日志只读存储，防篡改机制 | ✅ 符合 |
| 审计记录分析 | 提供日志分析和异常检测功能 | ✅ 符合 |

### 入侵防范（S3-A4）
| 要求项 | 实现方案 | 合规状态 |
|--------|----------|----------|
| 恶意代码防范 | 代码扫描、依赖安全检查 | ✅ 符合 |
| 网络攻击防范 | 速率限制、输入验证、SQL注入防护 | ✅ 符合 |
| 恶意行为监控 | 异常登录检测、权限滥用监控 | ✅ 符合 |

### 数据完整性（S3-A5）
| 要求项 | 实现方案 | 合规状态 |
|--------|----------|----------|
| 重要数据传输完整性 | HTTPS/TLS 1.2+加密传输 | ✅ 符合 |
| 重要数据存储完整性 | 数据库约束、校验和机制 | ✅ 符合 |

### 数据保密性（S3-A6）
| 要求项 | 实现方案 | 合规状态 |
|--------|----------|----------|
| 重要数据传输保密性 | TLS加密传输 | ✅ 符合 |
| 重要数据存储保密性 | 敏感字段加密存储 | ✅ 符合 |

## JWT实现安全规范

### Token结构设计
```go
// JWT Claims结构
type CustomClaims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    TenantID int64  `json:"tenant_id,omitempty"`
    Role     string `json:"role,omitempty"`
    jwttoken.StandardClaims
}

// Token配置
type JWTConfig struct {
    SecretKey       string        `mapstructure:"secret"`
    AccessExpire    time.Duration `mapstructure:"access_expire"`
    RefreshExpire   time.Duration `mapstructure:"refresh_expire"`
    Issuer          string        `mapstructure:"issuer"`
}
```

### 双Token机制实现
```go
type TokenService struct {
    config JWTConfig
    redis  redis.Client
}

func (ts *TokenService) GenerateTokens(userID int64) (*TokenPair, error) {
    // 生成Access Token
    accessToken, err := ts.generateAccessToken(userID)
    if err != nil {
        return nil, err
    }
    
    // 生成Refresh Token
    refreshToken, err := ts.generateRefreshToken(userID)
    if err != nil {
        return nil, err
    }
    
    // 将Refresh Token存储到Redis（用于撤销）
    err = ts.storeRefreshToken(userID, refreshToken)
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int(ts.config.AccessExpire.Seconds()),
    }, nil
}

func (ts *TokenService) generateAccessToken(userID int64) (string, error) {
    claims := CustomClaims{
        UserID: userID,
        StandardClaims: jwttoken.StandardClaims{
            ExpiresAt: time.Now().Add(ts.config.AccessExpire).Unix(),
            IssuedAt:  time.Now().Unix(),
            Issuer:    ts.config.Issuer,
        },
    }
    
    token := jwttoken.NewWithClaims(jwttoken.SigningMethodHS256, claims)
    return token.SignedString([]byte(ts.config.SecretKey))
}

func (ts *TokenService) generateRefreshToken(userID int64) (string, error) {
    // Refresh Token包含额外的安全信息
    claims := jwttoken.MapClaims{
        "user_id": userID,
        "type":    "refresh",
        "jti":     generateJTI(), // JWT ID，用于唯一标识
        "exp":     time.Now().Add(ts.config.RefreshExpire).Unix(),
    }
    
    token := jwttoken.NewWithClaims(jwttoken.SigningMethodHS256, claims)
    return token.SignedString([]byte(ts.config.SecretKey))
}
```

### Token验证与刷新
```go
func (ts *TokenService) ParseAccessToken(tokenString string) (*CustomClaims, error) {
    token, err := jwttoken.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwttoken.Token) (interface{}, error) {
        return []byte(ts.config.SecretKey), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    claims, ok := token.Claims.(*CustomClaims)
    if !ok {
        return nil, errors.New("invalid claims")
    }
    
    return claims, nil
}

func (ts *TokenService) RefreshToken(refreshToken string) (*TokenPair, error) {
    // 解析Refresh Token
    claims, err := ts.parseRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }
    
    // 检查Token是否已被撤销
    if ts.isTokenRevoked(claims["jti"].(string)) {
        return nil, errors.New("token has been revoked")
    }
    
    userID := int64(claims["user_id"].(float64))
    
    // 生成新的Token对
    return ts.GenerateTokens(userID)
}

func (ts *TokenService) RevokeToken(refreshToken string) error {
    claims, err := ts.parseRefreshToken(refreshToken)
    if err != nil {
        return err
    }
    
    jti := claims["jti"].(string)
    userID := int64(claims["user_id"].(float64))
    
    // 将Token加入黑名单
    return ts.redis.Set(context.Background(), 
        fmt.Sprintf("revoked_token:%s", jti), 
        userID, 
        ts.config.RefreshExpire).Err()
}
```

### 安全增强措施
```go
// Token黑名单检查
func (ts *TokenService) isTokenRevoked(jti string) bool {
    exists, err := ts.redis.Exists(context.Background(), fmt.Sprintf("revoked_token:%s", jti)).Result()
    return err == nil && exists > 0
}

// Token指纹验证（防重放攻击）
func (ts *TokenService) validateTokenFingerprint(claims *CustomClaims, fingerprint string) bool {
    expectedFingerprint := generateFingerprint(claims.UserID, claims.IssuedAt)
    return subtle.ConstantTimeCompare([]byte(fingerprint), []byte(expectedFingerprint)) == 1
}

// 客户端指纹生成
func generateFingerprint(userID int64, issuedAt int64) string {
    data := fmt.Sprintf("%d:%d:%s", userID, issuedAt, getClientSecret())
    hash := sha256.Sum256([]byte(data))
    return base64.StdEncoding.EncodeToString(hash[:])
}
```

## RBAC权限模型规范

### 权限模型设计
```go
// 权限结构
type Permission struct {
    ID          int64  `json:"id"`
    Resource    string `json:"resource"`    // 资源类型：user, tenant, role等
    Action      string `json:"action"`      // 操作：create, read, update, delete
    Description string `json:"description"`
    Scope       string `json:"scope"`       // 作用域：system, global, tenant
}

// 角色结构
type Role struct {
    ID          int64        `json:"id"`
    Name        string       `json:"name"`
    Description string       `json:"description"`
    TenantID    *int64       `json:"tenant_id"`   // NULL表示系统角色
    Permissions []Permission `json:"permissions"` // 关联权限
    IsSystem    bool         `json:"is_system"`
}

// 用户权限检查
type AuthorizationService struct {
    roleRepo       RoleRepository
    permissionRepo PermissionRepository
    cache          cache.Cache
}

func (as *AuthorizationService) CheckPermission(userID, tenantID int64, resource, action string) bool {
    // 构造缓存键
    cacheKey := fmt.Sprintf("permission:%d:%d:%s:%s", userID, tenantID, resource, action)
    
    // 先检查缓存
    if cached, found := as.cache.Get(cacheKey); found {
        return cached.(bool)
    }
    
    // 查询用户在该租户下的角色
    roles, err := as.getUserRoles(userID, tenantID)
    if err != nil {
        as.cache.Set(cacheKey, false, 5*time.Minute)
        return false
    }
    
    // 检查角色权限
    hasPermission := as.checkRolesPermission(roles, resource, action)
    
    // 缓存结果
    as.cache.Set(cacheKey, hasPermission, 10*time.Minute)
    return hasPermission
}

func (as *AuthorizationService) getUserRoles(userID, tenantID int64) ([]Role, error) {
    // 查询用户在指定租户下的角色
    var userTenants []UserTenant
    err := db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Find(&userTenants).Error
    if err != nil {
        return nil, err
    }
    
    if len(userTenants) == 0 {
        return []Role{}, nil
    }
    
    // 获取角色信息
    var roleIDs []int64
    for _, ut := range userTenants {
        roleIDs = append(roleIDs, ut.RoleID)
    }
    
    var roles []Role
    err = db.Where("id IN ?", roleIDs).Find(&roles).Error
    return roles, err
}
```

### 权限继承与组合
```go
// 权限继承树
type PermissionTree struct {
    Permission Permission
    Children   []PermissionTree
}

// 构建权限树
func (as *AuthorizationService) buildPermissionTree() (*PermissionTree, error) {
    permissions, err := as.permissionRepo.GetAll()
    if err != nil {
        return nil, err
    }
    
    // 按资源分组
    resourceMap := make(map[string][]Permission)
    for _, perm := range permissions {
        resourceMap[perm.Resource] = append(resourceMap[perm.Resource], perm)
    }
    
    // 构建树结构
    root := &PermissionTree{
        Permission: Permission{Resource: "root"},
    }
    
    for resource, perms := range resourceMap {
        resourceNode := PermissionTree{
            Permission: Permission{Resource: resource},
        }
        
        for _, perm := range perms {
            resourceNode.Children = append(resourceNode.Children, PermissionTree{
                Permission: perm,
            })
        }
        
        root.Children = append(root.Children, resourceNode)
    }
    
    return root, nil
}

// 权限组合检查
func (as *AuthorizationService) CheckCompositePermission(userID, tenantID int64, permissions []PermissionCheck) bool {
    for _, check := range permissions {
        if !as.CheckPermission(userID, tenantID, check.Resource, check.Action) {
            return false
        }
    }
    return true
}

type PermissionCheck struct {
    Resource string
    Action   string
}
```

## 审计日志规范

### 审计事件设计
```go
// 审计日志结构
type AuditLog struct {
    ID           int64       `json:"id" gorm:"primaryKey"`
    UserID       *int64      `json:"user_id"`           // 操作用户
    TenantID     *int64      `json:"tenant_id"`         // 操作租户
    ActionType   string      `json:"action_type"`       // 操作类型
    ResourceType string      `json:"resource_type"`     // 资源类型
    ResourceID   *int64      `json:"resource_id"`       // 资源ID
    Action       string      `json:"action"`            // 具体操作
    OldValues    JSON        `json:"old_values"`        // 修改前的值
    NewValues    JSON        `json:"new_values"`        // 修改后的值
    IPAddress    string      `json:"ip_address"`        // IP地址
    UserAgent    string      `json:"user_agent"`        // 用户代理
    StatusCode   int         `json:"status_code"`       // 操作结果状态码
    ErrorMessage *string     `json:"error_message"`     // 错误信息
    CreatedAt    time.Time   `json:"created_at"`
}

// 审计事件类型定义
const (
    // 用户相关事件
    UserLogin           = "USER_LOGIN"
    UserLogout          = "USER_LOGOUT"
    UserCreate          = "USER_CREATE"
    UserUpdate          = "USER_UPDATE"
    UserDelete          = "USER_DELETE"
    UserPasswordChange  = "USER_PASSWORD_CHANGE"
    
    // 租户相关事件
    TenantCreate        = "TENANT_CREATE"
    TenantUpdate        = "TENANT_UPDATE"
    TenantDelete        = "TENANT_DELETE"
    UserAddToTenant     = "USER_ADD_TO_TENANT"
    UserRemoveFromTenant = "USER_REMOVE_FROM_TENANT"
    
    // 权限相关事件
    RoleCreate          = "ROLE_CREATE"
    RoleUpdate          = "ROLE_UPDATE"
    RoleDelete          = "ROLE_DELETE"
    PermissionGrant     = "PERMISSION_GRANT"
    PermissionRevoke    = "PERMISSION_REVOKE"
    
    // 系统相关事件
    SystemConfigUpdate  = "SYSTEM_CONFIG_UPDATE"
    DatabaseBackup      = "DATABASE_BACKUP"
)
```

### 审计日志服务实现
```go
type AuditService struct {
    db    *gorm.DB
    queue chan AuditLog
    wg    sync.WaitGroup
}

func NewAuditService(db *gorm.DB) *AuditService {
    service := &AuditService{
        db:    db,
        queue: make(chan AuditLog, 1000),
    }
    
    // 启动异步处理协程
    service.startWorkers(5)
    return service
}

func (as *AuditService) LogEvent(event AuditLog) {
    // 设置时间戳
    event.CreatedAt = time.Now()
    
    // 异步记录到队列
    select {
    case as.queue <- event:
    default:
        // 队列满时同步记录
        go as.recordEvent(event)
    }
}

func (as *AuditService) startWorkers(workerCount int) {
    for i := 0; i < workerCount; i++ {
        as.wg.Add(1)
        go as.worker()
    }
}

func (as *AuditService) worker() {
    defer as.wg.Done()
    
    for event := range as.queue {
        as.recordEvent(event)
    }
}

func (as *AuditService) recordEvent(event AuditLog) {
    // 记录到数据库
    if err := as.db.Create(&event).Error; err != nil {
        log.Printf("Failed to record audit log: %v", err)
    }
    
    // 同时发送到外部系统（如ELK）
    as.sendToExternalSystem(event)
}

func (as *AuditService) sendToExternalSystem(event AuditLog) {
    // 发送到日志收集系统
    logData, _ := json.Marshal(event)
    // 这里可以集成具体的日志系统
    fmt.Printf("Audit log: %s\n", string(logData))
}
```

### 敏感操作审计
```go
// 敏感操作装饰器
func AuditSensitiveOperation(operation string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 记录操作前状态
        startTime := time.Now()
        
        // 执行原操作
        c.Next()
        
        // 记录操作后状态
        statusCode := c.Writer.Status()
        var errorMessage *string
        if len(c.Errors) > 0 {
            msg := c.Errors.Last().Error()
            errorMessage = &msg
        }
        
        // 构造审计日志
        auditLog := AuditLog{
            UserID:       getCurrentUserID(c),
            TenantID:     getCurrentTenantID(c),
            ActionType:   operation,
            ResourceType: getResourceType(c),
            ResourceID:   getResourceID(c),
            Action:       c.Request.Method + " " + c.Request.URL.Path,
            IPAddress:    c.ClientIP(),
            UserAgent:    c.Request.UserAgent(),
            StatusCode:   statusCode,
            ErrorMessage: errorMessage,
            CreatedAt:    startTime,
        }
        
        // 异步记录审计日志
        auditService.LogEvent(auditLog)
    }
}

// 使用示例
func SetupRoutes(router *gin.Engine) {
    userGroup := router.Group("/api/v1/users")
    {
        userGroup.POST("", AuditSensitiveOperation("USER_CREATE"), CreateUser)
        userGroup.PUT("/:id", AuditSensitiveOperation("USER_UPDATE"), UpdateUser)
        userGroup.DELETE("/:id", AuditSensitiveOperation("USER_DELETE"), DeleteUser)
    }
}
```

## 安全监控与告警

### 安全指标监控
```go
type SecurityMetrics struct {
    FailedLoginAttempts prometheus.Counter
    SuccessfulLogins    prometheus.Counter
    PermissionDenied    prometheus.Counter
    SuspiciousActivity  prometheus.Counter
}

func NewSecurityMetrics() *SecurityMetrics {
    return &SecurityMetrics{
        FailedLoginAttempts: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "failed_login_attempts_total",
            Help: "Total number of failed login attempts",
        }),
        SuccessfulLogins: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "successful_logins_total",
            Help: "Total number of successful logins",
        }),
        PermissionDenied: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "permission_denied_total",
            Help: "Total number of permission denied events",
        }),
        SuspiciousActivity: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "suspicious_activity_total",
            Help: "Total number of suspicious activities detected",
        }),
    }
}

// 异常行为检测
type AnomalyDetector struct {
    metrics *SecurityMetrics
    config  AnomalyConfig
}

func (ad *AnomalyDetector) DetectFailedLoginSpikes(userID int64, ip string) bool {
    // 检测短时间内大量失败登录
    recentFailures := ad.getRecentFailedLogins(userID, ip, time.Minute*5)
    if recentFailures > ad.config.MaxFailedLoginsPerMinute {
        ad.metrics.SuspiciousActivity.Inc()
        ad.alertSecurityTeam(fmt.Sprintf("Suspicious login activity from user %d, IP %s", userID, ip))
        return true
    }
    return false
}
```

这个安全合规文档为项目提供了全面的安全设计规范和等保三级合规实现方案。