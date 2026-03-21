package service

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/tenant/valueobject"
	useraggregate "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	userrepo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	uservo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// TenantService 租户领域服务
type TenantService struct {
	tenantRepo repository.TenantRepository
	userRepo   userrepo.UserRepository
}

// NewTenantService 创建租户服务
func NewTenantService(tenantRepo repository.TenantRepository, userRepo userrepo.UserRepository) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
	}
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, code, name string, ownerID uservo.UserID) (*aggregate.Tenant, error) {
	// 1. 检查租户编码是否已存在
	if _, err := s.tenantRepo.FindByCode(ctx, code); err == nil {
		return nil, kernel.NewBusinessError(aggregate.CodeTenantCodeExists, "tenant code already exists")
	}

	// 2. 检查所有者是否存在
	if _, err := s.userRepo.FindByID(ctx, ownerID); err != nil {
		return nil, kernel.NewBusinessError(useraggregate.CodeUserNotFound, "owner user not found")
	}

	// 3. 创建租户
	tenant, err := aggregate.NewTenant(code, name, ownerID)
	if err != nil {
		return nil, err
	}

	// 4. 保存租户
	if err := s.tenantRepo.Save(ctx, tenant); err != nil {
		return nil, err
	}

	// 5. 添加所有者为成员
	member := &valueobject.TenantMember{
		UserID:   ownerID,
		TenantID: tenant.ID().(valueobject.TenantID),
		Role:     valueobject.TenantRoleOwner,
		JoinedAt: time.Now().Format(time.RFC3339),
	}
	if err := s.tenantRepo.AddMember(ctx, tenant.ID().(valueobject.TenantID), member); err != nil {
		return nil, err
	}

	return tenant, nil
}

// AddMember 添加成员到租户
func (s *TenantService) AddMember(ctx context.Context, tenantID valueobject.TenantID, userID, addedBy uservo.UserID, role valueobject.TenantRole) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 检查租户状态
	if !tenant.IsActive() {
		return kernel.NewBusinessError(aggregate.CodeTenantNotActive, "tenant is not active")
	}

	// 3. 检查用户是否存在
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return kernel.NewBusinessError(useraggregate.CodeUserNotFound, "user not found")
	}

	// 4. 检查用户是否已经是成员
	if _, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, userID); err == nil {
		return kernel.NewBusinessError(aggregate.CodeAlreadyMember, "user is already a member of this tenant")
	}

	// 5. 检查成员数量限制
	members, err := s.tenantRepo.FindMembers(ctx, tenantID)
	if err != nil {
		return err
	}

	if !tenant.CanAddMember(len(members)) {
		return kernel.NewBusinessError(aggregate.CodeTenantMaxMembersReached, "tenant has reached maximum member limit")
	}

	// 6. 检查操作权限（只有 Owner 和 Admin 可以添加成员）
	addedByMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, addedBy)
	if err != nil {
		return kernel.NewBusinessError(aggregate.CodeOperatorNotMember, "operator is not a member of this tenant")
	}

	if addedByMember.Role != valueobject.TenantRoleOwner && addedByMember.Role != valueobject.TenantRoleAdmin {
		return kernel.NewBusinessError(aggregate.CodeInsufficientPermissions, "insufficient permissions to add members")
	}

	// 7. 添加成员
	member := &valueobject.TenantMember{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		JoinedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.tenantRepo.AddMember(ctx, tenantID, member); err != nil {
		return err
	}

	// 8. 发布领域事件
	evt := event.NewTenantMemberAddedEvent(tenantID, userID, addedBy, role)
	tenant.ApplyEvent(evt)

	return s.tenantRepo.Save(ctx, tenant)
}

