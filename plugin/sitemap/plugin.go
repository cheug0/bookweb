package sitemap

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Plugin sitemap 插件
type Plugin struct {
	config   *Config
	stopChan chan struct{}
}

// New 创建 sitemap 插件实例
func New() *Plugin {
	return &Plugin{
		stopChan: make(chan struct{}),
	}
}

// Name 返回插件名称
func (p *Plugin) Name() string {
	return "sitemap"
}

// Init 初始化插件
func (p *Plugin) Init(cfg map[string]interface{}) error {
	p.config = ParseConfig(cfg)
	SetConfig(p.config)

	if !p.config.Enabled {
		return nil
	}

	// 启动时生成 sitemap
	go func() {
		if err := GenerateSitemap(); err != nil {
			log.Printf("Sitemap: 生成失败: %v", err)
		} else {
			log.Printf("Sitemap: 初始生成完成")
		}
	}()

	// 启动定时更新
	if p.config.RegenerateHours > 0 {
		go p.startScheduler()
	}

	log.Printf("Sitemap plugin enabled: output=%s, regenerate=%dh",
		p.config.OutputPath, p.config.RegenerateHours)

	return nil
}

// startScheduler 启动定时更新任务
func (p *Plugin) startScheduler() {
	ticker := time.NewTicker(time.Duration(p.config.RegenerateHours) * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := GenerateSitemap(); err != nil {
				log.Printf("Sitemap: 定时更新失败: %v", err)
			} else {
				log.Printf("Sitemap: 定时更新完成")
			}
		case <-p.stopChan:
			return
		}
	}
}

// GetRoutes 获取插件路由
func (p *Plugin) GetRoutes() map[string]http.HandlerFunc {
	if p.config == nil || !p.config.Enabled {
		return nil
	}

	return map[string]http.HandlerFunc{
		"/sitemap/*filepath": p.serveSitemap,
	}
}

// serveSitemap 提供 sitemap 文件服务
func (p *Plugin) serveSitemap(w http.ResponseWriter, r *http.Request) {
	// 从 URL 中提取文件名
	path := r.URL.Path
	filename := strings.TrimPrefix(path, "/sitemap/")
	if filename == "" {
		filename = "sitemap.xml"
	}

	// 安全检查：防止路径遍历
	filename = filepath.Base(filename)
	if !strings.HasSuffix(filename, ".xml") {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(p.config.OutputPath, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 尝试生成
		if err := GenerateSitemap(); err != nil {
			http.Error(w, "Sitemap generation failed", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	http.ServeFile(w, r, filePath)
}

// Shutdown 关闭插件
func (p *Plugin) Shutdown() error {
	close(p.stopChan)
	log.Println("Sitemap plugin shutdown")
	return nil
}
