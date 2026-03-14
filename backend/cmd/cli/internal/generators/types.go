package generators

import (
	"fmt"
)

// InitOptions 项目初始化选项
type InitOptions struct {
	ProjectName  string
	ModulePath   string
	Author       string
	Email        string
	License      string
	Template     string
	SkipFrontend bool
	WithDocker   bool
	WithK8s      bool
}

// EntityOptions 实体生成选项
type EntityOptions struct {
	Name          string
	Fields        string
	Methods       string
	Package       string
	WithVO        bool
	WithAggregate bool
}

// RepositoryOptions Repository 生成选项
type RepositoryOptions struct {
	Name      string
	Domain    string
	OutputDir string
}

// ServiceOptions Service 生成选项
type ServiceOptions struct {
	Name         string
	Type         string
	Methods      string
	Dependencies []string
}

// HandlerOptions Handler 生成选项
type HandlerOptions struct {
	Name   string
	Type   string
	Domain string
}

// DTOOptions DTO 生成选项
type DTOOptions struct {
	Name           string
	Type           string
	Fields         string
	WithValidation bool
}

// CleanGeneratedFiles 清理生成的文件
func CleanGeneratedFiles(path string, dryRun bool) error {
	fmt.Printf("Cleaning generated files in %s (dry-run: %v)\n", path, dryRun)
	// TODO: 实现清理逻辑
	return nil
}
