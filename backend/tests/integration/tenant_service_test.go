package integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/application/tenant/dto"
	"go-ddd-scaffold/internal/application/tenant/service"
	user_entity "go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/infrastructure/auth"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/repo"
	"go-ddd-scaffold/internal/infrastructure/transaction"
)

// TestTenantService_CreateTenant_WithUnitOfWork 测试使用 UnitOfWork 创建租户的原子性
func TestTenantService_CreateTenant_WithUnitOfWork(t *testing.T) {
	db := setupTestDB(t)

	uow := transaction.NewGormUnitOfWork(db)
	tenantRepo := repo.NewTenantDAORepository(db)
	memberRepo := repo.NewTenantMemberDAORepository(db)

	casbinService, err := auth.NewCasbinServiceForTest(db)
	assert.NoError(t, err)

	tenantSvc := service.NewTenantService(tenantRepo, memberRepo, casbinService, uow)

	t.Run("成功创建租户及成员关系", func(t *testing.T) {
		ctx := context.Background()
		ownerID := uuid.New()
		tenantName := "测试租户"
		tenantDesc := "这是一个测试租户"

		req := &dto.CreateTenantRequest{
			Name:        tenantName,
			Description: &tenantDesc,
			MaxMembers:  10,
		}
		createdTenant, err := tenantSvc.CreateTenant(ctx, req, ownerID)

		assert.NoError(t, err)
		assert.NotNil(t, createdTenant)
		assert.Equal(t, tenantName, createdTenant.Name)

		var tenantCount int64
		err = db.Model(&user_entity.Tenant{}).Where("id = ?", createdTenant.ID).Count(&tenantCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), tenantCount)

		var memberCount int64
		err = db.Model(&user_entity.TenantMember{}).Where("tenant_id = ?", createdTenant.ID).Count(&memberCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), memberCount)

		var member user_entity.TenantMember
		err = db.Where("tenant_id = ? AND user_id = ?", createdTenant.ID, ownerID).First(&member).Error
		assert.NoError(t, err)
		assert.Equal(t, user_entity.RoleOwner, member.Role)
	})
}

// TestTenantService_GetUserTenants 测试查询用户租户
func TestTenantService_GetUserTenants(t *testing.T) {
	db := setupTestDB(t)

	uow := transaction.NewGormUnitOfWork(db)
	tenantRepo := repo.NewTenantDAORepository(db)
	memberRepo := repo.NewTenantMemberDAORepository(db)
	casbinService, _ := auth.NewCasbinServiceForTest(db)

	tenantSvc := service.NewTenantService(tenantRepo, memberRepo, casbinService, uow)

	t.Run("获取用户的所有租户", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()

		tenantDesc1 := "描述 1"
		tenantDesc2 := "描述 2"
		req1 := &dto.CreateTenantRequest{
			Name:        "租户 1",
			Description: &tenantDesc1,
			MaxMembers:  10,
		}
		req2 := &dto.CreateTenantRequest{
			Name:        "租户 2",
			Description: &tenantDesc2,
			MaxMembers:  10,
		}
		tenant1, _ := tenantSvc.CreateTenant(ctx, req1, userID)
		tenant2, _ := tenantSvc.CreateTenant(ctx, req2, userID)

		tenants, err := tenantSvc.GetUserTenants(ctx, userID)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(tenants), 2)

		foundTenant1 := false
		foundTenant2 := false
		for _, tn := range tenants {
			if tn.ID == tenant1.ID {
				foundTenant1 = true
				assert.Equal(t, string(user_entity.RoleOwner), tn.Role)
			}
			if tn.ID == tenant2.ID {
				foundTenant2 = true
				assert.Equal(t, string(user_entity.RoleOwner), tn.Role)
			}
		}

		assert.True(t, foundTenant1, "应该找到租户 1")
		assert.True(t, foundTenant2, "应该找到租户 2")
	})
}

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=test_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("无法连接测试数据库，跳过测试")
	}

	err = db.AutoMigrate(&user_entity.Tenant{}, &user_entity.TenantMember{})
	if err != nil {
		t.Fatalf("迁移失败：%v", err)
	}

	db.Exec("DELETE FROM tenant_members")
	db.Exec("DELETE FROM tenants")

	return db
}
