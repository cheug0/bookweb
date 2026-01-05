package controller

import (
	"bookweb/dao"
	"bookweb/utils"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// funcMap 定义模版辅助函数
var funcMap = template.FuncMap{
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
		return time.Unix(t, 0).Format("2006-01-02")
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
		currentSort, err := dao.GetSortByID(sortID)
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
	tPath, ok := GetTplPathOrError(w, "sort.html")
	if !ok {
		return
	}
	t := template.New("sort.html").Funcs(funcMap)
	t, err = t.ParseFiles(tPath, TplPath("head.html"), TplPath("foot.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
