package ads

import (
	"bookweb/utils"
	"net/http"
)

// Plugin 广告插件
type Plugin struct {
	config *Config
}

// New 创建广告插件实例
func New() *Plugin {
	return &Plugin{}
}

// Name 返回插件名称
func (p *Plugin) Name() string {
	return "ads"
}

// Init 初始化插件
func (p *Plugin) Init(cfg map[string]interface{}) error {
	p.config = ParseConfig(cfg)
	SetConfig(p.config)

	if p.config.Enabled {
		utils.LogInfo("Ads", "Ads plugin enabled with %d slots", len(p.config.Slots))
	}

	return nil
}

// GetRoutes 获取插件路由（广告插件无需额外路由）
func (p *Plugin) GetRoutes() map[string]http.HandlerFunc {
	return nil
}

// Shutdown 关闭插件
func (p *Plugin) Shutdown() error {
	utils.LogInfo("Ads", "Ads plugin shutdown")
	return nil
}

// GetConfig 获取插件配置
func (p *Plugin) GetConfig() *Config {
	return p.config
}
