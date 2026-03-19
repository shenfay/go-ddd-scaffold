package helpers

import (
	"github.com/shenfay/go-ddd-scaffold/internal/domain/loginlog"
	"github.com/shenfay/go-ddd-scaffold/pkg/useragent"
)

// UAParserAdapter User-Agent 解析器适配器
// 将通用的 useragent.DeviceInfo 转换为 domain 特定的 loginlog.DeviceInfo
type UAParserAdapter struct {
	parser *useragent.Parser
}

// NewUAParserAdapter 创建 User-Agent 解析器适配器
func NewUAParserAdapter() *UAParserAdapter {
	return &UAParserAdapter{
		parser: useragent.NewParser(),
	}
}

// Parse 解析 User-Agent 字符串
func (a *UAParserAdapter) Parse(ua string) loginlog.DeviceInfo {
	info := a.parser.Parse(ua)
	return loginlog.DeviceInfo{
		DeviceType: info.DeviceType,
		OS:         info.OS,
		Browser:    info.Browser,
	}
}
