package errors

import (
	"net/http"
)

// ============================================
// 网络错误码定义
// ============================================
var (
	// Network 错误
	ErrNetworkTimeout        = NewCategorized("Network", "Network.Timeout", "网络连接超时，请检查网络后重试")
	ErrNetworkUnavailable    = NewCategorized("Network", "Network.Unavailable", "网络不可用，请检查网络连接")
	ErrConnectionRefused     = NewCategorized("Network", "Network.ConnectionRefused", "连接被拒绝，请稍后重试")
	ErrDNSResolutionFailed   = NewCategorized("Network", "Network.DNSResolutionFailed", "DNS 解析失败")
	ErrSSLHandshakeFailed    = NewCategorized("Network", "Network.SSLHandshakeFailed", "SSL 握手失败")
	ErrRequestEntityTooLarge = NewCategorized("Network", "Network.RequestEntityTooLarge", "请求数据过大")
	ErrServiceUnavailable    = NewCategorized("Network", "Network.ServiceUnavailable", "服务暂时不可用，请稍后重试")
)

// ============================================
// HTTP 状态码映射表
// ============================================

// HTTPStatusMapper HTTP 状态码映射器
type HTTPStatusMapper struct {
	statusCode int
	message    string
}

// GetHTTPStatus 根据业务错误获取对应的 HTTP 状态码
func GetHTTPStatus(err error) (int, string) {
	if err == nil {
		return http.StatusOK, "Success"
	}

	appErr, ok := err.(*AppError)
	if !ok {
		// 非 AppError，返回 500
		return http.StatusInternalServerError, "Internal Server Error"
	}

	switch appErr.GetCategory() {
	case "Common":
		return mapCommonErrorToHTTP(appErr)
	case "User":
		return mapUserErrorToHTTP(appErr)
	case "KG":
		return mapKGErrorToHTTP(appErr)
	case "System":
		return mapSystemErrorToHTTP(appErr)
	case "Network":
		return mapNetworkErrorToHTTP(appErr)
	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

// mapCommonErrorToHTTP 通用错误映射
func mapCommonErrorToHTTP(err *AppError) (int, string) {
	switch err.GetCode() {
	case "Success":
		return http.StatusOK, "OK"
	case "InvalidParameter", "MissingParameter", "ValidationFailed":
		return http.StatusBadRequest, "Bad Request"
	case "Unauthorized":
		return http.StatusUnauthorized, "Unauthorized"
	case "Forbidden":
		return http.StatusForbidden, "Forbidden"
	case "NotFound":
		return http.StatusNotFound, "Not Found"
	case "MethodNotAllowed":
		return http.StatusMethodNotAllowed, "Method Not Allowed"
	case "TooManyRequests":
		return http.StatusTooManyRequests, "Too Many Requests"
	case "ResourceConflict":
		return http.StatusConflict, "Conflict"
	case "UnsupportedMediaType":
		return http.StatusUnsupportedMediaType, "Unsupported Media Type"
	default:
		return http.StatusBadRequest, "Bad Request"
	}
}

// mapUserErrorToHTTP 用户模块错误映射
func mapUserErrorToHTTP(err *AppError) (int, string) {
	switch err.GetCode() {
	case "User.Exists":
		return http.StatusConflict, "Conflict"
	case "User.NotFound":
		return http.StatusNotFound, "Not Found"
	case "User.InvalidPassword", "User.Unauthorized":
		return http.StatusUnauthorized, "Unauthorized"
	case "User.TenantLimitExceed":
		return http.StatusForbidden, "Forbidden"
	case "User.InvalidEmail":
		return http.StatusBadRequest, "Bad Request"
	case "User.InvalidRole":
		return http.StatusBadRequest, "Bad Request"
	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

// mapKGErrorToHTTP 知识图谱错误映射
func mapKGErrorToHTTP(err *AppError) (int, string) {
	switch err.GetCode() {
	case "KG.Domain.NotFound", "KG.Trunk.NotFound", "KG.Node.NotFound", "KG.Relationship.NotFound":
		return http.StatusNotFound, "Not Found"
	case "KG.Domain.AlreadyExists", "KG.Trunk.AlreadyExists", "KG.Node.AlreadyExists", "KG.Relationship.AlreadyExists":
		return http.StatusConflict, "Conflict"
	case "KG.Domain.InvalidData", "KG.Trunk.InvalidData", "KG.Node.InvalidData", "KG.Relationship.InvalidData":
		return http.StatusBadRequest, "Bad Request"
	case "KG.Node.InvalidType":
		return http.StatusBadRequest, "Bad Request"
	case "KG.Trunk.NotInDomain", "KG.Node.NotInTrunk":
		return http.StatusBadRequest, "Bad Request"
	case "KG.Relationship.InvalidType":
		return http.StatusBadRequest, "Bad Request"
	case "KG.Relationship.CycleDetected":
		return http.StatusBadRequest, "Bad Request"
	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

// mapSystemErrorToHTTP 系统错误映射
func mapSystemErrorToHTTP(err *AppError) (int, string) {
	switch err.GetCode() {
	case "System.InternalError":
		return http.StatusInternalServerError, "Internal Server Error"
	case "System.DatabaseError":
		return http.StatusServiceUnavailable, "Service Unavailable"
	case "System.CacheUnavailable":
		return http.StatusServiceUnavailable, "Service Unavailable"
	case "System.ExternalServiceError":
		return http.StatusBadGateway, "Bad Gateway"
	case "System.Timeout":
		return http.StatusGatewayTimeout, "Gateway Timeout"
	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

// mapNetworkErrorToHTTP 网络错误映射
func mapNetworkErrorToHTTP(err *AppError) (int, string) {
	switch err.GetCode() {
	case "Network.Timeout":
		return http.StatusGatewayTimeout, "Gateway Timeout"
	case "Network.Unavailable":
		return http.StatusServiceUnavailable, "Service Unavailable"
	case "Network.ConnectionRefused":
		return http.StatusServiceUnavailable, "Service Unavailable"
	case "Network.DNSResolutionFailed":
		return http.StatusServiceUnavailable, "Service Unavailable"
	case "Network.SSLHandshakeFailed":
		return http.StatusBadGateway, "Bad Gateway"
	case "Network.RequestEntityTooLarge":
		return http.StatusRequestEntityTooLarge, "Request Entity Too Large"
	case "Network.ServiceUnavailable":
		return http.StatusServiceUnavailable, "Service Unavailable"
	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

// GetUserFriendlyMessage 获取用户友好的错误消息
func GetUserFriendlyMessage(err error) string {
	if err == nil {
		return "操作成功"
	}

	appErr, ok := err.(*AppError)
	if !ok {
		return "系统繁忙，请稍后重试"
	}

	// 如果是开发环境，可以返回更详细的信息
	// if config.IsDev() {
	//     return fmt.Sprintf("%s [%s]", appErr.GetMessage(), appErr.GetCode())
	// }

	// 生产环境，返回友好的用户提示
	return appErr.GetMessage()
}

// ShouldRetry 判断是否应该重试
func ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	appErr, ok := err.(*AppError)
	if !ok {
		// 未知错误，不重试
		return false
	}

	// 网络相关错误建议重试
	retryCodes := []string{
		"Network.Timeout",
		"Network.Unavailable",
		"Network.ConnectionRefused",
		"Network.DNSResolutionFailed",
		"System.Timeout",
		"System.DatabaseError",
		"System.CacheUnavailable",
		"System.ExternalServiceError",
		"System.ServiceUnavailable",
	}

	code := appErr.GetCode()
	for _, retryCode := range retryCodes {
		if code == retryCode {
			return true
		}
	}

	return false
}

// IsClientError 判断是否是客户端错误（4xx）
func IsClientError(err error) bool {
	if err == nil {
		return false
	}

	statusCode, _ := GetHTTPStatus(err)
	return statusCode >= 400 && statusCode < 500
}

// IsServerError 判断是否是服务端错误（5xx）
func IsServerError(err error) bool {
	if err == nil {
		return false
	}

	statusCode, _ := GetHTTPStatus(err)
	return statusCode >= 500
}

// IsNetworkError 判断是否是网络错误
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	appErr, ok := err.(*AppError)
	if !ok {
		return false
	}

	return appErr.GetCategory() == "Network"
}
