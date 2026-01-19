package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/model"
	"bookweb/utils"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// indexFuncMap 首页模版函数
var indexFuncMap = template.FuncMap{
	"plus":  func(a, b int) int { return a + b },
	"minus": func(a, b int) int { return a - b },
	"formatSize": func(size int) string {
		if size >= 10000 {
			return fmt.Sprintf("%.1f万", float64(size)/10000.0)
		}
		return strconv.Itoa(size)
	},
	"formatDate": func(t int64) string {
		if t == 0 {
			return "-"
		}
		// 首页通常只显示 月-日
		return time.Unix(t, 0).Format("01-02")
	},
	"cover": func(id int) string {
		return utils.GetCoverPath(id)
	},
	"bookUrl": func(id int) string {
		return utils.BookUrl(id)
	},
	"readUrl": func(aid, cid int) string {
		return utils.ReadUrl(aid, cid)
	},
}

// categoryBlock 分类块数据结构
type categoryBlock struct {
	Sort     *model.Sort
	Articles []*model.Article
}

const indexCacheKey = "page_cache_index"

// Index 处理首页请求
func Index(w http.ResponseWriter, r *http.Request) {
	// 尝试从缓存获取整页HTML (如果开启)
	if config.GetGlobalConfig().Site.IndexCache {
		if cached, err := utils.CacheGet(indexCacheKey); err == nil && cached != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(cached))
			return
		}
	}

	data := GetCommonData(r).ApplySeo("index", nil)

	// 1. 大神小说 (取 allvisit 前 6，带缓存)
	topArticles, _ := dao.GetVisitArticlesCached(6)

	// 2. 热门小说 (取 allvisit 前 12，带缓存)
	hotArticles, _ := dao.GetVisitArticlesCached(12)

	// 3. 分类展示 (前 6 个分类，每个取 13 本，带缓存)
	allSorts, _ := dao.GetAllSortsCached()
	sortMap := make(map[int]string)
	var blocks []categoryBlock
	for i, s := range allSorts {
		sortMap[s.SortID] = s.Caption
		if i < 6 {
			arts, _ := dao.GetArticlesBySortIDCached(s.SortID, 0, 13)
			blocks = append(blocks, categoryBlock{
				Sort:     s,
				Articles: arts,
			})
		}
	}

	// 4. 最新更新 (取 lastupdate 前 30，带缓存)
	latestUpdates, _ := dao.GetArticlesBySortIDCached(0, 0, 30)

	// 5. 最新入库 (复用最新更新数据，避免重复查询)
	newArticles := latestUpdates

	data.Add("TopArticles", topArticles).
		Add("HotArticles", hotArticles).
		Add("Blocks", blocks).
		Add("LatestUpdates", latestUpdates).
		Add("NewArticles", newArticles).
		Add("SortMap", sortMap).
		Add("Links", config.GetGlobalConfig().Links)

	t := utils.GetTemplate("index.html")
	if t == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// 渲染到缓冲区并缓存
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := buf.String()
	// 缓存整页HTML（1分钟过期，如果开启）
	if config.GetGlobalConfig().Site.IndexCache {
		utils.CacheSet(indexCacheKey, html, 1*time.Minute)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
