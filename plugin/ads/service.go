// service.go (ads)
// 广告服务逻辑
// 提供获取广告内容的具体实现
package ads

import (
	"html/template"
	"sync"
)

var (
	currentConfig *Config
	configMu      sync.RWMutex
)

// SetConfig 设置当前配置
func SetConfig(cfg *Config) {
	configMu.Lock()
	defer configMu.Unlock()
	currentConfig = cfg
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return currentConfig
}

// GetAdContent 获取指定广告位的内容
func GetAdContent(slotID string) template.HTML {
	configMu.RLock()
	defer configMu.RUnlock()

	if currentConfig == nil || !currentConfig.Enabled {
		return ""
	}

	slot, ok := currentConfig.Slots[slotID]
	if !ok || !slot.Enabled {
		return ""
	}

	return template.HTML(slot.Content)
}

// GetAllSlots 获取所有广告位
func GetAllSlots() map[string]*AdSlot {
	configMu.RLock()
	defer configMu.RUnlock()

	if currentConfig == nil {
		return make(map[string]*AdSlot)
	}
	return currentConfig.Slots
}
