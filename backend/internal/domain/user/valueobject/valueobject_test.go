package vo_test

import (
	"testing"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// UserName 测试
// ============================================================================

func TestNewUserName_Success(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     string
	}{
		{"有效用户名 - 字母", "testuser", "testuser"},
		{"有效用户名 - 数字", "user123", "user123"},
		{"有效用户名 - 下划线", "test_user", "test_user"},
		{"有效用户名 - 连字符", "test-user", "test-user"},
		{"有效用户名 - 混合", "Test_User-123", "Test_User-123"},
		{"有效用户名 - 最小长度", "abc", "abc"},
		{"有效用户名 - 最大长度", "abcdefghijklmnopqrstuvwxyz1234567890ABCD", "abcdefghijklmnopqrstuvwxyz1234567890ABCD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userName, err := vo.NewUserName(tt.username)
			assert.NoError(t, err)
			assert.NotNil(t, userName)
			assert.Equal(t, tt.want, userName.Value())
		})
	}
}

func TestNewUserName_Failure(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		expectError string
	}{
		{"空用户名", "", "username cannot be empty"},
		{"空格用户名", "   ", "username cannot be empty"},
		{"长度不足", "ab", "username must be at least 3 characters long"},
		{"长度超标", "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNO", "username cannot exceed 50 characters"},
		{"包含特殊字符@", "test@user", "username can only contain letters, numbers, underscores and hyphens"},
		{"包含特殊字符!", "test!user", "username can only contain letters, numbers, underscores and hyphens"},
		{"包含空格", "test user", "username can only contain letters, numbers, underscores and hyphens"},
		{"包含中文字符", "测试用户", "username can only contain letters, numbers, underscores and hyphens"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userName, err := vo.NewUserName(tt.username)
			assert.Error(t, err)
			assert.Nil(t, userName)
			if err != nil {
				assert.Contains(t, err.Error(), tt.expectError)
			}
		})
	}
}

func TestUserName_Equals(t *testing.T) {
	tests := []struct {
		name     string
		value1   string
		value2   string
		expected bool
	}{
		{"相同用户名", "testuser", "testuser", true},
		{"不同大小写", "TestUser", "testuser", true},
		{"不同用户名", "user1", "user2", false},
		{"与 nil 比较", "testuser", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			un1, _ := vo.NewUserName(tt.value1)

			var un2 *vo.UserName
			if tt.value2 != "" {
				un2, _ = vo.NewUserName(tt.value2)
			}

			assert.Equal(t, tt.expected, un1.Equals(un2))
		})
	}
}

// ============================================================================
// Email 测试
// ============================================================================

func TestNewEmail_Success(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  string
	}{
		{"标准邮箱", "test@example.com", "test@example.com"},
		{"大写字母邮箱", "TEST@EXAMPLE.COM", "test@example.com"},
		{"带点号邮箱", "test.user@example.com", "test.user@example.com"},
		{"带加号邮箱", "test+tag@example.com", "test+tag@example.com"},
		{"带下划线邮箱", "test_user@example.com", "test_user@example.com"},
		{"子域名邮箱", "test@mail.example.com", "test@mail.example.com"},
		{"短后缀邮箱", "test@example.cn", "test@example.cn"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := vo.NewEmail(tt.email)
			assert.NoError(t, err)
			assert.NotNil(t, email)
			assert.Equal(t, tt.want, email.Value())
		})
	}
}

func TestNewEmail_Failure(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError string
	}{
		{"空邮箱", "", "email cannot be empty"},
		{"空格邮箱", "   ", "email cannot be empty"},
		{"缺少@符号", "testexample.com", "invalid email format"},
		{"缺少域名", "test@", "invalid email format"},
		{"缺少用户名", "@example.com", "invalid email format"},
		{"缺少后缀", "test@example", "invalid email format"},
		{"无效字符", "test!@example.com", "invalid email format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := vo.NewEmail(tt.email)
			assert.Error(t, err)
			assert.Nil(t, email)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	tests := []struct {
		name     string
		value1   string
		value2   string
		expected bool
	}{
		{"相同邮箱", "test@example.com", "test@example.com", true},
		{"不同大小写", "Test@Example.Com", "test@example.com", true},
		{"不同邮箱", "user1@example.com", "user2@example.com", false},
		{"与 nil 比较", "test@example.com", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e1, _ := vo.NewEmail(tt.value1)

			var e2 *vo.Email
			if tt.value2 != "" {
				e2, _ = vo.NewEmail(tt.value2)
			}

			assert.Equal(t, tt.expected, e1.Equals(e2))
		})
	}
}

