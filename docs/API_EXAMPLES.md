# API 使用示例

本文档提供完整的 API 使用示例，包括注册、登录、用户信息管理等常见场景。

---

## 快速开始

### 1. 健康检查

**请求**:
```bash
curl -X GET http://localhost:8080/api/health
```

**响应**:
```json
{
  "code": "Success",
  "message": "成功",
  "data": {
    "status": "healthy",
    "timestamp": "2026-03-08T09:00:00Z"
  }
}
```

---

## 认证相关

### 2. 用户注册

**请求**:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!",
    "nickname": "测试用户"
  }'
```

**请求参数**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码（最少 8 位，包含数字和字母） |
| nickname | string | 是 | 昵称（支持中文） |

**响应**:
```json
{
  "code": "Success",
  "message": "注册成功",
  "data": {
    "user_id": "341d982d-7deb-4906-a7c7-3e6afa43396c",
    "email": "test@example.com",
    "nickname": "测试用户",
    "created_at": "2026-03-08T09:00:00Z"
  }
}
```

**错误响应 - 邮箱已存在**:
```json
{
  "code": "BadRequest",
  "message": "邮箱已被注册",
  "details": {
    "field": "email"
  }
}
```

**错误响应 - 密码格式不正确**:
```json
{
  "code": "BadRequest",
  "message": "密码格式不正确",
  "details": {
    "field": "password"
  }
}
```

---

### 3. 用户登录

**请求**:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!"
  }'
```

**响应**:
```json
{
  "code": "Success",
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 7200,
    "user_info": {
      "user_id": "341d982d-7deb-4906-a7c7-3e6afa43396c",
      "email": "test@example.com",
      "nickname": "测试用户",
      "avatar": null
    }
  }
}
```

**错误响应 - 密码错误**:
```json
{
  "code": "InvalidCredential",
  "message": "用户名或密码错误",
  "details": {
    "field": "password"
  }
}
```

**错误响应 - 用户不存在**:
```json
{
  "code": "NotFound",
  "message": "用户不存在",
  "details": {
    "field": "email"
  }
}
```

---

### 4. 用户登出

**请求**:
```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应**:
```json
{
  "code": "Success",
  "message": "登出成功"
}
```

**错误响应 - Token 无效**:
```json
{
  "code": "Unauthorized",
  "message": "无效的访问令牌"
}
```

---

## 用户信息管理

### 5. 获取当前用户信息

**请求**:
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**响应**:
```json
{
  "code": "Success",
  "message": "成功",
  "data": {
    "user_id": "341d982d-7deb-4906-a7c7-3e6afa43396c",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": null,
    "phone": null,
    "bio": null,
    "status": "active",
    "created_at": "2026-03-08T09:00:00Z",
    "updated_at": "2026-03-08T09:00:00Z"
  }
}
```

---

### 6. 更新用户资料

**请求**:
```bash
curl -X PUT http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "新昵称",
    "phone": "13800138000",
    "bio": "个人简介"
  }'
```

**请求参数**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| nickname | string | 否 | 新昵称 |
| phone | string | 否 | 手机号 |
| bio | string | 否 | 个人简介 |

**响应**:
```json
{
  "code": "Success",
  "message": "更新成功",
  "data": {
    "user_id": "341d982d-7deb-4906-a7c7-3e6afa43396c",
    "email": "test@example.com",
    "nickname": "新昵称",
    "phone": "13800138000",
    "bio": "个人简介",
    "status": "active",
    "updated_at": "2026-03-08T10:00:00Z"
  }
}
```

**错误响应 - 昵称格式不正确**:
```json
{
  "code": "BadRequest",
  "message": "昵称格式不正确",
  "details": {
    "field": "nickname"
  }
}
```

---

## 错误处理

### 常见错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|------------|------|
| Success | 200 | 成功 |
| BadRequest | 400 | 请求参数错误 |
| Unauthorized | 401 | 未授权 |
| Forbidden | 403 | 禁止访问 |
| NotFound | 404 | 资源不存在 |
| InvalidCredential | 401 | 凭证无效 |
| InternalError | 500 | 服务器内部错误 |

### 错误响应格式

```json
{
  "code": "ErrorCode",
  "message": "错误描述信息",
  "details": {
    "field": "错误字段"
  }
}
```

---

## 完整示例

### 使用 JavaScript Fetch

```javascript
// 用户注册
const register = async () => {
  const response = await fetch('http://localhost:8080/api/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: 'test@example.com',
      password: 'Password123!',
      nickname: '测试用户',
    }),
  });
  
  const data = await response.json();
  if (response.ok) {
    console.log('注册成功:', data.data);
  } else {
    console.error('注册失败:', data.message);
  }
};

