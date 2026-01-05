package dao

import (
	"bookweb/model"
	"bookweb/utils"
)

// GetArticleByID 根据ArticleID获取小说信息
func GetArticleByID(id int) (*model.Article, error) {
	sqlStr := "select articleid, siteid, postdate, lastupdate, articlename, keywords, initial, authorid, author, posterid, poster, agentid, agent, sortid, typeid, intro, notice, setting, lastvolumeid, lastvolume, lastchapterid, lastchapter, chapters, size, lastvisit, dayvisit, weekvisit, monthvisit, allvisit, lastvote, dayvote, weekvote, monthvote, allvote, fullflag, imgflag from jieqi_article_article where articleid = ?"
	row := utils.Db.QueryRow(sqlStr, id)
	art := &model.Article{}
	err := row.Scan(&art.ArticleID, &art.SiteID, &art.PostDate, &art.LastUpdate, &art.ArticleName, &art.Keywords, &art.Initial, &art.AuthorID, &art.Author, &art.PosterID, &art.Poster, &art.AgentID, &art.Agent, &art.SortID, &art.TypeID, &art.Intro, &art.Notice, &art.Setting, &art.LastVolumeID, &art.LastVolume, &art.LastChapterID, &art.LastChapter, &art.Chapters, &art.Size, &art.LastVisit, &art.DayVisit, &art.WeekVisit, &art.MonthVisit, &art.AllVisit, &art.LastVote, &art.DayVote, &art.WeekVote, &art.MonthVote, &art.AllVote, &art.FullFlag, &art.ImgFlag)
	if err != nil {
		return nil, err
	}
	return art, nil
}

// GetChaptersByArticleID 根据ArticleID获取章节列表(不含内容)
func GetChaptersByArticleID(articleID int) ([]*model.Chapter, error) {
	sqlStr := "select chapterid, siteid, articleid, articlename, volumeid, posterid, poster, postdate, lastupdate, chaptername, chapterorder, size, saleprice, salenum, totalcost, isvip, chaptertype, power, display from jieqi_article_chapter where articleid = ? order by chapterorder asc"
	rows, err := utils.Db.Query(sqlStr, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []*model.Chapter
	for rows.Next() {
		ch := &model.Chapter{}
		// 注意：这里scan的字段必须和select语句中的字段一一对应且顺序一致
		// 这里的select语句没有包含 attachment 字段，因为获取列表通常不需要内容
		err := rows.Scan(&ch.ChapterID, &ch.SiteID, &ch.ArticleID, &ch.ArticleName, &ch.VolumeID, &ch.PosterID, &ch.Poster, &ch.PostDate, &ch.LastUpdate, &ch.ChapterName, &ch.ChapterOrder, &ch.Size, &ch.SalePrice, &ch.SaleNum, &ch.TotalCost, &ch.IsVIP, &ch.ChapterType, &ch.Power, &ch.Display)
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
	sqlStr := "select articleid, articlename, author, intro, size, lastupdate, sortid, fullflag, imgflag, lastchapterid, lastchapter from jieqi_article_article order by allvisit desc limit ?"
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

// GetArticlesBySortAndOrder 获取指定分类下按指定字段排序的小说列表
// 如果 sortID 为 0，则获取全部分类，orderBy 为排序字段，limit 为返回数量
// orderBy 可选值：allvisit-总点击量，postdate-发布时间
func GetArticlesBySortAndOrder(sortID int, orderBy string, limit int) ([]*model.Article, error) {
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

	sqlStr := "select articleid, articlename, author, intro, size, lastupdate, sortid, fullflag, imgflag, lastchapterid, lastchapter from jieqi_article_article where sortid = ? order by " + orderBy + " desc limit ?"
	rows, err := utils.Db.Query(sqlStr, sortID, limit)
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

// IncArticleVisit 增加文章点击量
func IncArticleVisit(id int) error {
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
	// 注意：这里需要处理并发更新问题，但在高并发场景下，通常会使用 Redis 计数然后定期回写 DB
	// 这里采用简单的 SQL 更新
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
	return err
}
