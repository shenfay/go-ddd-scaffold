// @description Go DDD Scaffold Worker - 异步任务处理器
//
// 负责处理 Asynq 任务队列中的异步任务
// 包括：领域事件处理、邮件发送、数据同步等
//
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
package main

import (
	"github.com/shenfay/go-ddd-scaffold/cmd/app"
)

func main() {
	// 初始化基础设施
	infra, logger, cleanup := app.Initialize("worker")
	defer cleanup()

	// 创建并运行 Worker
	processor := NewProcessor(infra, logger)
	processor.Run()
}
