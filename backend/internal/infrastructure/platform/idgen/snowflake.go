package idgen

import (
	idgen "github.com/yitter/idgenerator-go/idgen"
)

// Initialize 初始化全局 ID 生成器
// 只需在应用启动时调用一次
func Initialize(workerId uint64, workerIdBitLength uint8) {
	options := idgen.NewIdGeneratorOptions(uint16(workerId))
	options.WorkerIdBitLength = workerIdBitLength
	idgen.SetIdGenerator(options)
}

// Generate 生成雪花 ID
// 线程安全，可直接并发使用
func Generate() int64 {
	return idgen.NextId()
}

// ParseSnowflakeID 解析雪花 ID
// 返回时间戳、WorkerId、序列号
func ParseSnowflakeID(id int64) (timestamp int64, workerId uint64, sequence int64) {
	// yitter/idgenerator-go 的 ID 结构：
	// timestamp(41 位) | workerId(默认 6 位，最大 22 位) | sequence(默认 12 位)
	// 使用 ExtractTime 获取时间戳
	t := idgen.ExtractTime(id)
	timestamp = t.UnixMilli()

	// WorkerId 和 Sequence 需要根据实际配置计算
	// 默认配置：WorkerIdBitLength=6, SequenceBitLength=12
	const (
		workerIdBits  = uint8(6)
		sequenceBits  = uint8(12)
		workerIdShift = sequenceBits
		sequenceMask  = int64(-1 ^ (-1 << sequenceBits))
	)

	workerId = uint64((id >> workerIdShift) & ((1 << workerIdBits) - 1))
	sequence = id & sequenceMask
	return
}
