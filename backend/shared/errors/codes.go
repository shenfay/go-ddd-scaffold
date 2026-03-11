package apperrors

import (
	"fmt"
)

// ============================================
// 0 - 成功
// ============================================
const (
	// CodeSuccess 成功（统一为0）
	CodeSuccess = 0
)

// ============================================
// 错误码信息结构
// ============================================

// CodeInfo 错误码信息
type CodeInfo struct {
	Code       int    // 错误码
	HTTPStatus int    // HTTP状态码
	Message    string // 错误消息
	Module     string // 模块名称
}

// GetCodeInfo 获取错误码信息
func GetCodeInfo(code int) CodeInfo {
	info, ok := errorCodes[code]
	if !ok {
		return CodeInfo{
			Code:       code,
			HTTPStatus: 500,
			Message:    "未知错误",
			Module:     "unknown",
		}
	}
	return info
}

// GetHTTPStatus 获取HTTP状态码
func GetHTTPStatus(code int) int {
	info := GetCodeInfo(code)
	return info.HTTPStatus
}

// GetMessage 获取错误消息
func GetMessage(code int) string {
	info := GetCodeInfo(code)
	return info.Message
}

// errorCodes 错误码映射表
var errorCodes = map[int]CodeInfo{
	// 通用错误
	CodeSuccess:      {Code: 0, HTTPStatus: 200, Message: "success", Module: "system"},
	CodeUnknownError: {Code: 10000, HTTPStatus: 500, Message: "未知错误", Module: "system"},
	CodeInvalidParam: {Code: 10001, HTTPStatus: 400, Message: "参数无效", Module: "system"},
	CodeNotFound:     {Code: 10002, HTTPStatus: 404, Message: "资源不存在", Module: "system"},
	CodeConflict:     {Code: 10003, HTTPStatus: 409, Message: "资源冲突", Module: "system"},
	CodeUnauthorized: {Code: 10004, HTTPStatus: 401, Message: "未授权", Module: "system"},
	CodeForbidden:    {Code: 10005, HTTPStatus: 403, Message: "禁止访问", Module: "system"},
	CodeBadRequest:   {Code: 10009, HTTPStatus: 400, Message: "请求格式错误", Module: "system"},

	// 系统内部错误
	CodeInternalError: {Code: 10010, HTTPStatus: 500, Message: "内部错误", Module: "system"},
	CodeDatabaseError: {Code: 10011, HTTPStatus: 500, Message: "数据库错误", Module: "system"},
	CodeCacheError:    {Code: 10012, HTTPStatus: 500, Message: "缓存错误", Module: "system"},
	CodeExternalError: {Code: 10013, HTTPStatus: 502, Message: "外部服务错误", Module: "system"},
	CodeConcurrency:   {Code: 10014, HTTPStatus: 409, Message: "并发冲突", Module: "system"},

	// 用户模块
	CodeUserNotFound:    {Code: 20001, HTTPStatus: 404, Message: "用户不存在", Module: "user"},
	CodeUserExists:      {Code: 20002, HTTPStatus: 409, Message: "用户已存在", Module: "user"},
	CodeInvalidPassword: {Code: 21001, HTTPStatus: 401, Message: "密码错误", Module: "user"},
	CodeAccountLocked:   {Code: 21004, HTTPStatus: 403, Message: "账户已锁定", Module: "user"},
	CodeInvalidEmail:    {Code: 22002, HTTPStatus: 400, Message: "邮箱格式无效", Module: "user"},
	CodeEmailExists:     {Code: 22004, HTTPStatus: 409, Message: "邮箱已存在", Module: "user"},

	// 租户模块
	CodeTenantNotFound:  {Code: 30001, HTTPStatus: 404, Message: "租户不存在", Module: "tenant"},
	CodeTenantExists:    {Code: 30002, HTTPStatus: 409, Message: "租户已存在", Module: "tenant"},
	CodeNotTenantMember: {Code: 31001, HTTPStatus: 403, Message: "不是租户成员", Module: "tenant"},
	CodeNotTenantOwner:  {Code: 32003, HTTPStatus: 403, Message: "不是租户所有者", Module: "tenant"},

	// 认证授权
	CodeTokenExpired:     {Code: 40001, HTTPStatus: 401, Message: "Token已过期", Module: "auth"},
	CodeTokenInvalid:     {Code: 40002, HTTPStatus: 401, Message: "Token无效", Module: "auth"},
	CodePermissionDenied: {Code: 41001, HTTPStatus: 403, Message: "权限不足", Module: "auth"},
}

// Error 实现error接口
func (ci CodeInfo) Error() string {
	return fmt.Sprintf("[%05d] %s", ci.Code, ci.Message)
}

