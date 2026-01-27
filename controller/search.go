// search.go
// 搜索控制器
// 处理小说搜索请求，支持按书名或作者搜索
package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"net/http"
	"strconv"
	"time"
)

// Search 处理小说搜索请求
func Search(w http.ResponseWriter, r *http.Request) {
	// 1. 获取关键词和页码
	keyword := r.URL.Query().Get("key")
	if keyword == "" {
		// 如果关键词为空，重定向到首页或显示错误
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	pageStr := r.URL.Query().Get("page")
	currentPage, err := strconv.Atoi(pageStr)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}

	// 限制搜索频率
	limit := config.GetGlobalConfig().Site.SearchLimit
	if limit > 0 {
		cookie, err := r.Cookie("last_search_time")
		if err == nil {
			lastTime, _ := strconv.ParseInt(cookie.Value, 10, 64)
			if time.Now().Unix()-lastTime < int64(limit) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Write([]byte(`<script>alert("搜索过于频繁，请稍后再试");window.history.back();</script>`))
				return
			}
		}
		// 设置新的搜索时间 Cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "last_search_time",
			Value: strconv.FormatInt(time.Now().Unix(), 10),
			Path:  "/",
		})
	}

	// 2. 准备分页数据
	pageSize := 20
	offset := (currentPage - 1) * pageSize

	// 3. 执行搜索查询 (使用缓存)
	articles, err := dao.SearchArticlesCached(keyword, offset, pageSize)
	if err != nil {
		http.Error(w, "搜索请求失败", http.StatusInternalServerError)
		return
	}

	totalCount, err := dao.GetSearchCountCached(keyword)
	if err != nil {
		totalCount = 0
	}

	// 4. 计算总页数
	totalPage := (totalCount + pageSize - 1) / pageSize
	if totalPage == 0 {
		totalPage = 1
	}

	// 5. 应用 SEO 和通用数据
	tags := map[string]string{
		"keyword": keyword,
		"page":    strconv.Itoa(currentPage),
	}
	data := GetCommonData(r).
		ApplySeo("search", tags).
		Add("Keyword", keyword).
		Add("Articles", articles).
		Add("Page", currentPage).
		Add("TotalPage", totalPage).
		Add("TotalCount", totalCount).
		Add("Pages", generatePageList(currentPage, totalPage))

	// 6. 渲染页面
	t := GetRenderTemplate(w, r, "search.html")
	if t == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
