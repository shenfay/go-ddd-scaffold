package idgen

// Node ID 生成器节点接口（保持与 infrastructure 一致）
type Node interface {
	Generate() (int64, error)
}

// Generator ID 生成器端口
type Generator interface {
	// Generate 生成唯一 ID
	Generate() (int64, error)
}
