package service

import (
	"bookweb/config"
	"bookweb/utils"
	"log"
	"os"
	"time"
)

// ConfigWatcher 监听配置文件变更并自动重载
// 使用回调函数来避免 service -> router -> service 的循环依赖
func ConfigWatcher(onRouterReload func()) {
	// 监控的文件列表及其最后的修改时间
	files := map[string]time.Time{
		"config/config.conf": {},
		"config/link.conf":   {},
		"config/router.conf": {},
	}

	// 首次运行，记录当前时间
	for path := range files {
		if info, err := os.Stat(path); err == nil {
			files[path] = info.ModTime()
		}
	}

	log.Println("Config watcher service started.")

	for {
		time.Sleep(2 * time.Second) // 每 2 秒检查一次

		changed := false
		routerChanged := false
		for path, lastMod := range files {
			info, err := os.Stat(path)
			if err != nil {
				continue
			}

			if info.ModTime().After(lastMod) {
				log.Printf("Config file detected change: %s", path)
				files[path] = info.ModTime()
				changed = true
				if path == "config/router.conf" {
					routerChanged = true
				}
			}
		}

		if changed {
			reloadConfigs(routerChanged, onRouterReload)
		}
	}
}

// reloadConfigs 执行配置重载
func reloadConfigs(routerChanged bool, onRouterReload func()) {
	log.Println("Service: Reloading configurations...")

	// 1. 重载主体配置
	if newCfg, err := config.LoadAppConfig("config/config.conf"); err != nil {
		log.Printf("Error reloading config.conf: %v", err)
	} else {
		// 数据库热重载
		// 注意：这里简单实现，直接重连。高并发下建议加锁或使用连接池管理。
		utils.InitDB(&newCfg.Db)
		// 如果主配置变化，可能涉及 AdminPath 变化，也需要重载路由
		onRouterReload()
	}

	// 2. 重载友情链接
	if err := config.LoadLinkConfig("config/link.conf"); err != nil {
		log.Printf("Error reloading link.conf: %v", err)
	}

	// 3. 如果路由配置变了，触发外部传入的回调
	if routerChanged && onRouterReload != nil {
		onRouterReload()
	}

	log.Println("Service: All configurations reloaded.")
}
