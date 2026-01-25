package router

import (
	"bookweb/admin"
	"bookweb/config"
	"bookweb/controller"
	"bookweb/model"
	"bookweb/plugin"
	"bookweb/utils"
	"context"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
)

// RouterManager 动态路由管理器，支持热载
type RouterManager struct {
	router atomic.Value // 存储 *httprouter.Router
}

var currentManager *RouterManager

// NewRouterManager 创建并初始化路由管理器
func NewRouterManager(cfg *config.RouterConfig) *RouterManager {
	m := &RouterManager{}
	m.router.Store(SetupRouter(cfg))
	currentManager = m
	return m
}

// Reload 重新加载路由配置并替换 Router
func (m *RouterManager) Reload(cfg *config.RouterConfig) {
	m.router.Store(SetupRouter(cfg))
}

// ServeHTTP 实现 http.Handler 接口
func (m *RouterManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 注入开始时间
	ctx := context.WithValue(r.Context(), model.StartTimeKey, time.Now())
	m.router.Load().(*httprouter.Router).ServeHTTP(w, r.WithContext(ctx))
}

// GetManager 获取全局路由管理器
func GetManager() *RouterManager {
	return currentManager
}

// complexRoute 内部记录复杂路由的正则和处理函数
type complexRoute struct {
	re         *regexp.Regexp
	paramNames []string
	handler    http.HandlerFunc
	methods    []string
}

var complexRoutes []complexRoute

