package snowflake

import (
	"errors"
	"sync"
	"time"

	"github.com/shenfay/go-ddd-scaffold/pkg/util"
)

// Snowflake ID 组成部分
const (
	epoch     = int64(1704067200000) // 2024-01-01 00:00:00 UTC
	nodeBits  = uint(10)             // 节点 ID 位数（支持 1024 个节点）
	stepBits  = uint(12)             // 序列号位数（每毫秒 4096 个 ID）
	nodeMax   = int64(-1 ^ (-1 << nodeBits))
	stepMask  = int64(-1 ^ (-1 << stepBits))
	timeShift = uint(nodeBits + stepBits)
	nodeShift = uint(stepBits)
)

var (
	// ErrInvalidNodeID 节点 ID 无效
	ErrInvalidNodeID = errors.New("node number must be between 0 and 1023")
	// ErrTimeExhausted 时间耗尽（同一毫秒内生成的 ID 过多）
	ErrTimeExhausted = errors.New("timestamp is exhausted")
)

// Node Snowflake 节点
type Node struct {
	mu        sync.Mutex
	timestamp int64
	node      int64
	step      int64
}

// NewNode 创建新的 Snowflake 节点
func NewNode(nodeID int64) (*Node, error) {
	if nodeID < 0 || nodeID > nodeMax {
		return nil, ErrInvalidNodeID
	}

	return &Node{
		timestamp: 0,
		node:      nodeID,
		step:      0,
	}, nil
}

// Generate 生成 Snowflake ID
func (n *Node) Generate() (int64, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := util.Now().TimestampMilli()

	if now == n.timestamp {
		// 同一毫秒内，序列号递增
		n.step = (n.step + 1) & stepMask
		if n.step == 0 {
			// 序列号溢出，等待下一毫秒
			for now <= n.timestamp {
				now = util.Now().TimestampMilli()
			}
		}
	} else {
		// 不同毫秒，重置序列号
		n.step = 0
	}

	n.timestamp = now

	// 组合 ID: timestamp(41 位) | nodeID(10 位) | step(12 位)
	result := (now-epoch)<<timeShift |
		(n.node << nodeShift) |
		n.step

	return result, nil
}

// ParseSnowflakeID 解析 Snowflake ID
func ParseSnowflakeID(id int64) (timestamp time.Time, nodeID, sequence int64) {
	timestamp = time.UnixMilli((id >> timeShift) + epoch).UTC()
	nodeID = (id >> nodeShift) & nodeMax
	sequence = id & stepMask
	return
}
