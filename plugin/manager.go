package plugin

import (
	"bookweb/utils"
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

// Manager 插件管理器
type Manager struct {
	plugins  map[string]Plugin
	configs  map[string]map[string]interface{}
	mu       sync.RWMutex
	initDone bool
}

var (
	globalManager *Manager
	once          sync.Once
)

// GetManager 获取全局插件管理器
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			plugins: make(map[string]Plugin),
			configs: make(map[string]map[string]interface{}),
		}
	})
	return globalManager
}

// Register 注册插件
func (m *Manager) Register(p Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := p.Name()
	if _, exists := m.plugins[name]; exists {
		utils.LogWarn("Plugin", "Plugin %s already registered, skipping", name)
		return
	}
	m.plugins[name] = p
	utils.LogInfo("Plugin", "Plugin registered: %s", name)
}

// LoadConfig 加载插件配置文件
func (m *Manager) LoadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&m.configs); err != nil {
		return err
	}
	return nil
}

// InitAll 初始化所有已注册的插件
func (m *Manager) InitAll(configPath string) error {
	if err := m.LoadConfig(configPath); err != nil {
		utils.LogWarn("Plugin", "Failed to load plugin config: %v", err)
		// 继续执行，使用空配置
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for name, p := range m.plugins {
		cfg := m.configs[name]
		if cfg == nil {
			cfg = make(map[string]interface{})
		}

		// 检查是否启用
		if enabled, ok := cfg["enabled"].(bool); ok && !enabled {
			utils.LogInfo("Plugin", "Plugin %s is disabled, skipping init", name)
			continue
		}

		if err := p.Init(cfg); err != nil {
			utils.LogError("Plugin", "Failed to init plugin %s: %v", name, err)
			continue
		}
		utils.LogInfo("Plugin", "Plugin initialized: %s", name)
	}

	m.initDone = true
	return nil
}

// GetAllRoutes 获取所有插件的路由
func (m *Manager) GetAllRoutes() map[string]http.HandlerFunc {
	m.mu.RLock()
	defer m.mu.RUnlock()

	routes := make(map[string]http.HandlerFunc)
	for name, p := range m.plugins {
		// 检查是否启用
		cfg := m.configs[name]
		if cfg != nil {
			if enabled, ok := cfg["enabled"].(bool); ok && !enabled {
				continue
			}
		}

		pluginRoutes := p.GetRoutes()
		for pattern, handler := range pluginRoutes {
			routes[pattern] = handler
		}
	}
	return routes
}

// GetPlugin 获取指定插件
func (m *Manager) GetPlugin(name string) Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.plugins[name]
}

// GetConfig 获取指定插件的配置
func (m *Manager) GetConfig(name string) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.configs[name]
}

// ShutdownAll 关闭所有插件
func (m *Manager) ShutdownAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, p := range m.plugins {
		if err := p.Shutdown(); err != nil {
			utils.LogError("Plugin", "Error shutting down plugin %s: %v", name, err)
		}
	}
}

// IsEnabled 检查插件是否启用
func (m *Manager) IsEnabled(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cfg := m.configs[name]
	if cfg == nil {
		return false
	}
	if enabled, ok := cfg["enabled"].(bool); ok {
		return enabled
	}
	return false
}

// GetAllConfigs 获取所有插件配置
func (m *Manager) GetAllConfigs() map[string]map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]map[string]interface{})
	for name := range m.plugins {
		if cfg, ok := m.configs[name]; ok {
			result[name] = cfg
		} else {
			result[name] = map[string]interface{}{"enabled": false}
		}
	}
	return result
}

// ReloadConfig 重新加载并应用插件配置（热更新）
func (m *Manager) ReloadConfig(configPath string) error {
	if err := m.LoadConfig(configPath); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 重新初始化所有插件
	for name, p := range m.plugins {
		cfg := m.configs[name]
		if cfg == nil {
			cfg = make(map[string]interface{})
		}

		// 重新初始化插件（插件内部应处理配置更新）
		if err := p.Init(cfg); err != nil {
			utils.LogError("Plugin", "Failed to reload plugin %s: %v", name, err)
			continue
		}
		utils.LogInfo("Plugin", "Plugin config reloaded: %s", name)
	}

	return nil
}

// UpdatePluginConfig 更新单个插件配置并重新初始化
func (m *Manager) UpdatePluginConfig(name string, cfg map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.configs[name] = cfg

	// 重新初始化该插件
	if p, ok := m.plugins[name]; ok {
		if err := p.Init(cfg); err != nil {
			return err
		}
		utils.LogInfo("Plugin", "Plugin config updated and reloaded: %s", name)
	}

	return nil
}
