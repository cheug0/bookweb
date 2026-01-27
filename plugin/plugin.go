// plugin.go
// 插件接口定义
// 定义所有插件必须实现的接口规范
package plugin

import "net/http"

// Plugin 插件接口定义
type Plugin interface {
	// Name 返回插件名称（唯一标识符）
	Name() string

	// Init 初始化插件，传入插件配置
	Init(cfg map[string]interface{}) error

	// GetRoutes 获取插件需要注册的路由
	// 返回格式: map[路由模式]handler
	GetRoutes() map[string]http.HandlerFunc

	// Shutdown 插件关闭时的清理工作
	Shutdown() error
}

// RouteInfo 路由信息结构
type RouteInfo struct {
	Pattern string
	Handler http.HandlerFunc
	Methods []string // GET, POST, etc.
}
