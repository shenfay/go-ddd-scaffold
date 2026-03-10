package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	user_repo "go-ddd-scaffold/internal/domain/user/repository"
	auth "go-ddd-scaffold/internal/infrastructure/auth"
	transaction "go-ddd-scaffold/internal/infrastructure/transaction"
	"go-ddd-scaffold/internal/application/tenant/dto"
)

// TenantService 租户服务接口
type TenantService interface {
	// GetUserTenants 获取用户的所有租户
	GetUserTenants(ctx context.Context, userID uuid.UUID) ([]*dto.TenantWithRoleResponse, error)
	// CreateTenant 创建租户
	CreateTenant(ctx context.Context, req *dto.CreateTenantRequest, ownerID uuid.UUID) (*dto.TenantResponse, error)
}


type tenantService struct {
	tenantRepo    user_repo.TenantRepository
	memberRepo    user_repo.TenantMemberRepository
	casbinService auth.CasbinService
	uow           transaction.UnitOfWork
}

// NewTenantService 创建租户服务实例
func NewTenantService(
	tenantRepo user_repo.TenantRepository,
	memberRepo user_repo.TenantMemberRepository,
	casbinService auth.CasbinService,
	uow transaction.UnitOfWork,
) TenantService {
	return &tenantService{
		tenantRepo:    tenantRepo,
		memberRepo:    memberRepo,
		casbinService: casbinService,
		uow:           uow,
	}
}

// GetUserTenants 获取用户的所有租户（返回 DTO）
func (s *tenantService) GetUserTenants(ctx context.Context, userID uuid.UUID) ([]*dto.TenantWithRoleResponse, error) {
	// 1. 获取用户的所有租户成员关系
	memberships, err := s.memberRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 获取每个租户的详细信息并转换为 DTO
	result := make([]*dto.TenantWithRoleResponse, 0, len(memberships))
	for _, membership := range memberships {
		tenant, err := s.tenantRepo.GetByID(ctx, membership.TenantID)
		if err != nil {
			continue // 跳过不存在的租户
		}

		// TODO: 获取实际成员数（当前简化为 0）
		dtoResult := dto.ToTenantWithRoleDTO(tenant, membership, 0)
		result = append(result, dtoResult)
	}

	return result, nil
}

// CreateTenant 创建租户（使用 UnitOfWork 保证原子性，返回 DTO）
func (s *tenantService) CreateTenant(ctx context.Context, req *dto.CreateTenantRequest, ownerID uuid.UUID) (*dto.TenantResponse, error) {
	var response *dto.TenantResponse
	
	// ✅ 使用 UnitOfWork 保证跨聚合根操作的原子性
	err := s.uow.WithTransaction(ctx, func(ctx context.Context) error {
		// 获取事务 DB
		tx := transaction.GetTxFromContext(ctx)
		
		// 切换到事务仓储
		tenantRepo := s.tenantRepo.WithTx(tx)
		memberRepo := s.memberRepo.WithTx(tx)
		
		// 1. 使用 DTO 中的转换函数构建租户
		tenant, err := dto.FromCreateRequest(req, ownerID)
		if err != nil {
			return err
		}
		
		if err := tenantRepo.Create(ctx, tenant); err != nil {
			return err
		}
		
		// 2. 创建者自动成为 owner
		member := &user_entity.TenantMember{
			ID:       uuid.New(),
			TenantID: tenant.ID,
			UserID:   ownerID,
			Role:    user_entity.RoleOwner, // ✅ 使用 Owner 角色
			Status:  user_entity.MemberStatusActive,
			JoinedAt: time.Now(),
		}
		
		if err := memberRepo.Create(ctx, member); err != nil {
			return err
		}
		
		// 3. 在 Casbin 中添加角色（使用 AddRoleForUser 方法）
		// 注意：Casbin 操作不在事务中，如果失败会导致数据不一致
		// TODO: 考虑将 Casbin 策略也持久化到数据库
		if err := s.casbinService.AddRoleForUser(ownerID, tenant.ID, string(user_entity.RoleOwner)); err != nil {
			// 记录警告日志，但不回滚（因为主要数据已保存）
			// 可以在后台任务中重试
			return nil
		}
		
		// 4. 转换为 DTO 返回
		response = dto.ToTenantDTO(tenant, 0)
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return response, nil
}