// RemoveMember 从租户移除成员
func (s *TenantService) RemoveMember(ctx context.Context, tenantID valueobject.TenantID, userID, removedBy uservo.UserID) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 检查成员是否存在
	member, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, userID)
	if err != nil {
		return kernel.NewBusinessError(aggregate.CodeNotTenantMember, "user is not a member of this tenant")
	}

	// 3. 检查操作权限
	removedByMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, removedBy)
	if err != nil {
		return kernel.NewBusinessError(aggregate.CodeOperatorNotMember, "operator is not a member of this tenant")
	}

	// 4. 权限检查规则
	if member.Role == valueobject.TenantRoleOwner {
		// 不能移除所有者
		return kernel.NewBusinessError(aggregate.CodeCannotRemoveOwner, "cannot remove tenant owner")
	}

	if removedByMember.Role != valueobject.TenantRoleOwner && removedByMember.Role != valueobject.TenantRoleAdmin {
		return kernel.NewBusinessError(aggregate.CodeInsufficientPermissions, "insufficient permissions to remove members")
	}

	// Admin 不能移除其他 Admin
	if member.Role == valueobject.TenantRoleAdmin && removedByMember.Role != valueobject.TenantRoleOwner {
		return kernel.NewBusinessError(aggregate.CodeCannotRemoveAdmin, "only owner can remove admin members")
	}

	// 5. 移除成员
	if err := s.tenantRepo.RemoveMember(ctx, tenantID, userID); err != nil {
		return err
	}

	// 6. 发布领域事件
	evt := event.NewTenantMemberRemovedEvent(tenantID, userID, removedBy)
	tenant.ApplyEvent(evt)

	return s.tenantRepo.Save(ctx, tenant)
}

// ChangeMemberRole 更改成员角色
func (s *TenantService) ChangeMemberRole(ctx context.Context, tenantID valueobject.TenantID, userID, changedBy uservo.UserID, newRole valueobject.TenantRole) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 检查成员是否存在
	member, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, userID)
	if err != nil {
		return kernel.NewBusinessError(aggregate.CodeNotTenantMember, "user is not a member of this tenant")
	}

	// 3. 检查操作权限
	changedByMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, changedBy)
	if err != nil {
		return kernel.NewBusinessError(aggregate.CodeOperatorNotMember, "operator is not a member of this tenant")
	}

	// 4. 权限检查规则
	if changedByMember.Role != valueobject.TenantRoleOwner {
		return kernel.NewBusinessError(aggregate.CodeInsufficientPermissions, "only owner can change member roles")
	}

	// 5. 不能更改自己的角色
	if userID.Equals(changedBy) {
		return kernel.NewBusinessError(aggregate.CodeCannotChangeOwnRole, "cannot change your own role")
	}

	// 6. 保存旧角色用于事件
	oldRole := member.Role

	// 7. 更新角色（这里简化处理，实际应该通过仓储更新）
	member.Role = newRole

	// 8. 发布领域事件
	evt := event.NewTenantMemberRoleChangedEvent(tenantID, userID, changedBy, oldRole, newRole)
	tenant.ApplyEvent(evt)

	return s.tenantRepo.Save(ctx, tenant)
}

// TransferOwnership 转移租户所有权
func (s *TenantService) TransferOwnership(ctx context.Context, tenantID valueobject.TenantID, newOwnerID, currentOwnerID uservo.UserID) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 验证当前所有者是租户所有者
	if !tenant.OwnerID().Equals(currentOwnerID) {
		return kernel.NewBusinessError(aggregate.CodeNotTenantOwner, "only current owner can transfer ownership")
	}

	// 3. 检查新所有者是否是租户成员
	newOwnerMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, newOwnerID)
	if err != nil {
		return kernel.NewBusinessError(aggregate.CodeNotTenantMember, "new owner must be a member of the tenant")
	}

	// 4. 更新租户所有者（通过领域方法）
	if err := tenant.TransferOwnership(newOwnerID); err != nil {
		return err
	}

	// 5. 更新新所有者为 Owner
	newOwnerMember.Role = valueobject.TenantRoleOwner

	// 6. 保存
	return s.tenantRepo.Save(ctx, tenant)
}

// DeactivateTenant 停用租户
func (s *TenantService) DeactivateTenant(ctx context.Context, tenantID valueobject.TenantID, reason string, operatorID uservo.UserID) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 验证操作权限
	if !tenant.OwnerID().Equals(operatorID) {
		return kernel.NewBusinessError(aggregate.CodeInsufficientPermissions, "only owner can deactivate tenant")
	}

	// 3. 停用租户
	return tenant.Deactivate(reason)
}
