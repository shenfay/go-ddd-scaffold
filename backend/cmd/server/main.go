// Package main 服务端主程序
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
		log.Fatalf("创建应用失败: %v", err)
	}
	defer application.Close()

	// 运行应用
	if err := application.Run(); err != nil {
		log.Fatalf("应用运行失败: %v", err)
	}
}
