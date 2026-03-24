package idgen

// GeneratorAdapter ID 生成器端口适配器
type GeneratorAdapter struct {
	node *Node
}

// NewGeneratorAdapter 创建 ID 生成器适配器
func NewGeneratorAdapter(node *Node) *GeneratorAdapter {
	return &GeneratorAdapter{
		node: node,
	}
}

// Generate 生成唯一 ID
func (a *GeneratorAdapter) Generate() (int64, error) {
	return a.node.Generate()
}
