package factory_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-ddd-scaffold/internal/domain/tenant/factory"
)

func TestNewTenantBuilder_DefaultValues(t *testing.T) {
	ownerID := uuid.New()
	builder := factory.NewTenantBuilder(ownerID, "测试租户")

	assert.NotNil(t, builder)
}

func TestTenantBuilder_WithDescription(t *testing.T) {
	ownerID := uuid.New()
	description := "这是一个测试租户"

	tenant := factory.NewTenantBuilder(ownerID, "测试租户").
		WithDescription(description).
		MustBuild()

	assert.Equal(t, description, tenant.Description)
}

func TestTenantBuilder_WithMaxMembers(t *testing.T) {
	ownerID := uuid.New()
	maxMembers := 50

	tenant := factory.NewTenantBuilder(ownerID, "测试租户").
		WithMaxMembers(maxMembers).
		MustBuild()

	assert.Equal(t, maxMembers, tenant.MaxMembers)
}

func TestTenantBuilder_WithExpiredAt(t *testing.T) {
	ownerID := uuid.New()
	expiredAt := time.Now().AddDate(2, 0, 0)

	tenant := factory.NewTenantBuilder(ownerID, "测试租户").
		WithExpiredAt(expiredAt).
		MustBuild()

	assert.Equal(t, expiredAt.Year(), tenant.ExpiredAt.Year())
	assert.Equal(t, expiredAt.Month(), tenant.ExpiredAt.Month())
	assert.Equal(t, expiredAt.Day(), tenant.ExpiredAt.Day())
}

func TestTenantBuilder_Build_Success(t *testing.T) {
	ownerID := uuid.New()
	name := "成功租户"
	description := "描述信息"
	maxMembers := 100
	expiredAt := time.Now().AddDate(1, 6, 0)

	tenant, err := factory.NewTenantBuilder(ownerID, name).
		WithDescription(description).
		WithMaxMembers(maxMembers).
		WithExpiredAt(expiredAt).
		Build()

	require.NoError(t, err)
	require.NotNil(t, tenant)
	assert.Equal(t, name, tenant.Name)
	assert.Equal(t, description, tenant.Description)
	assert.Equal(t, maxMembers, tenant.MaxMembers)
	assert.True(t, tenant.ExpiredAt.After(time.Now()))
	assert.NotEqual(t, uuid.Nil, tenant.ID)
}

func TestTenantBuilder_Build_EmptyName(t *testing.T) {
	ownerID := uuid.New()

	tenant, err := factory.NewTenantBuilder(ownerID, "").Build()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant name is required")
	assert.Nil(t, tenant)
}

func TestTenantBuilder_Build_NilOwnerID(t *testing.T) {
	nilUUID := uuid.Nil

	tenant, err := factory.NewTenantBuilder(nilUUID, "测试租户").Build()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "owner ID is required")
	assert.Nil(t, tenant)
}

func TestTenantBuilder_Build_ExpiredTimePast(t *testing.T) {
	ownerID := uuid.New()
	pastTime := time.Now().AddDate(-1, 0, 0)

	tenant, err := factory.NewTenantBuilder(ownerID, "测试租户").
		WithExpiredAt(pastTime).
		Build()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired_at must be in the future")
	assert.Nil(t, tenant)
}

func TestTenantBuilder_MustBuild_Panic(t *testing.T) {
	ownerID := uuid.New()

	assert.Panics(t, func() {
		factory.NewTenantBuilder(ownerID, "").MustBuild()
	})
}

func TestTenantBuilder_ChainCall(t *testing.T) {
	ownerID := uuid.New()
	description := "完整链式调用测试"
	maxMembers := 75
	expiredAt := time.Now().AddDate(1, 3,15)

	tenant := factory.NewTenantBuilder(ownerID, "链式租户").
		WithDescription(description).
		WithMaxMembers(maxMembers).
		WithExpiredAt(expiredAt).
		MustBuild()

	assert.NotNil(t, tenant)
	assert.Equal(t, "链式租户", tenant.Name)
	assert.Equal(t, description, tenant.Description)
	assert.Equal(t, maxMembers, tenant.MaxMembers)
	assert.True(t, tenant.ExpiredAt.After(time.Now()))
	assert.NotEqual(t, uuid.Nil, tenant.ID)
}
