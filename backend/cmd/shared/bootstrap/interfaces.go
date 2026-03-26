package bootstrap

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// Module 基础接口：所有模块必须实现
type Module interface {
	Name() string
}

// HTTPModule 可选能力：支持 HTTP 路由注册
type HTTPModule interface {
	Module
	RegisterHTTP(group *gin.RouterGroup)
}

// EventModule 可选能力：支持事件订阅注册
type EventModule interface {
	Module
	RegisterSubscriptions(bus kernel.EventBus)
}

// GRPCModule 可选能力：支持 gRPC 服务注册（预留）
// type GRPCModule interface {
//     Module
//     RegisterGRPC(srv *grpc.Server)
// }
