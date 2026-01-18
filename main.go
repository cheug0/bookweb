package main

import (
	"bookweb/config"
	"bookweb/plugin"
	"bookweb/plugin/langtail"
	"bookweb/router"
	"bookweb/service"
	"bookweb/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 加载应用配置
	appCfg, err := config.LoadAppConfig("config/config.conf")
	if err != nil {
		log.Fatalf("Failed to load app config: %v", err)
	}

	// 初始化数据库
	utils.InitDB(&appCfg.Db)

	// 初始化 Redis 缓存 (如果启用)
	if appCfg.Redis.Enabled {
		if err := utils.InitRedis(&appCfg.Redis); err != nil {
			log.Printf("Warning: Failed to init Redis cache: %v", err)
		} else {
			log.Println("Redis cache enabled and connected.")
		}
	}

	// 加载友情链接配置
	if err := config.LoadLinkConfig("config/link.conf"); err != nil {
		log.Printf("Warning: Failed to load link config: %v", err)
	}

	// 加载 SEO 规则配置
	if err := config.LoadSeoConfig("config/seo.conf"); err != nil {
		log.Printf("Warning: Failed to load SEO config: %v", err)
	}

	// 加载路由配置
	routerCfg, err := config.LoadRouterConfig("config/router.conf")
	if err != nil {
		log.Fatalf("Failed to load router config: %v", err)
	}

	fmt.Println("Router configuration loaded successfully")

	// 初始化插件系统
	pluginManager := plugin.GetManager()
	pluginManager.Register(langtail.New())
	if err := pluginManager.InitAll("config/plugins.conf"); err != nil {
		log.Printf("Warning: Failed to init plugins: %v", err)
	}

	// 初始化动态路由管理器
	rm := router.NewRouterManager(routerCfg)

	// 启动配置监听协程 (解耦回调以打破循环依赖)
	go service.ConfigWatcher(func() {
		if newRouterCfg, err := config.LoadRouterConfig("config/router.conf"); err == nil {
			rm.Reload(newRouterCfg)
			log.Println("Router hot-swapped successfully.")
		} else {
			log.Printf("Failed to hot-swap router: %v", err)
		}
	})

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Server.Host, appCfg.Server.Port)
	fmt.Printf("\nServer starting on %s (Hot reload enabled)...\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, rm))
}
