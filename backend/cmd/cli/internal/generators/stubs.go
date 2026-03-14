package generators

import "fmt"

// EntityGenerator 实体生成器
type EntityGenerator struct{ opts EntityOptions }

func NewEntityGenerator(opts EntityOptions) *EntityGenerator { return &EntityGenerator{opts: opts} }
func (g *EntityGenerator) Generate() error {
	fmt.Printf("Generating entity: %s\n", g.opts.Name)
	return nil
}

// RepositoryGenerator Repository 生成器
type RepositoryGenerator struct{ opts RepositoryOptions }

func NewRepositoryGenerator(opts RepositoryOptions) *RepositoryGenerator {
	return &RepositoryGenerator{opts: opts}
}
func (g *RepositoryGenerator) Generate() error {
	fmt.Printf("Generating repository: %s\n", g.opts.Name)
	return nil
}

// ServiceGenerator Service 生成器
type ServiceGenerator struct{ opts ServiceOptions }

func NewServiceGenerator(opts ServiceOptions) *ServiceGenerator { return &ServiceGenerator{opts: opts} }
func (g *ServiceGenerator) Generate() error {
	fmt.Printf("Generating service: %s\n", g.opts.Name)
	return nil
}

// HandlerGenerator Handler 生成器
type HandlerGenerator struct{ opts HandlerOptions }

func NewHandlerGenerator(opts HandlerOptions) *HandlerGenerator { return &HandlerGenerator{opts: opts} }
func (g *HandlerGenerator) Generate() error {
	fmt.Printf("Generating handler: %s\n", g.opts.Name)
	return nil
}

// DTOGenerator DTO 生成器
type DTOGenerator struct{ opts DTOOptions }

func NewDTOGenerator(opts DTOOptions) *DTOGenerator { return &DTOGenerator{opts: opts} }
func (g *DTOGenerator) Generate() error             { fmt.Printf("Generating DTO: %s\n", g.opts.Name); return nil }
