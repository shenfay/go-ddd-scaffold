# 安全规范

本文档定义了 Go DDD Scaffold 项目的安全规范和最佳实践。

## 📋 安全原则

### 核心原则

1. **纵深防御** - 多层安全防护
2. **最小权限** - 只授予必要的权限
3. **默认安全** - 默认配置应该是安全的
4. **不信任输入** - 验证所有外部输入
5. **安全审计** - 记录所有敏感操作

---

## 🔐 认证安全

### JWT 令牌安全

#### 令牌生成

```go
// infrastructure/platform/auth/jwt_service.go
func (s *JWTService) GenerateTokenPair(userID int64, username, email string) (*TokenPair, error) {
    // Access Token
    now := time.Now()
    accessClaims := &TokenClaims{
        UserID:   userID,
        Username: username,
        Email:    email,
        JTI:      generateJTI(),  // 唯一标识符
        IssuedAt: now.Unix(),
        ExpiresAt: now.Add(s.accessExpire).Unix(),
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString([]byte(s.secret))
    if err != nil {
        return nil, fmt.Errorf("generate access token failed: %w", err)
    }
    
    // Refresh Token（更长有效期）
    refreshClaims := &TokenClaims{
        UserID:   userID,
        JTI:      generateJTI(),
        IssuedAt: now.Unix(),
        ExpiresAt: now.Add(s.refreshExpire).Unix(),
    }
    
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString([]byte(s.secret))
    if err != nil {
        return nil, fmt.Errorf("generate refresh token failed: %w", err)
    }
    
    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
        ExpiresAt:    accessClaims.ExpiresAt,
    }, nil
}

// 生成唯一的 JTI
func generateJTI() string {
    return uuid.New().String()
}
```

#### 令牌验证

```go
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
    // 检查是否在黑名单中
    if s.isBlacklisted(tokenString) {
        return nil, kernel.NewBusinessError(
            response.CodeTokenInvalid,
            "令牌已失效",
        )
    }
    
    // 解析令牌
    claims, err := s.ParseAccessToken(tokenString)
    if err != nil {
        return nil, err
    }
    
    // 检查过期时间
    if time.Now().Unix() > claims.ExpiresAt {
        return nil, kernel.NewBusinessError(
            response.CodeTokenExpired,
            "令牌已过期",
        )
    }
    
    return claims, nil
}

func (s *JWTService) ParseAccessToken(tokenString string) (*TokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 验证签名算法
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.secret), nil
    })
    
    if err != nil {
        return nil, kernel.NewBusinessError(
            response.CodeTokenInvalid,
            "无效的令牌",
        )
    }
    
    claims, ok := token.Claims.(*TokenClaims)
    if !ok || !token.Valid {
        return nil, kernel.NewBusinessError(
            response.CodeTokenInvalid,
            "无效的令牌",
        )
    }
    
    return claims, nil
}
```

#### 令牌黑名单（Redis）

```go
func (s *JWTService) BlacklistToken(token string, expiresAt time.Time) error {
    key := fmt.Sprintf("token:blacklist:%s", token)
    ttl := time.Until(expiresAt)
    
    // 如果已经过期，不需要加入黑名单
    if ttl <= 0 {
        return nil
    }
    
    // 加入黑名单，设置 TTL
    err := s.redis.Set(context.Background(), key, "1", ttl).Err()
    if err != nil {
        return fmt.Errorf("blacklist token failed: %w", err)
    }
    
    return nil
}

func (s *JWTService) isBlacklisted(token string) bool {
    key := fmt.Sprintf("token:blacklist:%s", token)
    exists, err := s.redis.Exists(context.Background(), key).Result()
    if err != nil {
        return false
    }
    return exists > 0
}
```

### 密码安全

#### 密码策略

