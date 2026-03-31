package common

import (
	"fmt"
)

// ============================================
// 业务错误辅助函数
// ============================================

// WrapBusinessError 包装业务错误，添加额外信息
func WrapBusinessError(err error, format string, args ...interface{}) error {
	if be := AsBusinessError(err); be != nil {
		return &BusinessError{
			Code:    be.Code,
			Message: fmt.Sprintf(format+" "+be.Message, args...),
			Details: be.Details,
		}
	}
	return err
}
