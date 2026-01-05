package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/model"
	"bookweb/utils"
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

// Index 处理首页请求
func Index(w http.ResponseWriter, r *http.Request) {
	data := GetCommonData(r).ApplySeo("index", nil)

	// 1. 大神小说 (取 allvisit 前 6)
	topArticles, _ := dao.GetVisitArticles(6)

	// 2. 热门小说 (取 allvisit 前 12)
	hotArticles, _ := dao.GetVisitArticles(12)

	// 3. 分类展示 (前 6 个分类，每个取 13 本)
	allSorts, _ := dao.GetAllSorts()
	sortMap := make(map[int]string)
	var blocks []categoryBlock
	for i, s := range allSorts {
		sortMap[s.SortID] = s.Caption
		if i < 6 {
			arts, _ := dao.GetArticlesBySortID(s.SortID, 0, 13)
			blocks = append(blocks, categoryBlock{
				Sort:     s,
				Articles: arts,
			})
		}
	}

	// 4. 最新章节 (取 lastupdate 前 30)
	latestUpdates, _ := dao.GetArticlesBySortID(0, 0, 30)

	// 5. 最新入库 (取 postdate 前 30)
	newArticles, _ := dao.GetArticlesBySortID(0, 0, 30)

	data.Add("TopArticles", topArticles).
		Add("HotArticles", hotArticles).
		Add("Blocks", blocks).
		Add("LatestUpdates", latestUpdates).
		Add("NewArticles", newArticles).
		Add("SortMap", sortMap).
		Add("Links", config.GetGlobalConfig().Links)

	tPath, ok := GetTplPathOrError(w, "index.html")
	if !ok {
		return
	}
	t := template.New("index.html").Funcs(indexFuncMap)
	t, err := t.ParseFiles(tPath, TplPath("head.html"), TplPath("foot.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