```go
// infrastructure/auth/password_policy.go
type PasswordPolicy struct {
    minLength         int
    requireUppercase  bool
    requireLowercase  bool
    requireNumber     bool
    requireSpecial    bool
}

func NewDefaultPasswordPolicy() *PasswordPolicy {
    return &PasswordPolicy{
        minLength:         8,
        requireUppercase:  true,
        requireLowercase:  true,
        requireNumber:     true,
        requireSpecial:    true,
    }
}

func (p *PasswordPolicy) Validate(password string) error {
    if len(password) < p.minLength {
        return kernel.FieldError("password", 
            fmt.Sprintf("密码长度至少为 %d 个字符", p.minLength), 
            password)
    }
    
    if p.requireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return kernel.FieldError("password", "密码必须包含大写字母", password)
    }
    
    if p.requireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
        return kernel.FieldError("password", "密码必须包含小写字母", password)
    }
    
    if p.requireNumber && !regexp.MustCompile(`[0-9]`).MatchString(password) {
        return kernel.FieldError("password", "密码必须包含数字", password)
    }
    
    if p.requireSpecial && !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
        return kernel.FieldError("password", "密码必须包含特殊字符", password)
    }
    
    return nil
}
```

#### 密码哈希

```go
// domain/user/service/password_hasher.go
type BcryptPasswordHasher struct {
    cost int
}

func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
    if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
        cost = bcrypt.DefaultCost
    }
    return &BcryptPasswordHasher{cost: cost}
}

func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
    if err != nil {
        return "", fmt.Errorf("hash password failed: %w", err)
    }
    return string(bytes), nil
}

func (h *BcryptPasswordHasher) Verify(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

#### 防止暴力破解

```go
// application/auth/service.go
func (s *AuthServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthResult, error) {
    // 检查登录失败次数
    failedAttempts, err := s.getLoginFailedAttempts(cmd.Identifier)
    if err != nil {
        return nil, fmt.Errorf("get failed attempts: %w", err)
    }
    
    if failedAttempts >= MaxLoginAttempts {
        // 账户锁定
        return nil, kernel.NewBusinessError(
            response.CodeUserLocked,
            "账户已被锁定，请稍后再试",
            kernel.WithMetadata("identifier", cmd.Identifier),
        )
    }
    
    // 查找用户并验证密码
    user, err := s.findUserByIdentifier(ctx, cmd.Identifier)
    if err != nil {
        s.recordLoginFailure(cmd.Identifier)
        return nil, application_shared.ErrInvalidCredentials
    }
    
    err = user.Login(cmd.Password, cmd.IP, cmd.UserAgent)
    if err != nil {
        s.recordLoginFailure(cmd.Identifier)
        return nil, err
    }
    
    // 登录成功，清除失败记录
    s.clearLoginFailure(cmd.Identifier)
    
    // 保存用户状态
    err = s.userRepo.Save(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("save user failed: %w", err)
    }
    
    return &AuthResult{
        UserID:    user.ID().String(),
        Username:  user.Username().String(),
        Email:     user.Email().String(),
    }, nil
}

func (s *AuthServiceImpl) recordLoginFailure(identifier string) {
    key := fmt.Sprintf("login:failed:%s", identifier)
    s.redis.Incr(s.ctx, key)
    s.redis.Expire(s.ctx, key, 30*time.Minute)  // 30 分钟后重置
}

func (s *AuthServiceImpl) getLoginFailedAttempts(identifier string) (int, error) {
    key := fmt.Sprintf("login:failed:%s", identifier)
    val, err := s.redis.Get(s.ctx, key).Int()
    if err == redis.Nil {
        return 0, nil
    }
    return val, err
}
```

---

## 🔒 授权安全

### RBAC 权限模型

```go
// domain/tenant/service/permission_checker.go
type PermissionChecker struct {
    roleRepo repository.RoleRepository
    userRepo repository.UserRepository
}

