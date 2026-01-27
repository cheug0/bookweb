// index.go
// 首页控制器
// 处理网站首页展示，包括推荐、分类、最新更新等数据的聚合
package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/model"
	"bookweb/utils"
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

	// 1. 大神小说 (使用配置)
	topCfg := config.GetGlobalConfig().Recommend.Top
	topArticles := getRecommendedArticles(topCfg.Picks, topCfg.Sort, topCfg.Limit)

	// 2. 热门小说 (使用配置)
	hotCfg := config.GetGlobalConfig().Recommend.Hot
	hotArticles := getRecommendedArticles(hotCfg.Picks, hotCfg.Sort, hotCfg.Limit)

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

	t := GetRenderTemplate(w, r, "index.html")
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

// getRecommendedArticles 获取推荐小说（合并 Picks 和 Sort）
func getRecommendedArticles(picksStr string, sortBy string, limit int) []*model.Article {
	var result []*model.Article
	existingIDs := make(map[int]bool)

	// 1. 获取 Picks
	if picksStr != "" {
		var ids []int
		parts := strings.Split(picksStr, ",")
		for _, p := range parts {
			// 支持去除首尾空格
			p = strings.TrimSpace(p)
			if id, err := strconv.Atoi(p); err == nil && id > 0 {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			pickedArts, _ := dao.GetArticlesByIDs(ids)
			// 按输入的 ID 顺序添加到结果
			for _, id := range ids {
				for _, art := range pickedArts {
					if art.ArticleID == id {
						result = append(result, art)
						existingIDs[art.ArticleID] = true
						break
					}
				}
			}
		}
	}

	// 如果已达到限制，直接返回
	if len(result) >= limit {
		return result[:limit]
	}

	// 2. 获取补充的排序文章
	// 稍微多取几个防止重复过滤后不足 (假设最多取 limit 个补充)
	sortedArts, _ := dao.GetArticlesBySortAndOrderCached(0, sortBy, limit)

	for _, art := range sortedArts {
		if !existingIDs[art.ArticleID] {
			result = append(result, art)
			existingIDs[art.ArticleID] = true
			if len(result) >= limit {
				break
			}
		}
	}

	return result
}
