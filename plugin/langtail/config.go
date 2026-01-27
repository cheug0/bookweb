// config.go (langtail)
// 长尾词配置
// 定义长尾词插件的采集周期、显示数量等配置
package langtail

import "strconv"

// Config 长尾词插件配置
type Config struct {
	Enabled        bool   `json:"enabled"`
	FetchCycleDays int    `json:"fetch_cycle_days"` // 抓取周期（天）
	CacheSeconds   int    `json:"cache_seconds"`    // 缓存时间（秒）
	RoutePattern   string `json:"route_pattern"`    // 路由模式，如 /langtail/:lid
	ShowCount      int    `json:"show_count"`       // 显示数量
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:        false,
		FetchCycleDays: 7,
		CacheSeconds:   3600,
		RoutePattern:   "/langtail/:lid",
		ShowCount:      50,
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
	if count, ok := cfg["show_count"].(float64); ok {
		c.ShowCount = int(count)
	} else if count, ok := cfg["show_count"].(int); ok {
		c.ShowCount = count
	} else if countStr, ok := cfg["show_count"].(string); ok {
		// 尝试解析字符串
		if val, err := strconv.Atoi(countStr); err == nil {
			c.ShowCount = val
		}
	}

	return c
}
