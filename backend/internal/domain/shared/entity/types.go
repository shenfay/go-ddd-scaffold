// Package entity 提供共享的领域实体类型定义
package entity

// UserRole 用户角色枚举（从 user/entity 迁移过来，避免循环依赖）
type UserRole string

const (
	// 基础角色
	RoleOwner  UserRole = "owner"  // 所有者
	RoleAdmin  UserRole = "admin"  // 管理员
	RoleMember UserRole = "member" // 普通成员
	RoleGuest  UserRole = "guest"  // 访客

	// 管理角色
	RoleSuperAdmin   UserRole = "super_admin"   // 超级管理员
	RoleContentAdmin UserRole = "content_admin" // 内容管理员
	RoleOpsAdmin     UserRole = "ops_admin"     // 运营管理员
)
