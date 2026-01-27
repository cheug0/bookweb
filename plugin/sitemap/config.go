// config.go (sitemap)
// Sitemap 配置
// 定义 Sitemap 的输出目录及生成频率等配置
package sitemap

// Config sitemap 插件配置
type Config struct {
	Enabled         bool    `json:"enabled"`
	OutputPath      string  `json:"output_path"`
	RegenerateHours int     `json:"regenerate_hours"`
	MaxURLsPerFile  int     `json:"max_urls_per_file"`
	IncludeBooks    bool    `json:"include_books"`
	IncludeChapters bool    `json:"include_chapters"`
	ChangeFreq      string  `json:"changefreq"`
	Priority        float64 `json:"priority"`
}

var globalConfig *Config

// SetConfig 设置全局配置
func SetConfig(cfg *Config) {
	globalConfig = cfg
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// ParseConfig 解析配置
func ParseConfig(cfg map[string]interface{}) *Config {
	config := &Config{
		Enabled:         true,
		OutputPath:      "static/sitemap",
		RegenerateHours: 24,
		MaxURLsPerFile:  50000,
		IncludeBooks:    true,
		IncludeChapters: false,
		ChangeFreq:      "daily",
		Priority:        0.8,
	}

	if v, ok := cfg["enabled"].(bool); ok {
		config.Enabled = v
	}
	if v, ok := cfg["output_path"].(string); ok {
		config.OutputPath = v
	}
	if v, ok := cfg["regenerate_hours"].(float64); ok {
		config.RegenerateHours = int(v)
	}
	if v, ok := cfg["max_urls_per_file"].(float64); ok {
		config.MaxURLsPerFile = int(v)
	}
	if v, ok := cfg["include_books"].(bool); ok {
		config.IncludeBooks = v
	}
	if v, ok := cfg["include_chapters"].(bool); ok {
		config.IncludeChapters = v
	}
	if v, ok := cfg["changefreq"].(string); ok {
		config.ChangeFreq = v
	}
	if v, ok := cfg["priority"].(float64); ok {
		config.Priority = v
	}

	return config
}