// ============================================================================
// HashedPassword 测试
// ============================================================================

func TestNewHashedPassword(t *testing.T) {
	tests := []struct {
		name   string
		hashed string
		want   string
	}{
		{"bcrypt 哈希", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"},
		{"空哈希", "", ""},
		{"普通字符串", "hashed_password", "hashed_password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := vo.NewHashedPassword(tt.hashed)
			assert.NotNil(t, hp)
			assert.Equal(t, tt.want, hp.Value())
		})
	}
}

// ============================================================================
// UserStatus 测试
// ============================================================================

func TestUserStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status vo.UserStatus
		want   string
	}{
		{"Pending", vo.UserStatusPending, "pending"},
		{"Active", vo.UserStatusActive, "active"},
		{"Inactive", vo.UserStatusInactive, "inactive"},
		{"Locked", vo.UserStatusLocked, "locked"},
		{"Unknown", vo.UserStatus(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.String())
		})
	}
}

// ============================================================================
// UserGender 测试
// ============================================================================

func TestUserGender_String(t *testing.T) {
	tests := []struct {
		name   string
		gender vo.UserGender
		want   string
	}{
		{"Unknown", vo.UserGenderUnknown, "unknown"},
		{"Male", vo.UserGenderMale, "male"},
		{"Female", vo.UserGenderFemale, "female"},
		{"Other", vo.UserGenderOther, "other"},
		{"Invalid", vo.UserGender(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.gender.String())
		})
	}
}

// ============================================================================
// UserProfile 测试
// ============================================================================

func TestNewUserProfile_Success(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		firstName   string
		lastName    string
		gender      vo.UserGender
		phoneNumber string
		avatarURL   string
	}{
		{"完整资料", "John Doe", "John", "Doe", vo.UserGenderMale, "1234567890", "https://example.com/avatar.jpg"},
		{"最小资料", "", "", "", vo.UserGenderUnknown, "", ""},
		{"仅显示名称", "John", "", "", vo.UserGenderUnknown, "", ""},
		{"仅姓名", "", "John", "Doe", vo.UserGenderUnknown, "", ""},
		{"带手机号", "", "", "", vo.UserGenderUnknown, "+86-138-0000-0000", ""},
		{"带头像", "", "", "", vo.UserGenderUnknown, "", "https://example.com/avatar.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := vo.NewUserProfile(
				tt.displayName,
				tt.firstName,
				tt.lastName,
				tt.gender,
				tt.phoneNumber,
				tt.avatarURL,
			)
			assert.NoError(t, err)
			assert.NotNil(t, profile)
			assert.Equal(t, tt.displayName, profile.DisplayName())
			assert.Equal(t, tt.firstName, profile.FirstName())
			assert.Equal(t, tt.lastName, profile.LastName())
			assert.Equal(t, tt.gender, profile.Gender())
			assert.Equal(t, tt.phoneNumber, profile.PhoneNumber())
			assert.Equal(t, tt.avatarURL, profile.AvatarURL())
		})
	}
}

func TestNewUserProfile_Failure(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		firstName   string
		lastName    string
		gender      vo.UserGender
		phoneNumber string
		avatarURL   string
		expectField string
	}{
		{
			"显示名称超长",
			"abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ1234",
			"", "", vo.UserGenderUnknown, "", "",
			"display_name",
		},
		{
			"名字超长",
			"",
			"abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWX",
			"", vo.UserGenderUnknown, "", "",
			"first_name",
		},
		{
			"姓氏超长",
			"", "",
			"abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWX",
			vo.UserGenderUnknown, "", "",
			"last_name",
		},
		{
			"手机号超长",
			"", "", "", vo.UserGenderUnknown,
			"123456789012345678901",
			"",
			"phone_number",
		},
		{
			"头像 URL 超长",
			"", "", "", vo.UserGenderUnknown, "",
			"https://example.com/" + string(make([]byte, 500)),
			"avatar_url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := vo.NewUserProfile(
				tt.displayName,
				tt.firstName,
				tt.lastName,
				tt.gender,
				tt.phoneNumber,
				tt.avatarURL,
			)
			assert.Error(t, err)
			assert.Nil(t, profile)
			if err != nil {
				valErr, ok := err.(*common.ValidationError)
				if assert.True(t, ok) {
					assert.Equal(t, tt.expectField, valErr.Field)
				}
			}
		})
	}
}

