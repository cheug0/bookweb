package controller

import (
	"bookweb/config"
	"bookweb/dao"
	"bookweb/model"
	"bookweb/plugin"
	"bookweb/utils"
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// LangtailUpdateFunc 长尾词更新回调函数（由插件设置）
var LangtailUpdateFunc func(sourceID int, sourceName string, cycleDays int)

// BookInfoData 书籍信息页面数据结构
type BookInfoData struct {
	Article        *model.Article
	SortName       string
	Chapters       []*model.Chapter
	LatestChapters []*model.Chapter
	LatestArticles []*model.Article
	HotArticles    []*model.Article
	Langtails      []*model.Langtail
}

// GetBookInfoData 获取书籍信息页面数据（供插件复用）
// 优化：使用并行查询减少总延迟
func GetBookInfoData(articleID int, articleName string) (*BookInfoData, error) {
	// 1. 先获取小说基本信息（必须先执行以获取 SortID）
	article, err := dao.GetArticleByIDCached(articleID)
	if err != nil {
		return nil, err
	}

	// 2. 并行执行剩余的独立查询
	var wg sync.WaitGroup
	var sortName = "全部分类"
	var chapters []*model.Chapter
	var latestArticles, hotArticles []*model.Article
	var langtails []*model.Langtail

	// 获取分类名称
	wg.Add(1)
	go func() {
		defer wg.Done()
		if sort, err := dao.GetSortByIDCached(article.SortID); err == nil {
			sortName = sort.Caption
		}
	}()

	// 获取章节目录
	wg.Add(1)
	go func() {
		defer wg.Done()
		chapters, _ = dao.GetChaptersByArticleIDCached(articleID)
		if chapters == nil {
			chapters = []*model.Chapter{}
		}
	}()

	// 获取最新文章
	wg.Add(1)
	go func() {
		defer wg.Done()
		latestArticles, _ = dao.GetArticlesBySortAndOrderCached(article.SortID, "postdate", 10)
	}()

	// 获取热门文章
	wg.Add(1)
	go func() {
		defer wg.Done()
		hotArticles, _ = dao.GetArticlesBySortAndOrderCached(article.SortID, "allvisit", 10)
	}()

	// 获取长尾词（如果插件启用）
	if plugin.GetManager().IsEnabled("langtail") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			langtails, _ = dao.GetLangtailsBySourceIDCached(articleID)
			// 如果没有长尾词或数据较少，异步抓取
			if len(langtails) < 3 && LangtailUpdateFunc != nil {
				go LangtailUpdateFunc(articleID, articleName, 7)
			}
		}()
	}

	wg.Wait()

	// 截取最新 12 条章节记录用于展示
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

	return &BookInfoData{
		Article:        article,
		SortName:       sortName,
		Chapters:       chapters,
		LatestChapters: latestChapters,
		LatestArticles: latestArticles,
		HotArticles:    hotArticles,
		Langtails:      langtails,
	}, nil
}

