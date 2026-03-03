// Package i18n 提供多语言文件加载功能
package i18n

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// LocaleData 语言数据映射
type LocaleData map[string]interface{}

// Loader 语言文件加载器
type Loader struct {
	mu       sync.RWMutex
	data     map[string]LocaleData // language -> key-value map
	embedded bool
	fs       *embed.FS
}

var (
	defaultLoader *Loader
	once          sync.Once
)

// NewLoader 创建新的加载器
func NewLoader() *Loader {
	return &Loader{
		data: make(map[string]LocaleData),
	}
}

// GetDefaultLoader 获取默认加载器
func GetDefaultLoader() *Loader {
	once.Do(func() {
		defaultLoader = NewLoader()
		// 默认加载内置翻译
		defaultLoader.loadEmbedded()
	})
	return defaultLoader
}

// loadEmbedded 加载内置翻译（从代码中的map）
func (l *Loader) loadEmbedded() {
	// 这里可以保留一些内置的fallback翻译
	// 实际生产环境会从YAML文件加载
}

// LoadFromFile 从文件加载翻译
func (l *Loader) LoadFromFile(lang string, filePath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read locale file: %w", err)
	}

	var locale LocaleData
	if err := yaml.Unmarshal(data, &locale); err != nil {
		return fmt.Errorf("failed to parse locale file: %w", err)
	}

	// 扁平化嵌套的YAML结构为点分隔的key
	flatData := l.flatten("", locale)
	l.data[lang] = flatData

	return nil
}

// LoadFromEmbedFS 从嵌入的文件系统加载
func (l *Loader) LoadFromEmbedFS(fs *embed.FS, lang string, filePath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := fs.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read embedded locale file: %w", err)
	}

	var locale LocaleData
	if err := yaml.Unmarshal(data, &locale); err != nil {
		return fmt.Errorf("failed to parse embedded locale file: %w", err)
	}

	flatData := l.flatten("", locale)
	l.data[lang] = flatData

	return nil
}

// LoadFromDirectory 加载目录下所有语言文件
func (l *Loader) LoadFromDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read locale directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		// 从文件名提取语言代码 (zh-CN.yaml -> zh-CN)
		lang := strings.TrimSuffix(file.Name(), ".yaml")
		filePath := filepath.Join(dirPath, file.Name())

		if err := l.LoadFromFile(lang, filePath); err != nil {
			return fmt.Errorf("failed to load locale %s: %w", lang, err)
		}
	}

	return nil
}

// flatten 将嵌套的map扁平化为点分隔的key
func (l *Loader) flatten(prefix string, data map[string]interface{}) LocaleData {
	result := make(LocaleData)

	for key, value := range data {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]interface{}); ok {
			// 递归处理嵌套结构
			nestedResult := l.flatten(newKey, nested)
			for k, v := range nestedResult {
				result[k] = v
			}
		} else {
			result[newKey] = value
		}
	}

	return result
}

// Get 获取翻译值
func (l *Loader) Get(lang string, key string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if langData, ok := l.data[lang]; ok {
		if value, ok := langData[key]; ok {
			if str, ok := value.(string); ok {
				return str, true
			}
		}
	}

	// 尝试默认语言
	if lang != "zh-CN" {
		if langData, ok := l.data["zh-CN"]; ok {
			if value, ok := langData[key]; ok {
				if str, ok := value.(string); ok {
					return str, true
				}
			}
		}
	}

	return "", false
}

// GetAll 获取所有翻译数据
func (l *Loader) GetAll(lang string) LocaleData {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if langData, ok := l.data[lang]; ok {
		return langData
	}

	return make(LocaleData)
}

// HasLanguage 检查是否已加载指定语言
func (l *Loader) HasLanguage(lang string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	_, ok := l.data[lang]
	return ok
}

// Languages 获取已加载的语言列表
func (l *Loader) Languages() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	langs := make([]string, 0, len(l.data))
	for lang := range l.data {
		langs = append(langs, lang)
	}

	return langs
}