// SetupRouter 设置路由并处理复杂路径模式
func SetupRouter(cfg *config.RouterConfig) *httprouter.Router {
	router := httprouter.New()
	complexRoutes = nil // 清空历史

	// 静态文件 - 使用绝对路径确保文件可以正确加载
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	// 模板静态文件 - 根据当前模板动态提供静态资源
	// 路径: /tpl_static/*filepath -> template/{current_template}/static/*filepath
	router.GET("/tpl_static/*filepath", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filePath := ps.ByName("filepath")
		tpl := config.GlobalConfig.Site.Template
		mobileTpl := config.GlobalConfig.Site.MobileTemplate

		// 检测是否为移动端访问
		isMobile := controller.IsMobile(r)
		templateName := tpl
		if isMobile && mobileTpl != "" {
			templateName = mobileTpl
		}

		// 构建模板静态文件完整路径
		fullPath := "template/" + templateName + "/static" + filePath

		// 检查文件是否存在，不存在则尝试回退到 PC 模板
		if _, err := http.Dir(".").Open(fullPath); err != nil && isMobile && tpl != "" {
			fullPath = "template/" + tpl + "/static" + filePath
		}

		http.ServeFile(w, r, fullPath)
	})

	// 后台管理路由
	adminPath := config.GlobalConfig.Site.AdminPath
	if adminPath == "" {
		adminPath = "/admin"
	}

	router.GET(adminPath+"/login", adaptHandlerFunc(admin.Login))
	router.POST(adminPath+"/login", adaptHandlerFunc(admin.Login))
	router.GET(adminPath+"/logout", adaptHandlerFunc(admin.Logout))
	router.GET(adminPath, adaptHandlerFunc(admin.AuthMiddleware(admin.Dashboard)))
	router.GET(adminPath+"/settings", adaptHandlerFunc(admin.AuthMiddleware(admin.Settings)))
	router.POST(adminPath+"/settings", adaptHandlerFunc(admin.AuthMiddleware(admin.Settings)))
	router.GET(adminPath+"/modules", adaptHandlerFunc(admin.AuthMiddleware(admin.Modules)))
	router.GET(adminPath+"/articles", adaptHandlerFunc(admin.AuthMiddleware(admin.Articles)))
	router.GET(adminPath+"/article/edit", adaptHandlerFunc(admin.AuthMiddleware(admin.ArticleEdit)))
	router.POST(adminPath+"/article/edit", adaptHandlerFunc(admin.AuthMiddleware(admin.ArticleEdit)))
	router.POST(adminPath+"/article/delete", adaptHandlerFunc(admin.AuthMiddleware(admin.ArticleDelete)))
	router.GET(adminPath+"/users", adaptHandlerFunc(admin.AuthMiddleware(admin.Users)))
	router.POST(adminPath+"/user/delete", adaptHandlerFunc(admin.AuthMiddleware(admin.UserDelete)))
	router.GET(adminPath+"/user/edit", adaptHandlerFunc(admin.AuthMiddleware(admin.UserEdit)))
	router.POST(adminPath+"/user/edit", adaptHandlerFunc(admin.AuthMiddleware(admin.UserEdit)))
	router.GET(adminPath+"/user/books", adaptHandlerFunc(admin.AuthMiddleware(admin.UserBooks)))
	router.POST(adminPath+"/user/bookcase/delete", adaptHandlerFunc(admin.AuthMiddleware(admin.UserBookcaseDelete)))
	router.POST(adminPath+"/user/bookmark/delete", adaptHandlerFunc(admin.AuthMiddleware(admin.UserBookmarkDelete)))
	router.GET(adminPath+"/links", adaptHandlerFunc(admin.AuthMiddleware(admin.Links)))
	router.POST(adminPath+"/link/add", adaptHandlerFunc(admin.AuthMiddleware(admin.LinkAdd)))
	router.POST(adminPath+"/link/delete", adaptHandlerFunc(admin.AuthMiddleware(admin.LinkDelete)))
	router.GET(adminPath+"/link/edit", adaptHandlerFunc(admin.AuthMiddleware(admin.LinkEdit)))
	router.POST(adminPath+"/link/edit", adaptHandlerFunc(admin.AuthMiddleware(admin.LinkEdit)))
	router.GET(adminPath+"/analytics", adaptHandlerFunc(admin.AuthMiddleware(admin.Analytics)))
	router.POST(adminPath+"/analytics", adaptHandlerFunc(admin.AuthMiddleware(admin.Analytics)))
	router.POST(adminPath+"/db/test", adaptHandlerFunc(admin.AuthMiddleware(admin.TestDBConnection)))
	router.GET(adminPath+"/security", adaptHandlerFunc(admin.AuthMiddleware(admin.Security))) // 新增安全设置
	router.POST(adminPath+"/security/password", adaptHandlerFunc(admin.AuthMiddleware(admin.SecurityPassword)))
	router.POST(adminPath+"/security/path", adaptHandlerFunc(admin.AuthMiddleware(admin.SecurityPath)))
	router.POST(adminPath+"/redis/test", adaptHandlerFunc(admin.AuthMiddleware(admin.TestRedisConnection)))
	router.POST(adminPath+"/cache/clear", adaptHandlerFunc(admin.AuthMiddleware(admin.ClearCache)))
	router.POST(adminPath+"/template/clear", adaptHandlerFunc(admin.AuthMiddleware(admin.ClearTemplates)))

	// 模块设置更新接口
	router.POST(adminPath+"/modules/routes", adaptHandlerFunc(admin.AuthMiddleware(admin.ModuleRoutesUpdate)))
	router.POST(adminPath+"/modules/sorts", adaptHandlerFunc(admin.AuthMiddleware(admin.ModuleSortsUpdate)))
	router.POST(adminPath+"/modules/seo", adaptHandlerFunc(admin.AuthMiddleware(admin.ModuleSeoUpdate)))

	// 插件管理路由
	router.GET(adminPath+"/plugins", adaptHandlerFunc(admin.AuthMiddleware(admin.Plugins)))
	router.POST(adminPath+"/plugins/toggle", adaptHandlerFunc(admin.AuthMiddleware(admin.PluginToggle)))
	router.POST(adminPath+"/plugins/config", adaptHandlerFunc(admin.AuthMiddleware(admin.PluginConfigUpdate)))

	// 遍历配置并分类注册
	for name, pattern := range cfg.Routes {
		handler := getHandler(name)
		if handler == nil {
			continue
		}

		methods := []string{"GET"}
		if name == "login" || name == "register" || name == "user_update" ||
			name == "bookcase_add" || name == "bookcase_delete" ||
			name == "bookmark_add" || name == "bookmark_delete" {
			methods = append(methods, "POST")
		}

		if isComplexPattern(pattern) {
			// 复杂路由：加入正则匹配列表
			addComplexRoute(pattern, handler, methods)
		} else {
			// 简单路由：直接注册到 httprouter
			for _, method := range methods {
				router.Handle(method, pattern, adaptHandler(handler))
			}
		}
	}

	// 注册插件路由
	pluginRoutes := plugin.GetManager().GetAllRoutes()
	pluginMethods := []string{"GET", "POST"}
	for pattern, handler := range pluginRoutes {
		if isComplexPattern(pattern) {
			addComplexRoute(pattern, handler, pluginMethods)
		} else {
			for _, method := range pluginMethods {
				router.Handle(method, pattern, adaptHandler(handler))
			}
		}
		utils.LogInfo("Router", "Plugin route registered: %s", pattern)
	}

	// 设置 NotFound 拦截器来处理复杂正则路由
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 遍历所有复杂路由进行匹配
		for _, cr := range complexRoutes {
			// 检查方法是否允许 (如果允许 GET，也应允许 HEAD)
			methodAllowed := false
			for _, m := range cr.methods {
				if r.Method == m || (r.Method == "HEAD" && m == "GET") {
					methodAllowed = true
					break
				}
			}
			if !methodAllowed {
				continue
			}

			// 执行正则匹配
			matches := cr.re.FindStringSubmatch(r.URL.Path)
			if matches != nil {
				// 提取参数
				params := make(model.Params, 0)
				for i, val := range matches[1:] {
					params = append(params, httprouter.Param{
						Key:   cr.paramNames[i],
						Value: val,
					})
				}

				// 注入 Context 并执行
				ctx := context.WithValue(r.Context(), model.ParamsKey, params)
				cr.handler(w, r.WithContext(ctx))
				return
			}
		}

		// 如果正则也没匹配到，返回 404
		controller.NotFound(w, r)
	})

	return router
}

