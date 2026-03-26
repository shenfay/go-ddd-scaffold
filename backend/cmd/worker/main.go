package main

import (
	"github.com/shenfay/go-ddd-scaffold/cmd/pkg"
)

func main() {
	// 1. 初始化基础设施
	infra, logger, cleanup := pkg.InitInfrastructure("worker")
	defer cleanup()

	// 2. 创建并运行 Worker
	processor := NewProcessor(infra, logger)
	processor.Run()
}