// ============================================
// 10000-10999 - 系统级通用错误
// ============================================
const (
	// 通用错误 (10000-10099)
	CodeUnknownError    = 10000 // 未知错误
	CodeInvalidParam    = 10001 // 参数无效
	CodeNotFound        = 10002 // 资源不存在
	CodeConflict        = 10003 // 资源冲突
	CodeUnauthorized    = 10004 // 未授权
	CodeForbidden       = 10005 // 禁止访问
	CodeTooManyRequests = 10006 // 请求过于频繁
	CodeTimeout         = 10007 // 请求超时
	CodeNotImplemented  = 10008 // 功能未实现
	CodeBadRequest      = 10009 // 请求格式错误

	// 系统内部错误 (10010-10099)
	CodeInternalError      = 10010 // 内部错误
	CodeDatabaseError      = 10011 // 数据库错误
	CodeCacheError         = 10012 // 缓存错误
	CodeExternalError      = 10013 // 外部服务错误
	CodeConcurrency        = 10014 // 并发冲突
	CodeNetworkError       = 10015 // 网络错误
	CodeSerializeError     = 10016 // 序列化错误
	CodeConfigError        = 10017 // 配置错误
	CodeServiceUnavailable = 10018 // 服务不可用
)

// ============================================
// 20000-29999 - 用户模块
// ============================================

// 用户基础 (20000-20999)
const (
	CodeUserNotFound     = 20001 // 用户不存在
	CodeUserExists       = 20002 // 用户已存在
	CodeUserCreateFailed = 20003 // 用户创建失败
	CodeUserUpdateFailed = 20004 // 用户更新失败
	CodeUserDeleteFailed = 20005 // 用户删除失败
	CodeUserDisabled     = 20006 // 用户已禁用
)

// 用户认证 (21000-21999) - 子模块21
const (
	CodeInvalidPassword     = 21001 // 密码错误
	CodePasswordTooWeak     = 21002 // 密码强度不足
	CodePasswordExpired     = 21003 // 密码已过期
	CodeAccountLocked       = 21004 // 账户已锁定
	CodeAccountDisabled     = 21005 // 账户已禁用
	CodeAccountNotActivated = 21006 // 账户未激活
	CodeLoginFailed         = 21007 // 登录失败
	CodeLogoutFailed        = 21008 // 登出失败
	CodeInvalidCredentials  = 21009 // 凭证无效
)

// 用户资料 (22000-22999) - 子模块22
const (
	CodeInvalidUsername    = 22001 // 用户名无效
	CodeInvalidEmail       = 22002 // 邮箱格式无效
	CodeInvalidPhone       = 22003 // 手机号无效
	CodeEmailExists        = 22004 // 邮箱已存在
	CodePhoneExists        = 22005 // 手机号已存在
	CodeAvatarUploadFailed = 22006 // 头像上传失败
	CodeUsernameTooLong    = 22007 // 用户名过长
	CodeUsernameTooShort   = 22008 // 用户名过短
)

// 用户关系 (23000-23999) - 子模块23
const (
	CodeUserNotFollowed  = 23001 // 未关注该用户
	CodeAlreadyFollowing = 23002 // 已关注该用户
	CodeCannotFollowSelf = 23003 // 不能关注自己
)

// ============================================
// 30000-39999 - 租户模块
// ============================================

// 租户基础 (30000-30999)
const (
	CodeTenantNotFound     = 30001 // 租户不存在
	CodeTenantExists       = 30002 // 租户已存在
	CodeTenantCreateFailed = 30003 // 租户创建失败
	CodeTenantUpdateFailed = 30004 // 租户更新失败
	CodeTenantDeleteFailed = 30005 // 租户删除失败
	CodeTenantDisabled     = 30006 // 租户已禁用
)

// 租户成员 (31000-31999) - 子模块31
const (
	CodeNotTenantMember      = 31001 // 不是租户成员
	CodeAlreadyMember        = 31002 // 已是租户成员
	CodeMemberAddFailed      = 31003 // 添加成员失败
	CodeMemberRemoveFailed   = 31004 // 移除成员失败
	CodeMemberUpdateFailed   = 31005 // 更新成员失败
	CodeLastOwnerCannotLeave = 31006 // 最后一位所有者不能离开
	CodeMemberNotFound       = 31007 // 成员不存在
)

// 租户角色权限 (32000-32999) - 子模块32
const (
	CodeTenantRoleNotFound = 32001 // 租户角色不存在
	CodeTenantRoleInvalid  = 32002 // 租户角色无效
	CodeNotTenantOwner     = 32003 // 不是租户所有者
	CodeNotTenantAdmin     = 32004 // 不是租户管理员
	CodeRoleAssignFailed   = 32005 // 角色分配失败
	CodeRoleRevokeFailed   = 32006 // 角色撤销失败
)