// BookInfo 小说信息页面
// BookInfo 小说信息页面
func BookInfo(w http.ResponseWriter, r *http.Request) {
	// 获取并校验参数
	articleID, ok := GetIDOr404(w, r, "aid")
	if !ok {
		return
	}

	// 增加点击量（排除爬虫）- 移到最前以确保即使命中缓存也能统计
	userAgent := r.UserAgent()
	if !utils.IsBot(userAgent) {
		go func() {
			dao.IncArticleVisit(articleID)
		}()
	}

	// 尝试从缓存获取整页 HTML (5分钟过期)
	// 使用 Redis 缓存页面，极大提升并发能力
	// 优先尝试 GZIP 缓存
	cacheKey := fmt.Sprintf("page_cache_book_%d", articleID)
	gzipCacheKey := cacheKey + "_gzip"

	useGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	// 检查是否开启小说信息页缓存
	if config.GetGlobalConfig().Site.BookCache && utils.IsRedisEnabled() {
		if useGzip {
			if cached, err := utils.CacheGet(gzipCacheKey); err == nil && cached != "" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("Content-Encoding", "gzip")
				w.Write([]byte(cached))
				return
			}
		}
		// 降级尝试普通缓存
		if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			// middleware 会自动压缩
			w.Write([]byte(cached))
			return
		}
	}

	// 获取书籍数据
	bookData, err := GetBookInfoData(articleID, "")
	if err != nil {
		NotFound(w, r)
		return
	}

	// 准备模版数据
	tags := map[string]string{
		"articlename": bookData.Article.ArticleName,
		"author":      bookData.Article.Author,
		"sortname":    bookData.SortName,
	}
	data := GetCommonData(r).
		ApplySeo("book_info", tags).
		Add("Article", bookData.Article).
		Add("SortName", bookData.SortName).
		Add("Chapters", bookData.Chapters).
		Add("LatestChapters", bookData.LatestChapters).
		Add("ChapterCount", len(bookData.Chapters)).
		Add("LatestArticles", bookData.LatestArticles).
		Add("HotArticles", bookData.HotArticles).
		Add("Langtails", bookData.Langtails)

	// 使用 Buffer 捕获渲染结果以便缓存
	var buf bytes.Buffer
	t := utils.GetTemplate("book_info.html")
	if t == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(&buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := buf.String()
	// 写入缓存 (5分钟)
	if config.GetGlobalConfig().Site.BookCache && utils.IsRedisEnabled() {
		utils.CacheSet(cacheKey, html, 5*time.Minute)

		// 同时预生成 GZIP 缓存
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write([]byte(html)); err == nil {
			if err := gz.Close(); err == nil {
				utils.CacheSet(gzipCacheKey, b.String(), 5*time.Minute)
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
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
	chapter, err := dao.GetChapterByIDCached(chapterID)
	if err != nil {
		NotFound(w, r)
		return
	}

	// 2. 获取小说信息 (用于展示名称、作者)
	article, err := dao.GetArticleByIDCached(articleID)
	if err != nil {
		NotFound(w, r)
		return
	}

	// 3. 上下页逻辑
	prevID, _ := dao.GetPrevChapterIDCached(articleID, chapter.ChapterOrder)
	nextID, _ := dao.GetNextChapterIDCached(articleID, chapter.ChapterOrder)

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

	t := utils.GetTemplate("book_reader.html")
	if t == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
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

	// 优先尝试 GZIP 缓存
	cacheKey := fmt.Sprintf("page_cache_index_%d", articleID)
	gzipCacheKey := cacheKey + "_gzip"

	useGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	// 检查是否开启小说目录页缓存
	if config.GetGlobalConfig().Site.BookIndexCache && utils.IsRedisEnabled() {
		if useGzip {
			if cached, err := utils.CacheGet(gzipCacheKey); err == nil && cached != "" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("Content-Encoding", "gzip")
				w.Write([]byte(cached))
				return
			}
		}
		// 降级尝试普通缓存
		if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			// middleware 会自动压缩
			w.Write([]byte(cached))
			return
		}
	}

	// 1. 获取小说基本信息（带缓存）
	article, err := dao.GetArticleByIDCached(articleID)
	if err != nil {
		NotFound(w, r)
		return
	}

	// 2. 获取章节目录（带缓存）
	chapters, err := dao.GetChaptersByArticleIDCached(articleID)
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

	// 5. 渲染页面
	var buf bytes.Buffer
	t := utils.GetTemplate("book_list.html")
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
	if config.GetGlobalConfig().Site.BookIndexCache && utils.IsRedisEnabled() {
		utils.CacheSet(cacheKey, html, 10*time.Minute)

		// 同时预生成 GZIP 缓存
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write([]byte(html)); err == nil {
			if err := gz.Close(); err == nil {
				utils.CacheSet(gzipCacheKey, b.String(), 10*time.Minute)
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
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
