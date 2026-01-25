package langtail

import (
	"bookweb/service"
	"bookweb/utils"
	"log"
	"net/http"
)

// Plugin 长尾词插件实现
type Plugin struct {
	config *Config
}

// New 创建长尾词插件实例
func New() *Plugin {
	return &Plugin{}
}

// Name 返回插件名称
func (p *Plugin) Name() string {
	return "langtail"
}

// Init 初始化插件
func (p *Plugin) Init(cfg map[string]interface{}) error {
	p.config = ParseConfig(cfg)
	SetConfig(p.config)

	// 注册长尾词更新回调函数
	service.LangtailUpdateFunc = func(sourceID int, sourceName string, cycleDays int) {
		_ = UpdateLangtailsIfNeeded(sourceID, sourceName, cycleDays)
	}

	if p.config.Enabled {
		utils.LogInfo("Langtail", "Langtail plugin enabled: cycle=%d days, cache=%d seconds, route=%s",
			p.config.FetchCycleDays, p.config.CacheSeconds, p.config.RoutePattern)
	}

	return nil
}

// GetRoutes 获取插件路由
func (p *Plugin) GetRoutes() map[string]http.HandlerFunc {
	if p.config == nil || !p.config.Enabled {
		return nil
	}

	routes := make(map[string]http.HandlerFunc)
	routes[p.config.RoutePattern] = LangtailInfo
	return routes
}

// Shutdown 关闭插件
func (p *Plugin) Shutdown() error {
	log.Println("Langtail plugin shutdown")
	return nil
}

// GetConfig 获取插件配置
func (p *Plugin) GetConfig() *Config {
	return p.config
}
