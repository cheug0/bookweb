// sort.go
// 分类控制器
// 处理小说分类列表页的展示与分页
package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/utils"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// SortList 处理小说分类页面请求
func SortList(w http.ResponseWriter, r *http.Request) {
	// 1. 获取并校验动态参数
	sortID, ok := GetIDOr404(w, r, "sid")
	if !ok {
		return
	}
	currentPage, ok := GetIDOr404(w, r, "page")
	if !ok {
		return
	}

	if currentPage < 1 {
		NotFound(w, r)
		return
	}

	// 尝试从缓存获取 (10分钟)
	cacheKey := fmt.Sprintf("page_cache_sort_%d_%d", sortID, currentPage)
	if config.GetGlobalConfig().Site.SortCache && utils.IsRedisEnabled() {
		if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(cached))
			return
		}
	}

	// 2. 准备基础数据
	pageSize := 20
	offset := (currentPage - 1) * pageSize

	// 获取所有分类 (侧边栏内容，带缓存)
	allSorts, err := dao.GetAllSortsCached()
	if err != nil {
		NotFound(w, r)
		return
	}

	// 3. 获取当前分类及对应的小说列表
	articles, err := dao.GetArticlesBySortID(sortID, offset, pageSize)
	if err != nil {
		NotFound(w, r)
		return
	}

	totalCount, err := dao.GetArticleCountBySortID(sortID)
	if err != nil {
		totalCount = 0
	}

	caption := "全部分类"
	if sortID > 0 {
		currentSort, err := dao.GetSortByIDCached(sortID)
		if err == nil {
			caption = currentSort.Caption
		} else {
			// 如果 ID 不存在，返回 404
			NotFound(w, r)
			return
		}
	}

	// 4. 构建模板数据
	totalPage := (totalCount + pageSize - 1) / pageSize
	if totalPage == 0 {
		totalPage = 1
	}

	// 严格校验：请求页码不能大于总页数
	if currentPage > totalPage {
		NotFound(w, r)
		return
	}

	// 应用标签化 SEO
	tags := map[string]string{
		"sortname": caption,
		"page":     strconv.Itoa(currentPage),
	}

	data := GetCommonData(r).
		ApplySeo("sort", tags).
		Add("Sorts", allSorts).
		Add("CurrentSID", sortID).
		Add("Caption", caption).
		Add("Articles", articles).
		Add("Page", currentPage).
		Add("TotalPage", totalPage).
		Add("TotalCount", totalCount).
		Add("Pages", generatePageList(currentPage, totalPage))

	// 5. 渲染页面
	var buf bytes.Buffer
	t := GetRenderTemplate(w, r, "sort.html")
	if t == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(&buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := buf.String()
	// 写入缓存 (10分钟)
	if config.GetGlobalConfig().Site.SortCache && utils.IsRedisEnabled() {
		utils.CacheSet(cacheKey, html, 10*time.Minute)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
