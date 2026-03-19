package useragent

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceType string
	OS         string
	Browser    string
}

// Parser User-Agent解析器
type Parser struct{}

// NewParser 创建User-Agent解析器
func NewParser() *Parser {
	return &Parser{}
}

// Parse 解析User-Agent字符串
func (p *Parser) Parse(ua string) DeviceInfo {
	info := DeviceInfo{
		DeviceType: "desktop",
		OS:         "Unknown",
		Browser:    "Unknown",
	}

	if ua == "" {
		return info
	}

	// 简单的关键词匹配
	if contains(ua, "Mobile") {
		info.DeviceType = "mobile"
	} else if contains(ua, "Tablet") {
		info.DeviceType = "tablet"
	}

	if contains(ua, "Windows") {
		info.OS = "Windows"
	} else if contains(ua, "Macintosh") || contains(ua, "Mac OS") {
		info.OS = "macOS"
	} else if contains(ua, "Android") {
		info.OS = "Android"
	} else if contains(ua, "iOS") || contains(ua, "iPhone") || contains(ua, "iPad") {
		info.OS = "iOS"
	} else if contains(ua, "Linux") {
		info.OS = "Linux"
	}

	if contains(ua, "Chrome") && !contains(ua, "Edg") {
		info.Browser = "Chrome"
	} else if contains(ua, "Safari") && !contains(ua, "Chrome") {
		info.Browser = "Safari"
	} else if contains(ua, "Firefox") {
		info.Browser = "Firefox"
	} else if contains(ua, "Edg") {
		info.Browser = "Edge"
	}

	return info
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)
}

// findSubstring 查找子串
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
