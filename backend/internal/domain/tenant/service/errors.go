// Package service 租户领域服务相关错误定义
package service

import "errors"

var (
	ErrCannotChangeOwnerRole = errors.New("cannot change owner role")
	ErrCannotPromoteToOwner  = errors.New("cannot promote to owner")
)