func (c *PermissionChecker) HasPermission(
    ctx context.Context,
    userID int64,
    tenantID int64,
    resource string,
    action string,
) (bool, error) {
    // 获取用户角色
    roles, err := c.userRepo.GetUserRoles(ctx, userID, tenantID)
    if err != nil {
        return false, fmt.Errorf("get user roles failed: %w", err)
    }
    
    // 检查每个角色的权限
    for _, role := range roles {
        permissions, err := c.roleRepo.GetPermissions(ctx, role.ID)
        if err != nil {
            return false, fmt.Errorf("get role permissions failed: %w", err)
        }
        
        for _, perm := range permissions {
            if perm.Resource == resource && perm.Action == action {
                return true, nil
            }
        }
    }
    
    return false, nil
}

// 权限检查中间件
func RequirePermission(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从上下文获取用户信息
        userID := GetUserIDFromContext(c)
        tenantID := GetTenantIDFromContext(c)
        
        checker := c.MustGet("permissionChecker").(*PermissionChecker)
        
        hasPerm, err := checker.HasPermission(c.Request.Context(), userID, tenantID, resource, action)
        if err != nil || !hasPerm {
            c.JSON(http.StatusForbidden, gin.H{
                "code":    response.CodePermissionDenied,
                "message": "权限不足",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

## 🛡️ 输入验证

### XSS 防护

```go
// pkg/util/sanitize.go
package util

import (
    "github.com/microcosm-cc/bluemonday"
)

var (
    // 严格模式 - 只允许纯文本
    strictPolicy = bluemonday.StrictPolicy()
    
    // 宽松模式 - 允许部分 HTML 标签
    defaultPolicy = bluemonday.UGCPolicy()
)

// SanitizeString 清理字符串，防止 XSS
func SanitizeString(input string) string {
    return strictPolicy.Sanitize(input)
}

// SanitizeHTML 清理 HTML，保留允许的标签
func SanitizeHTML(html string) string {
    return defaultPolicy.Sanitize(html)
}

// SanitizeEmail 验证并清理邮箱
func SanitizeEmail(email string) (string, error) {
    email = strings.TrimSpace(strings.ToLower(email))
    
    if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`).MatchString(email) {
        return "", kernel.FieldError("email", "无效的邮箱格式", email)
    }
    
    return email, nil
}
```

### SQL 注入防护

```go
// ✅ 正确：使用参数化查询
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*aggregate.User, error) {
    dao, err := r.daoQuery.User().
        WithContext(ctx).
        Where(r.daoQuery.User().Email.Eq(email)).  // 参数化
        First()
    if err != nil {
        return nil, err
    }
    return r.toDomain(dao)
}

// ❌ 错误：字符串拼接（禁止！）
func (r *UserRepository) FindByEmailUnsafe(email string) (*aggregate.User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)  // 危险！
    // ...
}
```

### 请求验证

```go
// interfaces/http/auth/dto.go
type LoginRequest struct {
    Identifier string `json:"identifier" validate:"required"`
    Password   string `json:"password" validate:"required,min=8"`
}

func (r *LoginRequest) Validate() error {
    // 验证标识符（邮箱或用户名）
    if r.Identifier == "" {
        return kernel.FieldError("identifier", "标识符不能为空", r.Identifier)
    }
    
    // 验证密码强度
    policy := auth.NewDefaultPasswordPolicy()
    if err := policy.Validate(r.Password); err != nil {
        return err
    }
    
    return nil
}
```

---

## 🔍 审计日志

### 审计追踪

```go
// domain/shared/audit_log.go
type AuditLog struct {
    ID          int64
    UserID      int64
    Action      string  // CREATE, UPDATE, DELETE, LOGIN, LOGOUT
    Resource    string  // User, Tenant, Role
    ResourceID  int64
    OldValues   map[string]interface{}
    NewValues   map[string]interface{}
    IPAddress   string
    UserAgent   string
    CreatedAt   time.Time
}

// 审计日志服务
type AuditLogService interface {
    Log(ctx context.Context, log *AuditLog) error
    Query(ctx context.Context, criteria AuditLogCriteria) ([]*AuditLog, error)
}

