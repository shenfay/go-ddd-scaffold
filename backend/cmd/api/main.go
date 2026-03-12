package main

import (
	"net/http"

	httpiface "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	apperrors "github.com/shenfay/go-ddd-scaffold/shared/errors"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("Starting API server...")

	// 初始化组件
	errorMapper := apperrors.NewErrorMapper()
	handler := httpiface.NewHandler(errorMapper)
	router := httpiface.GetRouter()

	// 构建依赖注入容器
	deps := httpiface.NewDependencies(handler)

	// 构建完整路由（自动触发所有已注册的领域路由）
	ginEngine := router.Build(deps)

	// 启动服务器
	logger.Info("Server listening on :8080")
	if err := ginEngine.Run(":8080"); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
