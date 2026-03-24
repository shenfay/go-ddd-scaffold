package idgen

// Generator ID 生成器端口
type Generator interface {
	// Generate 生成唯一 ID
	Generate() (int64, error)
}
