package shared

import "github.com/shenfay/go-ddd-scaffold/pkg/utils/ulid"

// GenerateUserID 生成用户ID
func GenerateUserID() string {
	return ulid.GenerateUserID()
}

// GenerateAuditLogID 生成审计日志ID
func GenerateAuditLogID() string {
	return ulid.GenerateAuditLogID()
}

// GenerateActivityLogID 生成活动日志ID
func GenerateActivityLogID() string {
	return ulid.GenerateActivityLogID()
}