// isComplexPattern 检查路径是否包含 httprouter 不直接支持的多参数或特殊格式
func isComplexPattern(pattern string) bool {
	segments := strings.Split(pattern, "/")
	for _, seg := range segments {
		// httprouter 一个段只允许一个 ':' 开始的参数
		if strings.Count(seg, ":") > 1 {
			return true
		}
		// 如果 ':' 不是该段的开头（例如 /book_:aid）或其后有非斜杠字符（例如 /:aid.html）
		if strings.Contains(seg, ":") {
			if !strings.HasPrefix(seg, ":") || containsSpecialChars(seg) {
				return true
			}
		}
	}
	return false
}

// containsSpecialChars 检查参数段是否包含非字母数字的特殊字符（这会导致 httprouter 匹配宽泛）
func containsSpecialChars(seg string) bool {
	// 去掉开头的 ':'
	paramPart := strings.TrimPrefix(seg, ":")
	// 如果剩余部分包含任何非标识符字符，判定为复杂
	return strings.ContainsAny(paramPart, "._-")
}

// addComplexRoute 编译复杂路由为正则并存储
func addComplexRoute(pattern string, handler http.HandlerFunc, methods []string) {
	re, paramNames := patternToRegexp(pattern)
	complexRoutes = append(complexRoutes, complexRoute{
		re:         re,
		paramNames: paramNames,
		handler:    handler,
		methods:    methods,
	})
}

