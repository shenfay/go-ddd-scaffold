package ulid

import (
	"crypto/rand"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropy io.Reader
	once    sync.Once
)

// init 初始化熵源
func init() {
	once.Do(func() {
		entropy = &lockedReader{r: rand.Reader}
	})
}

// lockedReader 线程安全的读取器
type lockedReader struct {
	r  io.Reader
	mu sync.Mutex
}

func (l *lockedReader) Read(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.r.Read(p)
}

// GenerateUserID 生成用户 ID
// 格式：纯 ULID（不带前缀）
func GenerateUserID() string {
	t := time.Now()
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

// GenerateTokenID 生成 Token ID
// 格式：tok_{ulid}
func GenerateTokenID() string {
	t := time.Now()
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return fmt.Sprintf("tok_%s", id.String())
}

// GenerateSessionID 生成会话 ID
// 格式：ses_{ulid}
func GenerateSessionID() string {
	t := time.Now()
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return fmt.Sprintf("ses_%s", id.String())
}

// GenerateAuditLogID 生成审计日志 ID
// 格式：aud_{ulid}
func GenerateAuditLogID() string {
	t := time.Now()
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return fmt.Sprintf("aud_%s", id.String())
}

// ParseULID 解析 ULID 字符串
func ParseULID(id string) (ulid.ULID, error) {
	// 移除前缀（如果有）
	var ulidStr string
	if len(id) > 4 && id[3] == '_' {
		ulidStr = id[4:]
	} else {
		ulidStr = id
	}

	return ulid.Parse(ulidStr)
}

// GetTimestamp 从 ID 中提取时间戳
func GetTimestamp(id string) (time.Time, error) {
	parsed, err := ParseULID(id)
	if err != nil {
		return time.Time{}, err
	}

	t := time.UnixMilli(int64(parsed.Time()))
	return t, nil
}
