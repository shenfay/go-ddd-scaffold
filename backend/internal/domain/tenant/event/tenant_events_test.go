package event_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go-ddd-scaffold/internal/domain/tenant/event"
)

func TestTenantCreatedEvent_Fields(t *testing.T) {
	tenantID := uuid.New()
	ownerID := uuid.New()
	tenantName := "测试租户"

	e := event.NewTenantCreatedEvent(tenantID, ownerID, tenantName)

	assert.Equal(t, "TenantCreated", e.EventType)
	assert.Equal(t, tenantID, e.TenantID)
	assert.Equal(t, ownerID, e.OwnerID)
	assert.Equal(t, tenantName, e.TenantName)
	assert.NotEmpty(t, e.EventID)
	assert.Equal(t, tenantID, e.AggregateID)
	assert.Equal(t, 1, e.Version)
}

func TestMemberJoinedEvent_Fields(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	role := "member"

	e := event.NewMemberJoinedEvent(tenantID, userID, role)

	assert.Equal(t, "MemberJoined", e.EventType)
	assert.Equal(t, tenantID, e.TenantID)
	assert.Equal(t, userID, e.UserID)
	assert.Equal(t, role, e.Role)
	assert.NotEmpty(t, e.EventID)
	assert.Equal(t, 1, e.Version)
}

func TestMemberLeftEvent_Fields(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()

	e := event.NewMemberLeftEvent(tenantID, userID)

	assert.Equal(t, "MemberLeft", e.EventType)
	assert.Equal(t, tenantID, e.TenantID)
	assert.Equal(t, userID, e.UserID)
	assert.NotEmpty(t, e.EventID)
	assert.Equal(t, 1, e.Version)
}

func TestRoleChangedEvent_Fields(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	oldRole := "member"
	newRole := "admin"

	e := event.NewRoleChangedEvent(tenantID, userID, oldRole, newRole)

	assert.Equal(t, "RoleChanged", e.EventType)
	assert.Equal(t, tenantID, e.TenantID)
	assert.Equal(t, userID, e.UserID)
	assert.Equal(t, oldRole, e.OldRole)
	assert.Equal(t, newRole, e.NewRole)
	assert.NotEmpty(t, e.EventID)
	assert.Equal(t, 1, e.Version)
}

func TestDomainEvent_Interface(t *testing.T) {
	tenantID := uuid.New()
	ownerID := uuid.New()
	e := event.NewTenantCreatedEvent(tenantID, ownerID, "测试")

	assert.Implements(t, (*interface{ GetEventType() string })(nil), e)
	assert.Implements(t, (*interface{ GetEventID() string })(nil), e)
	assert.Implements(t, (*interface{ GetAggregateID() uuid.UUID })(nil), e)
	assert.Implements(t, (*interface{ GetOccurredAt() time.Time })(nil), e)
	assert.Implements(t, (*interface{ GetVersion() int })(nil), e)
}

func TestEvent_JSONSerialization(t *testing.T) {
	tenantID := uuid.New()
	ownerID := uuid.New()
	e := event.NewTenantCreatedEvent(tenantID, ownerID, "测试租户")

	data, err := json.Marshal(e)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded event.TenantCreatedEvent
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, e.EventType, decoded.EventType)
	assert.Equal(t, e.TenantName, decoded.TenantName)
}
