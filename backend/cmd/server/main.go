// @title go-ddd-scaffold API
// @version 1.0
// @description DDD 架构脚手架项目 API文档
// @description 
// @description ## 认证模块
// @description - POST /api/auth/register 用户注册
// @description - POST /api/auth/login 用户登录
// @description - POST /api/auth/logout 用户登出
// @description 
// @description ## 用户管理
// @description - GET /api/users/:id 获取用户信息
// @description - PUT /api/users/:id 更新用户信息
// @description 
// @description ## 个人资料
// @description - GET /api/users/info 获取当前用户信息
// @description - PUT /api/users/profile 更新个人资料
// @description 
// @description ## 租户管理
// @description - POST /api/tenants 创建租户
// @description - GET /api/tenants/my-tenants 获取用户的租户列表
// @description 
// @BasePath/api
package main

import (
	"log"

	_ "go-ddd-scaffold/docs" // Swagger 文档
	"go-ddd-scaffold/internal/infrastructure/app"
)

func main() {
	// 创建应用实例
	application, err := app.NewApplication()
	if err != nil {
		log.Fatalf("创建应用失败：%v", err)
	}
	defer application.Close()

	// 运行应用
	if err := application.Run(); err != nil {
		log.Fatalf("应用运行失败：%v", err)
	}
}
