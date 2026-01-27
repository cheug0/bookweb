// controller.go (langtail)
// 长尾词控制器
// 处理长尾词落地页的展示请求
package langtail

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/model"
	"bookweb/utils"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var pluginConfig *Config

// SetConfig 设置插件配置
func SetConfig(cfg *Config) {
	pluginConfig = cfg
}

// GetConfig 获取插件配置
func GetConfig() *Config {
	return pluginConfig
}

// LangtailInfo 长尾词信息页面
func LangtailInfo(w http.ResponseWriter, r *http.Request) {
	// 获取长尾词ID
	langID := getLangIDFromRequest(r)
	if langID <= 0 {
		http.NotFound(w, r)
		return
	}

	// 1. 获取长尾词信息
	langtailItem, err := dao.GetLangtailByID(langID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 2. 获取小说信息
	article, err := dao.GetArticleByIDCached(langtailItem.SourceID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 3. 异步更新长尾词
	if pluginConfig != nil {
		go func() {
			UpdateLangtailsIfNeeded(langtailItem.SourceID, langtailItem.SourceName, pluginConfig.FetchCycleDays)
		}()
	}

	// 4. 获取分类名称
	sortName := "全部分类"
	if sort, err := dao.GetSortByID(article.SortID); err == nil {
		sortName = sort.Caption
	}

	// 5. 获取章节
	chapters, _ := dao.GetChaptersByArticleIDCached(langtailItem.SourceID)
	if chapters == nil {
		chapters = []*model.Chapter{}
	}

	// 截取最新章节
	latestChapters := chapters
	if len(chapters) > 12 {
		latestChapters = chapters[len(chapters)-12:]
		reversed := make([]*model.Chapter, len(latestChapters))
		for i, v := range latestChapters {
			reversed[len(latestChapters)-1-i] = v
		}
		latestChapters = reversed
	}

	// 6. 获取长尾词列表和推荐文章
	if pluginConfig == nil {
		utils.LogWarn("Langtail", "pluginConfig is nil in LangtailInfo")
	} else {
		utils.LogInfo("Langtail", "LangtailInfo: show_count=%d", pluginConfig.ShowCount)
	}
	langtails, _ := dao.GetLangtailsBySourceID(langtailItem.SourceID)
	if pluginConfig != nil && pluginConfig.ShowCount > 0 && len(langtails) > pluginConfig.ShowCount {
		langtails = langtails[:pluginConfig.ShowCount]
	}
	latestArticles, _ := dao.GetArticlesBySortAndOrder(article.SortID, "postdate", 10)
	hotArticles, _ := dao.GetArticlesBySortAndOrder(article.SortID, "allvisit", 10)

	// 7. 获取导航栏分类链接
	sorts, _ := dao.GetAllSorts()
	var sortLinks []map[string]string
	sortMap := make(map[int]string)
	sortRoute := config.GetRouterConfig().GetRoute("sort")
	for _, s := range sorts {
		sortMap[s.SortID] = s.Caption
		url := sortRoute
		url = strings.ReplaceAll(url, ":sid", fmt.Sprintf("%d", s.SortID))
		url = strings.ReplaceAll(url, ":page", "1")
		sortLinks = append(sortLinks, map[string]string{
			"Caption": s.Caption,
			"Url":     url,
		})
	}

	// 准备模版数据（使用正确的字段名）
	siteName := config.GlobalConfig.Site.SiteName
	siteDomain := config.GlobalConfig.Site.Domain
	topRoute := config.GetRouterConfig().GetRoute("top")

	// 处理耗时函数
	startTime := time.Now()
	processingComment := func() template.HTML {
		duration := time.Since(startTime).Seconds()
		return template.HTML(fmt.Sprintf("<!-- Processed in %.6f second(s) -->", duration))
	}

	// 安全获取配置
	showCount := 50
	if pluginConfig != nil {
		showCount = pluginConfig.ShowCount
	}

	data := map[string]interface{}{
		// SEO 字段
		"CurrentTitle":    langtailItem.LangName + " - " + siteName,
		"CurrentKeywords": langtailItem.LangName + "," + article.ArticleName,
		"CurrentDesc":     langtailItem.LangName + " - " + article.Intro,
		// 站点信息
		"SiteName":          siteName,
		"SiteDomain":        siteDomain,
		"SortLinks":         sortLinks,
		"SortMap":           sortMap,
		"TopUrl":            topRoute,
		"IsLogin":           false,
		"Analytics":         template.HTML(config.GlobalConfig.Analytics),
		"ProcessingComment": processingComment,
		// 文章信息
		"Article":        article,
		"ArticleName":    langtailItem.LangName,
		"IsLangtail":     true,
		"Langtail":       langtailItem,
		"Langtails":      langtails,
		"SortName":       sortName,
		"Chapters":       chapters,
		"LatestChapters": latestChapters,
		"ChapterCount":   len(chapters),
		"LatestArticles": latestArticles,
		"HotArticles":    hotArticles,
		"ShowCount":      showCount,
	}

	// 检查用户登录状态
	isLogin, sess := dao.IsLogin(r)
	if isLogin {
		data["IsLogin"] = true
		data["Username"] = sess.Username
	}

	tplPath := getTplPath("book_info.html")
	t := template.New("book_info.html").Funcs(utils.CommonFuncMap)
	t, err = t.ParseFiles(tplPath, getTplPath("head.html"), getTplPath("foot.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// getLangIDFromRequest 从请求中获取长尾词ID
func getLangIDFromRequest(r *http.Request) int {
	if ps := r.Context().Value(model.ParamsKey); ps != nil {
		params := ps.(httprouter.Params)
		if lid := params.ByName("lid"); lid != "" {
			id, _ := strconv.Atoi(lid)
			return id
		}
	}
	return 0
}

// getTplPath 获取模板路径
func getTplPath(name string) string {
	tpl := config.GlobalConfig.Site.Template
	return "template/" + tpl + "/" + name
}
