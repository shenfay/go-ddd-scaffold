// Package docs API 文档
package docs

import "github.com/swaggo/swag"

// @title Go DDD Scaffold API
// @version 1.0
// @description Go DDD Scaffold 通用脚手架 API

// @host localhost:8080
// @BasePath /api

// SwaggerInfo 存储 Swagger 信息
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{"http"},
	Title:            "Go DDD Scaffold API",
	Description:      "Go DDD Scaffold 通用脚手架 API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

const docTemplate = ``
