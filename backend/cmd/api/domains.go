package main

import (
	// 用户领域 - 触发 init() 自动注册路由
	_ "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/user"
	// 新增领域时，在这里添加导入：
	// _ "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/order"
	// _ "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/product"
)