// 租户限制 (33000-33999) - 子模块33
const (
	CodeTenantLimitReached = 33001 // 租户数量已达上限
	CodeTenantMemberLimit  = 33002 // 租户成员数已达上限
	CodeTenantExpired      = 33003 // 租户已过期
	CodeTenantSuspended    = 33004 // 租户已暂停
)

// ============================================
// 40000-49999 - 认证授权模块
// ============================================

// Token相关 (40000-40999)
const (
	CodeTokenExpired        = 40001 // Token已过期
	CodeTokenInvalid        = 40002 // Token无效
	CodeTokenRevoked        = 40003 // Token已撤销
	CodeTokenMissing        = 40004 // 缺少Token
	CodeTokenFormatInvalid  = 40005 // Token格式无效
	CodeRefreshTokenExpired = 40006 // 刷新Token已过期
	CodeRefreshTokenInvalid = 40007 // 刷新Token无效
)

// 权限相关 (41000-41999) - 子模块41
const (
	CodePermissionDenied  = 41001 // 权限不足
	CodeRoleNotFound      = 41002 // 角色不存在
	CodeResourceNotFound  = 41003 // 资源不存在
	CodeActionNotAllowed  = 41004 // 操作不允许
	CodeInsufficientQuota = 41005 // 配额不足
)

// 验证码安全 (42000-42999) - 子模块42
const (
	CodeCaptchaRequired  = 42001 // 需要验证码
	CodeCaptchaInvalid   = 42002 // 验证码错误
	CodeCaptchaExpired   = 42003 // 验证码已过期
	CodeCaptchaTooMany   = 42004 // 验证码尝试次数过多
	CodeIPBlocked        = 42005 // IP已被封禁
	CodeDeviceNotTrusted = 42006 // 设备不受信任
	CodeMFARequired      = 42007 // 需要多因素认证
	CodeMFACodeInvalid   = 42008 // MFA验证码错误
	CodeMFASetupFailed   = 42009 // MFA设置失败
)

// ============================================
// 50000-59999 - 内容资源模块 (预留)
// ============================================

// 内容基础 (50000-50999)
const (
	CodeContentNotFound     = 50001 // 内容不存在
	CodeContentExists       = 50002 // 内容已存在
	CodeContentCreateFailed = 50003 // 内容创建失败
	CodeContentUpdateFailed = 50004 // 内容更新失败
	CodeContentDeleteFailed = 50005 // 内容删除失败
)

// 文件媒体 (51000-51999)
const (
	CodeFileNotFound       = 51001 // 文件不存在
	CodeFileTooLarge       = 51002 // 文件过大
	CodeFileTypeInvalid    = 51003 // 文件类型无效
	CodeFileUploadFailed   = 51004 // 文件上传失败
	CodeFileDownloadFailed = 51005 // 文件下载失败
	CodeFileNotAllowed     = 51006 // 文件不允许上传
)

// ============================================
// 60000-69999 - 交易订单模块 (预留)
// ============================================

const (
	// 订单基础
	CodeOrderNotFound     = 60001 // 订单不存在
	CodeOrderCreateFailed = 60002 // 订单创建失败
	CodeOrderUpdateFailed = 60003 // 订单更新失败
	CodeOrderCancelFailed = 60004 // 订单取消失败

	// 订单状态
	CodeOrderAlreadyPaid   = 60010 // 订单已支付
	CodeOrderAlreadyClosed = 60011 // 订单已关闭
	CodeOrderCannotCancel  = 60012 // 订单不可取消

	// 支付相关
	CodePaymentFailed  = 60101 // 支付失败
	CodePaymentTimeout = 60102 // 支付超时
	CodeRefundFailed   = 60103 // 退款失败
)

// ============================================
// 70000-79999 - 消息通知模块 (预留)
// ============================================

const (
	CodeMessageNotFound    = 70001 // 消息不存在
	CodeMessageSendFailed  = 70002 // 消息发送失败
	CodeNotificationFailed = 70003 // 通知发送失败
)

// ============================================
// 80000-89999 - 工作流模块 (预留)
// ============================================

const (
	CodeWorkflowNotFound      = 80001 // 工作流不存在
	CodeWorkflowStartFailed   = 80002 // 工作流启动失败
	CodeWorkflowApproveFailed = 80003 // 审批失败
	CodeWorkflowRejectFailed  = 80004 // 拒绝失败
)

// ============================================
// 90000-99999 - 扩展预留
// ============================================
