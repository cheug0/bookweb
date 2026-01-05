package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/model"
	"bookweb/utils"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

// bookFuncMap 书籍页面的模版函数
var bookFuncMap = template.FuncMap{
	"formatSize": func(size int) string {
		if size >= 10000 {
			return fmt.Sprintf("%.1f万", float64(size)/10000.0)
		}
		return fmt.Sprintf("%d", size)
	},
	"formatDate": func(t int64) string {
		if t == 0 {
			return "-"
		}
		return time.Unix(t, 0).Format("2006-01-02")
	},
	"safe": func(s string) template.HTML {
		return template.HTML(s)
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

// BookInfo 小说信息页面
func BookInfo(w http.ResponseWriter, r *http.Request) {
	// 获取并校验参数
	articleID, ok := GetIDOr404(w, r, "aid")
	if !ok {
		return
	}

	// 1. 获取小说基本信息
	article, err := dao.GetArticleByID(articleID)
	if err != nil {
		NotFound(w, r)
		return
	}

	// 2. 获取分类名称
	sortName := "全部分类"
	sort, err := dao.GetSortByID(article.SortID)
	if err == nil {
		sortName = sort.Caption
	}

	// 3. 获取章节目录
	chapters, err := dao.GetChaptersByArticleID(articleID)
	if err != nil {
		chapters = []*model.Chapter{}
	}

	// 截取最新 12 条记录用于展示
	latestChapters := chapters
	if len(chapters) > 12 {
		latestChapters = chapters[len(chapters)-12:]
		// 逆序排列，让最新的在最上面
		reversedLatest := make([]*model.Chapter, len(latestChapters))
		for i, v := range latestChapters {
			reversedLatest[len(latestChapters)-1-i] = v
		}
		latestChapters = reversedLatest
	}

	latestArticles, _ := dao.GetArticlesBySortAndOrder(article.SortID, "postdate", 10)
	hotArticles, _ := dao.GetArticlesBySortAndOrder(article.SortID, "allvisit", 10)

	// 准备模版数据
	// 应用标签化 SEO
	tags := map[string]string{
		"articlename": article.ArticleName,
		"author":      article.Author,
		"sortname":    sortName,
	}
	data := GetCommonData(r).
		ApplySeo("book_info", tags).
		Add("Article", article).
		Add("SortName", sortName).
		Add("Chapters", chapters).
		Add("LatestChapters", latestChapters).
		Add("ChapterCount", len(chapters)).
		Add("LatestArticles", latestArticles).
		Add("HotArticles", hotArticles)

	tPath, ok := GetTplPathOrError(w, "book_info.html")
	if !ok {
		return
	}
	t := template.New("book_info.html").Funcs(bookFuncMap)
	t, err = t.ParseFiles(tPath, TplPath("head.html"), TplPath("foot.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// ChapterRead 章节阅读页面
func ChapterRead(w http.ResponseWriter, r *http.Request) {
	// 获取并校验参数
	articleID, ok := GetIDOr404(w, r, "aid")
	if !ok {
		return
	}
	chapterID, ok := GetIDOr404(w, r, "cid")
	if !ok {
		return
	}

	// 1. 获取章节内容
	chapter, err := dao.GetChapterByID(chapterID)
	if err != nil {
		NotFound(w, r)
		return
	}

	// 2. 获取小说信息 (用于展示名称、作者)
	article, err := dao.GetArticleByID(articleID)
	if err != nil {
		NotFound(w, r)
		return
	}

	// 3. 上下页逻辑
	prevID, _ := dao.GetPrevChapterID(articleID, chapter.ChapterOrder)
	nextID, _ := dao.GetNextChapterID(articleID, chapter.ChapterOrder)

	// 准备数据
	// 应用标签化 SEO
	tags := map[string]string{
		"articlename": article.ArticleName,
		"chaptername": chapter.ChapterName,
	}
	data := GetCommonData(r).
		ApplySeo("book_reader", tags).
		Add("Article", article).
		Add("Chapter", chapter).
		Add("PrevID", prevID).
		Add("NextID", nextID)

	tPath, ok := GetTplPathOrError(w, "book_reader.html")
	if !ok {
		return
	}
	t := template.New("book_reader.html").Funcs(bookFuncMap)
	t, err = t.ParseFiles(tPath, TplPath("head.html"), TplPath("foot.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// BookIndex 小说目录页
func BookIndex(w http.ResponseWriter, r *http.Request) {
	// 获取并校验参数
	articleID, ok := GetIDOr404(w, r, "aid")
	if !ok {
		return
	}

	// 1. 获取小说基本信息
	article, err := dao.GetArticleByID(articleID)
	if err != nil {
		NotFound(w, r)
		return
	}

	// 2. 获取章节目录
	chapters, err := dao.GetChaptersByArticleID(articleID)
	if err != nil {
		chapters = []*model.Chapter{}
	}

	// 准备数据
	tags := map[string]string{
		"articlename": article.ArticleName,
		"author":      article.Author,
	}
	data := GetCommonData(r).
		ApplySeo("book_index", tags).
		Add("Article", article).
		Add("Chapters", chapters).
		Add("ChapterCount", len(chapters))

	tPath, ok := GetTplPathOrError(w, "book_list.html")
	if !ok {
		return
	}
	t := template.New("book_list.html").Funcs(bookFuncMap)
	t, err = t.ParseFiles(tPath, TplPath("head.html"), TplPath("foot.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// CoverImage 处理封面图片请求
// 访问路径: /img/:aid.jpg
func CoverImage(w http.ResponseWriter, r *http.Request) {
	articleID, ok := GetIDOr404(w, r, "aid")
	if !ok {
		return
	}

	cfg := config.GetGlobalConfig()
	// 如果配置了 OSS 且有域名，重定向到远程地址
	if cfg != nil && cfg.Storage.Type == "oss" && cfg.Storage.Oss.Domain != "" {
		http.Redirect(w, r, utils.GetCoverPath(articleID), http.StatusFound)
		return
	}

	// 本地/NFS 模式：拼接物理基准路径
	relPath := utils.GetPhysicalCoverPath(articleID)
	basePath := "files"
	if cfg != nil && cfg.Storage.Local.Path != "" {
		basePath = cfg.Storage.Local.Path
	}
	fullPath := filepath.Join(basePath, relPath)

	// 设置强缓存和 Content-Type
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Content-Type", "image/jpeg")

	// 直接发送文件
	http.ServeFile(w, r, fullPath)
}
