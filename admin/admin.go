package admin

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/plugin"
	"bookweb/utils"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
)

// 模板路径
func tplPath(name string) string {
	return "admin/template/" + name
}

// adminFuncMap 后台模板函数
var adminFuncMap = template.FuncMap{
	"minus": func(a, b int) int { return a - b },
	"plus":  func(a, b int) int { return a + b },
}

// parseTpl 解析后台模板
func parseTpl(names ...string) (*template.Template, error) {
	paths := make([]string, len(names))
	for i, name := range names {
		paths[i] = tplPath(name)
	}
	return template.New(names[0]).Funcs(adminFuncMap).ParseFiles(paths...)
}

// getAdminData 获取后台通用模板数据
func getAdminData(r *http.Request, active string, title string) map[string]interface{} {
	session, _ := IsAdminLoggedIn(r)
	cfg := config.GetGlobalConfig()
	adminPath := cfg.Site.AdminPath
	if adminPath == "" {
		adminPath = "/admin"
	}
	return map[string]interface{}{
		"Title":     title,
		"Username":  session.Username,
		"Active":    active,
		"AdminPath": adminPath,
	}
}

// jsonResponse 返回 JSON 响应
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Login 后台登录页面
func Login(w http.ResponseWriter, r *http.Request) {
	adminPath := config.GlobalConfig.Site.AdminPath
	if adminPath == "" {
		adminPath = "/admin"
	}

	if r.Method == "GET" {
		t, err := parseTpl("login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]interface{}{"AdminPath": adminPath})
		return
	}

	// POST 处理登录
	username := r.FormValue("username")
	password := r.FormValue("password")

	admin, err := dao.VerifyAdminPassword(username, password)
	if err != nil {
		t, _ := parseTpl("login.html")
		t.Execute(w, map[string]interface{}{"Error": err.Error(), "AdminPath": adminPath})
		return
	}

	// 创建会话
	session := CreateAdminSession(admin.Id, admin.Username)
	SetAdminSessionCookie(w, session)

	http.Redirect(w, r, adminPath, http.StatusFound)
}

// ClearCache 清理 Redis 缓存
func ClearCache(w http.ResponseWriter, r *http.Request) {
	_, ok := IsAdminLoggedIn(r)
	if !ok {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "未登录"})
		return
	}

	if err := utils.CacheFlush(); err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "清理失败: " + err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "Redis 缓存已清空"})
}

// ClearTemplates 清理模板缓存 (重新加载)
func ClearTemplates(w http.ResponseWriter, r *http.Request) {
	_, ok := IsAdminLoggedIn(r)
	if !ok {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "未登录"})
		return
	}

	if err := utils.InitTemplates(); err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "重载失败: " + err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "模板缓存已重载"})
}

// Logout 后台注销
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(AdminSessionCookieName)
	if err == nil {
		DeleteAdminSession(cookie.Value)
	}
	ClearAdminSessionCookie(w)

	adminPath := config.GlobalConfig.Site.AdminPath
	if adminPath == "" {
		adminPath = "/admin"
	}
	http.Redirect(w, r, adminPath+"/login", http.StatusFound)
}

