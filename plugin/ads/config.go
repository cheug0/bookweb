// config.go (ads)
// 广告配置
// 定义广告插件的配置结构
package ads

// AdSlot 广告位结构
type AdSlot struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Enabled bool   `json:"enabled"`
}

// Config 广告插件配置
type Config struct {
	Enabled bool               `json:"enabled"`
	Slots   map[string]*AdSlot `json:"slots"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled: false,
		Slots:   make(map[string]*AdSlot),
	}
}

// ParseConfig 从 map 解析配置
func ParseConfig(cfg map[string]interface{}) *Config {
	config := DefaultConfig()

	if enabled, ok := cfg["enabled"].(bool); ok {
		config.Enabled = enabled
	}

	if slots, ok := cfg["slots"].(map[string]interface{}); ok {
		for slotID, slotData := range slots {
			if slotMap, ok := slotData.(map[string]interface{}); ok {
				slot := &AdSlot{}
				if name, ok := slotMap["name"].(string); ok {
					slot.Name = name
				}
				if content, ok := slotMap["content"].(string); ok {
					slot.Content = content
				}
				if enabled, ok := slotMap["enabled"].(bool); ok {
					slot.Enabled = enabled
				}
				config.Slots[slotID] = slot
			}
		}
	}

	return config
}