// patternToRegexp 将自定义路径模式转换为正则表达式和参数列表
func patternToRegexp(pattern string) (*regexp.Regexp, []string) {
	var paramNames []string

	// 1. 查找命名的参数 :name 并提取名称
	paramRegex := regexp.MustCompile(`:([a-zA-Z0-9]+)`)

	// 2. 将参数替换为临时占位符，避免在 QuoteMeta 时被转义
	tmpPattern := paramRegex.ReplaceAllStringFunc(pattern, func(m string) string {
		name := m[1:] // 去掉 :
		paramNames = append(paramNames, name)
		return "___PARAM_PLACEHOLDER___"
	})

	// 3. 对整体进行正则转义（处理 . _ - 等字符）
	rePattern := regexp.QuoteMeta(tmpPattern)

	// 4. 将占位符替换为捕获组
	finalPattern := strings.ReplaceAll(rePattern, "___PARAM_PLACEHOLDER___", "([^/_.]+)")

	return regexp.MustCompile("^" + finalPattern + "$"), paramNames
}

// getHandler 根据路由名称返回对应的 Handler
func getHandler(name string) http.HandlerFunc {
	handlers := map[string]http.HandlerFunc{
		"index":           controller.Index,
		"login":           controller.Login,
		"register":        controller.Regist,
		"user_center":     controller.UserCenter,
		"user_update":     controller.UpdateUser,
		"logout":          controller.Logout,
		"book":            controller.BookInfo,
		"book_index":      controller.BookIndex,
		"book_index_page": controller.BookIndex,
		"read":            controller.ChapterRead,
		"sort":            controller.SortList,
		"top":             controller.Top,
		"image":           controller.CoverImage,
		"search":          controller.Search,
		"bookcase_add":    controller.AddBookcase,
		"bookcase_delete": controller.DeleteBookcase,
		"bookmark_add":    controller.AddBookmark,
		"bookmark_delete": controller.DeleteBookmark,
	}
	return handlers[name]
}

// adaptHandler 标准适配逻辑
func adaptHandler(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !checkDBConnection(w) {
			return
		}
		if !checkDomain(w, r) {
			return
		}
		ctx := context.WithValue(r.Context(), model.ParamsKey, ps)
		h(w, r.WithContext(ctx))
	}
}

// adaptHandlerFunc 简化版适配（无参数注入）
func adaptHandlerFunc(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if !checkDBConnection(w) {
			return
		}
		if !checkDomain(w, r) {
			return
		}
		h(w, r)
	}
}

// checkDBConnection 检查数据库连接，失败则渲染错误页面
func checkDBConnection(w http.ResponseWriter) bool {
	if err := utils.Db.Ping(); err != nil {
		utils.LogError("Router", "Database connection failed: %v", err)
		t, parseErr := template.ParseFiles(controller.TplPath("db_error.html"))
		if parseErr != nil {
			// 如果模版加载失败，返回简单文本
			http.Error(w, "目前无法访问数据库，请稍后刷新重试。", http.StatusServiceUnavailable)
			return false
		}
		t.Execute(w, nil)
		return false
	}
	return true
}

// checkDomain 检查访问域名是否合法
func checkDomain(w http.ResponseWriter, r *http.Request) bool {
	cfg := config.GetGlobalConfig()
	if cfg != nil && cfg.Site.ForceDomain && cfg.Site.Domain != "" {
		// 如果是后台管理路径，不做域名限制
		adminPath := cfg.Site.AdminPath
		if adminPath == "" {
			adminPath = "/admin"
		}
		if strings.HasPrefix(r.URL.Path, adminPath) {
			return true
		}

		// 检查是否匹配 PC 域名或移动端域名
		host := r.Host
		if host != cfg.Site.Domain && host != cfg.Site.MobileDomain {
			http.Error(w, "Forbidden: Domain not allowed (Invalid Host)", http.StatusForbidden)
			return false
		}
	}
	return true
}

// GetParams 获取路由参数（供外部调用）
func GetParams(r *http.Request) httprouter.Params {
	if ps := r.Context().Value(model.ParamsKey); ps != nil {
		return ps.(httprouter.Params)
	}
	return nil
}
