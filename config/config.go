package config

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	configLock   sync.RWMutex
	Config       *RouterConfig
	GlobalConfig *AppConfig
)

// GetGlobalConfig 安全获取全局配置
func GetGlobalConfig() *AppConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return GlobalConfig
}

// GetRouterConfig 安全获取路由配置
func GetRouterConfig() *RouterConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return Config
}

// RouterConfig 路由配置结构
type RouterConfig struct {
	Routes map[string]string `json:"routes"`
}

// DbConfig 数据库配置结构
type DbConfig struct {
	Driver          string `json:"driver"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	User            string `json:"user"`
	Password        string `json:"password"`
	DbName          string `json:"dbname"`
	MaxOpenConns    int    `json:"max_open_conns"`    // 最大打开连接数
	MaxIdleConns    int    `json:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime int    `json:"conn_max_lifetime"` // 连接最大生命周期（秒）
}

// AppConfig 应用全局配置结构
type AppConfig struct {
	Db        DbConfig           `json:"db"`
	Server    ServerConfig       `json:"server"`
	Site      SiteConfig         `json:"site"`
	Storage   StorageConfig      `json:"storage"`
	SeoRules  map[string]SeoRule `json:"-"`
	Links     []LinkConfig       `json:"-"`
	Analytics string             `json:"analytics"`
	Redis     RedisConfig        `json:"redis"`
}

// RedisConfig Redis缓存配置
type RedisConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type  string      `json:"type"` // local, oss
	Local LocalConfig `json:"local"`
	Oss   OssConfig   `json:"oss"`
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	Path string `json:"path"`
}

// OssConfig 对象存储配置
type OssConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
}

// ServerConfig 服务器配置结构
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// LinkConfig 友情链接配置结构
type LinkConfig struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Order int    `json:"order"`
}

// SiteConfig 站点展示与 SEO 配置结构
type SiteConfig struct {
	SiteName       string `json:"sitename"`
	Domain         string `json:"domain"`
	MobileDomain   string `json:"mobile_domain"` // 移动端域名
	Template       string `json:"template"`
	MobileTemplate string `json:"mobile_template"`  // 移动端模板
	AdminPath      string `json:"admin_path"`       // 后台管理路径，默认 /admin
	SearchLimit    int    `json:"search_limit"`     // 搜索限制时间（秒）
	IndexCache     bool   `json:"index_cache"`      // 是否开启首页缓存
	BookCache      bool   `json:"book_cache"`       // 开启小说信息页缓存
	BookIndexCache bool   `json:"book_index_cache"` // 开启小说目录页缓存
	ReadCache      bool   `json:"read_cache"`       // 开启章节阅读页缓存
	SortCache      bool   `json:"sort_cache"`       // 开启分类页缓存
	TopCache       bool   `json:"top_cache"`        // 开启排行榜缓存
	ForceDomain    bool   `json:"force_domain"`     // 是否强制域名访问
	IdTransRule    string `json:"id_trans_rule"`    // 小说ID转换规则 (e.g. "*2,+100")
	GzipEnabled    bool   `json:"gzip_enabled"`     // 开启 GZIP 压缩
}

// SeoRule 定义单个页面的 SEO 模板
type SeoRule struct {
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
}

// LoadRouterConfig 加载路由配置
func LoadRouterConfig(configPath string) (*RouterConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg RouterConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	configLock.Lock()
	Config = &cfg
	configLock.Unlock()
	return &cfg, nil
}

// LoadAppConfig 加载应用全局配置
func LoadAppConfig(configPath string) (*AppConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	configLock.Lock()
	if cfg.Site.AdminPath == "" {
		cfg.Site.AdminPath = "/admin"
	}
	GlobalConfig = &cfg
	configLock.Unlock()
	return &cfg, nil
}

// LoadLinkConfig 加载友情链接配置
func LoadLinkConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var links []LinkConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&links); err != nil {
		return err
	}

	configLock.Lock()
	if GlobalConfig != nil {
		GlobalConfig.Links = links
	}
	configLock.Unlock()
	return nil
}

// LoadSeoConfig 加载 SEO 规则配置
func LoadSeoConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var seoRules map[string]SeoRule
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&seoRules); err != nil {
		return err
	}

	configLock.Lock()
	if GlobalConfig != nil {
		GlobalConfig.SeoRules = seoRules
	}
	configLock.Unlock()
	return nil
}

// SaveSeoConfig 保存 SEO 规则配置到文件
func SaveSeoConfig(configPath string) error {
	configLock.RLock()
	defer configLock.RUnlock()

	if GlobalConfig == nil {
		return nil
	}

	data, err := json.MarshalIndent(GlobalConfig.SeoRules, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// GetRoute 获取路由路径
func (c *RouterConfig) GetRoute(name string) string {
	if route, ok := c.Routes[name]; ok {
		return route
	}
	return ""
}

// SaveAppConfig 保存应用配置到文件
func SaveAppConfig(configPath string) error {
	configLock.RLock()
	defer configLock.RUnlock()

	if GlobalConfig == nil {
		return nil
	}

	data, err := json.MarshalIndent(GlobalConfig, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// SaveLinkConfig 保存友情链接配置到文件
func SaveLinkConfig(configPath string) error {
	configLock.RLock()
	defer configLock.RUnlock()

	if GlobalConfig == nil {
		return nil
	}

	data, err := json.MarshalIndent(GlobalConfig.Links, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// SaveRouterConfig 保存路由配置到文件
func SaveRouterConfig(configPath string) error {
	configLock.RLock()
	defer configLock.RUnlock()

	if Config == nil {
		return nil
	}

	data, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// GetPluginConfig 获取指定插件的配置
func GetPluginConfig(pluginName string) map[string]interface{} {
	file, err := os.Open("config/plugins.conf")
	if err != nil {
		return nil
	}
	defer file.Close()

	var configs map[string]map[string]interface{}
	if err := json.NewDecoder(file).Decode(&configs); err != nil {
		return nil
	}
	return configs[pluginName]
}
