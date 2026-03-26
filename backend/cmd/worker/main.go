package main

import (
	"github.com/shenfay/go-ddd-scaffold/cmd/shared"
)

func main() {
	// 初始化基础设施
	infra, logger, cleanup := shared.Initialize("worker")
	defer cleanup()

	// 2. 创建并运行 Worker
	processor := NewProcessor(infra, logger)
	processor.Run()
}
