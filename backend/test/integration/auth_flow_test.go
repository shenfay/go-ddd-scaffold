package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/suite"
)

// TestAuthIntegration 运行集成测试套件
func TestAuthIntegration(t *testing.T) {
	suite.Run(t, new(AuthIntegrationSuite))
}

// TestRegisterFlow 测试用户注册流程
func (s *AuthIntegrationSuite) TestRegisterFlow() {
	// 准备注册数据
	registerData := map[string]string{
		"email":    "test@example.com",
		"password": "Password123!",
	}
	body, _ := json.Marshal(registerData)

	// 发送注册请求
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// 验证响应
	s.Equal(http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	s.Contains(response, "user")
	s.Contains(response, "access_token")
	s.Contains(response, "refresh_token")
	s.Contains(response, "expires_in")

	// 验证用户数据
	user := response["user"].(map[string]interface{})
	s.Equal("test@example.com", user["email"])
	s.False(user["email_verified"].(bool))
}

// TestLoginFlow 测试用户登录流程
func (s *AuthIntegrationSuite) TestLoginFlow() {
	// 先注册用户
	s.createUser("login_test@example.com", "Password123!")

	// 准备登录数据
	loginData := map[string]string{
		"email":    "login_test@example.com",
		"password": "Password123!",
	}
	body, _ := json.Marshal(loginData)

	// 发送登录请求
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// 验证响应
	s.Equal(http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	s.Contains(response, "user")
	s.Contains(response, "access_token")
	s.Contains(response, "refresh_token")
}

// TestLoginWithWrongPassword 测试登录失败场景
func (s *AuthIntegrationSuite) TestLoginWithWrongPassword() {
	// 先注册用户
	s.createUser("wrong_password@example.com", "Password123!")

	// 准备错误的密码
	loginData := map[string]string{
		"email":    "wrong_password@example.com",
		"password": "WrongPassword",
	}
	body, _ := json.Marshal(loginData)

	// 发送登录请求
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// 验证响应（应该是 401）
	s.Equal(http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	s.Equal("INVALID_CREDENTIALS", response["code"])
}

// TestRefreshToken 测试刷新 Token
func (s *AuthIntegrationSuite) TestRefreshToken() {
	// 先注册用户并获取 refresh token
	s.createUser("refresh_test@example.com", "Password123!")
	
	// 登录获取 refresh token
	loginData := map[string]string{
		"email":    "refresh_test@example.com",
		"password": "Password123!",
	}
	body, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	refreshToken := loginResponse["refresh_token"].(string)

	// 使用 refresh token 刷新
	refreshData := map[string]string{
		"refresh_token": refreshToken,
	}
	body, _ = json.Marshal(refreshData)

	req, _ = http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// 验证响应
	s.Equal(http.StatusOK, w.Code)

	var refreshResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &refreshResponse)

	s.Contains(refreshResponse, "access_token")
	s.Contains(refreshResponse, "refresh_token")
	s.NotEqual(loginResponse["access_token"], refreshResponse["access_token"])
}

// TestAccountLockAfterFailedAttempts 测试账户锁定机制
func (s *AuthIntegrationSuite) TestAccountLockAfterFailedAttempts() {
	email := "lock_test@example.com"
	password := "Password123!"

	// 先注册用户
	s.createUser(email, password)

	// 连续失败 5 次
	for i := 0; i < 5; i++ {
		loginData := map[string]string{
			"email":    email,
			"password": "WrongPassword",
		}
		body, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusUnauthorized, w.Code)
	}

	// 第 6 次尝试应该返回账户锁定
	loginData := map[string]string{
		"email":    email,
		"password": "WrongPassword",
	}
	body, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// 验证账户被锁定
	s.Equal(http.StatusLocked, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	s.Equal("ACCOUNT_LOCKED", response["code"])
}

// createUser 辅助函数：创建测试用户
func (s *AuthIntegrationSuite) createUser(email, password string) {
	registerData := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(registerData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)
}
