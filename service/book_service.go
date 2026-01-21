package service

import (
	"bookweb/dao"
	"bookweb/model"
	"bookweb/plugin"
	"sync"
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

// GetBookInfoData 获取书籍信息页面数据
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
