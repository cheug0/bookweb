package main

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/plugin"
	"bookweb/plugin/ads"
	"bookweb/plugin/db_optimizer"
	"bookweb/plugin/langtail"
	"bookweb/plugin/sitemap"
	"bookweb/router"
	"bookweb/service"
	"bookweb/utils"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// 加载应用配置
	appCfg, err := config.LoadAppConfig("config/config.conf")
	if err != nil {
		log.Fatalf("Failed to load app config: %v", err)
	}

	// 初始化日志 (使用配置中的新参数)
	if err := utils.InitLogger(
		appCfg.Log.Level,
		appCfg.Log.Output,
		appCfg.Log.FilePath,
		appCfg.Log.EnableHTTP,
		appCfg.Log.MaxSize,
		appCfg.Log.MaxAge,
	); err != nil {
		log.Printf("Warning: Failed to init logger: %v", err)
	}
	utils.LogInfo("System", "Logger initialized with level=%s, output=%s", appCfg.Log.Level, appCfg.Log.Output)

	// 初始化 ID 转换规则
	if err := utils.ParseIdTransRule(appCfg.Site.IdTransRule); err != nil {
		utils.LogWarn("System", "Failed to parse ID trans rule: %v", err)
	}

	// 初始化数据库
	utils.InitDB(&appCfg.Db)

	// 初始化预编译SQL语句
	if err := dao.InitPreparedStatements(); err != nil {
		utils.LogWarn("DAO", "Failed to init prepared statements: %v", err)
	}

	// 初始化 Redis 缓存 (如果启用)
	if appCfg.Redis.Enabled {
		if err := utils.InitRedis(&appCfg.Redis); err != nil {
			utils.LogWarn("Redis", "Failed to init Redis cache: %v", err)
		} else {
			utils.LogInfo("Redis", "Redis cache enabled and connected.")
		}
	}

	// 加载友情链接配置
	if err := config.LoadLinkConfig("config/link.conf"); err != nil {
		utils.LogWarn("Config", "Failed to load link config: %v", err)
	}

	// 加载 SEO 规则配置
	if err := config.LoadSeoConfig("config/seo.conf"); err != nil {
		utils.LogWarn("Config", "Failed to load SEO config: %v", err)
	}

	// 加载路由配置
	routerCfg, err := config.LoadRouterConfig("config/router.conf")
	if err != nil {
		log.Fatalf("Failed to load router config: %v", err)
	}

	utils.LogInfo("Router", "Router configuration loaded successfully")

	// 初始化插件系统
	pluginManager := plugin.GetManager()
	pluginManager.Register(langtail.New())
	pluginManager.Register(ads.New())
	pluginManager.Register(db_optimizer.New())
	pluginManager.Register(sitemap.New())
	if err := pluginManager.InitAll("config/plugins.conf"); err != nil {
		utils.LogWarn("Plugin", "Failed to init plugins: %v", err)
	}

	// 注入广告获取函数到模板工具中，解决循环依赖
	utils.GetAdContentFunc = ads.GetAdContent

	// 初始化模板缓存
	if err := utils.InitTemplates(); err != nil {
		utils.LogError("Template", "Failed to init templates: %v", err)
		os.Exit(1)
	}

	// 初始化动态路由管理器
	rm := router.NewRouterManager(routerCfg)

	// 启动配置监听协程 (解耦回调以打破循环依赖)
	go service.ConfigWatcher(func() {
		if newRouterCfg, err := config.LoadRouterConfig("config/router.conf"); err == nil {
			rm.Reload(newRouterCfg)
			utils.LogInfo("Router", "Router hot-swapped successfully.")
		} else {
			utils.LogWarn("Router", "Failed to hot-swap router: %v", err)
		}
	})

	// 缓存预热 - 预填充常用数据缓存
	go warmupCache()

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Server.Host, appCfg.Server.Port)
	utils.LogInfo("Server", "Server starting on %s (Hot reload enabled)...", serverAddr)
	fmt.Printf("\nServer starting on %s (Hot reload enabled)...\n", serverAddr)

	// 使用中间件包装路由: Logging -> GZIP -> Router
	handler := router.LoggingMiddleware(utils.GzipMiddleware(rm))
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}

// warmupCache 缓存预热 - 启动时预填充常用数据
func warmupCache() {
	utils.LogInfo("Cache", "Cache warmup starting...")

	// 预热首页需要的数据
	dao.GetAllSortsCached()
	dao.GetVisitArticlesCached(6)
	dao.GetVisitArticlesCached(12)
	dao.GetArticlesBySortIDCached(0, 0, 30)

	// 预热各分类数据
	sorts, _ := dao.GetAllSortsCached()
	for i, s := range sorts {
		if i < 6 {
			dao.GetArticlesBySortIDCached(s.SortID, 0, 13)
		}
	}

	utils.LogInfo("Cache", "Cache warmup completed.")
}