// Dashboard 仪表板
func Dashboard(w http.ResponseWriter, r *http.Request) {
	// 获取统计数据
	stats, _ := dao.GetDashboardStats()

	t, err := parseTpl("layout.html", "dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "dashboard", "仪表板")
	data["Stats"] = stats // 追加额外数据
	t.ExecuteTemplate(w, "layout", data)
}

// Settings 系统设置页面
func Settings(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetGlobalConfig()

	if r.Method == "POST" {
		updateType := r.FormValue("update_type")

		if updateType == "basic" {
			// 保存基本设置和存储设置
			cfg.Site.SiteName = r.FormValue("sitename")
			cfg.Site.Domain = r.FormValue("domain")
			cfg.Site.Template = r.FormValue("template")
			if limit, err := strconv.Atoi(r.FormValue("search_limit")); err == nil {
				cfg.Site.SearchLimit = limit
			}
			cfg.Site.IndexCache = r.FormValue("index_cache") == "on"
			cfg.Site.BookCache = r.FormValue("book_cache") == "on"
			cfg.Site.BookIndexCache = r.FormValue("book_index_cache") == "on"
			cfg.Site.ReadCache = r.FormValue("read_cache") == "on"
			cfg.Site.SortCache = r.FormValue("sort_cache") == "on"
			cfg.Site.TopCache = r.FormValue("top_cache") == "on"
			cfg.Site.TopCache = r.FormValue("top_cache") == "on"
			cfg.Site.ForceDomain = r.FormValue("force_domain") == "on"
			cfg.Site.IdTransRule = r.FormValue("id_trans_rule")
			cfg.Site.GzipEnabled = r.FormValue("gzip_enabled") == "on"

			// 更新 ID 转换规则
			utils.ParseIdTransRule(cfg.Site.IdTransRule)

			cfg.Storage.Type = r.FormValue("storage_type")

			// 保存本地存储配置
			cfg.Storage.Local.Path = r.FormValue("storage_path")

			// 保存 OSS 配置
			cfg.Storage.Oss.Endpoint = r.FormValue("oss_endpoint")
			cfg.Storage.Oss.AccessKey = r.FormValue("oss_access_key")
			cfg.Storage.Oss.SecretKey = r.FormValue("oss_secret_key")
			cfg.Storage.Oss.Bucket = r.FormValue("oss_bucket")
			cfg.Storage.Oss.Domain = r.FormValue("oss_domain")
		} else if updateType == "db" {
			// 保存数据库配置
			cfg.Db.Driver = r.FormValue("db_driver")
			cfg.Db.Host = r.FormValue("db_host")
			cfg.Db.Port, _ = strconv.Atoi(r.FormValue("db_port"))
			cfg.Db.User = r.FormValue("db_user")
			cfg.Db.Password = r.FormValue("db_password")
			cfg.Db.DbName = r.FormValue("db_dbname")
			if maxOpen, err := strconv.Atoi(r.FormValue("db_max_open")); err == nil {
				cfg.Db.MaxOpenConns = maxOpen
			}
			if maxIdle, err := strconv.Atoi(r.FormValue("db_max_idle")); err == nil {
				cfg.Db.MaxIdleConns = maxIdle
			}
			if maxLife, err := strconv.Atoi(r.FormValue("db_max_life")); err == nil {
				cfg.Db.ConnMaxLifetime = maxLife
			}
		} else if updateType == "redis" {
			// 保存 Redis 配置
			cfg.Redis.Enabled = r.FormValue("redis_enabled") == "on"
			cfg.Redis.Host = r.FormValue("redis_host")
			cfg.Redis.Port, _ = strconv.Atoi(r.FormValue("redis_port"))
			cfg.Redis.Password = r.FormValue("redis_password")
			cfg.Redis.DB, _ = strconv.Atoi(r.FormValue("redis_db"))
		}

		err := config.SaveAppConfig("config/config.conf")
		if err != nil {
			jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
			return
		}
		jsonResponse(w, map[string]interface{}{"success": true, "message": "保存成功"})
		return
	}

	t, err := parseTpl("layout.html", "settings.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "settings", "系统设置")
	data["Config"] = cfg
	t.ExecuteTemplate(w, "layout", data)
}

// Articles 小说管理页面
func Articles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	keyword := r.URL.Query().Get("keyword")

	articles, total, err := dao.GetArticleListAdmin(page, 20, keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPage := (total + 19) / 20

	t, err := parseTpl("layout.html", "articles.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "articles", "小说管理")
	data["Articles"] = articles
	data["Page"] = page
	data["TotalPage"] = totalPage
	data["Total"] = total
	data["Keyword"] = keyword
	t.ExecuteTemplate(w, "layout", data)
}

// ArticleEdit 编辑小说
func ArticleEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("articlename")
		author := r.FormValue("author")
		sortID, _ := strconv.Atoi(r.FormValue("sortid"))
		fullFlag, _ := strconv.Atoi(r.FormValue("fullflag"))
		intro := r.FormValue("intro")

		err := dao.UpdateArticleAdmin(id, name, author, sortID, fullFlag, intro)
		if err != nil {
			jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
			return
		}
		jsonResponse(w, map[string]interface{}{"success": true, "message": "保存成功"})
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	article, err := dao.GetArticleByIDAdmin(id)
	if err != nil {
		http.Error(w, "小说不存在", http.StatusNotFound)
		return
	}
	sorts, _ := dao.GetAllSorts()

	t, err := parseTpl("layout.html", "article_edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "articles", "编辑小说")
	data["Article"] = article
	data["Sorts"] = sorts
	t.ExecuteTemplate(w, "layout", data)
}

// ArticleDelete 删除小说
func ArticleDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	err := dao.DeleteArticleAdmin(id)
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	jsonResponse(w, map[string]interface{}{"success": true, "message": "删除成功"})
}

// Users 用户管理页面
func Users(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	users, total, err := dao.GetUserListAdmin(page, 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPage := (total + 19) / 20

	t, err := parseTpl("layout.html", "users.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "users", "用户管理")
	data["Users"] = users
	data["Page"] = page
	data["TotalPage"] = totalPage
	data["Total"] = total
	t.ExecuteTemplate(w, "layout", data)
}

// UserDelete 删除用户
func UserDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	err := dao.DeleteUserAdmin(id)
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	jsonResponse(w, map[string]interface{}{"success": true, "message": "删除成功"})
}

// UserEdit 编辑用户
func UserEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		err := dao.UpdateUserAdmin(id, username, password, email)
		if err != nil {
			jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
			return
		}
		jsonResponse(w, map[string]interface{}{"success": true, "message": "保存成功"})
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	user, err := dao.GetUserByIDAdmin(id)
	if err != nil {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	t, err := parseTpl("layout.html", "user_edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "users", "编辑用户")
	data["User"] = user
	t.ExecuteTemplate(w, "layout", data)
}

// UserBooks 用户书籍管理（书架/书签）
func UserBooks(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	user, err := dao.GetUserByIDAdmin(id)
	if err != nil {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	bookcases, _ := dao.GetBookcaseList(id, 0, 100)
	bookmarks, _ := dao.GetBookmarkList(id, 0, 100)

	t, err := parseTpl("layout.html", "user_books.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "users", "用户书籍管理")
	data["User"] = user
	data["Bookcases"] = bookcases
	data["Bookmarks"] = bookmarks
	t.ExecuteTemplate(w, "layout", data)
}

// UserBookcaseDelete 删除书架记录
func UserBookcaseDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	err := dao.DeleteBookcase(id)
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	jsonResponse(w, map[string]interface{}{"success": true, "message": "删除成功"})
}

// UserBookmarkDelete 删除书签记录
func UserBookmarkDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	err := dao.DeleteBookmark(id)
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	jsonResponse(w, map[string]interface{}{"success": true, "message": "删除成功"})
}

// Links 友情链接管理
func Links(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetGlobalConfig()

	// 排序
	// 注意：这里是对 Slice 的引用进行排序，会影响全局配置的显示顺序
	// 但不会持久化到文件，直到调用 SaveLinkConfig
	// 为了不影响原始 Slice 的顺序直到保存，这里可以暂不处理，
	// 或者就是直接对 Config.Links 排序。
	// 鉴于这是一个管理操作，直接排序 GlobalConfig.Links 并在后续保存时持久化也是合理的。
	// 这里为了简单，我们每次显示前都排一次序。
	// 但更好的做法是：在 Add/Edit/Delete 后保存前排序。
	// 这里仅做显示排序，不修改 Config，避免并发问题（虽然这里并没有锁保护 slice 读写在 controller 层）
	// 不过 config 里面有锁。GetGlobalConfig 返回的是指针。
	// 安全起见，我们在 template 渲染前拷贝一份或者直接在 Save 时排序。
	// 让我们在 SaveLinkConfig 前排序。
	// 那显示的时候呢？
	// 让我们在 Links handler 里构造一个副本排序用于显示。
	links := make([]config.LinkConfig, len(cfg.Links))
	copy(links, cfg.Links)

	// 简单的冒泡排序或者自定义排序，这里因为结构体在 config 包，不能直接定义方法
	// 除非把 LinkConfig 移出去或者使用 sort.Slice
	// 我们需要 import "sort"
	// 假设已经 import "sort"
	sort.Slice(links, func(i, j int) bool {
		return links[i].Order > links[j].Order // 降序：权重大的在前
	})

	t, err := parseTpl("layout.html", "links.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "links", "友情链接")
	data["Links"] = links
	t.ExecuteTemplate(w, "layout", data)
}

// LinkAdd 添加友情链接
func LinkAdd(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	url := r.FormValue("url")
	order, _ := strconv.Atoi(r.FormValue("order"))

	cfg := config.GetGlobalConfig()
	cfg.Links = append(cfg.Links, config.LinkConfig{Name: name, Url: url, Order: order})

	err := saveAndSortLinks()
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	jsonResponse(w, map[string]interface{}{"success": true, "message": "添加成功"})
}

// LinkEdit 编辑友情链接
func LinkEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		index, _ := strconv.Atoi(r.FormValue("index"))
		name := r.FormValue("name")
		url := r.FormValue("url")
		order, _ := strconv.Atoi(r.FormValue("order"))

		cfg := config.GetGlobalConfig()
		if index >= 0 && index < len(cfg.Links) {
			cfg.Links[index] = config.LinkConfig{Name: name, Url: url, Order: order}
			err := saveAndSortLinks()
			if err != nil {
				jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
				return
			}
			jsonResponse(w, map[string]interface{}{"success": true, "message": "修改成功"})
		} else {
			jsonResponse(w, map[string]interface{}{"success": false, "message": "索引无效"})
		}
		return
	}

	index, _ := strconv.Atoi(r.URL.Query().Get("index"))
	cfg := config.GetGlobalConfig()
	var link *config.LinkConfig
	if index >= 0 && index < len(cfg.Links) {
		link = &cfg.Links[index]
	}

	if link == nil {
		http.Error(w, "链接不存在", http.StatusNotFound)
		return
	}

	t, err := parseTpl("layout.html", "link_edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "links", "编辑友情链接")
	data["Link"] = link
	data["Index"] = index
	t.ExecuteTemplate(w, "layout", data)
}

// LinkDelete 删除友情链接
func LinkDelete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(r.FormValue("index"))

	cfg := config.GetGlobalConfig()
	if index >= 0 && index < len(cfg.Links) {
		cfg.Links = append(cfg.Links[:index], cfg.Links[index+1:]...)
	}

	err := saveAndSortLinks()
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	jsonResponse(w, map[string]interface{}{"success": true, "message": "删除成功"})
}

// saveAndSortLinks 保存并排序链接
func saveAndSortLinks() error {
	cfg := config.GetGlobalConfig()
	// 保存时统一排序：权重从大到小
	sort.Slice(cfg.Links, func(i, j int) bool {
		return cfg.Links[i].Order > cfg.Links[j].Order
	})
	return config.SaveLinkConfig("config/link.conf")
}

// Modules 模块设置页面
func Modules(w http.ResponseWriter, r *http.Request) {
	routerCfg := config.GetRouterConfig()
	appCfg := config.GetGlobalConfig()
	sorts, _ := dao.GetAllSorts()

	// 过滤路由：只显示允许自定义的路由
	allowedKeys := []string{"book", "book_index", "read", "sort", "top"}
	displayRoutes := make(map[string]string)
	for _, key := range allowedKeys {
		if val, ok := routerCfg.Routes[key]; ok {
			displayRoutes[key] = val
		}
	}

	t, err := parseTpl("layout.html", "modules.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "modules", "模块设置")
	data["Routes"] = displayRoutes
	data["SeoRules"] = appCfg.SeoRules
	data["Sorts"] = sorts
	t.ExecuteTemplate(w, "layout", data)
}

// Analytics 统计代码管理
func Analytics(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetGlobalConfig()

	if r.Method == "POST" {
		cfg.Analytics = r.FormValue("analytics")
		err := config.SaveAppConfig("config/config.conf")
		if err != nil {
			jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
			return
		}
		jsonResponse(w, map[string]interface{}{"success": true, "message": "保存成功"})
		return
	}

	t, err := parseTpl("layout.html", "analytics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "analytics", "统计代码")
	data["Analytics"] = cfg.Analytics
	t.ExecuteTemplate(w, "layout", data)
}

// ModuleRoutesUpdate 更新路由配置
func ModuleRoutesUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表单数据
	r.ParseForm()

	// 更新配置
	// 修正逻辑：只更新允许修改的路由，保留其他系统路由
	cfg := config.GetRouterConfig()

	allowedKeys := map[string]bool{
		"book":       true,
		"book_index": true,
		"read":       true,
		"sort":       true,
		"top":        true,
	}

	for key, values := range r.Form {
		if len(values) > 0 && allowedKeys[key] {
			cfg.Routes[key] = values[0]
		}
	}

	err := config.SaveRouterConfig("config/router.conf")
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	jsonResponse(w, map[string]interface{}{"success": true, "message": "路由配置保存成功"})
}

// ModuleSortsUpdate 更新分类配置
func ModuleSortsUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 接收 JSON 数据比较方便处理批量更新
	var sorts []struct {
		SortID    int    `json:"sortid"`
		Caption   string `json:"caption"`
		ShortName string `json:"shortname"`
		Weight    int    `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sorts); err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "无效的数据格式"})
		return
	}

	for _, s := range sorts {
		if err := dao.UpdateSort(s.SortID, s.Caption, s.ShortName, s.Weight); err != nil {
			jsonResponse(w, map[string]interface{}{"success": false, "message": "更新分类 " + s.Caption + " 失败: " + err.Error()})
			return
		}
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "分类保存成功"})
}

// ModuleSeoUpdate 更新 SEO 配置
func ModuleSeoUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 接收 JSON 数据
	var seoRules map[string]config.SeoRule
	if err := json.NewDecoder(r.Body).Decode(&seoRules); err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "无效的数据格式"})
		return
	}

	cfg := config.GetGlobalConfig()
	cfg.SeoRules = seoRules

	err := config.SaveSeoConfig("config/seo.conf")
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	jsonResponse(w, map[string]interface{}{"success": true, "message": "SEO 配置保存成功"})
}

// TestDBConnection 测试数据库连接
func TestDBConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	driver := r.FormValue("driver")
	host := r.FormValue("host")
	portStr := r.FormValue("port")
	user := r.FormValue("user")
	password := r.FormValue("password")
	dbname := r.FormValue("dbname")

	port, _ := strconv.Atoi(portStr)

	// 尝试连接
	err := dao.TestConnection(driver, host, port, user, password, dbname)
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "连接失败: " + err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "连接成功！"})
}

// Security 安全设置页面
func Security(w http.ResponseWriter, r *http.Request) {
	t, err := parseTpl("layout.html", "security.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "", "安全设置")
	t.ExecuteTemplate(w, "layout", data)
}

// SecurityPassword 修改密码
func SecurityPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := IsAdminLoggedIn(r)
	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	if newPassword != confirmPassword {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "两次输入的密码不一致"})
		return
	}

	// 验证原密码
	_, err := dao.VerifyAdminPassword(session.Username, oldPassword)
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "原密码错误"})
		return
	}

	// 更新密码
	err = dao.UpdateAdminPassword(session.AdminID, newPassword)
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "修改失败: " + err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "密码修改成功，请重新登录"})
}

// SecurityPath 修改后台入口
func SecurityPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	newPath := r.FormValue("admin_path")
	if newPath == "" {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "路径不能为空"})
		return
	}

	cfg := config.GetGlobalConfig()
	cfg.Site.AdminPath = newPath

	err := config.SaveAppConfig("config/config.conf")
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "保存配置失败: " + err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "入口已修改，正在跳转...", "new_path": newPath})
}

// TestRedisConnection 测试 Redis 连接
func TestRedisConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	host := r.FormValue("host")
	portStr := r.FormValue("port")
	password := r.FormValue("password")
	dbStr := r.FormValue("db")

	port, _ := strconv.Atoi(portStr)
	db, _ := strconv.Atoi(dbStr)

	// 创建测试配置
	testCfg := &config.RedisConfig{
		Enabled:  true,
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}

	// 尝试连接
	err := utils.InitRedis(testCfg)
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "连接失败: " + err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "连接成功！"})
}

// Plugins 插件管理页面
func Plugins(w http.ResponseWriter, r *http.Request) {
	pluginConfigs := plugin.GetManager().GetAllConfigs()

	t, err := parseTpl("layout.html", "plugins.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := getAdminData(r, "plugins", "插件管理")
	data["Plugins"] = pluginConfigs
	t.ExecuteTemplate(w, "layout", data)
}

// PluginToggle 启用/禁用插件
func PluginToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "无效的数据格式"})
		return
	}

	// 读取配置文件
	pluginConfigs, err := loadPluginConfigFile()
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "读取配置失败: " + err.Error()})
		return
	}

	// 更新启用状态
	if cfg, ok := pluginConfigs[req.Name]; ok {
		cfg["enabled"] = req.Enabled
		pluginConfigs[req.Name] = cfg
	} else {
		pluginConfigs[req.Name] = map[string]interface{}{"enabled": req.Enabled}
	}

	// 保存配置文件
	if err := savePluginConfigFile(pluginConfigs); err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "保存配置失败: " + err.Error()})
		return
	}

	// 热更新插件配置
	if err := plugin.GetManager().UpdatePluginConfig(req.Name, pluginConfigs[req.Name]); err != nil {
		jsonResponse(w, map[string]interface{}{"success": true, "message": "配置已保存，但热更新失败: " + err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "操作成功，已立即生效"})
}

// PluginConfigUpdate 更新插件配置
func PluginConfigUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name   string                 `json:"name"`
		Config map[string]interface{} `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "无效的数据格式"})
		return
	}

	// 读取配置文件
	pluginConfigs, err := loadPluginConfigFile()
	if err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "读取配置失败: " + err.Error()})
		return
	}

	// 更新配置
	if existing, ok := pluginConfigs[req.Name]; ok {
		// 保留 enabled 状态
		for key, value := range req.Config {
			existing[key] = value
		}
		pluginConfigs[req.Name] = existing
	} else {
		pluginConfigs[req.Name] = req.Config
	}

	// 保存配置文件
	if err := savePluginConfigFile(pluginConfigs); err != nil {
		jsonResponse(w, map[string]interface{}{"success": false, "message": "保存配置失败: " + err.Error()})
		return
	}

	// 热更新插件配置
	if err := plugin.GetManager().UpdatePluginConfig(req.Name, pluginConfigs[req.Name]); err != nil {
		jsonResponse(w, map[string]interface{}{"success": true, "message": "配置已保存，但热更新失败: " + err.Error()})
		return
	}

	jsonResponse(w, map[string]interface{}{"success": true, "message": "配置保存成功，已立即生效"})
}

// loadPluginConfigFile 读取插件配置文件
func loadPluginConfigFile() (map[string]map[string]interface{}, error) {
	file, err := os.Open("config/plugins.conf")
	if err != nil {
		return make(map[string]map[string]interface{}), nil
	}
	defer file.Close()

	var configs map[string]map[string]interface{}
	if err := json.NewDecoder(file).Decode(&configs); err != nil {
		return nil, err
	}
	return configs, nil
}

// savePluginConfigFile 保存插件配置文件
func savePluginConfigFile(configs map[string]map[string]interface{}) error {
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config/plugins.conf", data, 0644)
}
