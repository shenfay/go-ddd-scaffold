// @title Go DDD Scaffold API
// @version 1.0
// @description Go DDD Scaffold API 文档 - 基于 DDD 和 CQRS 的企业级脚手架
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在 Header 中输入：Bearer {token}

package main

import (
	"github.com/shenfay/go-ddd-scaffold/cmd/app"
)

func main() {
	// 初始化基础设施
	infra, logger, cleanup := app.Initialize("api")
	defer cleanup()

	// 创建并运行服务器
	server := NewServer(infra, logger)
	server.Run()
}
