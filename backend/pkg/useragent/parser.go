package useragent

import "strings"

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceType string
	OS         string
	Browser    string
}

// Parser User-Agent 解析器
type Parser struct{}

// NewParser 创建 User-Agent 解析器
func NewParser() *Parser {
	return &Parser{}
}

// Parse 解析 User-Agent 字符串
func (p *Parser) Parse(ua string) DeviceInfo {
	if ua == "" {
		return DeviceInfo{
			DeviceType: "desktop",
			OS:         "Unknown",
			Browser:    "Unknown",
		}
	}

	return DeviceInfo{
		DeviceType: p.detectDeviceType(ua),
		OS:         p.detectOS(ua),
		Browser:    p.detectBrowser(ua),
	}
}

// detectDeviceType 检测设备类型
func (p *Parser) detectDeviceType(ua string) string {
	if contains(ua, "Mobile") {
		return "mobile"
	}
	if contains(ua, "Tablet") {
		return "tablet"
	}
	return "desktop"
}

// detectOS 检测操作系统
func (p *Parser) detectOS(ua string) string {
	osMatchers := map[string]string{
		"Windows":   "Windows",
		"Macintosh": "macOS",
		"Mac OS":    "macOS",
		"Android":   "Android",
		"iOS":       "iOS",
		"iPhone":    "iOS",
		"iPad":      "iOS",
		"Linux":     "Linux",
	}

	for pattern, os := range osMatchers {
		if contains(ua, pattern) {
			return os
		}
	}
	return "Unknown"
}

// detectBrowser 检测浏览器
func (p *Parser) detectBrowser(ua string) string {
	// Edge 必须在 Chrome 之前检测（因为 Edge 也包含 Chrome）
	if contains(ua, "Edg") {
		return "Edge"
	}

	// Chrome 排除 Edge
	if contains(ua, "Chrome") && !contains(ua, "Edg") {
		return "Chrome"
	}

	// Safari 排除 Chrome
	if contains(ua, "Safari") && !contains(ua, "Chrome") {
		return "Safari"
	}

	if contains(ua, "Firefox") {
		return "Firefox"
	}

	return "Unknown"
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
