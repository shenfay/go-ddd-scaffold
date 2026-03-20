package tenant

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
)

// TenantService 租户领域服务
type TenantService struct {
	tenantRepo TenantRepository
	userRepo   repository.UserRepository
}

// NewTenantService 创建租户服务
func NewTenantService(tenantRepo TenantRepository, userRepo repository.UserRepository) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
	}
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, code, name string, ownerID vo.UserID) (*Tenant, error) {
	// 1. 检查租户编码是否已存在
	if _, err := s.tenantRepo.FindByCode(ctx, code); err == nil {
		return nil, kernel.NewBusinessError(kernel.CodeTenantCodeExists, "tenant code already exists")
	}

	// 2. 检查所有者是否存在
	if _, err := s.userRepo.FindByID(ctx, ownerID); err != nil {
		return nil, kernel.NewBusinessError(kernel.CodeUserNotFound, "owner user not found")
	}

	// 3. 创建租户
	tenant, err := NewTenant(code, name, ownerID)
	if err != nil {
		return nil, err
	}

	// 4. 保存租户
	if err := s.tenantRepo.Save(ctx, tenant); err != nil {
		return nil, err
	}

	// 5. 添加所有者为成员
	member := &TenantMember{
		UserID:   ownerID,
		TenantID: tenant.ID().(TenantID),
		Role:     TenantRoleOwner,
		JoinedAt: time.Now().Format(time.RFC3339),
	}
	if err := s.tenantRepo.AddMember(ctx, tenant.ID().(TenantID), member); err != nil {
		return nil, err
	}

	return tenant, nil
}

// AddMember 添加成员到租户
func (s *TenantService) AddMember(ctx context.Context, tenantID TenantID, userID, addedBy vo.UserID, role TenantRole) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 检查租户状态
	if !tenant.IsActive() {
		return kernel.NewBusinessError(kernel.CodeTenantNotActive, "tenant is not active")
	}

	// 3. 检查用户是否存在
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return kernel.NewBusinessError(kernel.CodeUserNotFound, "user not found")
	}

	// 4. 检查用户是否已经是成员
	if _, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, userID); err == nil {
		return kernel.NewBusinessError(kernel.CodeAlreadyMember, "user is already a member of this tenant")
	}

	// 5. 检查成员数量限制
	members, err := s.tenantRepo.FindMembers(ctx, tenantID)
	if err != nil {
		return err
	}

	if !tenant.CanAddMember(len(members)) {
		return kernel.NewBusinessError(kernel.CodeTenantMaxMembersReached, "tenant has reached maximum member limit")
	}

	// 6. 检查操作权限（只有 Owner 和 Admin 可以添加成员）
	addedByMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, addedBy)
	if err != nil {
		return kernel.NewBusinessError(kernel.CodeOperatorNotMember, "operator is not a member of this tenant")
	}

	if addedByMember.Role != TenantRoleOwner && addedByMember.Role != TenantRoleAdmin {
		return kernel.NewBusinessError(kernel.CodeInsufficientPermissions, "insufficient permissions to add members")
	}

	// 7. 添加成员
	member := &TenantMember{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		JoinedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.tenantRepo.AddMember(ctx, tenantID, member); err != nil {
		return err
	}

	// 8. 发布领域事件
	event := NewTenantMemberAddedEvent(tenantID, userID, addedBy, role)
	tenant.ApplyEvent(event)

	return s.tenantRepo.Save(ctx, tenant)
}

