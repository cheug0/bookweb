package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/utils"
	"html/template"
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

	// 6. 定义模版函数 (复用 sort.go 的逻辑)
	funcMap := template.FuncMap{
		"plus":  func(a, b int) int { return a + b },
		"minus": func(a, b int) int { return a - b },
		"formatSize": func(size int) string {
			if size >= 10000 {
				return strconv.FormatFloat(float64(size)/10000.0, 'f', 1, 64) + "万"
			}
			return strconv.Itoa(size)
		},
		"formatDate": func(t int64) string {
			return utils.FormatDate(t, "2006-01-02")
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

	// 7. 渲染页面
	tPath, ok := GetTplPathOrError(w, "search.html")
	if !ok {
		return
	}
	t := template.New("search.html").Funcs(funcMap)
	t, err = t.ParseFiles(tPath, TplPath("head.html"), TplPath("foot.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