// 用户登录
const login = async () => {
  const response = await fetch('http://localhost:8080/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: 'test@example.com',
      password: 'Password123!',
    }),
  });
  
  const data = await response.json();
  if (response.ok) {
    localStorage.setItem('accessToken', data.data.access_token);
    console.log('登录成功');
  } else {
    console.error('登录失败:', data.message);
  }
};

// 获取用户信息
const getProfile = async () => {
  const token = localStorage.getItem('accessToken');
  const response = await fetch('http://localhost:8080/api/user/profile', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  const data = await response.json();
  if (response.ok) {
    console.log('用户信息:', data.data);
  } else {
    console.error('获取失败:', data.message);
  }
};
```

---

### 使用 Python Requests

```python
import requests

BASE_URL = 'http://localhost:8080/api'

# 用户注册
def register():
    response = requests.post(f'{BASE_URL}/auth/register', json={
        'email': 'test@example.com',
        'password': 'Password123!',
        'nickname': '测试用户',
    })
    
    data = response.json()
    if response.status_code == 200:
        print('注册成功:', data['data'])
    else:
        print('注册失败:', data['message'])

# 用户登录
def login():
    response = requests.post(f'{BASE_URL}/auth/login', json={
        'email': 'test@example.com',
        'password': 'Password123!',
    })
    
    data = response.json()
    if response.status_code == 200:
        access_token = data['data']['access_token']
        print('登录成功，Token:', access_token)
        return access_token
    else:
        print('登录失败:', data['message'])
        return None

# 获取用户信息
def get_profile(access_token):
    headers = {'Authorization': f'Bearer {access_token}'}
    response = requests.get(f'{BASE_URL}/user/profile', headers=headers)
    
    data = response.json()
    if response.status_code == 200:
        print('用户信息:', data['data'])
    else:
        print('获取失败:', data['message'])

# 使用示例
if __name__ == '__main__':
    register()
    token = login()
    if token:
        get_profile(token)
```

---

## Postman 集合

导入以下 JSON 到 Postman：

```json
{
  "info": {
    "name": "DDD Scaffold API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Auth",
      "item": [
        {
          "name": "Register",
          "request": {
            "method": "POST",
            "header": [{"key": "Content-Type", "value": "application/json"}],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"Password123!\",\n  \"nickname\": \"测试用户\"\n}"
            },
            "url": {
              "raw": "{{baseUrl}}/auth/register",
              "host": ["{{baseUrl}}"],
              "path": ["auth", "register"]
            }
          }
        },
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "header": [{"key": "Content-Type", "value": "application/json"}],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"Password123!\"\n}"
            },
            "url": {
              "raw": "{{baseUrl}}/auth/login",
              "host": ["{{baseUrl}}"],
              "path": ["auth", "login"]
            }
          }
        }
      ]
    }
  ],
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080/api"
    }
  ]
}
```

---

## 测试工具

### 集成测试脚本

项目包含完整的集成测试脚本：

```bash
# 运行集成测试
cd backend
./scripts/integration_test.sh
```

测试内容包括：
- ✅ 健康检查
- ✅ 用户注册
- ✅ 用户登录
- ✅ 获取用户信息
- ✅ 更新用户资料
- ✅ 用户登出
- ✅ 错误密码验证

---

## 最佳实践

### 1. Token 管理

```javascript
// 自动刷新 Token
const refreshToken = async () => {
  const refreshToken = localStorage.getItem('refreshToken');
  const response = await fetch('http://localhost:8080/api/auth/refresh', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
  
  const data = await response.json();
  if (response.ok) {
    localStorage.setItem('accessToken', data.data.access_token);
    return data.data.access_token;
  } else {
    // Token 刷新失败，跳转到登录页
    window.location.href = '/login';
  }
};
```

### 2. 错误处理

```javascript
// 统一错误处理
const handleApiError = (error) => {
  switch (error.code) {
    case 'Unauthorized':
      // Token 过期，刷新或重新登录
      break;
    case 'BadRequest':
      // 显示表单验证错误
      break;
    case 'InternalError':
      // 显示系统错误提示
      break;
    default:
      console.error('未知错误:', error.message);
  }
};
```

### 3. 请求拦截器

```javascript
// Axios 拦截器示例
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Token 过期处理
      return refreshToken().then(newToken => {
        error.config.headers.Authorization = `Bearer ${newToken}`;
        return axios.request(error.config);
      });
    }
    return Promise.reject(error);
  }
);
```

---

## 相关文档

- [Swagger API 文档](./api_swagger_guide.md)
- [错误处理指南](./error_handling_guide.md)
- [部署指南](./deployment_guide.md)
