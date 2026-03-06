package service_test

import (
	"strings"
	"testing"
	"time"

	"go-ddd-scaffold/internal/domain/user/service"
)

func TestBcryptPasswordHasher_Hash(t *testing.T) {
	t.Run("成功哈希密码", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		plainPassword := "TestPassword123!"

		// Act
		hashed, err := hasher.Hash(plainPassword)

		// Assert
		if err != nil {
			t.Fatalf("Hash 失败：%v", err)
		}
		if hashed == "" {
			t.Fatal("Hash 结果不能为空")
		}
		if len(hashed) < 60 {
			t.Errorf("bcrypt hash 长度应该至少为 60，实际：%d", len(hashed))
		}
	})

	t.Run("相同密码产生不同 hash", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		password := "SamePassword123"

		// Act
		hash1, _ := hasher.Hash(password)
		hash2, _ := hasher.Hash(password)

		// Assert
		if hash1 == hash2 {
			t.Error("bcrypt 应该为相同密码生成不同的 hash（因为 salt）")
		}
	})

	t.Run("空密码哈希", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()

		// Act
		hashed, err := hasher.Hash("")

		// Assert
		if err != nil {
			t.Fatalf("空密码 Hash 不应该失败：%v", err)
		}
		if hashed == "" {
			t.Error("空密码也应该生成 hash")
		}
	})
}

func TestBcryptPasswordHasher_Verify(t *testing.T) {
	t.Run("验证正确密码", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		password := "CorrectPassword123!"
		hashed, _ := hasher.Hash(password)

		// Act
		valid := hasher.Verify(hashed, password)

		// Assert
		if !valid {
			t.Error("正确密码应该验证通过")
		}
	})

	t.Run("验证错误密码", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		password := "CorrectPassword123!"
		wrongPassword := "WrongPassword456!"
		hashed, _ := hasher.Hash(password)

		// Act
		valid := hasher.Verify(hashed, wrongPassword)

		// Assert
		if valid {
			t.Error("错误密码应该验证失败")
		}
	})

	t.Run("验证空密码", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		hashed, _ := hasher.Hash("")

		// Act
		validEmpty := hasher.Verify(hashed, "")
		validNonEmpty := hasher.Verify(hashed, "NotEmpty")

		// Assert
		if !validEmpty {
			t.Error("空密码应该验证通过")
		}
		if validNonEmpty {
			t.Error("非空密码应该验证失败")
		}
	})

	t.Run("验证无效 hash 格式", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		invalidHash := "$2a$10$invalid"
		password := "AnyPassword123"

		// Act
		valid := hasher.Verify(invalidHash, password)

		// Assert
		if valid {
			t.Error("无效 hash 格式应该验证失败")
		}
	})
}

func TestBcryptPasswordHasher_RoundTrip(t *testing.T) {
	t.Run("完整的 Hash 和 Verify 流程", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		testCases := []struct {
			name     string
			password string
		}{
			{"普通密码", "SimplePassword123"},
			{"复杂密码", "C0mpl3x!@#$%^&*()Pass"},
			{"长密码", "ThisIsAVeryLongPasswordWithManyCharacters123456789"},
			{"包含特殊字符", "Pass with spaces and 特殊字符！@#"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Act - Hash
				hashed, err := hasher.Hash(tc.password)
				if err != nil {
					t.Fatalf("Hash 失败：%v", err)
				}

				// Act - Verify 正确密码
				validCorrect := hasher.Verify(hashed, tc.password)
				
				// Act - Verify 错误密码
				validWrong := hasher.Verify(hashed, tc.password+"wrong")

				// Assert
				if !validCorrect {
					t.Error("正确密码验证失败")
				}
				if validWrong {
					t.Error("错误密码验证通过")
				}
			})
		}
	})
}

func TestNewDefaultBcryptPasswordHasher(t *testing.T) {
	t.Run("创建默认实例", func(t *testing.T) {
		// Arrange & Act
		hasher := service.NewDefaultBcryptPasswordHasher()

		// Assert
		if hasher == nil {
			t.Fatal("创建的 hasher 不应该为 nil")
		}
		
		// 验证返回类型
		var _ service.PasswordHasher = hasher // 编译时检查
	})
}

func TestNewBcryptPasswordHasher_CustomCost(t *testing.T) {
	t.Run("自定义 cost 值", func(t *testing.T) {
		// Arrange
		costs := []int{4, 8, 10, 12, 14}

		for _, cost := range costs {
			t.Run(string(rune(cost)), func(t *testing.T) {
				// Arrange
				hasher := service.NewBcryptPasswordHasher(cost)
				password := "TestPassword123"

				// Act
				hashed, err := hasher.Hash(password)
				if err != nil {
					t.Fatalf("Hash 失败：%v", err)
				}

				// Assert
				valid := hasher.Verify(hashed, password)
				if !valid {
					t.Error("密码验证失败")
				}
			})
		}
	})

	t.Run("cost 值影响 hash 时间", func(t *testing.T) {
		// Arrange
		hasher4 := service.NewBcryptPasswordHasher(4)
		hasher12 := service.NewBcryptPasswordHasher(12)
		password := "TestPassword123"

		// Act - 测量不同 cost 的时间
		start4 := time.Now()
		hasher4.Hash(password)
		time4 := time.Since(start4)

		start12 := time.Now()
		hasher12.Hash(password)
		time12 := time.Since(start12)

		// Assert
		if time12 <= time4 {
			t.Logf("警告：高 cost 没有明显更慢（cost=4: %v, cost=12: %v）", time4, time12)
		} else {
			t.Logf("✓ cost=12 确实更慢（cost=4: %v, cost=12: %v）", time4, time12)
		}
	})
}

// 性能测试
func BenchmarkBcryptPasswordHasher_Hash_Cost10(b *testing.B) {
	hasher := service.NewDefaultBcryptPasswordHasher()
	password := "BenchmarkPassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = hasher.Hash(password)
	}
}

func BenchmarkBcryptPasswordHasher_Verify_Cost10(b *testing.B) {
	hasher := service.NewDefaultBcryptPasswordHasher()
	password := "BenchmarkPassword123"
	hashed, _ := hasher.Hash(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hasher.Verify(hashed, password)
	}
}

// 安全性测试
func TestPasswordHasher_Security(t *testing.T) {
	t.Run("hash 不应该泄露密码信息", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		password := "SecretPassword123"

		// Act
		hashed, _ := hasher.Hash(password)

		// Assert
		if strings.Contains(hashed, password) {
			t.Error("hash 结果不应该包含原始密码")
		}
		if strings.Contains(hashed, "password") {
			t.Error("hash 结果不应该包含 'password' 字样")
		}
	})

	t.Run("不同密码的 hash 应该不同", func(t *testing.T) {
		// Arrange
		hasher := service.NewDefaultBcryptPasswordHasher()
		password1 := "Password1"
		password2 := "Password2"

		// Act
		hash1, _ := hasher.Hash(password1)
		hash2, _ := hasher.Hash(password2)

		// Assert
		if hash1 == hash2 {
			t.Error("不同密码的 hash 应该不同")
		}
	})
}
