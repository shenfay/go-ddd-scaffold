// Package valueobject_test 值对象测试
package valueobject_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go-ddd-scaffold/internal/domain/shared/valueobject"
)

// TestUserID_Creation 测试 UserID 创建
func TestUserID_Creation(t *testing.T) {
	id := uuid.New()
	userID := valueobject.NewUserID(id)

	assert.Equal(t, id, userID.Value())
	assert.Equal(t, id.String(), userID.String())
}

// TestUserID_Parse 测试 UserID 解析
func TestUserID_Parse(t *testing.T) {
	idStr := "550e8400-e29b-41d4-a716-446655440000"
	
	userID, err := valueobject.ParseUserID(idStr)
	
	assert.NoError(t, err)
	assert.Equal(t, idStr, userID.String())
}

// TestUserID_Parse_Invalid 测试无效 UserID 解析
func TestUserID_Parse_Invalid(t *testing.T) {
	invalidStr := "not-a-valid-uuid"
	
	_, err := valueobject.ParseUserID(invalidStr)
	
	assert.Error(t, err)
}

// TestUserID_Equals 测试 UserID 相等性
func TestUserID_Equals(t *testing.T) {
	id := uuid.New()
	userID1 := valueobject.NewUserID(id)
	userID2 := valueobject.NewUserID(id)
	userID3 := valueobject.NewUserID(uuid.New())

	assert.True(t, userID1.Equals(userID2))
	assert.False(t, userID1.Equals(userID3))
}

// TestUserID_JSON 测试 UserID JSON 序列化
func TestUserID_JSON(t *testing.T) {
	id := uuid.New()
	userID := valueobject.NewUserID(id)

	// Marshal
	data, err := json.Marshal(userID)
	assert.NoError(t, err)
	assert.Equal(t, `"`+id.String()+`"`, string(data))

	// Unmarshal
	var parsedUserID valueobject.UserID
	err = json.Unmarshal(data, &parsedUserID)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

// TestTenantID_Creation 测试 TenantID 创建
func TestTenantID_Creation(t *testing.T) {
	id := uuid.New()
	tenantID := valueobject.NewTenantID(id)

	assert.Equal(t, id, tenantID.Value())
	assert.Equal(t, id.String(), tenantID.String())
}

// TestTenantID_Equals 测试 TenantID 相等性
func TestTenantID_Equals(t *testing.T) {
	id := uuid.New()
	tenantID1 := valueobject.NewTenantID(id)
	tenantID2 := valueobject.NewTenantID(id)
	tenantID3 := valueobject.NewTenantID(uuid.New())

	assert.True(t, tenantID1.Equals(tenantID2))
	assert.False(t, tenantID1.Equals(tenantID3))
}

// TestEmail_ValidCreation 测试有效 Email 创建
func TestEmail_ValidCreation(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"admin@test.org",
		"TEST@EXAMPLE.COM", // 应该转为小写
	}

	for _, emailStr := range validEmails {
		email, err := valueobject.NewEmail(emailStr)
		assert.NoError(t, err)
		assert.Equal(t, emailStr, email.Value())
	}
}

// TestEmail_Normalization 测试 Email 标准化
func TestEmail_Normalization(t *testing.T) {
	email, err := valueobject.NewEmail("  TEST@Example.COM  ")
	
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", email.Value())
}

// TestEmail_InvalidFormat 测试无效 Email 格式
func TestEmail_InvalidFormat(t *testing.T) {
	invalidEmails := []string{
		"invalid",
		"@example.com",
		"user@",
		"user@.com",
		"user@com",
	}

	for _, emailStr := range invalidEmails {
		_, err := valueobject.NewEmail(emailStr)
		assert.Error(t, err, "Email %s should be invalid", emailStr)
	}
}

// TestEmail_Equals 测试 Email 相等性
func TestEmail_Equals(t *testing.T) {
	email1, _ := valueobject.NewEmail("test@example.com")
	email2, _ := valueobject.NewEmail("test@example.com")
	email3, _ := valueobject.NewEmail("other@example.com")

	assert.True(t, email1.Equals(email2))
	assert.False(t, email1.Equals(email3))
}

// TestEmail_JSON 测试 Email JSON 序列化
func TestEmail_JSON(t *testing.T) {
	email, _ := valueobject.NewEmail("test@example.com")

	// Marshal
	data, err := json.Marshal(email)
	assert.NoError(t, err)
	assert.Equal(t, `"test@example.com"`, string(data))

	// Unmarshal
	var parsedEmail valueobject.Email
	err = json.Unmarshal(data, &parsedEmail)
	assert.NoError(t, err)
	assert.Equal(t, email, parsedEmail)
}

// TestEmail_InvalidJSON 测试无效 Email JSON
func TestEmail_InvalidJSON(t *testing.T) {
	jsonData := []byte(`"invalid-email"`)
	
	var email valueobject.Email
	err := json.Unmarshal(jsonData, &email)
	
	assert.Error(t, err)
}

// TestIsValidEmail 测试邮箱验证函数
func TestIsValidEmail(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name+tag@domain.co.uk",
		"admin@test.org",
	}

	invalidEmails := []string{
		"invalid",
		"",
		"@example.com",
		"user@",
	}

	for _, email := range validEmails {
		assert.True(t, valueobject.IsValidEmail(email), "%s should be valid", email)
	}

	for _, email := range invalidEmails {
		assert.False(t, valueobject.IsValidEmail(email), "%s should be invalid", email)
	}
}
