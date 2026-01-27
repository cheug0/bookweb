// article_dao.go
// 文章 DAO
// 处理小说信息的增删改查及点击量统计
package dao

import (
	"bookweb/model"
	"bookweb/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// GetArticleByID 根据ArticleID获取小说信息
func GetArticleByID(id int) (*model.Article, error) {
	var row *sql.Row
	if stmtGetArticleByID != nil {
		row = stmtGetArticleByID.QueryRow(id)
	} else {
		sqlStr := "select articleid, siteid, postdate, lastupdate, articlename, keywords, initial, authorid, author, posterid, poster, agentid, agent, sortid, typeid, intro, notice, setting, lastvolumeid, lastvolume, lastchapterid, lastchapter, chapters, size, lastvisit, dayvisit, weekvisit, monthvisit, allvisit, lastvote, dayvote, weekvote, monthvote, allvote, fullflag, imgflag from jieqi_article_article where articleid = ?"
		row = utils.Db.QueryRow(sqlStr, id)
	}
	art := &model.Article{}
	err := row.Scan(&art.ArticleID, &art.SiteID, &art.PostDate, &art.LastUpdate, &art.ArticleName, &art.Keywords, &art.Initial, &art.AuthorID, &art.Author, &art.PosterID, &art.Poster, &art.AgentID, &art.Agent, &art.SortID, &art.TypeID, &art.Intro, &art.Notice, &art.Setting, &art.LastVolumeID, &art.LastVolume, &art.LastChapterID, &art.LastChapter, &art.Chapters, &art.Size, &art.LastVisit, &art.DayVisit, &art.WeekVisit, &art.MonthVisit, &art.AllVisit, &art.LastVote, &art.DayVote, &art.WeekVote, &art.MonthVote, &art.AllVote, &art.FullFlag, &art.ImgFlag)
	if err != nil {
		return nil, err
	}
	return art, nil
}

// GetChaptersByArticleID 根据ArticleID获取章节列表(不含内容)
func GetChaptersByArticleID(articleID int) ([]*model.Chapter, error) {
	var rows *sql.Rows
	var err error
	if stmtGetChaptersByArticle != nil {
		rows, err = stmtGetChaptersByArticle.Query(articleID)
	} else {
		sqlStr := "select chapterid, chaptername, chapterorder, isvip, size, lastupdate from jieqi_article_chapter where articleid = ? order by chapterorder asc"
		rows, err = utils.Db.Query(sqlStr, articleID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []*model.Chapter
	for rows.Next() {
		ch := &model.Chapter{}
		err := rows.Scan(&ch.ChapterID, &ch.ChapterName, &ch.ChapterOrder, &ch.IsVIP, &ch.Size, &ch.LastUpdate)
		if err != nil {
			return nil, err
		}
		chapters = append(chapters, ch)
	}
	return chapters, nil
}

// GetArticlesBySortID 分页获取指定分类的小说列表
// 如果 sortID 为 0，则获取全部分类
func GetArticlesBySortID(sortID int, offset, limit int) ([]*model.Article, error) {
	sqlStr := "select articleid, articlename, author, intro, size, lastupdate, sortid, fullflag, imgflag, lastchapterid, lastchapter from jieqi_article_article"
	var args []interface{}
	if sortID > 0 {
		sqlStr += " where sortid = ?"
		args = append(args, sortID)
	}
	sqlStr += " order by lastupdate desc limit ?, ?"
	args = append(args, offset, limit)

	rows, err := utils.Db.Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		art := &model.Article{}
		// 仅扫描列表页需要的字段以提升性能
		err := rows.Scan(&art.ArticleID, &art.ArticleName, &art.Author, &art.Intro, &art.Size, &art.LastUpdate, &art.SortID, &art.FullFlag, &art.ImgFlag, &art.LastChapterID, &art.LastChapter)
		if err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil
}

// GetArticlesBySortIDCached 带缓存分页获取指定分类的小说列表
func GetArticlesBySortIDCached(sortID int, offset, limit int) ([]*model.Article, error) {
	cacheKey := fmt.Sprintf("articles_sort_%d_%d_%d", sortID, offset, limit)

	if !utils.IsRedisEnabled() {
		return GetArticlesBySortID(sortID, offset, limit)
	}

	// 尝试从 Redis 缓存获取
	if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
		var articles []*model.Article
		if err := json.Unmarshal([]byte(cached), &articles); err == nil {
			return articles, nil
		}
	}

	// 从数据库获取
	articles, err := GetArticlesBySortID(sortID, offset, limit)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(articles); err == nil {
		utils.CacheSet(cacheKey, string(data), 5*time.Minute)
	}
	return articles, nil
}

// GetArticlesByIDs 批量获取小说信息
func GetArticlesByIDs(ids []int) ([]*model.Article, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// 构建 IN 查询
	sqlStr := "select articleid, articlename, author, intro, size, lastupdate, sortid, fullflag, imgflag, lastchapterid, lastchapter from jieqi_article_article where articleid in ("
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		if i > 0 {
			sqlStr += ","
		}
		sqlStr += "?"
		args[i] = id
	}
	sqlStr += ")"

	rows, err := utils.Db.Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 使用 Map 暂存结果以便按 ID 顺序返回
	artMap := make(map[int]*model.Article)
	for rows.Next() {
		art := &model.Article{}
		err := rows.Scan(&art.ArticleID, &art.ArticleName, &art.Author, &art.Intro, &art.Size, &art.LastUpdate, &art.SortID, &art.FullFlag, &art.ImgFlag, &art.LastChapterID, &art.LastChapter)
		if err != nil {
			continue
		}
		artMap[art.ArticleID] = art
	}

	var articles []*model.Article
	for _, id := range ids {
		if art, ok := artMap[id]; ok {
			articles = append(articles, art)
		}
	}
	return articles, nil
}

// GetArticleCountBySortID 获取指定分类的小说总数
func GetArticleCountBySortID(sortID int) (int, error) {
	sqlStr := "select count(*) from jieqi_article_article"
	var args []interface{}
	if sortID > 0 {
		sqlStr += " where sortid = ?"
		args = append(args, sortID)
	}

	var count int
	err := utils.Db.QueryRow(sqlStr, args...).Scan(&count)
	return count, err
}

// GetVisitArticles 获取按点击量排序的小说列表
func GetVisitArticles(limit int) ([]*model.Article, error) {
	var rows *sql.Rows
	var err error
	if stmtGetVisitArticles != nil {
		rows, err = stmtGetVisitArticles.Query(limit)
	} else {
		sqlStr := "select articleid, articlename, author, intro, size, lastupdate, sortid, fullflag, imgflag, lastchapterid, lastchapter from jieqi_article_article order by allvisit desc limit ?"
		rows, err = utils.Db.Query(sqlStr, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		art := &model.Article{}
		err := rows.Scan(&art.ArticleID, &art.ArticleName, &art.Author, &art.Intro, &art.Size, &art.LastUpdate, &art.SortID, &art.FullFlag, &art.ImgFlag, &art.LastChapterID, &art.LastChapter)
		if err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil
}

// GetVisitArticlesCached 带缓存获取按点击量排序的小说列表
func GetVisitArticlesCached(limit int) ([]*model.Article, error) {
	cacheKey := fmt.Sprintf("visit_articles_%d", limit)

	if !utils.IsRedisEnabled() {
		return GetVisitArticles(limit)
	}

	// 尝试从 Redis 缓存获取
	if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
		var articles []*model.Article
		if err := json.Unmarshal([]byte(cached), &articles); err == nil {
			return articles, nil
		}
	}

	articles, err := GetVisitArticles(limit)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(articles); err == nil {
		utils.CacheSet(cacheKey, string(data), 5*time.Minute)
	}
	return articles, nil
}

// GetArticlesBySortAndOrder 获取分类小说列表 (按指定字段排序)
func GetArticlesBySortAndOrder(sortID int, orderBy string, limit int) ([]*model.Article, error) {
	// 简单的白名单校验
	validFields := map[string]bool{
		"allvisit":   true,
		"monthvisit": true,
		"weekvisit":  true,
		"dayvisit":   true,
		"allvote":    true,
		"size":       true,
		"lastupdate": true,
		"postdate":   true,
	}
	if !validFields[orderBy] {
		orderBy = "allvisit"
	}

	sqlStr := fmt.Sprintf("select articleid, articlename, author, sortid, intro, size, lastupdate, postdate, allvisit, monthvisit, weekvisit, dayvisit from jieqi_article_article where sortid = ? and display = 0 order by %s desc limit ?", orderBy)

	rows, err := utils.Db.Query(sqlStr, sortID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		a := &model.Article{}
		err := rows.Scan(&a.ArticleID, &a.ArticleName, &a.Author, &a.SortID, &a.Intro, &a.Size, &a.LastUpdate, &a.PostDate, &a.AllVisit, &a.MonthVisit, &a.WeekVisit, &a.DayVisit)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, nil
}

// GetArticlesBySortAndOrderCached 带缓存获取分类小说列表
func GetArticlesBySortAndOrderCached(sortID int, order string, limit int) ([]*model.Article, error) {
	cacheKey := fmt.Sprintf("articles_sort_%d_%s_%d", sortID, order, limit)

	if !utils.IsRedisEnabled() {
		return GetArticlesBySortAndOrder(sortID, order, limit)
	}

	// 尝试从 Redis 缓存获取
	if cached, err := utils.CacheGet(cacheKey); err == nil && cached != "" {
		var articles []*model.Article
		if err := json.Unmarshal([]byte(cached), &articles); err == nil {
			return articles, nil
		}
	}

	articles, err := GetArticlesBySortAndOrder(sortID, order, limit)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(articles); err == nil {
		utils.CacheSet(cacheKey, string(data), 10*time.Minute)
	}
	return articles, nil
}

// GetRankArticles 获取按指定字段排序的小说列表
func GetRankArticles(orderBy string, limit int) ([]*model.Article, error) {
	// 简单的白名单校验
	validFields := map[string]bool{
		"allvisit":   true,
		"monthvisit": true,
		"weekvisit":  true,
		"dayvisit":   true,
		"allvote":    true,
		"monthvote":  true,
		"weekvote":   true,
		"dayvote":    true,
		"size":       true,
		"lastupdate": true,
		"postdate":   true,
	}
	if !validFields[orderBy] {
		orderBy = "allvisit"
	}

	sqlStr := "select articleid, articlename, author, intro, size, lastupdate, sortid, fullflag, imgflag, lastchapterid, lastchapter from jieqi_article_article order by " + orderBy + " desc limit ?"
	rows, err := utils.Db.Query(sqlStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		art := &model.Article{}
		err := rows.Scan(&art.ArticleID, &art.ArticleName, &art.Author, &art.Intro, &art.Size, &art.LastUpdate, &art.SortID, &art.FullFlag, &art.ImgFlag, &art.LastChapterID, &art.LastChapter)
		if err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil
}

// SearchArticles 模糊查询小说列表
func SearchArticles(keyword string, offset, limit int) ([]*model.Article, error) {
	sqlStr := "select articleid, articlename, author, intro, size, lastupdate, sortid, fullflag, imgflag, lastchapterid, lastchapter from jieqi_article_article where articlename like ? or author like ? order by lastupdate desc limit ?, ?"
	rows, err := utils.Db.Query(sqlStr, "%"+keyword+"%", "%"+keyword+"%", offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		art := &model.Article{}
		err := rows.Scan(&art.ArticleID, &art.ArticleName, &art.Author, &art.Intro, &art.Size, &art.LastUpdate, &art.SortID, &art.FullFlag, &art.ImgFlag, &art.LastChapterID, &art.LastChapter)
		if err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil
}

// GetSearchCount 获取模糊查询的结果总数
func GetSearchCount(keyword string) (int, error) {
	sqlStr := "select count(*) from jieqi_article_article where articlename like ? or author like ?"
	var count int
	err := utils.Db.QueryRow(sqlStr, "%"+keyword+"%", "%"+keyword+"%").Scan(&count)
	return count, err
}

// IncArticleVisit 增加文章点击量 (优化版：Redis 缓冲 + 批量回写)
func IncArticleVisit(id int) error {
	const FlushThreshold = 10

	// 如果开启了 Redis，先尝试在 Redis 中缓冲
	if utils.IsRedisEnabled() {
		bufferKey := fmt.Sprintf("article:visit_buffer:%d", id)

		// 1. 增加缓冲区计数
		val, err := utils.CacheIncr(bufferKey)
		if err != nil {
			// 如果 Redis 操作失败，记录错误并降级到直接写库 (这里简化处理，直接返回错误或继续)
			// 为保证数据一致性，如果 Redis 挂了，这里可以选择降级
			utils.LogError("Redis", "Redis Incr failed: %v", err)
		} else {
			// 如果未达到阈值，直接返回，不写库
			if val < FlushThreshold {
				return nil
			}

			// 达到阈值，获取并重置缓冲区
			// 使用 GetSet 原子性地重置为 0 并拿到旧值
			oldValStr, err := utils.CacheGetSet(bufferKey, 0)
			if err != nil {
				return err
			}

			// 解析增加的访问量
			delta := 0
			fmt.Sscanf(oldValStr, "%d", &delta)

			// 如果 delta <= 0 说明可能并发重置了，或者刚重置完
			if delta <= 0 {
				return nil
			}

			// 2. 准备回写数据库
			article, err := GetArticleByID(id)
			if err != nil {
				return err
			}

			now := utils.NowTime()
			lastVisit := article.LastVisit

			// 3. 判断是否需要重置 (按时间)
			resetDay := !utils.IsSameDay(lastVisit, now)
			resetWeek := !utils.IsSameWeek(lastVisit, now)
			resetMonth := !utils.IsSameMonth(lastVisit, now)

			// 4. 构建 SQL
			sqlStr := "update jieqi_article_article set allvisit=allvisit+?, lastvisit=?"
			args := []interface{}{delta, now}

			if resetDay {
				sqlStr += ", dayvisit=?"
				args = append(args, delta)
			} else {
				sqlStr += ", dayvisit=dayvisit+?"
				args = append(args, delta)
			}

			if resetWeek {
				sqlStr += ", weekvisit=?"
				args = append(args, delta)
			} else {
				sqlStr += ", weekvisit=weekvisit+?"
				args = append(args, delta)
			}

			if resetMonth {
				sqlStr += ", monthvisit=?"
				args = append(args, delta)
			} else {
				sqlStr += ", monthvisit=monthvisit+?"
				args = append(args, delta)
			}

			sqlStr += " where articleid=?"
			args = append(args, id)

			_, err = utils.Db.Exec(sqlStr, args...)

			// 5. 假如写入成功，清理相关页面缓存
			if err == nil {
				InvalidateArticleCache(id)
			}

			return err
		}
	}

	// === 降级处理 / 未开启 Redis 的原有逻辑 ===

	// 1. 获取当前点击量信息
	article, err := GetArticleByID(id)
	if err != nil {
		return err
	}

	now := utils.NowTime()
	lastVisit := article.LastVisit

	// 2. 判断是否需要重置
	resetDay := !utils.IsSameDay(lastVisit, now)
	resetWeek := !utils.IsSameWeek(lastVisit, now)
	resetMonth := !utils.IsSameMonth(lastVisit, now)

	// 3. 更新点击量
	sqlStr := "update jieqi_article_article set allvisit=allvisit+1, lastvisit=?"
	args := []interface{}{now}

	if resetDay {
		sqlStr += ", dayvisit=1"
	} else {
		sqlStr += ", dayvisit=dayvisit+1"
	}

	if resetWeek {
		sqlStr += ", weekvisit=1"
	} else {
		sqlStr += ", weekvisit=weekvisit+1"
	}

	if resetMonth {
		sqlStr += ", monthvisit=1"
	} else {
		sqlStr += ", monthvisit=monthvisit+1"
	}

	sqlStr += " where articleid=?"
	args = append(args, id)

	_, err = utils.Db.Exec(sqlStr, args...)

	// 清理缓存
	if err == nil && utils.IsRedisEnabled() {
		InvalidateArticleCache(id)
	}

	return err
}

// GetAllArticlesForSitemap 获取所有文章用于生成 sitemap
// 只返回必要的字段：ArticleID, LastUpdate
func GetAllArticlesForSitemap() ([]*model.Article, error) {
	// 仅选择 display=0 (显示) 的文章
	sqlStr := "SELECT articleid, lastupdate FROM jieqi_article_article WHERE display = 0"
	rows, err := utils.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		art := &model.Article{}
		err := rows.Scan(&art.ArticleID, &art.LastUpdate)
		if err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil
}