func TestUserProfile_FullName(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		firstName   string
		lastName    string
		want        string
	}{
		{"有姓名", "John Doe", "John", "Doe", "John Doe"},
		{"无姓", "John", "John", "", "John"},
		{"无名", "Doe", "", "Doe", "Doe"},
		{"无姓名", "John Doe", "", "", "John Doe"},
		{"全空", "", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, _ := vo.NewUserProfile(
				tt.displayName,
				tt.firstName,
				tt.lastName,
				vo.UserGenderUnknown,
				"",
				"",
			)
			assert.Equal(t, tt.want, profile.FullName())
		})
	}
}

func TestUserProfile_Equals(t *testing.T) {
	profile1, _ := vo.NewUserProfile("John", "John", "Doe", vo.UserGenderMale, "1234567890", "https://example.com/avatar.jpg")
	profile2, _ := vo.NewUserProfile("John", "John", "Doe", vo.UserGenderMale, "1234567890", "https://example.com/avatar.jpg")
	profile3, _ := vo.NewUserProfile("Jane", "Jane", "Doe", vo.UserGenderFemale, "0987654321", "https://example.com/jane.jpg")

	assert.True(t, profile1.Equals(profile2))
	assert.False(t, profile1.Equals(profile3))
	assert.False(t, profile1.Equals(nil))
}

func TestUserProfile_UpdateDisplayName(t *testing.T) {
	profile, _ := vo.NewUserProfile("John", "John", "Doe", vo.UserGenderMale, "", "")

	newProfile, err := profile.UpdateDisplayName("John Updated")
	assert.NoError(t, err)
	assert.NotNil(t, newProfile)
	assert.Equal(t, "John Updated", newProfile.DisplayName())
	assert.Equal(t, "John", profile.DisplayName()) // 原对象不变
}

func TestUserProfile_UpdateName(t *testing.T) {
	profile, _ := vo.NewUserProfile("John Doe", "John", "Doe", vo.UserGenderMale, "", "")

	newProfile, err := profile.UpdateName("John", "Updated")
	assert.NoError(t, err)
	assert.NotNil(t, newProfile)
	assert.Equal(t, "John", newProfile.FirstName())
	assert.Equal(t, "Updated", newProfile.LastName())
	assert.Equal(t, "John", profile.FirstName()) // 原对象不变
}

func TestUserProfile_UpdateGender(t *testing.T) {
	profile, _ := vo.NewUserProfile("", "", "", vo.UserGenderUnknown, "", "")

	newProfile, err := profile.UpdateGender(vo.UserGenderMale)
	assert.NoError(t, err)
	assert.NotNil(t, newProfile)
	assert.Equal(t, vo.UserGenderMale, newProfile.Gender())
	assert.Equal(t, vo.UserGenderUnknown, profile.Gender()) // 原对象不变
}

func TestUserProfile_UpdatePhoneNumber(t *testing.T) {
	profile, _ := vo.NewUserProfile("", "", "", vo.UserGenderUnknown, "", "")

	newProfile, err := profile.UpdatePhoneNumber("13800000000")
	assert.NoError(t, err)
	assert.NotNil(t, newProfile)
	assert.Equal(t, "13800000000", newProfile.PhoneNumber())
	assert.Equal(t, "", profile.PhoneNumber()) // 原对象不变
}

func TestUserProfile_UpdateAvatarURL(t *testing.T) {
	profile, _ := vo.NewUserProfile("", "", "", vo.UserGenderUnknown, "", "")

	newProfile, err := profile.UpdateAvatarURL("https://example.com/new-avatar.jpg")
	assert.NoError(t, err)
	assert.NotNil(t, newProfile)
	assert.Equal(t, "https://example.com/new-avatar.jpg", newProfile.AvatarURL())
	assert.Equal(t, "", profile.AvatarURL()) // 原对象不变
}
