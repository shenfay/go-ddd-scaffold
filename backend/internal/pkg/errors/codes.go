package errors

// ============================================
// 通用错误码定义
// ============================================
var (
	// Success 成功
	Success = NewCategorized("Common", "Success", "操作成功")

	// InvalidParameter 参数错误
	InvalidParameter = NewCategorized("Common", "InvalidParameter", "无效请求，请检查输入参数")

	// MissingParameter 缺少必要参数
	MissingParameter = NewCategorized("Common", "MissingParameter", "缺少必要参数")

	// Unauthorized 未授权
	Unauthorized = NewCategorized("Common", "Unauthorized", "未授权，请先登录")

	// Forbidden 禁止访问
	Forbidden = NewCategorized("Common", "Forbidden", "禁止访问")

	// NotFound 资源不存在
	NotFound = NewCategorized("Common", "NotFound", "请求的资源不存在")

	// MethodNotAllowed 不支持的请求方法
	MethodNotAllowed = NewCategorized("Common", "MethodNotAllowed", "不支持的请求方法")

	// TooManyRequests 请求过于频繁
	TooManyRequests = NewCategorized("Common", "TooManyRequests", "请求过于频繁，请稍后再试")

	// ValidationFailed 参数校验失败
	ValidationFailed = NewCategorized("Common", "ValidationFailed", "参数校验失败")

	// ResourceConflict 资源冲突
	ResourceConflict = NewCategorized("Common", "ResourceConflict", "资源冲突")

	// UnsupportedMediaType 不支持的媒体类型
	UnsupportedMediaType = NewCategorized("Common", "UnsupportedMediaType", "不支持的媒体类型")
)

// ============================================
// 用户模块错误码
// ============================================
var (
	// User 模块错误
	ErrUserExists        = NewCategorized("User", "User.Exists", "用户已存在")
	ErrUserNotFound      = NewCategorized("User", "User.NotFound", "用户不存在")
	ErrInvalidPassword   = NewCategorized("User", "User.InvalidPassword", "密码错误")
	ErrUnauthorized      = NewCategorized("User", "User.Unauthorized", "未授权或账户已被禁用")
	ErrTenantLimitExceed = NewCategorized("User", "User.TenantLimitExceed", "超过租户用户数限制")

	// User 校验错误
	ErrInvalidEmail = NewCategorized("User", "User.InvalidEmail", "无效的邮箱格式")
	ErrInvalidRole  = NewCategorized("User", "User.InvalidRole", "无效的用户角色")
)

// ============================================
// 知识图谱错误码定义
// ============================================
var (
	// Domain 错误
	DomainNotFound      = NewCategorized("KG", "KG.Domain.NotFound", "知识领域不存在")
	DomainAlreadyExists = NewCategorized("KG", "KG.Domain.AlreadyExists", "知识领域已存在")
	DomainInvalidData   = New("KG.Domain.InvalidData", "知识领域数据无效").SetCategory("KG")

	// Trunk 错误
	TrunkNotFound      = NewCategorized("KG", "KG.Trunk.NotFound", "知识主线不存在")
	TrunkAlreadyExists = NewCategorized("KG", "KG.Trunk.AlreadyExists", "知识主线已存在")
	TrunkNotInDomain   = NewCategorized("KG", "KG.Trunk.NotInDomain", "知识主线不属于指定的领域")
	TrunkInvalidData   = New("KG.Trunk.InvalidData", "知识主线数据无效").SetCategory("KG")

	// Node 错误
	NodeNotFound      = NewCategorized("KG", "KG.Node.NotFound", "知识节点不存在")
	NodeAlreadyExists = NewCategorized("KG", "KG.Node.AlreadyExists", "知识节点已存在")
	NodeInvalidType   = NewCategorized("KG", "KG.Node.InvalidType", "无效的节点类型，必须是 C/S/T/P 之一")
	NodeNotInTrunk    = NewCategorized("KG", "KG.Node.NotInTrunk", "知识节点不属于指定的主线")
	NodeInvalidData   = New("KG.Node.InvalidData", "知识节点数据无效").SetCategory("KG")

	// Relationship 错误
	RelationshipNotFound      = NewCategorized("KG", "KG.Relationship.NotFound", "知识关系不存在")
	RelationshipAlreadyExists = NewCategorized("KG", "KG.Relationship.AlreadyExists", "知识关系已存在")
	RelationshipInvalidType   = NewCategorized("KG", "KG.Relationship.InvalidType", "无效的关系类型，必须是 PREREQ/SUP_SKILL/THINK_PAT 之一")
	RelationshipCycleDetected = NewCategorized("KG", "KG.Relationship.CycleDetected", "检测到循环引用，无法建立关系")
	RelationshipInvalidData   = New("KG.Relationship.InvalidData", "知识关系数据无效").SetCategory("KG")
)

// ============================================
// 系统错误码定义
// ============================================
var (
	InternalError        = NewCategorized("System", "System.InternalError", "系统内部错误，请稍后重试")
	DatabaseError        = NewCategorized("System", "System.DatabaseError", "数据库操作失败")
	CacheUnavailable     = NewCategorized("System", "System.CacheUnavailable", "缓存服务不可用")
	ExternalServiceError = NewCategorized("System", "System.ExternalServiceError", "外部服务调用失败")
	Timeout              = NewCategorized("System", "System.Timeout", "请求超时，请稍后重试")
)
