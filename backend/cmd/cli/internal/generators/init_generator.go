package generators

import (
	"fmt"
)

// InitGenerator 项目初始化生成器
type InitGenerator struct {
	opts InitOptions
}

// NewInitGenerator 创建初始化生成器
func NewInitGenerator(opts InitOptions) *InitGenerator {
	return &InitGenerator{opts: opts}
}

// Generate 生成项目结构
func (g *InitGenerator) Generate() error {
	fmt.Printf("Initializing project: %s\n", g.opts.ProjectName)
	fmt.Printf("Module path: %s\n", g.opts.ModulePath)
	fmt.Printf("Template: %s\n", g.opts.Template)

	// TODO: 实现项目初始化逻辑

	return nil
}