// RemoveMember 从租户移除成员
func (s *TenantService) RemoveMember(ctx context.Context, tenantID TenantID, userID, removedBy vo.UserID) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 检查成员是否存在
	member, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, userID)
	if err != nil {
		return kernel.NewBusinessError(kernel.CodeNotTenantMember, "user is not a member of this tenant")
	}

	// 3. 检查操作权限
	removedByMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, removedBy)
	if err != nil {
		return kernel.NewBusinessError(kernel.CodeOperatorNotMember, "operator is not a member of this tenant")
	}

	// 4. 权限检查规则
	if member.Role == TenantRoleOwner {
		// 不能移除所有者
		return kernel.NewBusinessError(kernel.CodeCannotRemoveOwner, "cannot remove tenant owner")
	}

	if removedByMember.Role != TenantRoleOwner && removedByMember.Role != TenantRoleAdmin {
		return kernel.NewBusinessError(kernel.CodeInsufficientPermissions, "insufficient permissions to remove members")
	}

	// Admin 不能移除其他 Admin
	if member.Role == TenantRoleAdmin && removedByMember.Role != TenantRoleOwner {
		return kernel.NewBusinessError(kernel.CodeCannotRemoveAdmin, "only owner can remove admin members")
	}

	// 5. 移除成员
	if err := s.tenantRepo.RemoveMember(ctx, tenantID, userID); err != nil {
		return err
	}

	// 6. 发布领域事件
	event := NewTenantMemberRemovedEvent(tenantID, userID, removedBy)
	tenant.ApplyEvent(event)

	return s.tenantRepo.Save(ctx, tenant)
}

// ChangeMemberRole 更改成员角色
func (s *TenantService) ChangeMemberRole(ctx context.Context, tenantID TenantID, userID, changedBy vo.UserID, newRole TenantRole) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 检查成员是否存在
	member, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, userID)
	if err != nil {
		return kernel.NewBusinessError(kernel.CodeNotTenantMember, "user is not a member of this tenant")
	}

	// 3. 检查操作权限
	changedByMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, changedBy)
	if err != nil {
		return kernel.NewBusinessError(kernel.CodeOperatorNotMember, "operator is not a member of this tenant")
	}

	// 4. 权限检查规则
	if changedByMember.Role != TenantRoleOwner {
		return kernel.NewBusinessError(kernel.CodeInsufficientPermissions, "only owner can change member roles")
	}

	// 5. 不能更改自己的角色
	if userID.Equals(changedBy) {
		return kernel.NewBusinessError(kernel.CodeCannotChangeOwnRole, "cannot change your own role")
	}

	// 6. 保存旧角色用于事件
	oldRole := member.Role

	// 7. 更新角色（这里简化处理，实际应该通过仓储更新）
	member.Role = newRole

	// 8. 发布领域事件
	event := NewTenantMemberRoleChangedEvent(tenantID, userID, changedBy, oldRole, newRole)
	tenant.ApplyEvent(event)

	return s.tenantRepo.Save(ctx, tenant)
}

// TransferOwnership 转移租户所有权
func (s *TenantService) TransferOwnership(ctx context.Context, tenantID TenantID, newOwnerID, currentOwnerID vo.UserID) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 验证当前所有者是租户所有者
	if !tenant.OwnerID().Equals(currentOwnerID) {
		return kernel.NewBusinessError(kernel.CodeNotTenantOwner, "only current owner can transfer ownership")
	}

	// 3. 检查新所有者是否是租户成员
	newOwnerMember, err := s.tenantRepo.FindMemberByUserID(ctx, tenantID, newOwnerID)
	if err != nil {
		return kernel.NewBusinessError(kernel.CodeNotTenantMember, "new owner must be a member of the tenant")
	}

	// 4. 更新所有者
	tenant.ownerID = newOwnerID
	tenant.updatedAt = time.Now()
	tenant.IncrementVersion()

	// 5. 更新原所有者为 Admin
	// 更新新所有者为 Owner
	newOwnerMember.Role = TenantRoleOwner

	// 6. 保存
	return s.tenantRepo.Save(ctx, tenant)
}

// DeactivateTenant 停用租户
func (s *TenantService) DeactivateTenant(ctx context.Context, tenantID TenantID, reason string, operatorID vo.UserID) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return kernel.ErrAggregateNotFound
	}

	// 2. 验证操作权限
	if !tenant.OwnerID().Equals(operatorID) {
		return kernel.NewBusinessError(kernel.CodeInsufficientPermissions, "only owner can deactivate tenant")
	}

	// 3. 停用租户
	return tenant.Deactivate(reason)
}
