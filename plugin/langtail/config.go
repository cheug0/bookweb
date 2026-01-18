package langtail

// Config 长尾词插件配置
type Config struct {
	Enabled        bool   `json:"enabled"`
	FetchCycleDays int    `json:"fetch_cycle_days"` // 抓取周期（天）
	CacheSeconds   int    `json:"cache_seconds"`    // 缓存时间（秒）
	RoutePattern   string `json:"route_pattern"`    // 路由模式，如 /langtail/:lid
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:        false,
		FetchCycleDays: 7,
		CacheSeconds:   3600,
		RoutePattern:   "/langtail/:lid",
	}
}

// ParseConfig 从 map 解析配置
func ParseConfig(cfg map[string]interface{}) *Config {
	c := DefaultConfig()

	if enabled, ok := cfg["enabled"].(bool); ok {
		c.Enabled = enabled
	}
	if days, ok := cfg["fetch_cycle_days"].(float64); ok {
		c.FetchCycleDays = int(days)
	}
	if seconds, ok := cfg["cache_seconds"].(float64); ok {
		c.CacheSeconds = int(seconds)
	}
	if pattern, ok := cfg["route_pattern"].(string); ok && pattern != "" {
		c.RoutePattern = pattern
	}

	return c
}