// 使用示例
func (s *UserService) CreateUser(ctx context.Context, cmd *CreateUserCommand) (*User, error) {
    user, err := aggregate.NewUser(cmd.Username, cmd.Email, cmd.Password)
    if err != nil {
        return nil, err
    }
    
    err = s.userRepo.Save(ctx, user)
    if err != nil {
        return nil, err
    }
    
    // 记录审计日志
    s.auditLogger.Log(ctx, &AuditLog{
        Action:     "CREATE",
        Resource:   "User",
        ResourceID: user.ID().Value(),
        NewValues: map[string]interface{}{
            "username": user.Username().String(),
            "email":    user.Email().String(),
        },
        IPAddress: getClientIP(ctx),
    })
    
    return user, nil
}
```

---

## 🔐 数据加密

### 敏感数据加密存储

```go
// infrastructure/crypto/encryption.go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

type Encryptor struct {
    key []byte
}

func NewEncryptor(secret string) (*Encryptor, error) {
    block, err := aes.NewCipher([]byte(secret))
    if err != nil {
        return nil, err
    }
    
    return &Encryptor{key: []byte(secret)}, nil
}

func (e *Encryptor) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", err
    }
    
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }
    
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
    
    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
    decoded, err := base64.URLEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", err
    }
    
    if len(decoded) < aes.BlockSize {
        return "", fmt.Errorf("ciphertext too short")
    }
    
    iv := decoded[:aes.BlockSize]
    decoded = decoded[aes.BlockSize:]
    
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(decoded, decoded)
    
    return string(decoded), nil
}
```

---

## 🚨 安全响应头

### HTTP 安全头

```go
// interfaces/http/middleware/security.go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 防止点击劫持
        c.Header("X-Frame-Options", "DENY")
        
        // 防止 MIME 类型嗅探
        c.Header("X-Content-Type-Options", "nosniff")
        
        // XSS 防护
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // 内容安全策略
        c.Header("Content-Security-Policy", "default-src 'self'")
        
        // 严格传输安全
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // Referrer 策略
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        c.Next()
    }
}
```

---

## 📊 安全监控

### 异常检测

```go
// 检测异常登录行为
func (s *AuthService) detectSuspiciousLogin(ip, userAgent string) bool {
    // 检查地理位置突变
    lastLogin := s.getLastLoginLocation()
    currentLocation := s.getGeoLocation(ip)
    
    if lastLogin != nil && s.isDistanceTooFar(lastLogin, currentLocation) {
        s.logger.Warn("suspicious login detected",
            zap.String("last_location", lastLogin),
            zap.String("current_location", currentLocation),
        )
        return true
    }
    
    // 检查设备变更
    if s.isNewDevice(userAgent) {
        s.logger.Info("new device login",
            zap.String("user_agent", userAgent),
        )
    }
    
    return false
}
```

---

## ✅ 安全检查清单

### 开发阶段

- [ ] 所有输入都经过验证
- [ ] 使用参数化查询，防止 SQL 注入
- [ ] 输出经过清理，防止 XSS
- [ ] 密码使用强哈希算法（bcrypt）
- [ ] 敏感数据加密存储
- [ ] 实现适当的访问控制
- [ ] 记录所有敏感操作

### 部署阶段

- [ ] 修改默认密码和密钥
- [ ] 禁用调试模式
- [ ] 配置 HTTPS
- [ ] 设置安全响应头
- [ ] 限制 API 访问频率
- [ ] 配置防火墙规则
- [ ] 启用日志审计

### 运维阶段

- [ ] 定期更新依赖
- [ ] 监控系统异常
- [ ] 审查访问日志
- [ ] 定期备份数据
- [ ] 进行安全扫描
- [ ] 执行渗透测试
- [ ] 制定应急响应计划

---

## 📚 参考资源

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go 语言安全编程](https://golang.org/doc/articles/go_command.html)
- [JWT 最佳实践](https://tools.ietf.org/html/rfc7519)
- [bcrypt 密码哈希](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
